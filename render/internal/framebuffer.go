package internal

import (
	"github.com/mokiat/lacking/log"
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/wasmgl"
)

func NewFramebuffer(info render.FramebufferInfo) *Framebuffer {
	raw := wasmgl.CreateFramebuffer()
	wasmgl.BindFramebuffer(wasmgl.FRAMEBUFFER, raw)

	var drawBuffers []int
	for i, attachment := range info.ColorAttachments {
		if colorAttachment, ok := attachment.(*Texture); ok {
			attachmentID := wasmgl.COLOR_ATTACHMENT0 + int(i)
			wasmgl.FramebufferTexture2D(wasmgl.FRAMEBUFFER, attachmentID, wasmgl.TEXTURE_2D, colorAttachment.raw, 0)
			drawBuffers = append(drawBuffers, attachmentID)
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
		log.Error("Framebuffer is incomplete")
	}

	return &Framebuffer{
		raw: raw,
	}
}

var DefaultFramebuffer = &Framebuffer{
	raw: wasmgl.NilFramebuffer,
}

type Framebuffer struct {
	raw wasmgl.Framebuffer
}

func (f *Framebuffer) Release() {
	wasmgl.DeleteFramebuffer(f.raw)
	f.raw = wasmgl.NilFramebuffer
}
