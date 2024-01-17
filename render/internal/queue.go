package internal

import (
	"fmt"

	"github.com/mokiat/gog/opt"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/wasmgl"
)

func NewQueue() *Queue {
	return &Queue{
		invalidateAttachments: make([]uint32, 0, 16),
	}
}

type Queue struct {
	render.QueueMarker

	currentProgram                     opt.T[uint32]
	currentTopology                    opt.T[uint32]
	currentIndexType                   opt.T[uint32]
	currentCullTest                    opt.T[bool]
	currentCullFace                    opt.T[uint32]
	currentFrontFace                   opt.T[uint32]
	currentDepthTest                   opt.T[bool]
	currentDepthWrite                  opt.T[bool]
	currentDepthComparison             opt.T[uint32]
	currentStencilTest                 opt.T[bool]
	currentStencilOpStencilFailFront   opt.T[uint32]
	currentStencilOpDepthFailFront     opt.T[uint32]
	currentStencilOpPassFront          opt.T[uint32]
	currentStencilOpStencilFailBack    opt.T[uint32]
	currentStencilOpDepthFailBack      opt.T[uint32]
	currentStencilOpPassBack           opt.T[uint32]
	currentStencilComparisonFuncFront  opt.T[uint32]
	currentStencilComparisonRefFront   opt.T[int32]
	currentStencilComparisonMaskFront  opt.T[uint32]
	currentStencilComparisonFuncBack   opt.T[uint32]
	currentStencilComparisonRefBack    opt.T[int32]
	currentStencilComparisonMaskBack   opt.T[uint32]
	currentStencilMaskFront            opt.T[uint32]
	currentStencilMaskBack             opt.T[uint32]
	currentColorMask                   opt.T[[4]bool]
	currentBlending                    opt.T[bool]
	currentBlendColor                  opt.T[[4]float32]
	currentBlendModeRGB                opt.T[uint32]
	currentBlendModeAlpha              opt.T[uint32]
	currentBlendSourceFactorRGB        opt.T[uint32]
	currentBlendDestinationFactorRGB   opt.T[uint32]
	currentBlendSourceFactorAlpha      opt.T[uint32]
	currentBlendDestinationFactorAlpha opt.T[uint32]

	invalidateAttachments []wasmgl.GLenum
}

func (q *Queue) Invalidate() {
	q.currentProgram = opt.Unspecified[uint32]()
	q.currentTopology = opt.Unspecified[uint32]()
	q.currentIndexType = opt.Unspecified[uint32]()
	q.currentCullTest = opt.Unspecified[bool]()
	q.currentCullFace = opt.Unspecified[uint32]()
	q.currentFrontFace = opt.Unspecified[uint32]()
	q.currentDepthTest = opt.Unspecified[bool]()
	q.currentDepthWrite = opt.Unspecified[bool]()
	q.currentDepthComparison = opt.Unspecified[uint32]()
	q.currentStencilTest = opt.Unspecified[bool]()
	q.currentStencilOpStencilFailFront = opt.Unspecified[uint32]()
	q.currentStencilOpDepthFailFront = opt.Unspecified[uint32]()
	q.currentStencilOpPassFront = opt.Unspecified[uint32]()
	q.currentStencilOpStencilFailBack = opt.Unspecified[uint32]()
	q.currentStencilOpDepthFailBack = opt.Unspecified[uint32]()
	q.currentStencilOpPassBack = opt.Unspecified[uint32]()
	q.currentStencilComparisonFuncFront = opt.Unspecified[uint32]()
	q.currentStencilComparisonRefFront = opt.Unspecified[int32]()
	q.currentStencilComparisonMaskFront = opt.Unspecified[uint32]()
	q.currentStencilComparisonFuncBack = opt.Unspecified[uint32]()
	q.currentStencilComparisonRefBack = opt.Unspecified[int32]()
	q.currentStencilComparisonMaskBack = opt.Unspecified[uint32]()
	q.currentStencilMaskFront = opt.Unspecified[uint32]()
	q.currentStencilMaskBack = opt.Unspecified[uint32]()
	q.currentColorMask = opt.Unspecified[[4]bool]()
	q.currentBlending = opt.Unspecified[bool]()
	q.currentBlendColor = opt.Unspecified[[4]float32]()
	q.currentBlendSourceFactorRGB = opt.Unspecified[uint32]()
	q.currentBlendDestinationFactorRGB = opt.Unspecified[uint32]()
	q.currentBlendSourceFactorAlpha = opt.Unspecified[uint32]()
	q.currentBlendDestinationFactorAlpha = opt.Unspecified[uint32]()
	q.currentBlendModeRGB = opt.Unspecified[uint32]()
	q.currentBlendModeAlpha = opt.Unspecified[uint32]()

	q.invalidateAttachments = q.invalidateAttachments[:0]
}

