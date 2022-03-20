package internal

import (
	"encoding/binary"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking-js/webgl"
	"github.com/mokiat/lacking/data/buffer"
	"github.com/mokiat/wasmgl"
)

const (
	textPositionAttribIndex = 0
	textTexCoordAttribIndex = 1

	textMeshVertexSize = 2*4 + 2*4
)

func newTextMesh(vertexCount int) *TextMesh {
	data := make([]byte, vertexCount*textMeshVertexSize)
	return &TextMesh{
		vertexData:    data,
		vertexPlotter: buffer.NewPlotter(data, binary.LittleEndian),
		vertexBuffer:  webgl.NewBuffer(),
		vertexArray:   webgl.NewVertexArray(),
	}
}

type TextMesh struct {
	vertexData    []byte
	vertexPlotter *buffer.Plotter
	vertexOffset  int
	vertexBuffer  *webgl.Buffer
	vertexArray   *webgl.VertexArray
}

func (m *TextMesh) Allocate() {
	m.vertexBuffer.Allocate(webgl.BufferAllocateInfo{
		BufferType: wasmgl.ARRAY_BUFFER,
		Dynamic:    true,
		Data:       m.vertexData,
	})

	m.vertexArray.Allocate(webgl.VertexArrayAllocateInfo{
		Attributes: []webgl.VertexArrayAttribute{
			{
				Buffer:         m.vertexBuffer,
				Index:          textPositionAttribIndex,
				ComponentCount: 2,
				ComponentType:  wasmgl.FLOAT,
				Normalized:     false,
				StrideBytes:    textMeshVertexSize,
				OffsetBytes:    0,
			},
			{
				Buffer:         m.vertexBuffer,
				Index:          textTexCoordAttribIndex,
				ComponentCount: 2,
				ComponentType:  wasmgl.FLOAT,
				Normalized:     false,
				StrideBytes:    textMeshVertexSize,
				OffsetBytes:    2 * 4,
			},
		},
	})
}

func (m *TextMesh) Release() {
	m.vertexArray.Release()
	m.vertexBuffer.Release()
}

func (m *TextMesh) Update() {
	if length := m.vertexPlotter.Offset(); length > 0 {
		m.vertexBuffer.Update(webgl.BufferUpdateInfo{
			Data:        m.vertexData[:m.vertexPlotter.Offset()],
			OffsetBytes: 0,
		})
	}
}

func (m *TextMesh) Reset() {
	m.vertexOffset = 0
	m.vertexPlotter.Rewind()
}

func (m *TextMesh) Offset() int {
	return m.vertexOffset
}

func (m *TextMesh) Append(vertex TextVertex) {
	m.vertexPlotter.PlotFloat32(vertex.position.X)
	m.vertexPlotter.PlotFloat32(vertex.position.Y)
	m.vertexPlotter.PlotFloat32(vertex.texCoord.X)
	m.vertexPlotter.PlotFloat32(vertex.texCoord.Y)
	m.vertexOffset++
}

type TextVertex struct {
	position sprec.Vec2
	texCoord sprec.Vec2
}
