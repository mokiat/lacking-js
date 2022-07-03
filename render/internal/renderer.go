package internal

import (
	"fmt"

	"github.com/mokiat/lacking/render"
	"github.com/mokiat/wasmgl"
)

func NewRenderer() *Renderer {
	result := &Renderer{
		framebuffer:   DefaultFramebuffer,
		isDirty:       true,
		isInvalidated: true,
		desiredState: &State{
			CullTest:                    false,
			CullFace:                    wasmgl.BACK,
			FrontFace:                   wasmgl.CCW,
			DepthTest:                   false,
			DepthMask:                   true,
			DepthComparison:             wasmgl.LESS,
			StencilTest:                 false,
			StencilOpStencilFailFront:   wasmgl.KEEP,
			StencilOpDepthFailFront:     wasmgl.KEEP,
			StencilOpPassFront:          wasmgl.KEEP,
			StencilOpStencilFailBack:    wasmgl.KEEP,
			StencilOpDepthFailBack:      wasmgl.KEEP,
			StencilOpPassBack:           wasmgl.KEEP,
			StencilComparisonFuncFront:  wasmgl.ALWAYS,
			StencilComparisonRefFront:   0x00,
			StencilComparisonMaskFront:  0xFF,
			StencilComparisonFuncBack:   wasmgl.ALWAYS,
			StencilComparisonRefBack:    0x00,
			StencilComparisonMaskBack:   0xFF,
			StencilMaskFront:            0xFF,
			StencilMaskBack:             0xFF,
			ColorMask:                   render.ColorMaskTrue,
			Blending:                    false,
			BlendModeRGB:                wasmgl.FUNC_ADD,
			BlendModeAlpha:              wasmgl.FUNC_ADD,
			BlendSourceFactorRGB:        wasmgl.ONE,
			BlendDestinationFactorRGB:   wasmgl.ZERO,
			BlendSourceFactorAlpha:      wasmgl.ONE,
			BlendDestinationFactorAlpha: wasmgl.ZERO,
		},
		actualState: &State{},
	}
	result.Invalidate()
	return result
}

type Renderer struct {
	framebuffer           *Framebuffer
	invalidateAttachments []int
	topology              int
	indexType             int

	isDirty       bool
	isInvalidated bool
	desiredState  *State
	actualState   *State
}

func (r *Renderer) BeginRenderPass(info render.RenderPassInfo) {
	r.validateState()

	r.framebuffer = info.Framebuffer.(*Framebuffer)
	isDefaultFramebuffer := r.framebuffer == DefaultFramebuffer

	wasmgl.BindFramebuffer(wasmgl.FRAMEBUFFER, r.framebuffer.raw)
	wasmgl.Viewport(
		info.Viewport.X,
		info.Viewport.Y,
		info.Viewport.Width,
		info.Viewport.Height,
	)

	oldColorMask := r.actualState.ColorMask

	var colorMaskChanged bool
	for i, attachment := range info.Colors {
		if r.framebuffer.activeDrawBuffers[i] && (attachment.LoadOp == render.LoadOperationClear) {
			if !colorMaskChanged {
				r.executeCommandColorWrite(CommandColorWrite{
					Mask: render.ColorMaskTrue,
				})
				r.validateColorMask(false)
				colorMaskChanged = true
			}
			wasmgl.ClearBufferfv(wasmgl.COLOR, i, attachment.ClearValue[:])
		}
	}
	if colorMaskChanged {
		r.executeCommandColorWrite(CommandColorWrite{
			Mask: oldColorMask,
		})
	}

	oldDepthMask := r.actualState.DepthMask
	oldStencilMaskFront := r.actualState.StencilMaskFront
	oldStencilMaskBack := r.actualState.StencilMaskBack

	clearDepth := info.DepthLoadOp == render.LoadOperationClear
	clearStencil := info.StencilLoadOp == render.LoadOperationClear

	if clearDepth && clearStencil {
		r.executeCommandDepthWrite(CommandDepthWrite{
			Enabled: true,
		})
		r.validateDepthMask(false)
		r.executeCommandStencilMask(CommandStencilMask{
			Face: wasmgl.FRONT_AND_BACK,
			Mask: 0xFF,
		})
		r.validateStencilMask(false)
		depthValue := info.DepthClearValue
		stencilValue := int32(info.StencilClearValue)
		wasmgl.ClearBufferfi(wasmgl.DEPTH_STENCIL, 0, depthValue, stencilValue)
	} else {
		if clearDepth {
			r.executeCommandDepthWrite(CommandDepthWrite{
				Enabled: true,
			})
			r.validateDepthMask(false)
			depthValues := [1]float32{info.DepthClearValue}
			wasmgl.ClearBufferfv(wasmgl.DEPTH, 0, depthValues[:])
		}
		if clearStencil {
			r.executeCommandStencilMask(CommandStencilMask{
				Face: wasmgl.FRONT_AND_BACK,
				Mask: 0xFF,
			})
			r.validateStencilMask(false)
			stencilValues := [1]int32{int32(info.StencilClearValue)}
			wasmgl.ClearBufferiv(wasmgl.STENCIL, 0, stencilValues[:])
		}
	}

	r.invalidateAttachments = r.invalidateAttachments[:0]

	for i, attachment := range info.Colors {
		if r.framebuffer.activeDrawBuffers[i] && (attachment.StoreOp == render.StoreOperationDontCare) {
			if isDefaultFramebuffer {
				if i == 0 {
					r.invalidateAttachments = append(r.invalidateAttachments, wasmgl.COLOR)
				}
			} else {
				r.invalidateAttachments = append(r.invalidateAttachments, wasmgl.COLOR_ATTACHMENT0+i)
			}
		}
	}

	invalidateDepth := info.DepthStoreOp == render.StoreOperationDontCare
	invalidateStencil := info.StencilStoreOp == render.StoreOperationDontCare

	if invalidateDepth && invalidateStencil && !isDefaultFramebuffer {
		r.invalidateAttachments = append(r.invalidateAttachments, wasmgl.DEPTH_STENCIL_ATTACHMENT)
	} else {
		if invalidateDepth {
			if isDefaultFramebuffer {
				r.invalidateAttachments = append(r.invalidateAttachments, wasmgl.DEPTH)
			} else {
				r.invalidateAttachments = append(r.invalidateAttachments, wasmgl.DEPTH_ATTACHMENT)
			}
		}
		if invalidateStencil {
			if isDefaultFramebuffer {
				r.invalidateAttachments = append(r.invalidateAttachments, wasmgl.STENCIL)
			} else {
				r.invalidateAttachments = append(r.invalidateAttachments, wasmgl.STENCIL_ATTACHMENT)
			}
		}
	}

	r.executeCommandDepthWrite(CommandDepthWrite{
		Enabled: oldDepthMask,
	})
	r.executeCommandStencilMask(CommandStencilMask{
		Face: wasmgl.FRONT,
		Mask: uint32(oldStencilMaskFront),
	})
	r.executeCommandStencilMask(CommandStencilMask{
		Face: wasmgl.BACK,
		Mask: uint32(oldStencilMaskBack),
	})
}