func (q *Queue) WriteBuffer(buffer render.Buffer, offset int, data []byte) {
	actualBuffer := buffer.(*Buffer)
	wasmgl.BindBuffer(actualBuffer.kind, actualBuffer.raw)
	wasmgl.BufferSubData(actualBuffer.kind, wasmgl.GLintptr(offset), data)
	wasmgl.BindBuffer(actualBuffer.kind, wasmgl.NilBuffer)
}

func (q *Queue) ReadBuffer(buffer render.Buffer, offset int, target []byte) {
	actualBuffer := buffer.(*Buffer)
	wasmgl.BindBuffer(actualBuffer.kind, actualBuffer.raw)
	wasmgl.GetBufferSubData(actualBuffer.kind, wasmgl.GLintptr(offset), target)
	wasmgl.BindBuffer(actualBuffer.kind, wasmgl.NilBuffer)
}

func (q *Queue) Submit(commands render.CommandBuffer) {
	commandBuffer := commands.(*CommandBuffer)
	for commandBuffer.HasMoreCommands() {
		header := readCommandChunk[CommandHeader](commandBuffer)
		switch header.Kind {
		case CommandKindCopyFramebufferToBuffer:
			command := readCommandChunk[CommandCopyFramebufferToBuffer](commandBuffer)
			q.executeCommandCopyFramebufferToBuffer(command)
		case CommandKindCopyFramebufferToTexture:
			command := readCommandChunk[CommandCopyFramebufferToTexture](commandBuffer)
			q.executeCommandCopyFramebufferToTexture(command)
		case CommandKindBeginRenderPass:
			command := readCommandChunk[CommandBeginRenderPass](commandBuffer)
			q.executeCommandBeginRenderPass(command)
		case CommandKindEndRenderPass:
			command := readCommandChunk[CommandEndRenderPass](commandBuffer)
			q.executeCommandEndRenderPass(command)
		case CommandKindBindPipeline:
			command := readCommandChunk[CommandBindPipeline](commandBuffer)
			q.executeCommandBindPipeline(command)
		case CommandKindTextureUnit:
			command := readCommandChunk[CommandTextureUnit](commandBuffer)
			q.executeCommandTextureUnit(command)
		case CommandKindUniformBufferUnit:
			command := readCommandChunk[CommandUniformBufferUnit](commandBuffer)
			q.executeCommandUniformBufferUnit(command)
		case CommandKindDraw:
			command := readCommandChunk[CommandDraw](commandBuffer)
			q.executeCommandDraw(command)
		case CommandKindDrawIndexed:
			command := readCommandChunk[CommandDrawIndexed](commandBuffer)
			q.executeCommandDrawIndexed(command)
		default:
			panic(fmt.Errorf("unknown command kind: %v", header.Kind))
		}
	}
	commandBuffer.Reset()
}

func (q *Queue) TrackSubmittedWorkDone() render.Fence {
	return NewFence()
}

