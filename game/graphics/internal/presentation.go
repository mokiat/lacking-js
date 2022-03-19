package internal

import (
	"github.com/mokiat/lacking-js/webgl"
	"github.com/mokiat/wasmgl"
)

type Presentation struct {
	Program *webgl.Program
}

func (p *Presentation) Delete() {
	p.Program.Release()
}

type PostprocessingPresentation struct {
	Presentation

	FramebufferDraw0Location wasmgl.UniformLocation

	ExposureLocation wasmgl.UniformLocation
}

func NewPostprocessingPresentation(vertexSrc, fragmentSrc string) *PostprocessingPresentation {
	program := buildProgram(vertexSrc, fragmentSrc)
	return &PostprocessingPresentation{
		Presentation: Presentation{
			Program: program,
		},
		FramebufferDraw0Location: program.UniformLocation("fbColor0TextureIn"),
		ExposureLocation:         program.UniformLocation("exposureIn"),
	}
}

type SkyboxPresentation struct {
	Presentation

	ProjectionMatrixLocation wasmgl.UniformLocation
	ViewMatrixLocation       wasmgl.UniformLocation

	AlbedoCubeTextureLocation wasmgl.UniformLocation
	AlbedoColorLocation       wasmgl.UniformLocation
}

func NewSkyboxPresentation(vertexSrc, fragmentSrc string) *SkyboxPresentation {
	program := buildProgram(vertexSrc, fragmentSrc)
	return &SkyboxPresentation{
		Presentation: Presentation{
			Program: program,
		},
		ProjectionMatrixLocation:  program.UniformLocation("projectionMatrixIn"),
		ViewMatrixLocation:        program.UniformLocation("viewMatrixIn"),
		AlbedoCubeTextureLocation: program.UniformLocation("albedoCubeTextureIn"),
		AlbedoColorLocation:       program.UniformLocation("albedoColorIn"),
	}
}

type ShadowPresentation struct {
	Presentation
}

func NewShadowPresentation(vertexSrc, fragmentSrc string) *ShadowPresentation {
	program := buildProgram(vertexSrc, fragmentSrc)
	return &ShadowPresentation{
		Presentation: Presentation{
			Program: program,
		},
	}
}

type GeometryPresentation struct {
	Presentation

	ProjectionMatrixLocation wasmgl.UniformLocation
	ModelMatrixLocation      wasmgl.UniformLocation
	ViewMatrixLocation       wasmgl.UniformLocation

	MetalnessLocation wasmgl.UniformLocation
	RoughnessLocation wasmgl.UniformLocation

	AlbedoColorLocation   wasmgl.UniformLocation
	AlbedoTextureLocation wasmgl.UniformLocation
}

func NewGeometryPresentation(vertexSrc, fragmentSrc string) *GeometryPresentation {
	program := buildProgram(vertexSrc, fragmentSrc)
	return &GeometryPresentation{
		Presentation: Presentation{
			Program: program,
		},
		ProjectionMatrixLocation: program.UniformLocation("projectionMatrixIn"),
		ModelMatrixLocation:      program.UniformLocation("modelMatrixIn"),
		ViewMatrixLocation:       program.UniformLocation("viewMatrixIn"),
		MetalnessLocation:        program.UniformLocation("metalnessIn"),
		RoughnessLocation:        program.UniformLocation("roughnessIn"),
		AlbedoColorLocation:      program.UniformLocation("albedoColorIn"),
		AlbedoTextureLocation:    program.UniformLocation("albedoTwoDTextureIn"),
	}
}

type LightingPresentation struct {
	Presentation

	FramebufferDraw0Location wasmgl.UniformLocation
	FramebufferDraw1Location wasmgl.UniformLocation
	FramebufferDepthLocation wasmgl.UniformLocation

	ProjectionMatrixLocation wasmgl.UniformLocation
	CameraMatrixLocation     wasmgl.UniformLocation
	ViewMatrixLocation       wasmgl.UniformLocation

	ReflectionTextureLocation wasmgl.UniformLocation
	RefractionTextureLocation wasmgl.UniformLocation

	LightDirection wasmgl.UniformLocation
	LightIntensity wasmgl.UniformLocation
}

func NewLightingPresentation(vertexSrc, fragmentSrc string) *LightingPresentation {
	program := buildProgram(vertexSrc, fragmentSrc)
	return &LightingPresentation{
		Presentation: Presentation{
			Program: program,
		},

		FramebufferDraw0Location: program.UniformLocation("fbColor0TextureIn"),
		FramebufferDraw1Location: program.UniformLocation("fbColor1TextureIn"),
		FramebufferDepthLocation: program.UniformLocation("fbDepthTextureIn"),

		ProjectionMatrixLocation: program.UniformLocation("projectionMatrixIn"),
		CameraMatrixLocation:     program.UniformLocation("cameraMatrixIn"),
		ViewMatrixLocation:       program.UniformLocation("viewMatrixIn"),

		ReflectionTextureLocation: program.UniformLocation("reflectionTextureIn"),
		RefractionTextureLocation: program.UniformLocation("refractionTextureIn"),

		LightDirection: program.UniformLocation("lightDirectionIn"),
		LightIntensity: program.UniformLocation("lightIntensityIn"),
	}
}

func buildProgram(vertSrc, fragSrc string) *webgl.Program {
	vertexShader := webgl.NewShader()
	vertexShader.Allocate(webgl.ShaderAllocateInfo{
		ShaderType: wasmgl.VERTEX_SHADER,
		SourceCode: vertSrc,
	})
	defer func() {
		vertexShader.Release()
	}()

	fragmentShader := webgl.NewShader()
	fragmentShader.Allocate(webgl.ShaderAllocateInfo{
		ShaderType: wasmgl.FRAGMENT_SHADER,
		SourceCode: fragSrc,
	})
	defer func() {
		fragmentShader.Release()
	}()

	program := webgl.NewProgram()
	program.Allocate(webgl.ProgramAllocateInfo{
		VertexShader:   vertexShader,
		FragmentShader: fragmentShader,
	})
	return program
}