func (r *Renderer) EndRenderPass() {
	if len(r.invalidateAttachments) > 0 {
		wasmgl.InvalidateFramebuffer(wasmgl.FRAMEBUFFER, r.invalidateAttachments)
	}
	r.framebuffer = DefaultFramebuffer
}

func (r *Renderer) Invalidate() {
	r.isDirty = true
	r.isInvalidated = true
}

func (r *Renderer) BindPipeline(pipeline render.Pipeline) {
	intPipeline := pipeline.(*Pipeline)
	r.executeCommandBindPipeline(CommandBindPipeline{
		ProgramID:        intPipeline.ProgramID,
		Topology:         intPipeline.Topology,
		CullTest:         intPipeline.CullTest,
		FrontFace:        intPipeline.FrontFace,
		DepthTest:        intPipeline.DepthTest,
		DepthWrite:       intPipeline.DepthWrite,
		DepthComparison:  intPipeline.DepthComparison,
		StencilTest:      intPipeline.StencilTest,
		StencilOpFront:   intPipeline.StencilOpFront,
		StencilOpBack:    intPipeline.StencilOpBack,
		StencilFuncFront: intPipeline.StencilFuncFront,
		StencilFuncBack:  intPipeline.StencilFuncBack,
		StencilMaskFront: intPipeline.StencilMaskFront,
		StencilMaskBack:  intPipeline.StencilMaskBack,
		ColorWrite:       intPipeline.ColorWrite,
		BlendEnabled:     intPipeline.BlendEnabled,
		BlendColor:       intPipeline.BlendColor,
		BlendEquation:    intPipeline.BlendEquation,
		BlendFunc:        intPipeline.BlendFunc,
		VertexArray:      intPipeline.VertexArray,
	})
}

func (r *Renderer) Uniform1f(location render.UniformLocation, value float32) {
	intLocation := location.(*UniformLocation)
	r.executeCommandUniform1f(CommandUniform1f{
		Location: intLocation.id,
		Value:    value,
	})
}

func (r *Renderer) Uniform1i(location render.UniformLocation, value int) {
	intLocation := location.(*UniformLocation)
	r.executeCommandUniform1i(CommandUniform1i{
		Location: intLocation.id,
		Value:    int32(value),
	})
}

func (r *Renderer) Uniform3f(location render.UniformLocation, values [3]float32) {
	intLocation := location.(*UniformLocation)
	r.executeCommandUniform3f(CommandUniform3f{
		Location: intLocation.id,
		Values:   values,
	})
}