func (q *Queue) executeCommandCopyFramebufferToBuffer(command CommandCopyFramebufferToBuffer) {
	buffer := buffers.Get(command.BufferID)
	wasmgl.BindBuffer(
		buffer.kind,
		buffer.raw,
	)
	wasmgl.ReadPixels(
		wasmgl.GLint(command.X),
		wasmgl.GLint(command.Y),
		wasmgl.GLsizei(command.Width),
		wasmgl.GLsizei(command.Height),
		wasmgl.GLenum(command.Format),
		wasmgl.GLenum(command.XType),
		wasmgl.GLintptr(command.BufferOffset),
	)
	wasmgl.BindBuffer(
		buffer.kind,
		wasmgl.NilBuffer,
	)
}

func (q *Queue) executeCommandCopyFramebufferToTexture(command CommandCopyFramebufferToTexture) {
	intTexture := textures.Get(command.TextureID)
	wasmgl.BindTexture(intTexture.kind, intTexture.raw)
	wasmgl.CopyTexSubImage2D(
		intTexture.kind,
		command.TextureLevel,
		command.TextureX,
		command.TextureY,
		command.FramebufferX,
		command.FramebufferY,
		command.Width,
		command.Height,
	)
	if command.GenerateMipmaps {
		wasmgl.GenerateMipmap(intTexture.kind)
	}
}

func (q *Queue) executeCommandBeginRenderPass(command CommandBeginRenderPass) {
	intFramebuffer := framebuffers.Get(command.FramebufferID)

	wasmgl.BindFramebuffer(wasmgl.FRAMEBUFFER, intFramebuffer.raw)
	wasmgl.Viewport(
		command.ViewportX,
		command.ViewportY,
		command.ViewportWidth,
		command.ViewportHeight,
	)

	var colorMaskChanged bool
	for i, attachment := range command.Colors {
		loadOp := CommandLoadOperationToRender(attachment.LoadOp)
		if intFramebuffer.activeDrawBuffers[i] && (loadOp == render.LoadOperationClear) {
			if !colorMaskChanged {
				q.executeCommandColorWrite(CommandColorWrite{
					Mask: render.ColorMaskTrue,
				})
				colorMaskChanged = true
			}
			wasmgl.ClearBufferfv(wasmgl.COLOR, wasmgl.GLint(i), attachment.ClearValue[:])
		}
	}

	clearDepth := CommandLoadOperationToRender(command.DepthLoadOp) == render.LoadOperationClear
	clearStencil := CommandLoadOperationToRender(command.StencilLoadOp) == render.LoadOperationClear

	if clearDepth && clearStencil {
		q.executeCommandDepthWrite(CommandDepthWrite{
			Enabled: true,
		})
		q.executeCommandStencilMask(CommandStencilMask{
			Face: wasmgl.FRONT_AND_BACK,
			Mask: 0xFF,
		})
		depthValue := command.DepthClearValue
		stencilValue := command.StencilClearValue
		wasmgl.ClearBufferfi(wasmgl.DEPTH_STENCIL, 0, depthValue, stencilValue)
	} else {
		if clearDepth {
			q.executeCommandDepthWrite(CommandDepthWrite{
				Enabled: true,
			})
			depthValues := [1]float32{command.DepthClearValue}
			wasmgl.ClearBufferfv(wasmgl.DEPTH, 0, depthValues[:])
		}
		if clearStencil {
			q.executeCommandStencilMask(CommandStencilMask{
				Face: wasmgl.FRONT_AND_BACK,
				Mask: 0xFF,
			})
			stencilValues := [1]int32{command.StencilClearValue}
			wasmgl.ClearBufferiv(wasmgl.STENCIL, 0, stencilValues[:])
		}
	}

	isDefaultFramebuffer := intFramebuffer == DefaultFramebuffer
	q.invalidateAttachments = q.invalidateAttachments[:0]

	for i, attachment := range command.Colors {
		if intFramebuffer.activeDrawBuffers[i] && (CommandStoreOperationToRender(attachment.StoreOp) == render.StoreOperationDontCare) {
			if isDefaultFramebuffer {
				if i == 0 {
					q.invalidateAttachments = append(q.invalidateAttachments, wasmgl.COLOR)
				}
			} else {
				q.invalidateAttachments = append(q.invalidateAttachments, wasmgl.COLOR_ATTACHMENT0+wasmgl.GLenum(i))
			}
		}
	}

	invalidateDepth := CommandStoreOperationToRender(command.DepthStoreOp) == render.StoreOperationDontCare
	invalidateStencil := CommandStoreOperationToRender(command.StencilStoreOp) == render.StoreOperationDontCare

	if invalidateDepth && invalidateStencil && !isDefaultFramebuffer {
		q.invalidateAttachments = append(q.invalidateAttachments, wasmgl.DEPTH_STENCIL_ATTACHMENT)
	} else {
		if invalidateDepth {
			if isDefaultFramebuffer {
				q.invalidateAttachments = append(q.invalidateAttachments, wasmgl.DEPTH)
			} else {
				q.invalidateAttachments = append(q.invalidateAttachments, wasmgl.DEPTH_ATTACHMENT)
			}
		}
		if invalidateStencil {
			if isDefaultFramebuffer {
				q.invalidateAttachments = append(q.invalidateAttachments, wasmgl.STENCIL)
			} else {
				q.invalidateAttachments = append(q.invalidateAttachments, wasmgl.STENCIL_ATTACHMENT)
			}
		}
	}
}

