package internal

import (
	"encoding/binary"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking-js/webgl"
	"github.com/mokiat/lacking/data/buffer"
	"github.com/mokiat/wasmgl"
)

const (
	shapePositionAttribIndex = 0

	shapeMeshVertexSize = 2 * 4
)

func newShapeMesh(vertexCount int) *ShapeMesh {
	data := make([]byte, vertexCount*shapeMeshVertexSize)
	return &ShapeMesh{
		vertexData:    data,
		vertexPlotter: buffer.NewPlotter(data, binary.LittleEndian),
		vertexBuffer:  webgl.NewBuffer(),
		vertexArray:   webgl.NewVertexArray(),
	}
}

type ShapeMesh struct {
	vertexData    []byte
	vertexPlotter *buffer.Plotter
	vertexOffset  int
	vertexBuffer  *webgl.Buffer
	vertexArray   *webgl.VertexArray
}

func (m *ShapeMesh) Allocate() {
	m.vertexBuffer.Allocate(webgl.BufferAllocateInfo{
		BufferType: wasmgl.ARRAY_BUFFER,
		Dynamic:    true,
		Data:       m.vertexData,
	})

	m.vertexArray.Allocate(webgl.VertexArrayAllocateInfo{
		Attributes: []webgl.VertexArrayAttribute{
			{
				Buffer:         m.vertexBuffer,
				Index:          shapePositionAttribIndex,
				ComponentCount: 2,
				ComponentType:  wasmgl.FLOAT,
				Normalized:     false,
				StrideBytes:    shapeMeshVertexSize,
				OffsetBytes:    0,
			},
		},
	})
}

func (m *ShapeMesh) Release() {
	m.vertexArray.Release()
	m.vertexBuffer.Release()
}

func (m *ShapeMesh) Update() {
	if length := m.vertexPlotter.Offset(); length > 0 {
		m.vertexBuffer.Update(webgl.BufferUpdateInfo{
			Data:        m.vertexData[:m.vertexPlotter.Offset()],
			OffsetBytes: 0,
		})
	}
}

func (m *ShapeMesh) Reset() {
	m.vertexOffset = 0
	m.vertexPlotter.Rewind()
}

func (m *ShapeMesh) Offset() int {
	return m.vertexOffset
}

func (m *ShapeMesh) Append(vertex ShapeVertex) {
	m.vertexPlotter.PlotFloat32(vertex.position.X)
	m.vertexPlotter.PlotFloat32(vertex.position.Y)
	m.vertexOffset++
}

type ShapeVertex struct {
	position sprec.Vec2
}