func (r *Renderer) Uniform4f(location render.UniformLocation, values [4]float32) {
	intLocation := location.(*UniformLocation)
	r.executeCommandUniform4f(CommandUniform4f{
		Location: intLocation.id,
		Values:   values,
	})
}

func (r *Renderer) UniformMatrix4f(location render.UniformLocation, values [16]float32) {
	intLocation := location.(*UniformLocation)
	r.executeCommandUniformMatrix4f(CommandUniformMatrix4f{
		Location: intLocation.id,
		Values:   values,
	})
}

func (r *Renderer) UniformBufferUnit(index int, buffer render.Buffer) {
	r.executeCommandUniformBufferUnit(CommandUniformBufferUnit{
		Index:    uint32(index),
		BufferID: buffer.(*Buffer).id,
	})
}

func (r *Renderer) UniformBufferUnitRange(index int, buffer render.Buffer, offset, size int) {
	r.executeCommandUniformBufferUnitRange(CommandUniformBufferUnitRange{
		Index:    uint32(index),
		BufferID: buffer.(*Buffer).id,
		Offset:   uint32(offset),
		Size:     uint32(size),
	})
}

func (r *Renderer) TextureUnit(index int, texture render.Texture) {
	r.executeCommandTextureUnit(CommandTextureUnit{
		Index:     uint32(index),
		TextureID: texture.(*Texture).id,
	})
}

func (r *Renderer) Draw(vertexOffset, vertexCount, instanceCount int) {
	r.executeCommandDraw(CommandDraw{
		VertexOffset:  int32(vertexOffset),
		VertexCount:   int32(vertexCount),
		InstanceCount: int32(instanceCount),
	})
}

func (r *Renderer) DrawIndexed(indexOffset, indexCount, instanceCount int) {
	r.executeCommandDrawIndexed(CommandDrawIndexed{
		IndexOffset:   int32(indexOffset),
		IndexCount:    int32(indexCount),
		InstanceCount: int32(instanceCount),
	})
}

func (r *Renderer) CopyContentToTexture(info render.CopyContentToTextureInfo) {
	intTexture := info.Texture.(*Texture)
	wasmgl.BindTexture(intTexture.kind, intTexture.raw)
	wasmgl.CopyTexSubImage2D(
		intTexture.kind,
		info.TextureLevel,
		info.TextureX,
		info.TextureY,
		info.FramebufferX,
		info.FramebufferY,
		info.Width,
		info.Height,
	)
	if info.GenerateMipmaps {
		wasmgl.GenerateMipmap(intTexture.kind)
	}
}

func (r *Renderer) SubmitQueue(queue *CommandQueue) {
	for MoreCommands(queue) {
		header := PopCommand[CommandHeader](queue)
		switch header.Kind {
		case CommandKindBindPipeline:
			command := PopCommand[CommandBindPipeline](queue)
			r.executeCommandBindPipeline(command)
		case CommandKindTopology:
			command := PopCommand[CommandTopology](queue)
			r.executeCommandTopology(command)
		case CommandKindCullTest:
			command := PopCommand[CommandCullTest](queue)
			r.executeCommandCullTest(command)
		case CommandKindFrontFace:
			command := PopCommand[CommandFrontFace](queue)
			r.executeCommandFrontFace(command)
		case CommandKindDepthTest:
			command := PopCommand[CommandDepthTest](queue)
			r.executeCommandDepthTest(command)
		case CommandKindDepthWrite:
			command := PopCommand[CommandDepthWrite](queue)
			r.executeCommandDepthWrite(command)
		case CommandKindDepthComparison:
			command := PopCommand[CommandDepthComparison](queue)
			r.executeCommandDepthComparison(command)
		case CommandKindUniform1f:
			command := PopCommand[CommandUniform1f](queue)
			r.executeCommandUniform1f(command)
		case CommandKindUniform1i:
			command := PopCommand[CommandUniform1i](queue)
			r.executeCommandUniform1i(command)
		case CommandKindUniform3f:
			command := PopCommand[CommandUniform3f](queue)
			r.executeCommandUniform3f(command)
		case CommandKindUniform4f:
			command := PopCommand[CommandUniform4f](queue)
			r.executeCommandUniform4f(command)
		case CommandKindUniformMatrix4f:
			command := PopCommand[CommandUniformMatrix4f](queue)
			r.executeCommandUniformMatrix4f(command)
		case CommandKindUniformBufferUnit:
			command := PopCommand[CommandUniformBufferUnit](queue)
			r.executeCommandUniformBufferUnit(command)
		case CommandKindUniformBufferUnitRange:
			command := PopCommand[CommandUniformBufferUnitRange](queue)
			r.executeCommandUniformBufferUnitRange(command)
		case CommandKindTextureUnit:
			command := PopCommand[CommandTextureUnit](queue)
			r.executeCommandTextureUnit(command)
		case CommandKindDraw:
			command := PopCommand[CommandDraw](queue)
			r.executeCommandDraw(command)
		case CommandKindDrawIndexed:
			command := PopCommand[CommandDrawIndexed](queue)
			r.executeCommandDrawIndexed(command)
		case CommandKindCopyContentToBuffer:
			command := PopCommand[CommandCopyContentToBuffer](queue)
			r.executeCommandCopyContentToBuffer(command)
		default:
			panic(fmt.Errorf("unknown command kind: %v", header.Kind))
		}
	}
	queue.Reset()
}