func (q *Queue) executeCommandEndRenderPass(command CommandEndRenderPass) {
	if len(q.invalidateAttachments) > 0 {
		wasmgl.InvalidateFramebuffer(wasmgl.FRAMEBUFFER, q.invalidateAttachments)
	}
}

func (q *Queue) executeCommandBindPipeline(command CommandBindPipeline) {
	program := programs.Get(command.ProgramID)
	if isDirty(q.currentProgram, command.ProgramID) {
		q.currentProgram = opt.V(command.ProgramID)
		wasmgl.UseProgram(program.raw)
	}
	q.executeCommandTopology(command.Topology)
	q.executeCommandCullTest(command.CullTest)
	if command.CullTest.Enabled {
		q.executeCommandCullFace(command.CullTest)
	}
	q.executeCommandFrontFace(command.FrontFace)
	q.executeCommandDepthTest(command.DepthTest)
	q.executeCommandDepthWrite(command.DepthWrite)
	if command.DepthTest.Enabled {
		q.executeCommandDepthComparison(command.DepthComparison)
	}
	q.executeCommandStencilTest(command.StencilTest)
	if command.StencilTest.Enabled {
		q.executeCommandStencilOperation(command.StencilOpFront)
		q.executeCommandStencilOperation(command.StencilOpBack)
		q.executeCommandStencilFunc(command.StencilFuncFront)
		q.executeCommandStencilFunc(command.StencilFuncBack)
		q.executeCommandStencilMask(command.StencilMaskFront)
		q.executeCommandStencilMask(command.StencilMaskBack)
	}
	q.executeCommandColorWrite(command.ColorWrite)
	if isDirty(q.currentBlending, command.BlendEnabled) {
		q.currentBlending = opt.V(command.BlendEnabled)
		if command.BlendEnabled {
			wasmgl.Enable(wasmgl.BLEND)
		} else {
			wasmgl.Disable(wasmgl.BLEND)
		}
	}
	if command.BlendEnabled {
		q.executeCommandBlendColor(command.BlendColor)
		q.executeCommandBlendEquation(command.BlendEquation)
		q.executeCommandBlendFunc(command.BlendFunc)
	}
	q.executeCommandBindVertexArray(command.VertexArray)
}

func (q *Queue) executeCommandTopology(command CommandTopology) {
	q.currentTopology = opt.V(command.Topology)
}

func (q *Queue) executeCommandCullTest(command CommandCullTest) {
	needsUpdate := isDirty(q.currentCullTest, command.Enabled)
	if needsUpdate {
		q.currentCullTest = opt.V(command.Enabled)
		if command.Enabled {
			wasmgl.Enable(wasmgl.CULL_FACE)
		} else {
			wasmgl.Disable(wasmgl.CULL_FACE)
		}
	}
}

