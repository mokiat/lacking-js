package internal

import (
	"fmt"

	"github.com/mokiat/lacking/render"
	"github.com/mokiat/wasmgl"
)

func NewRenderer() *Renderer {
	return &Renderer{
		framebuffer: DefaultFramebuffer,
	}
}

type Renderer struct {
	framebuffer           *Framebuffer
	invalidateAttachments []int
	primitive             int
	indexType             int
}

func (r *Renderer) BeginRenderPass(info render.RenderPassInfo) {
	r.framebuffer = info.Framebuffer.(*Framebuffer)

	wasmgl.BindFramebuffer(wasmgl.FRAMEBUFFER, r.framebuffer.raw)
	wasmgl.Viewport(
		info.Viewport.X,
		info.Viewport.Y,
		info.Viewport.Width,
		info.Viewport.Height,
	)

	for i, attachment := range info.Colors {
		if i == 0 && (attachment.LoadOp == render.LoadOperationClear) {
			wasmgl.ClearColor(attachment.ClearValue[0], attachment.ClearValue[1], attachment.ClearValue[2], attachment.ClearValue[3])
			wasmgl.Clear(wasmgl.COLOR_BUFFER_BIT)
		}
		// if r.framebuffer.activeDrawBuffers[i] && (attachment.LoadOp == render.LoadOperationClear) {
		// 	gl.ClearNamedFramebufferfv(r.framebuffer.id, gl.COLOR, int32(i), &attachment.ClearValue[0])
		// }
	}

	clearDepth := info.DepthLoadOp == render.LoadOperationClear
	clearStencil := info.StencilLoadOp == render.LoadOperationClear

	if clearDepth && clearStencil {
		wasmgl.ClearDepth(info.DepthClearValue)
		wasmgl.ClearStencil(info.StencilClearValue)
		wasmgl.Clear(wasmgl.DEPTH_BUFFER_BIT | wasmgl.STENCIL_BUFFER_BIT)
		// 	depthValue := info.DepthClearValue
		// 	stencilValue := int32(info.StencilClearValue)
		// 	gl.ClearNamedFramebufferfi(r.framebuffer.id, gl.DEPTH_STENCIL, 0, depthValue, stencilValue)
	} else {
		if clearDepth {
			wasmgl.ClearDepth(info.DepthClearValue)
			wasmgl.Clear(wasmgl.DEPTH_BUFFER_BIT)
			// 		depthValue := info.DepthClearValue
			// 		gl.ClearNamedFramebufferfv(r.framebuffer.id, gl.DEPTH, 0, &depthValue)
		}
		if clearStencil {
			wasmgl.ClearStencil(info.StencilClearValue)
			wasmgl.Clear(wasmgl.STENCIL_BUFFER_BIT)
			// 		stencilValue := uint32(info.StencilClearValue)
			// 		gl.ClearNamedFramebufferuiv(r.framebuffer.id, gl.STENCIL, 0, &stencilValue)
		}
	}

	r.invalidateAttachments = r.invalidateAttachments[:0]

	invalidateDepth := info.DepthStoreOp == render.StoreOperationDontCare
	invalidateStencil := info.StencilStoreOp == render.StoreOperationDontCare

	for i, attachment := range info.Colors {
		if r.framebuffer.activeDrawBuffers[i] && (attachment.StoreOp == render.StoreOperationDontCare) {
			r.invalidateAttachments = append(r.invalidateAttachments, wasmgl.COLOR_ATTACHMENT0+i)
		}
	}

	if invalidateDepth && invalidateStencil {
		r.invalidateAttachments = append(r.invalidateAttachments, wasmgl.DEPTH_STENCIL_ATTACHMENT)
	} else {
		if invalidateDepth {
			r.invalidateAttachments = append(r.invalidateAttachments, wasmgl.DEPTH_ATTACHMENT)
		}
		if invalidateStencil {
			r.invalidateAttachments = append(r.invalidateAttachments, wasmgl.STENCIL_ATTACHMENT)
		}
	}
}

func (r *Renderer) EndRenderPass() {
	if len(r.invalidateAttachments) > 0 {
		// TODO: When the viewport is just part of the framebuffer
		// we should use glInvalidateNamedFramebufferSubData
		wasmgl.InvalidateFramebuffer(wasmgl.FRAMEBUFFER, r.invalidateAttachments)
	}

	// FIXME
	wasmgl.Disable(wasmgl.BLEND)
	wasmgl.DepthMask(true)
	wasmgl.ColorMask(true, true, true, true)

	r.framebuffer = DefaultFramebuffer
}

