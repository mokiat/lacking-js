package internal

import (
	"fmt"

	"github.com/mokiat/lacking/render"
	"github.com/mokiat/wasmgl"
)

func NewVertexArray(info render.VertexArrayInfo) *VertexArray {
	if glLogger.IsDebugEnabled() {
		defer trackError("Error creating vertex array (%v)", info.Label)()
	}

	raw := wasmgl.CreateVertexArray()
	wasmgl.BindVertexArray(raw)
	for _, attribute := range info.Attributes {
		binding := info.Bindings[attribute.Binding]
		if vertexBuffer, ok := binding.VertexBuffer.(*Buffer); ok {
			wasmgl.BindBuffer(vertexBuffer.kind, vertexBuffer.raw)
		}
		wasmgl.EnableVertexAttribArray(wasmgl.GLuint(attribute.Location))
		count, compType, normalized, integer := glAttribParams(attribute.Format)
		if integer {
			wasmgl.VertexAttribIPointer(wasmgl.GLuint(attribute.Location), count, compType, wasmgl.GLsizei(binding.Stride), wasmgl.GLintptr(attribute.Offset))
		} else {
			wasmgl.VertexAttribPointer(wasmgl.GLuint(attribute.Location), count, compType, normalized, wasmgl.GLsizei(binding.Stride), wasmgl.GLintptr(attribute.Offset))
		}
	}
	if indexBuffer, ok := info.IndexBuffer.(*Buffer); ok {
		wasmgl.BindBuffer(indexBuffer.kind, indexBuffer.raw)
	}
	wasmgl.BindVertexArray(wasmgl.NilVertexArray)

	result := &VertexArray{
		label:       info.Label,
		raw:         raw,
		indexFormat: glIndexFormat(info.IndexFormat),
	}
	result.id = vertexArrays.Allocate(result)
	return result
}

type VertexArray struct {
	render.VertexArrayMarker

	label       string
	id          uint32
	raw         wasmgl.VertexArray
	indexFormat wasmgl.GLenum
}

func (a *VertexArray) Label() string {
	return a.label
}

func (a *VertexArray) Release() {
	vertexArrays.Release(a.id)
	wasmgl.DeleteVertexArray(a.raw)
	a.raw = wasmgl.NilVertexArray
	a.id = 0
}