func (q *Queue) executeCommandCullFace(command CommandCullTest) {
	needsUpdate := isDirty(q.currentCullFace, command.Face)
	if needsUpdate {
		q.currentCullFace = opt.V(command.Face)
		wasmgl.CullFace(command.Face)
	}
}

func (q *Queue) executeCommandFrontFace(command CommandFrontFace) {
	needsUpdate := isDirty(q.currentFrontFace, command.Orientation)
	if needsUpdate {
		q.currentFrontFace = opt.V(command.Orientation)
		wasmgl.FrontFace(command.Orientation)
	}
}

func (q *Queue) executeCommandDepthTest(command CommandDepthTest) {
	needsUpdate := isDirty(q.currentDepthTest, command.Enabled)
	if needsUpdate {
		q.currentDepthTest = opt.V(command.Enabled)
		if command.Enabled {
			wasmgl.Enable(wasmgl.DEPTH_TEST)
		} else {
			wasmgl.Disable(wasmgl.DEPTH_TEST)
		}
	}
}

func (q *Queue) executeCommandDepthWrite(command CommandDepthWrite) {
	needsUpdate := isDirty(q.currentDepthWrite, command.Enabled)
	if needsUpdate {
		q.currentDepthWrite = opt.V(command.Enabled)
		wasmgl.DepthMask(command.Enabled)
	}
}

func (q *Queue) executeCommandDepthComparison(command CommandDepthComparison) {
	needsUpdate := isDirty(q.currentDepthComparison, command.Mode)
	if needsUpdate {
		q.currentDepthComparison = opt.V(command.Mode)
		wasmgl.DepthFunc(command.Mode)
	}
}

func (q *Queue) executeCommandStencilTest(command CommandStencilTest) {
	needsUpdate := isDirty(q.currentStencilTest, command.Enabled)
	if needsUpdate {
		q.currentStencilTest = opt.V(command.Enabled)
		if command.Enabled {
			wasmgl.Enable(wasmgl.STENCIL_TEST)
		} else {
			wasmgl.Disable(wasmgl.STENCIL_TEST)
		}
	}
}

func (q *Queue) executeCommandStencilOperation(command CommandStencilOperation) {
	affectsFront := command.Face == wasmgl.FRONT || command.Face == wasmgl.FRONT_AND_BACK
	affectsBack := command.Face == wasmgl.BACK || command.Face == wasmgl.FRONT_AND_BACK

	frontNeedsUpdate := isDirty(q.currentStencilOpStencilFailFront, command.StencilFail) ||
		isDirty(q.currentStencilOpDepthFailFront, command.DepthFail) ||
		isDirty(q.currentStencilOpPassFront, command.Pass)
	if frontNeedsUpdate && affectsFront {
		q.currentStencilOpStencilFailFront = opt.V(command.StencilFail)
		q.currentStencilOpDepthFailFront = opt.V(command.DepthFail)
		q.currentStencilOpPassFront = opt.V(command.Pass)
	}

	backNeedsUpdate := isDirty(q.currentStencilOpStencilFailBack, command.StencilFail) ||
		isDirty(q.currentStencilOpDepthFailBack, command.DepthFail) ||
		isDirty(q.currentStencilOpPassBack, command.Pass)
	if backNeedsUpdate && affectsBack {
		q.currentStencilOpStencilFailBack = opt.V(command.StencilFail)
		q.currentStencilOpDepthFailBack = opt.V(command.DepthFail)
		q.currentStencilOpPassBack = opt.V(command.Pass)
	}

	switch {
	case affectsFront && affectsBack && frontNeedsUpdate && backNeedsUpdate:
		wasmgl.StencilOpSeparate(
			wasmgl.FRONT_AND_BACK,
			command.StencilFail,
			command.DepthFail,
			command.Pass,
		)
	case affectsFront && frontNeedsUpdate:
		wasmgl.StencilOpSeparate(
			wasmgl.FRONT,
			command.StencilFail,
			command.DepthFail,
			command.Pass,
		)
	case affectsBack && backNeedsUpdate:
		wasmgl.StencilOpSeparate(
			wasmgl.BACK,
			command.StencilFail,
			command.DepthFail,
			command.Pass,
		)
	}
}