func (r *Renderer) executeCommandBindPipeline(command CommandBindPipeline) {
	program := programs.Get(command.ProgramID)
	wasmgl.UseProgram(program.raw)
	r.executeCommandTopology(command.Topology)
	r.executeCommandCullTest(command.CullTest)
	r.executeCommandFrontFace(command.FrontFace)
	r.executeCommandDepthTest(command.DepthTest)
	r.executeCommandDepthWrite(command.DepthWrite)
	if command.DepthTest.Enabled {
		r.executeCommandDepthComparison(command.DepthComparison)
	}
	r.executeCommandStencilTest(command.StencilTest)
	if command.StencilTest.Enabled {
		r.executeCommandStencilOperation(command.StencilOpFront)
		r.executeCommandStencilOperation(command.StencilOpBack)
		r.executeCommandStencilFunc(command.StencilFuncFront)
		r.executeCommandStencilFunc(command.StencilFuncBack)
		r.executeCommandStencilMask(command.StencilMaskFront)
		r.executeCommandStencilMask(command.StencilMaskBack)
	}
	r.executeCommandColorWrite(command.ColorWrite)
	r.desiredState.Blending = command.BlendEnabled
	r.isDirty = true
	if command.BlendEnabled {
		r.executeCommandBlendColor(command.BlendColor)
		r.executeCommandBlendEquation(command.BlendEquation)
		r.executeCommandBlendFunc(command.BlendFunc)
	}
	r.executeCommandBindVertexArray(command.VertexArray)
}

func (r *Renderer) executeCommandTopology(command CommandTopology) {
	r.topology = int(command.Topology)
}

func (r *Renderer) executeCommandCullTest(command CommandCullTest) {
	r.desiredState.CullTest = command.Enabled
	if command.Enabled {
		r.desiredState.CullFace = int(command.Face)
	}
	r.isDirty = true
}

func (r *Renderer) executeCommandFrontFace(command CommandFrontFace) {
	r.desiredState.FrontFace = int(command.Orientation)
	r.isDirty = true
}

func (r *Renderer) executeCommandDepthTest(command CommandDepthTest) {
	r.desiredState.DepthTest = command.Enabled
	r.isDirty = true
}

func (r *Renderer) executeCommandDepthWrite(command CommandDepthWrite) {
	r.desiredState.DepthMask = command.Enabled
	r.isDirty = true
}

func (r *Renderer) executeCommandDepthComparison(command CommandDepthComparison) {
	r.desiredState.DepthComparison = int(command.Mode)
	r.isDirty = true
}

func (r *Renderer) executeCommandStencilTest(command CommandStencilTest) {
	r.desiredState.StencilTest = command.Enabled
	r.isDirty = true
}

func (r *Renderer) executeCommandStencilOperation(command CommandStencilOperation) {
	if int(command.Face) == wasmgl.FRONT || int(command.Face) == wasmgl.FRONT_AND_BACK {
		r.desiredState.StencilOpStencilFailFront = int(command.StencilFail)
		r.desiredState.StencilOpDepthFailFront = int(command.DepthFail)
		r.desiredState.StencilOpPassFront = int(command.Pass)
	}
	if int(command.Face) == wasmgl.BACK || int(command.Face) == wasmgl.FRONT_AND_BACK {
		r.desiredState.StencilOpStencilFailBack = int(command.StencilFail)
		r.desiredState.StencilOpDepthFailBack = int(command.DepthFail)
		r.desiredState.StencilOpPassBack = int(command.Pass)
	}
	r.isDirty = true
}

func (r *Renderer) executeCommandStencilFunc(command CommandStencilFunc) {
	if int(command.Face) == wasmgl.FRONT || int(command.Face) == wasmgl.FRONT_AND_BACK {
		r.desiredState.StencilComparisonFuncFront = int(command.Func)
		r.desiredState.StencilComparisonRefFront = int(command.Ref)
		r.desiredState.StencilComparisonMaskFront = int(command.Mask)
	}
	if int(command.Face) == wasmgl.BACK || int(command.Face) == wasmgl.FRONT_AND_BACK {
		r.desiredState.StencilComparisonFuncBack = int(command.Func)
		r.desiredState.StencilComparisonRefBack = int(command.Ref)
		r.desiredState.StencilComparisonMaskBack = int(command.Mask)
	}
	r.isDirty = true
}

