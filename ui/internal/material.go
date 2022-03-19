package internal

import (
	"github.com/mokiat/lacking-js/webgl"
	"github.com/mokiat/wasmgl"
)

func newMaterial(vertexSrc, fragmentSrc func() string) *Material {
	return &Material{
		vertexSrc:   vertexSrc,
		fragmentSrc: fragmentSrc,

		program: webgl.NewProgram(),
	}
}

type Material struct {
	vertexSrc   func() string
	fragmentSrc func() string

	program                        *webgl.Program
	transformMatrixLocation        wasmgl.UniformLocation
	textureTransformMatrixLocation wasmgl.UniformLocation
	projectionMatrixLocation       wasmgl.UniformLocation
	clipDistancesLocation          wasmgl.UniformLocation
	textureLocation                wasmgl.UniformLocation
	colorLocation                  wasmgl.UniformLocation
}

func (m *Material) Allocate() {
	vertexShader := webgl.NewShader()
	vertexShader.Allocate(webgl.ShaderAllocateInfo{
		ShaderType: wasmgl.VERTEX_SHADER,
		SourceCode: m.vertexSrc(),
	})
	defer func() {
		vertexShader.Release()
	}()

	fragmentShader := webgl.NewShader()
	fragmentShader.Allocate(webgl.ShaderAllocateInfo{
		ShaderType: wasmgl.FRAGMENT_SHADER,
		SourceCode: m.fragmentSrc(),
	})
	defer func() {
		fragmentShader.Release()
	}()

	m.program.Allocate(webgl.ProgramAllocateInfo{
		VertexShader:   vertexShader,
		FragmentShader: fragmentShader,
	})

	m.transformMatrixLocation = m.program.UniformLocation("transformMatrixIn")
	m.textureTransformMatrixLocation = m.program.UniformLocation("textureTransformMatrixIn")
	m.projectionMatrixLocation = m.program.UniformLocation("projectionMatrixIn")
	m.clipDistancesLocation = m.program.UniformLocation("clipDistancesIn")
	m.textureLocation = m.program.UniformLocation("textureIn")
	m.colorLocation = m.program.UniformLocation("colorIn")
}

func (m *Material) Release() {
	m.program.Release()
}
