package graphics

import (
	"encoding/binary"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking-js/webgl"
	"github.com/mokiat/lacking/data/buffer"
	"github.com/mokiat/wasmgl"
)

type quadMeshVertex struct {
	Position sprec.Vec2
}

func (v quadMeshVertex) Serialize(plotter *buffer.Plotter) {
	plotter.PlotFloat32(v.Position.X)
	plotter.PlotFloat32(v.Position.Y)
}

func newQuadMesh() *QuadMesh {
	return &QuadMesh{
		VertexBuffer: webgl.NewBuffer(),
		IndexBuffer:  webgl.NewBuffer(),
		VertexArray:  webgl.NewVertexArray(),
	}
}

type QuadMesh struct {
	VertexBuffer     *webgl.Buffer
	IndexBuffer      *webgl.Buffer
	VertexArray      *webgl.VertexArray
	Primitive        int
	IndexCount       int
	IndexOffsetBytes int
}

func (m *QuadMesh) Allocate() {
	m.Primitive = wasmgl.TRIANGLES // Do not do this in constructor
	m.IndexCount = 6
	m.IndexOffsetBytes = 0

	const vertexSize = 2 * 4
	vertexPlotter := buffer.NewPlotter(
		make([]byte, vertexSize*4),
		binary.LittleEndian,
	)

	quadMeshVertex{
		Position: sprec.NewVec2(-1.0, 1.0),
	}.Serialize(vertexPlotter)
	quadMeshVertex{
		Position: sprec.NewVec2(-1.0, -1.0),
	}.Serialize(vertexPlotter)
	quadMeshVertex{
		Position: sprec.NewVec2(1.0, -1.0),
	}.Serialize(vertexPlotter)
	quadMeshVertex{
		Position: sprec.NewVec2(1.0, 1.0),
	}.Serialize(vertexPlotter)

	const indexSize = 1 * 2
	indexPlotter := buffer.NewPlotter(
		make([]byte, indexSize*6),
		binary.LittleEndian,
	)

	indexPlotter.PlotUint16(0)
	indexPlotter.PlotUint16(1)
	indexPlotter.PlotUint16(2)

	indexPlotter.PlotUint16(0)
	indexPlotter.PlotUint16(2)
	indexPlotter.PlotUint16(3)

	m.VertexBuffer.Allocate(webgl.BufferAllocateInfo{
		BufferType: wasmgl.ARRAY_BUFFER,
		Dynamic:    false,
		Data:       vertexPlotter.Data(),
	})
	m.IndexBuffer.Allocate(webgl.BufferAllocateInfo{
		BufferType: wasmgl.ELEMENT_ARRAY_BUFFER,
		Dynamic:    false,
		Data:       indexPlotter.Data(),
	})
	m.VertexArray.Allocate(webgl.VertexArrayAllocateInfo{
		Attributes: []webgl.VertexArrayAttribute{
			{
				Buffer:         m.VertexBuffer,
				Index:          0,
				ComponentCount: 2,
				ComponentType:  wasmgl.FLOAT,
				Normalized:     false,
				StrideBytes:    vertexSize,
				OffsetBytes:    0,
			},
		},
		IndexBuffer: m.IndexBuffer,
	})
}

func (m *QuadMesh) Release() {
	m.VertexArray.Release()
	m.IndexBuffer.Release()
	m.VertexBuffer.Release()
}
