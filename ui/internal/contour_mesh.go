package internal

import (
	"encoding/binary"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking-js/webgl"
	"github.com/mokiat/lacking/data/buffer"
	"github.com/mokiat/wasmgl"
)

const (
	contourPositionAttribIndex = 0
	contourColorAttribIndex    = 2

	contourMeshVertexSize = 2*4 + 1*4
)

func newContourMesh(vertexCount int) *ContourMesh {
	data := make([]byte, vertexCount*contourMeshVertexSize)
	return &ContourMesh{
		vertexData:    data,
		vertexPlotter: buffer.NewPlotter(data, binary.LittleEndian),
		vertexBuffer:  webgl.NewBuffer(),
		vertexArray:   webgl.NewVertexArray(),
	}
}

type ContourMesh struct {
	vertexData    []byte
	vertexPlotter *buffer.Plotter
	vertexOffset  int
	vertexBuffer  *webgl.Buffer
	vertexArray   *webgl.VertexArray
}

func (m *ContourMesh) Allocate() {
	m.vertexBuffer.Allocate(webgl.BufferAllocateInfo{
		BufferType: wasmgl.ARRAY_BUFFER,
		Dynamic:    true,
		Data:       m.vertexData,
	})

	m.vertexArray.Allocate(webgl.VertexArrayAllocateInfo{
		Attributes: []webgl.VertexArrayAttribute{
			{
				Buffer:         m.vertexBuffer,
				Index:          contourPositionAttribIndex,
				ComponentCount: 2,
				ComponentType:  wasmgl.FLOAT,
				Normalized:     false,
				StrideBytes:    contourMeshVertexSize,
				OffsetBytes:    0,
			},
			{
				Buffer:         m.vertexBuffer,
				Index:          contourColorAttribIndex,
				ComponentCount: 4,
				ComponentType:  wasmgl.UNSIGNED_BYTE,
				Normalized:     true,
				StrideBytes:    contourMeshVertexSize,
				OffsetBytes:    2 * 4,
			},
		},
	})
}

func (m *ContourMesh) Release() {
	m.vertexArray.Release()
	m.vertexBuffer.Release()
}

func (m *ContourMesh) Update() {
	if length := m.vertexPlotter.Offset(); length > 0 {
		m.vertexBuffer.Update(webgl.BufferUpdateInfo{
			Data:        m.vertexData[:m.vertexPlotter.Offset()],
			OffsetBytes: 0,
		})
	}
}

func (m *ContourMesh) Reset() {
	m.vertexOffset = 0
	m.vertexPlotter.Rewind()
}

func (m *ContourMesh) Offset() int {
	return m.vertexOffset
}

func (m *ContourMesh) Append(vertex ContourVertex) {
	m.vertexPlotter.PlotFloat32(vertex.position.X)
	m.vertexPlotter.PlotFloat32(vertex.position.Y)
	m.vertexPlotter.PlotByte(byte(vertex.color.X * 255))
	m.vertexPlotter.PlotByte(byte(vertex.color.Y * 255))
	m.vertexPlotter.PlotByte(byte(vertex.color.Z * 255))
	m.vertexPlotter.PlotByte(byte(vertex.color.W * 255))
	m.vertexOffset++
}

type ContourVertex struct {
	position sprec.Vec2
	color    sprec.Vec4
}
