package webgl

import (
	"github.com/mokiat/wasmgl"
)

var (
	gl string // TODO: REMOVE
)

func NewVertexArray() *VertexArray {
	return &VertexArray{}
}

type VertexArray struct {
	raw wasmgl.VertexArray
}

func (a *VertexArray) Raw() wasmgl.VertexArray {
	return a.raw
}

func (a *VertexArray) Allocate(info VertexArrayAllocateInfo) {
	a.raw = wasmgl.CreateVertexArray()
	wasmgl.BindVertexArray(a.raw)

	var lastBuffer *Buffer
	for _, attribute := range info.Attributes {
		if attribute.Buffer != lastBuffer {
			attribute.Buffer.Use()
		}
		wasmgl.EnableVertexAttribArray(attribute.Index)
		wasmgl.VertexAttribPointer(attribute.Index, attribute.ComponentCount, attribute.ComponentType, attribute.Normalized, attribute.StrideBytes, attribute.OffsetBytes)
	}
	if info.IndexBuffer != nil {
		info.IndexBuffer.Use()
	}

	wasmgl.BindVertexArray(wasmgl.NilVertexArray)
}

func (a *VertexArray) Use() {
	wasmgl.BindVertexArray(a.raw)
}

func (a *VertexArray) Release() {
	wasmgl.DeleteVertexArray(a.raw)
	a.raw = wasmgl.VertexArray{}
}

type VertexArrayAllocateInfo struct {
	Attributes  []VertexArrayAttribute
	IndexBuffer *Buffer
}

func NewVertexArrayBufferBinding(buffer *Buffer, offsetBytes int, strideBytes int32) VertexArrayBufferBinding {
	return VertexArrayBufferBinding{
		VertexBuffer: buffer,
		OffsetBytes:  offsetBytes,
		StrideBytes:  strideBytes,
	}
}

type VertexArrayBufferBinding struct {
	VertexBuffer *Buffer
	OffsetBytes  int
	StrideBytes  int32
}

func NewVertexArrayAttribute(buffer *Buffer, index, compCount, compType int, norm bool, strideBytes, offsetBytes int) VertexArrayAttribute {
	return VertexArrayAttribute{
		Buffer:         buffer,
		Index:          index,
		ComponentCount: compCount,
		ComponentType:  compType,
		Normalized:     norm,
		StrideBytes:    strideBytes,
		OffsetBytes:    offsetBytes,
	}
}

type VertexArrayAttribute struct {
	Buffer         *Buffer
	Index          int
	ComponentCount int
	ComponentType  int
	Normalized     bool
	StrideBytes    int
	OffsetBytes    int
}