func (q *Queue) executeCommandStencilFunc(command CommandStencilFunc) {
	affectsFront := command.Face == wasmgl.FRONT || command.Face == wasmgl.FRONT_AND_BACK
	affectsBack := command.Face == wasmgl.BACK || command.Face == wasmgl.FRONT_AND_BACK

	frontNeedsUpdate := isDirty(q.currentStencilComparisonFuncFront, command.Func) ||
		isDirty(q.currentStencilComparisonRefFront, command.Ref) ||
		isDirty(q.currentStencilComparisonMaskFront, command.Mask)
	if frontNeedsUpdate && affectsFront {
		q.currentStencilComparisonFuncFront = opt.V(command.Func)
		q.currentStencilComparisonRefFront = opt.V(command.Ref)
		q.currentStencilComparisonMaskFront = opt.V(command.Mask)
	}

	backNeedsUpdate := isDirty(q.currentStencilComparisonFuncBack, command.Func) ||
		isDirty(q.currentStencilComparisonRefBack, command.Ref) ||
		isDirty(q.currentStencilComparisonMaskBack, command.Mask)
	if backNeedsUpdate && affectsBack {
		q.currentStencilComparisonFuncBack = opt.V(command.Func)
		q.currentStencilComparisonRefBack = opt.V(command.Ref)
		q.currentStencilComparisonMaskBack = opt.V(command.Mask)
	}

	switch {
	case affectsFront && affectsBack && frontNeedsUpdate && backNeedsUpdate:
		wasmgl.StencilFuncSeparate(
			wasmgl.FRONT_AND_BACK,
			command.Func,
			command.Ref,
			command.Mask,
		)
	case affectsFront && frontNeedsUpdate:
		wasmgl.StencilFuncSeparate(
			wasmgl.FRONT,
			command.Func,
			command.Ref,
			command.Mask,
		)
	case affectsBack && backNeedsUpdate:
		wasmgl.StencilFuncSeparate(
			wasmgl.BACK,
			command.Func,
			command.Ref,
			command.Mask,
		)
	}
}

func (q *Queue) executeCommandStencilMask(command CommandStencilMask) {
	affectsFront := command.Face == wasmgl.FRONT || command.Face == wasmgl.FRONT_AND_BACK
	affectsBack := command.Face == wasmgl.BACK || command.Face == wasmgl.FRONT_AND_BACK

	frontNeedsUpdate := isDirty(q.currentStencilMaskFront, command.Mask)
	if frontNeedsUpdate && affectsFront {
		q.currentStencilMaskFront = opt.V(command.Mask)
	}

	backNeedsUpdate := isDirty(q.currentStencilMaskBack, command.Mask)
	if backNeedsUpdate && affectsBack {
		q.currentStencilMaskBack = opt.V(command.Mask)
	}

	switch {
	case affectsFront && affectsBack && frontNeedsUpdate && backNeedsUpdate:
		wasmgl.StencilMaskSeparate(
			wasmgl.FRONT_AND_BACK,
			command.Mask,
		)
	case affectsFront && frontNeedsUpdate:
		wasmgl.StencilMaskSeparate(
			wasmgl.FRONT,
			command.Mask,
		)
	case affectsBack && backNeedsUpdate:
		wasmgl.StencilMaskSeparate(
			wasmgl.BACK,
			command.Mask,
		)
	}
}

func (q *Queue) executeCommandColorWrite(command CommandColorWrite) {
	needsUpdate := isDirty(q.currentColorMask, command.Mask)
	if needsUpdate {
		q.currentColorMask = opt.V(command.Mask)
		wasmgl.ColorMask(
			command.Mask[0],
			command.Mask[1],
			command.Mask[2],
			command.Mask[3],
		)
	}
}