func (r *Renderer) executeCommandStencilMask(command CommandStencilMask) {
	if int(command.Face) == wasmgl.FRONT || int(command.Face) == wasmgl.FRONT_AND_BACK {
		r.desiredState.StencilMaskFront = int(command.Mask)
	}
	if int(command.Face) == wasmgl.BACK || int(command.Face) == wasmgl.FRONT_AND_BACK {
		r.desiredState.StencilMaskBack = int(command.Mask)
	}
	r.isDirty = true
}

func (r *Renderer) executeCommandColorWrite(command CommandColorWrite) {
	r.desiredState.ColorMask = command.Mask
	r.isDirty = true
}

func (r *Renderer) executeCommandBlendColor(command CommandBlendColor) {
	r.desiredState.BlendColor = command.Color
	r.isDirty = true
}

func (r *Renderer) executeCommandBlendEquation(command CommandBlendEquation) {
	r.desiredState.BlendModeRGB = int(command.ModeRGB)
	r.desiredState.BlendModeAlpha = int(command.ModeAlpha)
	r.isDirty = true
}

func (r *Renderer) executeCommandBlendFunc(command CommandBlendFunc) {
	r.desiredState.BlendSourceFactorRGB = int(command.SourceFactorRGB)
	r.desiredState.BlendDestinationFactorRGB = int(command.DestinationFactorRGB)
	r.desiredState.BlendSourceFactorAlpha = int(command.SourceFactorAlpha)
	r.desiredState.BlendDestinationFactorAlpha = int(command.DestinationFactorAlpha)
	r.isDirty = true
}

func (r *Renderer) executeCommandBindVertexArray(command CommandBindVertexArray) {
	vertexArray := vertexArrays.Get(command.VertexArrayID)
	wasmgl.BindVertexArray(vertexArray.raw)
	r.indexType = int(command.IndexFormat)
}

func (r *Renderer) executeCommandUniform1f(command CommandUniform1f) {
	location := locations.Get(uint32(command.Location))
	wasmgl.Uniform1f(
		location.raw,
		command.Value,
	)
}

func (r *Renderer) executeCommandUniform1i(command CommandUniform1i) {
	location := locations.Get(uint32(command.Location))
	wasmgl.Uniform1i(
		location.raw,
		int(command.Value),
	)
}

func (r *Renderer) executeCommandUniform3f(command CommandUniform3f) {
	location := locations.Get(uint32(command.Location))
	wasmgl.Uniform3f(
		location.raw,
		command.Values[0],
		command.Values[1],
		command.Values[2],
	)
}

func (r *Renderer) executeCommandUniform4f(command CommandUniform4f) {
	location := locations.Get(uint32(command.Location))
	wasmgl.Uniform4f(
		location.raw,
		command.Values[0],
		command.Values[1],
		command.Values[2],
		command.Values[3],
	)
}

func (r *Renderer) executeCommandUniformMatrix4f(command CommandUniformMatrix4f) {
	location := locations.Get(uint32(command.Location))
	wasmgl.UniformMatrix4fv(
		location.raw,
		false,
		command.Values[:],
	)
}

func (r *Renderer) executeCommandUniformBufferUnit(command CommandUniformBufferUnit) {
	buffer := buffers.Get(command.BufferID)
	wasmgl.BindBufferBase(
		wasmgl.UNIFORM_BUFFER,
		int(command.Index),
		buffer.raw,
	)
}

func (r *Renderer) executeCommandUniformBufferUnitRange(command CommandUniformBufferUnitRange) {
	buffer := buffers.Get(command.BufferID)
	wasmgl.BindBufferRange(
		wasmgl.UNIFORM_BUFFER,
		int(command.Index),
		buffer.raw,
		int(command.Offset),
		int(command.Size),
	)
}

func (r *Renderer) executeCommandTextureUnit(command CommandTextureUnit) {
	texture := textures.Get(command.TextureID)
	wasmgl.ActiveTexture(wasmgl.TEXTURE0 + int(command.Index))
	wasmgl.BindTexture(texture.kind, texture.raw)
}

func (r *Renderer) executeCommandDraw(command CommandDraw) {
	r.validateState()
	wasmgl.DrawArraysInstanced(
		r.topology,
		int(command.VertexOffset),
		int(command.VertexCount),
		int(command.InstanceCount),
	)
}

func (r *Renderer) executeCommandDrawIndexed(command CommandDrawIndexed) {
	r.validateState()
	wasmgl.DrawElementsInstanced(
		r.topology,
		int(command.IndexCount),
		r.indexType,
		int(command.IndexOffset),
		int(command.InstanceCount),
	)
}

