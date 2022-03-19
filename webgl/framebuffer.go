package webgl

import (
	"fmt"

	"github.com/mokiat/wasmgl"
)

var defaultFramebuffer = &Framebuffer{
	raw: wasmgl.NullFramebuffer,
}

func DefaultFramebuffer() *Framebuffer {
	return defaultFramebuffer
}

func NewFramebuffer() *Framebuffer {
	return &Framebuffer{
		raw: wasmgl.NullFramebuffer,
	}
}

type Framebuffer struct {
	raw wasmgl.Framebuffer
}

func (b *Framebuffer) Raw() wasmgl.Framebuffer {
	return b.raw
}

func (b *Framebuffer) Allocate(info FramebufferAllocateInfo) {
	b.raw = wasmgl.CreateFramebuffer()
	wasmgl.BindFramebuffer(wasmgl.FRAMEBUFFER, b.raw)

	var drawBuffers []int
	for i, colorAttachment := range info.ColorAttachments {
		if colorAttachment != nil {
			attachmentID := wasmgl.COLOR_ATTACHMENT0 + int(i)
			wasmgl.FramebufferTexture2D(wasmgl.FRAMEBUFFER, attachmentID, wasmgl.TEXTURE_2D, colorAttachment.Raw(), 0)
			drawBuffers = append(drawBuffers, attachmentID)
		}
	}
	if info.DepthStencilAttachment != nil {
		wasmgl.FramebufferTexture2D(wasmgl.FRAMEBUFFER, wasmgl.DEPTH_STENCIL_ATTACHMENT, wasmgl.TEXTURE_2D, info.DepthStencilAttachment.Raw(), 0)
	} else {
		if info.DepthAttachment != nil {
			wasmgl.FramebufferTexture2D(wasmgl.FRAMEBUFFER, wasmgl.DEPTH_ATTACHMENT, wasmgl.TEXTURE_2D, info.DepthAttachment.Raw(), 0)
		}
		if info.StencilAttachment != nil {
			wasmgl.FramebufferTexture2D(wasmgl.FRAMEBUFFER, wasmgl.STENCIL_ATTACHMENT, wasmgl.TEXTURE_2D, info.StencilAttachment.Raw(), 0)
		}
	}
	wasmgl.DrawBuffers(drawBuffers)

	if wasmgl.CheckFramebufferStatus(wasmgl.FRAMEBUFFER) != wasmgl.FRAMEBUFFER_COMPLETE {
		panic(fmt.Errorf("framebuffer is incomplete"))
	}
}

func (b *Framebuffer) Use() {
	wasmgl.BindFramebuffer(wasmgl.FRAMEBUFFER, b.raw)
}

func (b *Framebuffer) Release() {
	wasmgl.DeleteFramebuffer(b.raw)
	b.raw = wasmgl.NullFramebuffer
}

type FramebufferAllocateInfo struct {
	ColorAttachments       []*Texture
	DepthAttachment        *Texture
	StencilAttachment      *Texture
	DepthStencilAttachment *Texture
}
