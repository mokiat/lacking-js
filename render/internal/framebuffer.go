package internal

import (
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/wasmgl"
)

func NewFramebuffer(info render.FramebufferInfo) *Framebuffer {
	raw := wasmgl.CreateFramebuffer()
	wasmgl.BindFramebuffer(wasmgl.FRAMEBUFFER, raw)

	var activeDrawBuffers [4]bool
	var drawBuffers []wasmgl.GLenum
	for i, attachment := range info.ColorAttachments {
		if colorAttachment, ok := attachment.(*Texture); ok {
			attachmentID := wasmgl.COLOR_ATTACHMENT0 + wasmgl.GLenum(i)
			wasmgl.FramebufferTexture2D(wasmgl.FRAMEBUFFER, attachmentID, wasmgl.TEXTURE_2D, colorAttachment.raw, 0)
			drawBuffers = append(drawBuffers, attachmentID)
			activeDrawBuffers[i] = true
		}
	}

	if depthStencilAttachment, ok := info.DepthStencilAttachment.(*Texture); ok {
		wasmgl.FramebufferTexture2D(wasmgl.FRAMEBUFFER, wasmgl.DEPTH_STENCIL_ATTACHMENT, wasmgl.TEXTURE_2D, depthStencilAttachment.raw, 0)
	} else {
		if depthAttachment, ok := info.DepthAttachment.(*Texture); ok {
			wasmgl.FramebufferTexture2D(wasmgl.FRAMEBUFFER, wasmgl.DEPTH_ATTACHMENT, wasmgl.TEXTURE_2D, depthAttachment.raw, 0)
		}
		if stencilAttachment, ok := info.StencilAttachment.(*Texture); ok {
			wasmgl.FramebufferTexture2D(wasmgl.FRAMEBUFFER, wasmgl.STENCIL_ATTACHMENT, wasmgl.TEXTURE_2D, stencilAttachment.raw, 0)
		}
	}

	wasmgl.DrawBuffers(drawBuffers)

	status := wasmgl.CheckFramebufferStatus(wasmgl.FRAMEBUFFER)
	if status != wasmgl.FRAMEBUFFER_COMPLETE {
		logger.Error("Framebuffer (%q) is incomplete!", info.Label)
	}

	result := &Framebuffer{
		raw:               raw,
		activeDrawBuffers: activeDrawBuffers,
	}
	result.id = framebuffers.Allocate(result)
	return result
}

var DefaultFramebuffer = &Framebuffer{
	raw:               wasmgl.NilFramebuffer,
	activeDrawBuffers: [4]bool{true, false, false, false},
}

func init() {
	DefaultFramebuffer.id = framebuffers.Allocate(DefaultFramebuffer)
}

type Framebuffer struct {
	render.FramebufferMarker
	id                uint32
	raw               wasmgl.Framebuffer
	activeDrawBuffers [4]bool
}

func (f *Framebuffer) Release() {
	framebuffers.Release(f.id)
	wasmgl.DeleteFramebuffer(f.raw)
	f.raw = wasmgl.NilFramebuffer
	f.id = 0
	f.activeDrawBuffers = [4]bool{}
}

func DetermineContentFormat(framebuffer render.Framebuffer) render.DataFormat {
	fb := framebuffer.(*Framebuffer)
	wasmgl.BindFramebuffer(wasmgl.FRAMEBUFFER, fb.raw)
	defer func() {
		wasmgl.BindFramebuffer(wasmgl.FRAMEBUFFER, wasmgl.NilFramebuffer)
	}()
	glFormat := wasmgl.GetParameter(
		wasmgl.IMPLEMENTATION_COLOR_READ_FORMAT,
	).GLenum()
	if glFormat != wasmgl.RGBA {
		return render.DataFormatUnsupported
	}
	glType := wasmgl.GetParameter(
		wasmgl.IMPLEMENTATION_COLOR_READ_TYPE,
	).GLenum()
	switch glType {
	case wasmgl.UNSIGNED_BYTE:
		return render.DataFormatRGBA8
	case wasmgl.HALF_FLOAT:
		return render.DataFormatRGBA16F
	case wasmgl.FLOAT:
		return render.DataFormatRGBA32F
	default:
		return render.DataFormatUnsupported
	}
}