func (r *Renderer) executeCommandCopyContentToBuffer(command CommandCopyContentToBuffer) {
	buffer := buffers.Get(command.BufferID)
	wasmgl.BindBuffer(
		buffer.kind,
		buffer.raw,
	)
	wasmgl.ReadPixels(
		int(command.X),
		int(command.Y),
		int(command.Width),
		int(command.Height),
		int(command.Format),
		int(command.XType),
		int(command.BufferOffset),
	)
	wasmgl.BindBuffer(
		buffer.kind,
		wasmgl.NilBuffer,
	)
}

func (r *Renderer) validateState() {
	if r.isDirty || r.isInvalidated {
		forcedUpdate := r.isInvalidated
		r.validateCullTest(forcedUpdate)
		r.validateCullFace(forcedUpdate)
		r.validateFrontFace(forcedUpdate)
		r.validateDepthTest(forcedUpdate)
		r.validateDepthMask(forcedUpdate)
		r.validateDepthComparison(forcedUpdate)
		r.validateStencilTest(forcedUpdate)
		r.validateStencilOperation(forcedUpdate)
		r.validateStencilComparison(forcedUpdate)
		r.validateStencilMask(forcedUpdate)
		r.validateColorMask(forcedUpdate)
		r.validateBlending(forcedUpdate)
		r.validateBlendEquation(forcedUpdate)
		r.validateBlendFunc(forcedUpdate)
	}
	r.isDirty = false
	r.isInvalidated = false
}

func (r *Renderer) validateCullTest(forcedUpdate bool) {
	needsUpdate := forcedUpdate ||
		(r.actualState.CullTest != r.desiredState.CullTest)

	if needsUpdate {
		r.actualState.CullTest = r.desiredState.CullTest
		if r.actualState.CullTest {
			wasmgl.Enable(wasmgl.CULL_FACE)
		} else {
			wasmgl.Disable(wasmgl.CULL_FACE)
		}
	}
}

func (r *Renderer) validateCullFace(forcedUpdate bool) {
	needsUpdate := forcedUpdate ||
		(r.actualState.CullFace != r.desiredState.CullFace)

	if needsUpdate {
		r.actualState.CullFace = r.desiredState.CullFace
		wasmgl.CullFace(r.actualState.CullFace)
	}
}

func (r *Renderer) validateFrontFace(forcedUpdate bool) {
	needsUpdate := forcedUpdate ||
		(r.actualState.FrontFace != r.desiredState.FrontFace)

	if needsUpdate {
		r.actualState.FrontFace = r.desiredState.FrontFace
		wasmgl.FrontFace(r.actualState.FrontFace)
	}
}

func (r *Renderer) validateDepthTest(forcedUpdate bool) {
	needsUpdate := forcedUpdate ||
		(r.actualState.DepthTest != r.desiredState.DepthTest)

	if needsUpdate {
		r.actualState.DepthTest = r.desiredState.DepthTest
		if r.actualState.DepthTest {
			wasmgl.Enable(wasmgl.DEPTH_TEST)
		} else {
			wasmgl.Disable(wasmgl.DEPTH_TEST)
		}
	}
}

func (r *Renderer) validateDepthMask(forcedUpdate bool) {
	needsUpdate := forcedUpdate ||
		(r.actualState.DepthMask != r.desiredState.DepthMask)

	if needsUpdate {
		r.actualState.DepthMask = r.desiredState.DepthMask
		wasmgl.DepthMask(r.actualState.DepthMask)
	}
}

func (r *Renderer) validateDepthComparison(forcedUpdate bool) {
	needsUpdate := forcedUpdate ||
		(r.actualState.DepthComparison != r.desiredState.DepthComparison)

	if needsUpdate {
		r.actualState.DepthComparison = r.desiredState.DepthComparison
		wasmgl.DepthFunc(r.actualState.DepthComparison)
	}
}

func (r *Renderer) validateStencilTest(forcedUpdate bool) {
	needsUpdate := forcedUpdate ||
		(r.actualState.StencilTest != r.desiredState.StencilTest)

	if needsUpdate {
		r.actualState.StencilTest = r.desiredState.StencilTest
		if r.actualState.StencilTest {
			wasmgl.Enable(wasmgl.STENCIL_TEST)
		} else {
			wasmgl.Disable(wasmgl.STENCIL_TEST)
		}
	}
}

