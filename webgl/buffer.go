package webgl

import "github.com/mokiat/wasmgl"

func NewBuffer() *Buffer {
	return &Buffer{}
}

type Buffer struct {
	raw  wasmgl.Buffer
	kind int
}

func (b *Buffer) Raw() wasmgl.Buffer {
	return b.raw
}

func (b *Buffer) Allocate(info BufferAllocateInfo) {
	b.kind = info.BufferType
	b.raw = wasmgl.CreateBuffer()
	wasmgl.BindBuffer(b.kind, b.raw)
	wasmgl.BufferData(b.kind, info.Data, info.usage())
}

func (b *Buffer) Update(info BufferUpdateInfo) {
	wasmgl.BindBuffer(b.kind, b.raw)
	wasmgl.BufferSubData(b.kind, info.OffsetBytes, info.Data)
}

func (b *Buffer) Use() {
	wasmgl.BindBuffer(b.kind, b.raw)
}

func (b *Buffer) Release() {
	wasmgl.DeleteBuffer(b.raw)
	b.raw = wasmgl.Buffer{}
	b.kind = 0
}

type BufferAllocateInfo struct {
	BufferType int
	Dynamic    bool
	Data       []byte
}

func (i BufferAllocateInfo) usage() int {
	if i.Dynamic {
		return wasmgl.DYNAMIC_DRAW
	} else {
		return wasmgl.STATIC_DRAW
	}
}

type BufferUpdateInfo struct {
	Data        []byte
	OffsetBytes int
}
