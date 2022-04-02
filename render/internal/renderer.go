package internal

import (
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

	// TODO
	// var rgba = [4]float32{
	// 	0.0,
	// 	0.3,
	// 	0.6,
	// 	1.0,
	// }
	// gl.ClearNamedFramebufferfv(r.framebuffer.id, gl.COLOR, 0, &rgba[0])

	// clearDepth := info.StencilLoadOp == render.LoadOperationClear
	// clearStencil := info.StencilLoadOp == render.LoadOperationClear

	// if clearDepth && clearStencil {
	// 	depthValue := info.DepthClearValue
	// 	stencilValue := int32(info.StencilClearValue)
	// 	gl.ClearNamedFramebufferfi(r.framebuffer.id, gl.DEPTH_STENCIL, 0, depthValue, stencilValue)
	// } else {
	// 	if clearDepth {
	// 		depthValue := info.DepthClearValue
	// 		gl.ClearNamedFramebufferfv(r.framebuffer.id, gl.DEPTH, 0, &depthValue)
	// 	}
	// 	if clearStencil {
	// 		stencilValue := uint32(info.StencilClearValue)
	// 		gl.ClearNamedFramebufferuiv(r.framebuffer.id, gl.STENCIL, 0, &stencilValue)
	// 	}
	// }

	invalidateDepth := info.DepthStoreOp == render.StoreOperationDontCare
	invalidateStencil := info.StencilStoreOp == render.StoreOperationDontCare

	r.invalidateAttachments = r.invalidateAttachments[:0]

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

func (r *Renderer) BindPipeline(pipeline render.Pipeline) {
	if pipeline, ok := pipeline.(*Pipeline); ok {
		r.primitive = wasmgl.TRIANGLE_FAN // FIXME

		switch pipeline.Culling {
		case render.CullModeNone:
			wasmgl.Disable(wasmgl.CULL_FACE)
		case render.CullModeBack:
			wasmgl.Enable(wasmgl.CULL_FACE)
			wasmgl.CullFace(wasmgl.BACK)
		case render.CullModeFront:
			wasmgl.Enable(wasmgl.CULL_FACE)
			wasmgl.CullFace(wasmgl.FRONT)
		case render.CullModeFrontAndBack:
			wasmgl.Enable(wasmgl.CULL_FACE)
			wasmgl.CullFace(wasmgl.FRONT_AND_BACK)
		}

		// switch pipeline.FrontFace {
		// case render.FaceOrientationCCW:
		// 	wasmgl.FrontFace(wasmgl.CCW)
		// case render.FaceOrientationCW:
		// 	wasmgl.FrontFace(wasmgl.CW)
		// }

		// wasmgl.LineWidth(pipeline.LineWidth)

		if pipeline.DepthTest {
			wasmgl.Enable(wasmgl.DEPTH_TEST)
			// gl.DepthFunc(xfunc uint32) // TODO
		} else {
			wasmgl.Disable(wasmgl.DEPTH_TEST)
		}
		if pipeline.DepthWrite {
			wasmgl.DepthMask(true)
		} else {
			wasmgl.DepthMask(false)
		}

		if pipeline.StencilTest {
			wasmgl.Enable(wasmgl.STENCIL_TEST)
		} else {
			wasmgl.Disable(wasmgl.STENCIL_TEST)
		}

		wasmgl.Enable(wasmgl.BLEND)
		wasmgl.BlendFunc(wasmgl.SRC_ALPHA, wasmgl.ONE_MINUS_SRC_ALPHA)

		wasmgl.ColorMask(pipeline.ColorWrite[0], pipeline.ColorWrite[1], pipeline.ColorWrite[2], pipeline.ColorWrite[3])

		if program, ok := pipeline.Program.(*Program); ok {
			wasmgl.UseProgram(program.raw)
		}

		if vertexArray, ok := pipeline.VertexArray.(*VertexArray); ok {
			wasmgl.BindVertexArray(vertexArray.raw)
			r.indexType = vertexArray.indexFormat
		}
	}
}

func (r *Renderer) Uniform4f(location render.UniformLocation, values [4]float32) {
	wasmgl.Uniform4f(location.(wasmgl.UniformLocation), values[0], values[1], values[2], values[3])
}

func (r *Renderer) Uniform1i(location render.UniformLocation, value int) {
	wasmgl.Uniform1i(location.(wasmgl.UniformLocation), value)
}

func (r *Renderer) UniformMatrix4f(location render.UniformLocation, values [16]float32) {
	wasmgl.UniformMatrix4fv(location.(wasmgl.UniformLocation), false, values[:])
}

func (r *Renderer) TextureUnit(index int, texture render.Texture) {
	if texture, ok := texture.(*Texture); ok {
		wasmgl.ActiveTexture(wasmgl.TEXTURE0 + index)
		wasmgl.BindTexture(texture.kind, texture.raw)
	}
}

func (r *Renderer) Draw(vertexOffset, vertexCount, instanceCount int) {
	wasmgl.DrawArrays(r.primitive, vertexOffset, vertexCount)
	// gl.DrawArraysInstanced(r.primitive, int32(vertexOffset), int32(vertexCount), int32(instanceCount))
}

func (r *Renderer) DrawIndexed(indexOffset, indexCount, instanceCount int) {
	wasmgl.DrawElements(r.primitive, indexCount, r.indexType, indexOffset)
	// gl.DrawElementsInstanced(r.primitive, int32(indexCount), r.indexType, gl.PtrOffset(indexOffset), int32(instanceCount))
}

func (r *Renderer) EndRenderPass() {
	if len(r.invalidateAttachments) > 0 {
		// TODO: When the viewport is just part of the framebuffer
		// we should use glInvalidateNamedFramebufferSubData
		// wasmgl.InvalidateNamedFramebufferData(r.framebuffer.id, 1, &r.invalidateAttachments[0])
	}

	// FIXME
	wasmgl.Disable(wasmgl.BLEND)
	wasmgl.DepthMask(true)

	r.framebuffer = DefaultFramebuffer
}