func (r *Renderer) validateStencilOperation(forcedUpdate bool) {
	frontNeedsUpdate := forcedUpdate ||
		(r.actualState.StencilOpStencilFailFront != r.desiredState.StencilOpStencilFailFront) ||
		(r.actualState.StencilOpDepthFailFront != r.desiredState.StencilOpDepthFailFront) ||
		(r.actualState.StencilOpPassFront != r.desiredState.StencilOpPassFront)

	backNeedsUpdate := forcedUpdate ||
		(r.actualState.StencilOpStencilFailBack != r.desiredState.StencilOpStencilFailBack) ||
		(r.actualState.StencilOpDepthFailBack != r.desiredState.StencilOpDepthFailBack) ||
		(r.actualState.StencilOpPassBack != r.desiredState.StencilOpPassBack)

	if frontNeedsUpdate {
		r.actualState.StencilOpStencilFailFront = r.desiredState.StencilOpStencilFailFront
		r.actualState.StencilOpDepthFailFront = r.desiredState.StencilOpDepthFailFront
		r.actualState.StencilOpPassFront = r.desiredState.StencilOpPassFront
	}

	if backNeedsUpdate {
		r.actualState.StencilOpStencilFailBack = r.desiredState.StencilOpStencilFailBack
		r.actualState.StencilOpDepthFailBack = r.desiredState.StencilOpDepthFailBack
		r.actualState.StencilOpPassBack = r.desiredState.StencilOpPassBack
	}

	frontEqualsBack := (r.desiredState.StencilOpStencilFailFront == r.desiredState.StencilOpStencilFailBack) &&
		(r.desiredState.StencilOpDepthFailFront == r.desiredState.StencilOpDepthFailBack) &&
		(r.desiredState.StencilOpPassFront == r.desiredState.StencilOpPassBack)

	if frontNeedsUpdate && backNeedsUpdate && frontEqualsBack {
		wasmgl.StencilOpSeparate(
			wasmgl.FRONT_AND_BACK,
			r.actualState.StencilOpStencilFailFront,
			r.actualState.StencilOpDepthFailFront,
			r.actualState.StencilOpPassFront,
		)
	} else {
		if frontNeedsUpdate {
			wasmgl.StencilOpSeparate(
				wasmgl.FRONT,
				r.actualState.StencilOpStencilFailFront,
				r.actualState.StencilOpDepthFailFront,
				r.actualState.StencilOpPassFront,
			)
		}
		if backNeedsUpdate {
			wasmgl.StencilOpSeparate(
				wasmgl.BACK,
				r.actualState.StencilOpStencilFailBack,
				r.actualState.StencilOpDepthFailBack,
				r.actualState.StencilOpPassBack,
			)
		}
	}
}

func (r *Renderer) validateStencilComparison(forcedUpdate bool) {
	frontNeedsUpdate := forcedUpdate ||
		(r.actualState.StencilComparisonFuncFront != r.desiredState.StencilComparisonFuncFront) ||
		(r.actualState.StencilComparisonRefFront != r.desiredState.StencilComparisonRefFront) ||
		(r.actualState.StencilComparisonMaskFront != r.desiredState.StencilComparisonMaskFront)

	backNeedsUpdate := forcedUpdate ||
		(r.actualState.StencilComparisonFuncBack != r.desiredState.StencilComparisonFuncBack) ||
		(r.actualState.StencilComparisonRefBack != r.desiredState.StencilComparisonRefBack) ||
		(r.actualState.StencilComparisonMaskBack != r.desiredState.StencilComparisonMaskBack)

	if frontNeedsUpdate {
		r.actualState.StencilComparisonFuncFront = r.desiredState.StencilComparisonFuncFront
		r.actualState.StencilComparisonRefFront = r.desiredState.StencilComparisonRefFront
		r.actualState.StencilComparisonMaskFront = r.desiredState.StencilComparisonMaskFront
	}

	if backNeedsUpdate {
		r.actualState.StencilComparisonFuncBack = r.desiredState.StencilComparisonFuncBack
		r.actualState.StencilComparisonRefBack = r.desiredState.StencilComparisonRefBack
		r.actualState.StencilComparisonMaskBack = r.desiredState.StencilComparisonMaskBack
	}

	frontEqualsBack := (r.desiredState.StencilComparisonFuncFront == r.desiredState.StencilComparisonFuncBack) &&
		(r.desiredState.StencilComparisonRefFront == r.desiredState.StencilComparisonRefBack) &&
		(r.desiredState.StencilComparisonMaskFront == r.desiredState.StencilComparisonMaskBack)

	if frontNeedsUpdate && backNeedsUpdate && frontEqualsBack {
		wasmgl.StencilFuncSeparate(
			wasmgl.FRONT_AND_BACK,
			r.actualState.StencilComparisonFuncFront,
			r.actualState.StencilComparisonRefFront,
			r.actualState.StencilComparisonMaskFront,
		)
	} else {
		if frontNeedsUpdate {
			wasmgl.StencilFuncSeparate(
				wasmgl.FRONT,
				r.actualState.StencilComparisonFuncFront,
				r.actualState.StencilComparisonRefFront,
				r.actualState.StencilComparisonMaskFront,
			)
		}
		if backNeedsUpdate {
			wasmgl.StencilFuncSeparate(
				wasmgl.BACK,
				r.actualState.StencilComparisonFuncBack,
				r.actualState.StencilComparisonRefBack,
				r.actualState.StencilComparisonMaskBack,
			)
		}
	}
}

