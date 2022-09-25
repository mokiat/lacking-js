package game

import (
	_ "embed"
	"fmt"

	"github.com/mokiat/lacking-js/internal"
	"github.com/mokiat/lacking/game/graphics"
)

var (
	//go:embed shaders/pbr_geometry.vert
	pbrGeometryVertexShader string

	//go:embed shaders/pbr_geometry.frag
	pbrGeometryFragmentShader string

	//go:embed shaders/dir_light.vert
	directionalLightVertexShader string

	//go:embed shaders/dir_light.frag
	directionalLightFragmentShader string

	//go:embed shaders/amb_light.vert
	ambientLightVertexShader string

	//go:embed shaders/amb_light.frag
	ambientLightFragmentShader string

	//go:embed shaders/point_light.vert
	pointLightVertexShader string

	//go:embed shaders/point_light.frag
	pointLightFragmentShader string

	//go:embed shaders/shadow.vert
	shadowMappingVertexShader string

	//go:embed shaders/shadow.frag
	shadowMappingFragmentShader string

	//go:embed shaders/skybox.vert
	cubeSkyboxVertexShader string

	//go:embed shaders/skybox.frag
	cubeSkyboxFragmentShader string

	//go:embed shaders/skycolor.vert
	colorSkyboxVertexShader string

	//go:embed shaders/skycolor.frag
	colorSkyboxFragmentShader string

	//go:embed shaders/exposure.vert
	exposureVertexShader string

	//go:embed shaders/exposure.frag
	exposureFragmentShader string

	//go:embed shaders/postprocess.vert
	tonePostprocessingVertexShader string

	//go:embed shaders/postprocess.frag
	tonePostprocessingFragmentShader string
)

func NewShaderCollection() graphics.ShaderCollection {
	return graphics.ShaderCollection{
		ShadowMappingSet:    newShadowMappingSet,
		PBRGeometrySet:      newPBRGeometrySet,
		DirectionalLightSet: newDirectionalLightShaderSet,
		AmbientLightSet:     newAmbientLightShaderSet,
		PointLightSet:       newPointLightShaderSet,
		SkyboxSet:           newSkyboxShaderSet,
		SkycolorSet:         newSkycolorShaderSet,
		ExposureSet:         newExposureShaderSet,
		PostprocessingSet:   newPostprocessingShaderSet,
	}
}

func newShadowMappingSet(cfg graphics.ShadowMappingShaderConfig) graphics.ShaderSet {
	vsBuilder := internal.NewShaderSourceBuilder(shadowMappingVertexShader)
	fsBuilder := internal.NewShaderSourceBuilder(shadowMappingFragmentShader)
	if cfg.HasArmature {
		vsBuilder.AddFeature("USES_BONES")
		fsBuilder.AddFeature("USES_BONES")
	}
	return graphics.ShaderSet{
		VertexShader:   vsBuilder.Build(),
		FragmentShader: fsBuilder.Build(),
	}
}

func newPBRGeometrySet(cfg graphics.PBRGeometryShaderConfig) graphics.ShaderSet {
	vsBuilder := internal.NewShaderSourceBuilder(pbrGeometryVertexShader)
	fsBuilder := internal.NewShaderSourceBuilder(pbrGeometryFragmentShader)
	if cfg.HasArmature {
		vsBuilder.AddFeature("USES_BONES")
		fsBuilder.AddFeature("USES_BONES")
	}
	if cfg.HasAlphaTesting {
		vsBuilder.AddFeature("USES_ALPHA_TEST")
		fsBuilder.AddFeature("USES_ALPHA_TEST")
	}
	if cfg.HasAlbedoTexture {
		vsBuilder.AddFeature("USES_ALBEDO_TEXTURE")
		fsBuilder.AddFeature("USES_ALBEDO_TEXTURE")
		vsBuilder.AddFeature("USES_TEX_COORD0")
		fsBuilder.AddFeature("USES_TEX_COORD0")
	}
	return graphics.ShaderSet{
		VertexShader:   vsBuilder.Build(),
		FragmentShader: fsBuilder.Build(),
	}
}

func newDirectionalLightShaderSet() graphics.ShaderSet {
	vsBuilder := internal.NewShaderSourceBuilder(directionalLightVertexShader)
	fsBuilder := internal.NewShaderSourceBuilder(directionalLightFragmentShader)
	return graphics.ShaderSet{
		VertexShader:   vsBuilder.Build(),
		FragmentShader: fsBuilder.Build(),
	}
}

func newAmbientLightShaderSet() graphics.ShaderSet {
	vsBuilder := internal.NewShaderSourceBuilder(ambientLightVertexShader)
	fsBuilder := internal.NewShaderSourceBuilder(ambientLightFragmentShader)
	return graphics.ShaderSet{
		VertexShader:   vsBuilder.Build(),
		FragmentShader: fsBuilder.Build(),
	}
}

func newPointLightShaderSet() graphics.ShaderSet {
	vsBuilder := internal.NewShaderSourceBuilder(pointLightVertexShader)
	fsBuilder := internal.NewShaderSourceBuilder(pointLightFragmentShader)
	return graphics.ShaderSet{
		VertexShader:   vsBuilder.Build(),
		FragmentShader: fsBuilder.Build(),
	}
}

func newSkyboxShaderSet() graphics.ShaderSet {
	vsBuilder := internal.NewShaderSourceBuilder(cubeSkyboxVertexShader)
	fsBuilder := internal.NewShaderSourceBuilder(cubeSkyboxFragmentShader)
	return graphics.ShaderSet{
		VertexShader:   vsBuilder.Build(),
		FragmentShader: fsBuilder.Build(),
	}
}

func newSkycolorShaderSet() graphics.ShaderSet {
	vsBuilder := internal.NewShaderSourceBuilder(colorSkyboxVertexShader)
	fsBuilder := internal.NewShaderSourceBuilder(colorSkyboxFragmentShader)
	return graphics.ShaderSet{
		VertexShader:   vsBuilder.Build(),
		FragmentShader: fsBuilder.Build(),
	}
}

func newExposureShaderSet() graphics.ShaderSet {
	vsBuilder := internal.NewShaderSourceBuilder(exposureVertexShader)
	fsBuilder := internal.NewShaderSourceBuilder(exposureFragmentShader)
	return graphics.ShaderSet{
		VertexShader:   vsBuilder.Build(),
		FragmentShader: fsBuilder.Build(),
	}
}

func newPostprocessingShaderSet(cfg graphics.PostprocessingShaderConfig) graphics.ShaderSet {
	vsBuilder := internal.NewShaderSourceBuilder(tonePostprocessingVertexShader)
	fsBuilder := internal.NewShaderSourceBuilder(tonePostprocessingFragmentShader)
	switch cfg.ToneMapping {
	case graphics.ReinhardToneMapping:
		fsBuilder.AddFeature("MODE_REINHARD")
	case graphics.ExponentialToneMapping:
		fsBuilder.AddFeature("MODE_EXPONENTIAL")
	default:
		panic(fmt.Errorf("unknown tone mapping mode: %s", cfg.ToneMapping))
	}
	return graphics.ShaderSet{
		VertexShader:   vsBuilder.Build(),
		FragmentShader: fsBuilder.Build(),
	}
}