func (q *Queue) executeCommandBlendColor(command CommandBlendColor) {
	needsUpdate := isDirty(q.currentBlendColor, command.Color)
	if needsUpdate {
		q.currentBlendColor = opt.V(command.Color)
		wasmgl.BlendColor(
			command.Color[0],
			command.Color[1],
			command.Color[2],
			command.Color[3],
		)
	}
}

func (q *Queue) executeCommandBlendEquation(command CommandBlendEquation) {
	needsUpdate := isDirty(q.currentBlendModeRGB, command.ModeRGB) ||
		isDirty(q.currentBlendModeAlpha, command.ModeAlpha)
	if needsUpdate {
		q.currentBlendModeRGB = opt.V(command.ModeRGB)
		q.currentBlendModeAlpha = opt.V(command.ModeAlpha)
		wasmgl.BlendEquationSeparate(
			command.ModeRGB,
			command.ModeAlpha,
		)
	}
}

func (q *Queue) executeCommandBlendFunc(command CommandBlendFunc) {
	needsUpdate := isDirty(q.currentBlendSourceFactorRGB, command.SourceFactorRGB) ||
		isDirty(q.currentBlendDestinationFactorRGB, command.DestinationFactorRGB) ||
		isDirty(q.currentBlendSourceFactorAlpha, command.SourceFactorAlpha) ||
		isDirty(q.currentBlendDestinationFactorAlpha, command.DestinationFactorAlpha)
	if needsUpdate {
		q.currentBlendSourceFactorRGB = opt.V(command.SourceFactorRGB)
		q.currentBlendDestinationFactorRGB = opt.V(command.DestinationFactorRGB)
		q.currentBlendSourceFactorAlpha = opt.V(command.SourceFactorAlpha)
		q.currentBlendDestinationFactorAlpha = opt.V(command.DestinationFactorAlpha)
		wasmgl.BlendFuncSeparate(
			command.SourceFactorRGB,
			command.DestinationFactorRGB,
			command.SourceFactorAlpha,
			command.DestinationFactorAlpha,
		)
	}
}

func (q *Queue) executeCommandBindVertexArray(command CommandBindVertexArray) {
	// NOTE: We don't cache the array since there is a risk that during creation
	// the array has been changed without the queue knowing about it.
	vertexArray := vertexArrays.Get(command.VertexArrayID)
	wasmgl.BindVertexArray(vertexArray.raw)
	q.currentIndexType = opt.V(command.IndexFormat)
}

func (q *Queue) executeCommandTextureUnit(command CommandTextureUnit) {
	texture := textures.Get(command.TextureID)
	wasmgl.ActiveTexture(wasmgl.TEXTURE0 + command.Index)
	wasmgl.BindTexture(texture.kind, texture.raw)
}

func (q *Queue) executeCommandUniformBufferUnit(command CommandUniformBufferUnit) {
	buffer := buffers.Get(command.BufferID)
	wasmgl.BindBufferRange(
		wasmgl.UNIFORM_BUFFER,
		wasmgl.GLuint(command.Index),
		buffer.raw,
		wasmgl.GLintptr(command.Offset),
		wasmgl.GLsizeiptr(command.Size),
	)
}

func (q *Queue) executeCommandDraw(command CommandDraw) {
	wasmgl.DrawArraysInstanced(
		q.currentTopology.Value,
		wasmgl.GLint(command.VertexOffset),
		wasmgl.GLsizei(command.VertexCount),
		wasmgl.GLsizei(command.InstanceCount),
	)
}

func (q *Queue) executeCommandDrawIndexed(command CommandDrawIndexed) {
	wasmgl.DrawElementsInstanced(
		q.currentTopology.Value,
		wasmgl.GLsizei(command.IndexCount),
		q.currentIndexType.Value,
		wasmgl.GLintptr(command.IndexOffset),
		wasmgl.GLsizei(command.InstanceCount),
	)
}

func isDirty[T comparable](cached opt.T[T], desired T) bool {
	return !cached.Specified || (cached.Value != desired)
}