func glAttribParams(format render.VertexAttributeFormat) (wasmgl.GLint, wasmgl.GLenum, wasmgl.GLboolean, bool) {
	switch format {
	case render.VertexAttributeFormatR32F:
		return 1, wasmgl.FLOAT, false, false
	case render.VertexAttributeFormatRG32F:
		return 2, wasmgl.FLOAT, false, false
	case render.VertexAttributeFormatRGB32F:
		return 3, wasmgl.FLOAT, false, false
	case render.VertexAttributeFormatRGBA32F:
		return 4, wasmgl.FLOAT, false, false

	case render.VertexAttributeFormatR16F:
		return 1, wasmgl.HALF_FLOAT, false, false
	case render.VertexAttributeFormatRG16F:
		return 2, wasmgl.HALF_FLOAT, false, false
	case render.VertexAttributeFormatRGB16F:
		return 3, wasmgl.HALF_FLOAT, false, false
	case render.VertexAttributeFormatRGBA16F:
		return 4, wasmgl.HALF_FLOAT, false, false

	case render.VertexAttributeFormatR16S:
		return 1, wasmgl.SHORT, false, false
	case render.VertexAttributeFormatRG16S:
		return 2, wasmgl.SHORT, false, false
	case render.VertexAttributeFormatRGB16S:
		return 3, wasmgl.SHORT, false, false
	case render.VertexAttributeFormatRGBA16S:
		return 4, wasmgl.SHORT, false, false

	case render.VertexAttributeFormatR16SN:
		return 1, wasmgl.SHORT, true, false
	case render.VertexAttributeFormatRG16SN:
		return 2, wasmgl.SHORT, true, false
	case render.VertexAttributeFormatRGB16SN:
		return 3, wasmgl.SHORT, true, false
	case render.VertexAttributeFormatRGBA16SN:
		return 4, wasmgl.SHORT, true, false

	case render.VertexAttributeFormatR16U:
		return 1, wasmgl.UNSIGNED_SHORT, false, false
	case render.VertexAttributeFormatRG16U:
		return 2, wasmgl.UNSIGNED_SHORT, false, false
	case render.VertexAttributeFormatRGB16U:
		return 3, wasmgl.UNSIGNED_SHORT, false, false
	case render.VertexAttributeFormatRGBA16U:
		return 4, wasmgl.UNSIGNED_SHORT, false, false

	case render.VertexAttributeFormatR16UN:
		return 1, wasmgl.UNSIGNED_SHORT, true, false
	case render.VertexAttributeFormatRG16UN:
		return 2, wasmgl.UNSIGNED_SHORT, true, false
	case render.VertexAttributeFormatRGB16UN:
		return 3, wasmgl.UNSIGNED_SHORT, true, false
	case render.VertexAttributeFormatRGBA16UN:
		return 4, wasmgl.UNSIGNED_SHORT, true, false

	case render.VertexAttributeFormatR8S:
		return 1, wasmgl.BYTE, false, false
	case render.VertexAttributeFormatRG8S:
		return 2, wasmgl.BYTE, false, false
	case render.VertexAttributeFormatRGB8S:
		return 3, wasmgl.BYTE, false, false
	case render.VertexAttributeFormatRGBA8S:
		return 4, wasmgl.BYTE, false, false

	case render.VertexAttributeFormatR8SN:
		return 1, wasmgl.BYTE, true, false
	case render.VertexAttributeFormatRG8SN:
		return 2, wasmgl.BYTE, true, false
	case render.VertexAttributeFormatRGB8SN:
		return 3, wasmgl.BYTE, true, false
	case render.VertexAttributeFormatRGBA8SN:
		return 4, wasmgl.BYTE, true, false

	case render.VertexAttributeFormatR8U:
		return 1, wasmgl.UNSIGNED_BYTE, false, false
	case render.VertexAttributeFormatRG8U:
		return 2, wasmgl.UNSIGNED_BYTE, false, false
	case render.VertexAttributeFormatRGB8U:
		return 3, wasmgl.UNSIGNED_BYTE, false, false
	case render.VertexAttributeFormatRGBA8U:
		return 4, wasmgl.UNSIGNED_BYTE, false, false

	case render.VertexAttributeFormatR8UN:
		return 1, wasmgl.UNSIGNED_BYTE, true, false
	case render.VertexAttributeFormatRG8UN:
		return 2, wasmgl.UNSIGNED_BYTE, true, false
	case render.VertexAttributeFormatRGB8UN:
		return 3, wasmgl.UNSIGNED_BYTE, true, false
	case render.VertexAttributeFormatRGBA8UN:
		return 4, wasmgl.UNSIGNED_BYTE, true, false

	case render.VertexAttributeFormatR8IU:
		return 1, wasmgl.UNSIGNED_BYTE, false, true
	case render.VertexAttributeFormatRG8IU:
		return 2, wasmgl.UNSIGNED_BYTE, false, true
	case render.VertexAttributeFormatRGB8IU:
		return 3, wasmgl.UNSIGNED_BYTE, false, true
	case render.VertexAttributeFormatRGBA8IU:
		return 4, wasmgl.UNSIGNED_BYTE, false, true

	default:
		panic(fmt.Errorf("unknown attribute format: %d", format))
	}
}

func glIndexFormat(format render.IndexFormat) wasmgl.GLenum {
	switch format {
	case render.IndexFormatUnsignedU16:
		return wasmgl.UNSIGNED_SHORT
	case render.IndexFormatUnsignedU32:
		return wasmgl.UNSIGNED_INT
	default:
		panic(fmt.Errorf("unknown index format: %d", format))
	}
}