func (r *Renderer) BindPipeline(pipeline render.Pipeline) {
	intPipeline := pipeline.(*Pipeline)
	r.executeCommandBindPipeline(CommandBindPipeline{
		ProgramID:        intPipeline.ProgramID,
		Topology:         intPipeline.Topology,
		CullTest:         intPipeline.CullTest,
		FrontFace:        intPipeline.FrontFace,
		LineWidth:        intPipeline.LineWidth,
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
	wasmgl.BindTexture(wasmgl.TEXTURE_2D, intTexture.raw)
	wasmgl.CopyTexSubImage2D(
		wasmgl.TEXTURE_2D,
		info.TextureLevel,
		info.TextureX,
		info.TextureY,
		info.FramebufferX,
		info.FramebufferY,
		info.Width,
		info.Height,
	)
	if info.GenerateMipmaps {
		wasmgl.GenerateMipmap(wasmgl.TEXTURE_2D)
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
		case CommandKindLineWidth:
			command := PopCommand[CommandLineWidth](queue)
			r.executeCommandLineWidth(command)
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
		case CommandKindTextureUnit:
			command := PopCommand[CommandTextureUnit](queue)
			r.executeCommandTextureUnit(command)
		case CommandKindDraw:
			command := PopCommand[CommandDraw](queue)
			r.executeCommandDraw(command)
		case CommandKindDrawIndexed:
			command := PopCommand[CommandDrawIndexed](queue)
			r.executeCommandDrawIndexed(command)
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
	r.executeCommandLineWidth(command.LineWidth)
	r.executeCommandDepthTest(command.DepthTest)
	r.executeCommandDepthWrite(command.DepthWrite)
	r.executeCommandDepthComparison(command.DepthComparison)
	r.executeCommandStencilTest(command.StencilTest)
	// TODO: Optimize if equal except for face
	r.executeCommandStencilFunc(command.StencilFuncFront)
	r.executeCommandStencilFunc(command.StencilFuncBack)
	// TODO: Optimize if equal except for face
	r.executeCommandStencilOperation(command.StencilOpFront)
	r.executeCommandStencilOperation(command.StencilOpBack)
	// TODO: Optimize if equal except for face
	r.executeCommandStencilMask(command.StencilMaskFront)
	r.executeCommandStencilMask(command.StencilMaskBack)
	r.executeCommandColorWrite(command.ColorWrite)
	if command.BlendEnabled {
		wasmgl.Enable(wasmgl.BLEND)
	} else {
		wasmgl.Disable(wasmgl.BLEND)
	}
	r.executeCommandBlendEquation(command.BlendEquation)
	r.executeCommandBlendFunc(command.BlendFunc)
	r.executeCommandBlendColor(command.BlendColor)
	r.executeCommandBindVertexArray(command.VertexArray)
}

func (r *Renderer) executeCommandTopology(command CommandTopology) {
	r.primitive = int(command.Topology)
}

func (r *Renderer) executeCommandCullTest(command CommandCullTest) {
	if command.Enabled {
		wasmgl.Enable(wasmgl.CULL_FACE)
		wasmgl.CullFace(int(command.Face))
	} else {
		wasmgl.Disable(wasmgl.CULL_FACE)
	}
}

func (r *Renderer) executeCommandFrontFace(command CommandFrontFace) {
	wasmgl.FrontFace(int(command.Orientation))
}

func (r *Renderer) executeCommandLineWidth(command CommandLineWidth) {
	if command.Width > 0.0 {
		wasmgl.LineWidth(command.Width)
	}
}

func (r *Renderer) executeCommandDepthTest(command CommandDepthTest) {
	if command.Enabled {
		wasmgl.Enable(wasmgl.DEPTH_TEST)
	} else {
		wasmgl.Disable(wasmgl.DEPTH_TEST)
	}
}

func (r *Renderer) executeCommandDepthWrite(command CommandDepthWrite) {
	wasmgl.DepthMask(command.Enabled)
}

func (r *Renderer) executeCommandDepthComparison(command CommandDepthComparison) {
	wasmgl.DepthFunc(int(command.Mode))
}

func (r *Renderer) executeCommandStencilTest(command CommandStencilTest) {
	if command.Enabled {
		wasmgl.Enable(wasmgl.STENCIL_TEST)
	} else {
		wasmgl.Disable(wasmgl.STENCIL_TEST)
	}
}

func (r *Renderer) executeCommandStencilOperation(command CommandStencilOperation) {
	wasmgl.StencilOpSeparate(
		int(command.Face),
		int(command.StencilFail),
		int(command.DepthFail),
		int(command.Pass),
	)
}

func (r *Renderer) executeCommandStencilFunc(command CommandStencilFunc) {
	wasmgl.StencilFuncSeparate(
		int(command.Face),
		int(command.Func),
		int(command.Ref),
		int(command.Mask),
	)
}

func (r *Renderer) executeCommandStencilMask(command CommandStencilMask) {
	wasmgl.StencilMaskSeparate(
		int(command.Face),
		int(command.Mask),
	)
}

func (r *Renderer) executeCommandColorWrite(command CommandColorWrite) {
	wasmgl.ColorMask(command.Mask[0], command.Mask[1], command.Mask[2], command.Mask[3])
}

func (r *Renderer) executeCommandBlendColor(command CommandBlendColor) {
	wasmgl.BlendColor(
		command.Color[0],
		command.Color[1],
		command.Color[2],
		command.Color[3],
	)
}

func (r *Renderer) executeCommandBlendEquation(command CommandBlendEquation) {
	wasmgl.BlendEquationSeparate(
		int(command.ModeRGB),
		int(command.ModeAlpha),
	)
}

func (r *Renderer) executeCommandBlendFunc(command CommandBlendFunc) {
	wasmgl.BlendFuncSeparate(
		int(command.SourceFactorRGB),
		int(command.DestinationFactorRGB),
		int(command.SourceFactorAlpha),
		int(command.DestinationFactorAlpha),
	)
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

func (r *Renderer) executeCommandTextureUnit(command CommandTextureUnit) {
	texture := textures.Get(command.TextureID)
	wasmgl.ActiveTexture(wasmgl.TEXTURE0 + int(command.Index))
	wasmgl.BindTexture(texture.kind, texture.raw)
}

func (r *Renderer) executeCommandDraw(command CommandDraw) {
	wasmgl.DrawArraysInstanced(
		r.primitive,
		int(command.VertexOffset),
		int(command.VertexCount),
		int(command.InstanceCount),
	)
}

func (r *Renderer) executeCommandDrawIndexed(command CommandDrawIndexed) {
	wasmgl.DrawElementsInstanced(
		r.primitive,
		int(command.IndexCount),
		r.indexType,
		int(command.IndexOffset),
		int(command.InstanceCount),
	)
}
