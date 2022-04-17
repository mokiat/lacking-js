package internal

import (
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/wasmgl"
)

func NewVertexBuffer(info render.BufferInfo) *Buffer {
	return newBuffer(info, wasmgl.ARRAY_BUFFER)
}

func NewIndexBuffer(info render.BufferInfo) *Buffer {
	return newBuffer(info, wasmgl.ELEMENT_ARRAY_BUFFER)
}

func NewPixelTransferBuffer(info render.BufferInfo) render.Buffer {
	return newBuffer(info, wasmgl.PIXEL_PACK_BUFFER)
}

func newBuffer(info render.BufferInfo, kind int) *Buffer {
	raw := wasmgl.CreateBuffer()
	wasmgl.BindBuffer(kind, raw)
	if info.Data != nil {
		wasmgl.BufferData(kind, len(info.Data), info.Data, glBufferUsage(info.Dynamic))
	} else {
		wasmgl.BufferData(kind, info.Size, nil, glBufferUsage(info.Dynamic))
	}
	result := &Buffer{
		raw:  raw,
		kind: kind,
	}
	result.id = buffers.Allocate(result)
	return result
}

type Buffer struct {
	render.BufferObject
	id   uint32
	raw  wasmgl.Buffer
	kind int
}

func (b *Buffer) Update(info render.BufferUpdateInfo) {
	wasmgl.BindBuffer(b.kind, b.raw)
	wasmgl.BufferSubData(b.kind, info.Offset, info.Data)
}

func (b *Buffer) Release() {
	buffers.Release(b.id)
	wasmgl.DeleteBuffer(b.raw)
	b.raw = wasmgl.NilBuffer
	b.kind = 0
	b.id = 0
}

func glBufferUsage(dynamic bool) int {
	if dynamic {
		return wasmgl.DYNAMIC_DRAW
	} else {
		return wasmgl.STATIC_DRAW
	}
}
