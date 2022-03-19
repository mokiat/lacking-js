package graphics

import (
	"encoding/binary"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking-js/webgl"
	"github.com/mokiat/lacking/data/buffer"
	"github.com/mokiat/wasmgl"
)

type skyboxMeshVertex struct {
	Position sprec.Vec3
}

func (v skyboxMeshVertex) Serialize(plotter *buffer.Plotter) {
	plotter.PlotFloat32(v.Position.X)
	plotter.PlotFloat32(v.Position.Y)
	plotter.PlotFloat32(v.Position.Z)
}

func newSkyboxMesh() *SkyboxMesh {
	return &SkyboxMesh{
		VertexBuffer: webgl.NewBuffer(),
		IndexBuffer:  webgl.NewBuffer(),
		VertexArray:  webgl.NewVertexArray(),
	}
}

type SkyboxMesh struct {
	VertexBuffer     *webgl.Buffer
	IndexBuffer      *webgl.Buffer
	VertexArray      *webgl.VertexArray
	Primitive        int
	IndexCount       int
	IndexOffsetBytes int
}

func (m *SkyboxMesh) Allocate() {
	m.Primitive = wasmgl.TRIANGLES // Do not do this in constructor!
	m.IndexCount = 36
	m.IndexOffsetBytes = 0

	const vertexSize = 3 * 4
	vertexPlotter := buffer.NewPlotter(
		make([]byte, vertexSize*8),
		binary.LittleEndian,
	)

	skyboxMeshVertex{
		Position: sprec.NewVec3(-1.0, 1.0, 1.0),
	}.Serialize(vertexPlotter)
	skyboxMeshVertex{
		Position: sprec.NewVec3(-1.0, -1.0, 1.0),
	}.Serialize(vertexPlotter)
	skyboxMeshVertex{
		Position: sprec.NewVec3(1.0, -1.0, 1.0),
	}.Serialize(vertexPlotter)
	skyboxMeshVertex{
		Position: sprec.NewVec3(1.0, 1.0, 1.0),
	}.Serialize(vertexPlotter)

	skyboxMeshVertex{
		Position: sprec.NewVec3(-1.0, 1.0, -1.0),
	}.Serialize(vertexPlotter)
	skyboxMeshVertex{
		Position: sprec.NewVec3(-1.0, -1.0, -1.0),
	}.Serialize(vertexPlotter)
	skyboxMeshVertex{
		Position: sprec.NewVec3(1.0, -1.0, -1.0),
	}.Serialize(vertexPlotter)
	skyboxMeshVertex{
		Position: sprec.NewVec3(1.0, 1.0, -1.0),
	}.Serialize(vertexPlotter)

	const indexSize = 1 * 2
	indexPlotter := buffer.NewPlotter(
		make([]byte, indexSize*36),
		binary.LittleEndian,
	)

	indexPlotter.PlotUint16(3)
	indexPlotter.PlotUint16(2)
	indexPlotter.PlotUint16(1)

	indexPlotter.PlotUint16(3)
	indexPlotter.PlotUint16(1)
	indexPlotter.PlotUint16(0)

	indexPlotter.PlotUint16(0)
	indexPlotter.PlotUint16(1)
	indexPlotter.PlotUint16(5)

	indexPlotter.PlotUint16(0)
	indexPlotter.PlotUint16(5)
	indexPlotter.PlotUint16(4)

	indexPlotter.PlotUint16(7)
	indexPlotter.PlotUint16(6)
	indexPlotter.PlotUint16(2)

	indexPlotter.PlotUint16(7)
	indexPlotter.PlotUint16(2)
	indexPlotter.PlotUint16(3)

	indexPlotter.PlotUint16(4)
	indexPlotter.PlotUint16(5)
	indexPlotter.PlotUint16(6)

	indexPlotter.PlotUint16(4)
	indexPlotter.PlotUint16(6)
	indexPlotter.PlotUint16(7)

	indexPlotter.PlotUint16(5)
	indexPlotter.PlotUint16(1)
	indexPlotter.PlotUint16(2)

	indexPlotter.PlotUint16(5)
	indexPlotter.PlotUint16(2)
	indexPlotter.PlotUint16(6)

	indexPlotter.PlotUint16(0)
	indexPlotter.PlotUint16(4)
	indexPlotter.PlotUint16(7)

	indexPlotter.PlotUint16(0)
	indexPlotter.PlotUint16(7)
	indexPlotter.PlotUint16(3)

	vertexBufferInfo := webgl.BufferAllocateInfo{
		BufferType: wasmgl.ARRAY_BUFFER,
		Dynamic:    false,
		Data:       vertexPlotter.Data(),
	}
	m.VertexBuffer.Allocate(vertexBufferInfo)

	indexBufferInfo := webgl.BufferAllocateInfo{
		BufferType: wasmgl.ELEMENT_ARRAY_BUFFER,
		Dynamic:    false,
		Data:       indexPlotter.Data(),
	}
	m.IndexBuffer.Allocate(indexBufferInfo)

	vertexArrayInfo := webgl.VertexArrayAllocateInfo{
		Attributes: []webgl.VertexArrayAttribute{
			{
				Buffer:         m.VertexBuffer,
				Index:          0,
				ComponentCount: 3,
				ComponentType:  wasmgl.FLOAT,
				Normalized:     false,
				StrideBytes:    3 * 4,
				OffsetBytes:    0,
			},
		},
		IndexBuffer: m.IndexBuffer,
	}
	m.VertexArray.Allocate(vertexArrayInfo)
}

func (m *SkyboxMesh) Release() {
	m.VertexArray.Release()
	m.IndexBuffer.Release()
	m.VertexBuffer.Release()
}