func (r *Renderer) validateStencilMask(forcedUpdate bool) {
	frontNeedsUpdate := forcedUpdate ||
		(r.actualState.StencilMaskFront != r.desiredState.StencilMaskFront)

	backNeedsUpdate := forcedUpdate ||
		(r.actualState.StencilMaskBack != r.desiredState.StencilMaskBack)

	if frontNeedsUpdate {
		r.actualState.StencilMaskFront = r.desiredState.StencilMaskFront
	}
	if backNeedsUpdate {
		r.actualState.StencilMaskBack = r.desiredState.StencilMaskBack
	}

	frontEqualsBack := (r.desiredState.StencilMaskFront == r.desiredState.StencilMaskBack)

	if frontNeedsUpdate && backNeedsUpdate && frontEqualsBack {
		wasmgl.StencilMaskSeparate(
			wasmgl.FRONT_AND_BACK,
			r.actualState.StencilMaskFront,
		)
	} else {
		if frontNeedsUpdate {
			wasmgl.StencilMaskSeparate(
				wasmgl.FRONT,
				r.actualState.StencilMaskFront,
			)
		}
		if backNeedsUpdate {
			wasmgl.StencilMaskSeparate(
				wasmgl.BACK,
				r.actualState.StencilMaskBack,
			)
		}
	}
}

func (r *Renderer) validateColorMask(forcedUpdate bool) {
	needsUpdate := forcedUpdate ||
		(r.actualState.ColorMask != r.desiredState.ColorMask)

	if needsUpdate {
		r.actualState.ColorMask = r.desiredState.ColorMask
		wasmgl.ColorMask(
			r.actualState.ColorMask[0],
			r.actualState.ColorMask[1],
			r.actualState.ColorMask[2],
			r.actualState.ColorMask[3],
		)
	}
}

func (r *Renderer) validateBlending(forcedUpdate bool) {
	needsUpdate := forcedUpdate ||
		(r.actualState.Blending != r.desiredState.Blending)

	if needsUpdate {
		r.actualState.Blending = r.desiredState.Blending
		if r.actualState.Blending {
			wasmgl.Enable(wasmgl.BLEND)
		} else {
			wasmgl.Disable(wasmgl.BLEND)
		}
	}
}

func (r *Renderer) validateBlendColor(forcedUpdate bool) {
	needsUpdate := forcedUpdate ||
		(r.actualState.BlendColor != r.desiredState.BlendColor)

	if needsUpdate {
		r.actualState.BlendColor = r.desiredState.BlendColor
		wasmgl.BlendColor(
			r.actualState.BlendColor[0],
			r.actualState.BlendColor[1],
			r.actualState.BlendColor[2],
			r.actualState.BlendColor[3],
		)
	}
}

func (r *Renderer) validateBlendEquation(forcedUpdate bool) {
	needsUpdate := forcedUpdate ||
		(r.actualState.BlendModeRGB != r.desiredState.BlendModeRGB) ||
		(r.actualState.BlendModeAlpha != r.desiredState.BlendModeAlpha)

	if needsUpdate {
		r.actualState.BlendModeRGB = r.desiredState.BlendModeRGB
		r.actualState.BlendModeAlpha = r.desiredState.BlendModeAlpha
		wasmgl.BlendEquationSeparate(
			r.actualState.BlendModeRGB,
			r.actualState.BlendModeAlpha,
		)
	}
}

func (r *Renderer) validateBlendFunc(forcedUpdate bool) {
	needsUpdate := forcedUpdate ||
		(r.actualState.BlendSourceFactorRGB != r.desiredState.BlendSourceFactorRGB) ||
		(r.actualState.BlendDestinationFactorRGB != r.desiredState.BlendDestinationFactorRGB) ||
		(r.actualState.BlendSourceFactorAlpha != r.desiredState.BlendSourceFactorAlpha) ||
		(r.actualState.BlendDestinationFactorAlpha != r.desiredState.BlendDestinationFactorAlpha)

	if needsUpdate {
		r.actualState.BlendSourceFactorRGB = r.desiredState.BlendSourceFactorRGB
		r.actualState.BlendDestinationFactorRGB = r.desiredState.BlendDestinationFactorRGB
		r.actualState.BlendSourceFactorAlpha = r.desiredState.BlendSourceFactorAlpha
		r.actualState.BlendDestinationFactorAlpha = r.desiredState.BlendDestinationFactorAlpha
		wasmgl.BlendFuncSeparate(
			r.actualState.BlendSourceFactorRGB,
			r.actualState.BlendDestinationFactorRGB,
			r.actualState.BlendSourceFactorAlpha,
			r.actualState.BlendDestinationFactorAlpha,
		)
	}
}
