package internal

import (
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/wasmgl"
)

func NewVertexBuffer(info render.BufferInfo) *Buffer {
	if glLogger.IsDebugEnabled() {
		defer trackError("Error creating vertex buffer (%v)", info.Label)()
	}
	return newBuffer(info, wasmgl.ARRAY_BUFFER)
}

func NewIndexBuffer(info render.BufferInfo) *Buffer {
	if glLogger.IsDebugEnabled() {
		defer trackError("Error creating index buffer (%v)", info.Label)()
	}
	return newBuffer(info, wasmgl.ELEMENT_ARRAY_BUFFER)
}

func NewPixelTransferBuffer(info render.BufferInfo) render.Buffer {
	if glLogger.IsDebugEnabled() {
		defer trackError("Error creating pixel transfer buffer (%v)", info.Label)()
	}

	raw := wasmgl.CreateBuffer()
	wasmgl.BindBuffer(wasmgl.PIXEL_PACK_BUFFER, raw)
	wasmgl.BufferData(wasmgl.PIXEL_PACK_BUFFER, wasmgl.GLintptr(info.Size), nil, wasmgl.DYNAMIC_READ)
	result := &Buffer{
		label: info.Label,
		raw:   raw,
		kind:  wasmgl.PIXEL_PACK_BUFFER,
	}
	result.id = buffers.Allocate(result)
	return result
}

func NewUniformBuffer(info render.BufferInfo) render.Buffer {
	if glLogger.IsDebugEnabled() {
		defer trackError("Error creating uniform buffer (%v)", info.Label)()
	}
	return newBuffer(info, wasmgl.UNIFORM_BUFFER)
}

func newBuffer(info render.BufferInfo, kind wasmgl.GLenum) *Buffer {
	raw := wasmgl.CreateBuffer()
	wasmgl.BindBuffer(kind, raw)
	if info.Data != nil {
		wasmgl.BufferData(kind, wasmgl.GLintptr(len(info.Data)), info.Data, glBufferUsage(info.Dynamic))
	} else {
		wasmgl.BufferData(kind, wasmgl.GLintptr(info.Size), nil, glBufferUsage(info.Dynamic))
	}
	result := &Buffer{
		label: info.Label,
		raw:   raw,
		kind:  kind,
	}
	result.id = buffers.Allocate(result)
	return result
}

type Buffer struct {
	render.BufferMarker

	label string
	id    uint32
	raw   wasmgl.Buffer
	kind  wasmgl.GLenum
}

func (b *Buffer) Label() string {
	return b.label
}

func (b *Buffer) Release() {
	buffers.Release(b.id)
	wasmgl.DeleteBuffer(b.raw)
	b.raw = wasmgl.NilBuffer
	b.kind = 0
	b.id = 0
}

func glBufferUsage(dynamic bool) wasmgl.GLenum {
	if dynamic {
		return wasmgl.DYNAMIC_DRAW
	} else {
		return wasmgl.STATIC_DRAW
	}
}
