package internal

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking-js/webgl"
)

type SubMesh struct {
	surface                Surface
	material               *Material
	vertexArray            *webgl.VertexArray
	transformMatrix        sprec.Mat4
	textureTransformMatrix sprec.Mat4
	texture                *webgl.TwoDTexture
	color                  sprec.Vec4
	vertexOffset           int
	vertexCount            int
	primitive              int
	culling                bool
	cullFace               int
	clipBounds             sprec.Vec4
	skipColor              bool
	stencil                bool
	stencilCfg             stencilConfig
}

type stencilConfig struct {
	stencilOpFront   stencilOp
	stencilOpBack    stencilOp
	stencilFuncFront stencilFunc
	stencilFuncBack  stencilFunc
}

type stencilOp struct {
	sfail  int
	dpfail int
	dppass int
}

type stencilFunc struct {
	fn   int
	ref  int
	mask int
}
