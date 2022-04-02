package render

import (
	"github.com/mokiat/lacking-js/render/internal"
	"github.com/mokiat/lacking/render"
)

func NewAPI() render.API {
	return &API{
		renderer: internal.NewRenderer(),
	}
}

type API struct {
	renderer *internal.Renderer
}

func (a *API) Capabilities() render.Capabilities {
	return render.Capabilities{
		Quality: render.QualityHigh,
	}
}

func (a *API) DefaultFramebuffer() render.Framebuffer {
	return internal.DefaultFramebuffer
}

func (a *API) CreateFramebuffer(info render.FramebufferInfo) render.Framebuffer {
	return internal.NewFramebuffer(info)
}

func (a *API) CreateColorTexture2D(info render.ColorTexture2DInfo) render.Texture {
	return internal.NewColorTexture2D(info)
}

func (a *API) CreateColorTextureCube(info render.ColorTextureCubeInfo) render.Texture {
	panic("TODO")
}

func (a *API) CreateDepthTexture2D(info render.DepthTexture2DInfo) render.Texture {
	return internal.NewDepthTexture2D(info)
}

func (a *API) CreateStencilTexture2D(info render.StencilTexture2DInfo) render.Texture {
	return internal.NewStencilTexture2D(info)
}

func (a *API) CreateDepthStencilTexture2D(info render.DepthStencilTexture2DInfo) render.Texture {
	return internal.NewDepthStencilTexture2D(info)
}

func (a *API) CreateVertexShader(info render.ShaderInfo) render.Shader {
	return internal.NewVertexShader(info)
}

func (a *API) CreateFragmentShader(info render.ShaderInfo) render.Shader {
	return internal.NewFragmentShader(info)
}

func (a *API) CreateProgram(info render.ProgramInfo) render.Program {
	return internal.NewProgram(info)
}

func (a *API) CreateVertexBuffer(info render.BufferInfo) render.Buffer {
	return internal.NewVertexBuffer(info)
}

func (a *API) CreateIndexBuffer(info render.BufferInfo) render.Buffer {
	return internal.NewIndexBuffer(info)
}

func (a *API) CreateVertexArray(info render.VertexArrayInfo) render.VertexArray {
	return internal.NewVertexArray(info)
}

func (a *API) CreatePipeline(info render.PipelineInfo) render.Pipeline {
	return internal.NewPipeline(info)
}

func (a *API) BeginRenderPass(info render.RenderPassInfo) {
	a.renderer.BeginRenderPass(info)
}

func (a *API) BindPipeline(pipeline render.Pipeline) {
	a.renderer.BindPipeline(pipeline)
}

func (a *API) Uniform4f(location render.UniformLocation, values [4]float32) {
	a.renderer.Uniform4f(location, values)
}

func (a *API) Uniform1i(location render.UniformLocation, value int) {
	a.renderer.Uniform1i(location, value)
}

func (a *API) UniformMatrix4f(location render.UniformLocation, values [16]float32) {
	a.renderer.UniformMatrix4f(location, values)
}

func (a *API) TextureUnit(index int, texture render.Texture) {
	a.renderer.TextureUnit(index, texture)
}

func (a *API) Draw(vertexOffset, vertexCount, instanceCount int) {
	a.renderer.Draw(vertexOffset, vertexCount, instanceCount)
}

func (a *API) DrawIndexed(indexOffset, indexCount, instanceCount int) {
	a.renderer.DrawIndexed(indexOffset, indexCount, instanceCount)
}

func (a *API) EndRenderPass() {
	a.renderer.EndRenderPass()
}
