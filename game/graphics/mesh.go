package graphics

import (
	"github.com/mokiat/lacking-js/game/graphics/internal"
	"github.com/mokiat/lacking-js/webgl"
	"github.com/mokiat/lacking/game/graphics"
)

var _ graphics.MeshTemplate = (*MeshTemplate)(nil)

type SubMeshTemplate struct {
	material         *Material
	primitive        int
	indexCount       int
	indexOffsetBytes int
	indexType        int
}

type MeshTemplate struct {
	vertexBuffer *webgl.Buffer
	indexBuffer  *webgl.Buffer
	vertexArray  *webgl.VertexArray
	subMeshes    []SubMeshTemplate
}

func (t *MeshTemplate) Delete() {
	t.vertexArray.Release()
	t.indexBuffer.Release()
	t.vertexBuffer.Release()
	t.subMeshes = nil
}

var _ graphics.Mesh = (*Mesh)(nil)

type Mesh struct {
	internal.Node

	scene *Scene
	prev  *Mesh
	next  *Mesh

	template *MeshTemplate
}

func (m *Mesh) Delete() {
	m.scene.detachMesh(m)
	m.scene.cacheMesh(m)
	m.scene = nil
}
