package game

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"

	"github.com/mokiat/lacking-js/internal"
	"github.com/mokiat/lacking/game/graphics"
)

//go:embed shaders/*
var sources embed.FS

var rootTemplate = template.Must(template.
	New("root").
	Delims("/*", "*/").
	ParseFS(sources, "shaders/*.glsl"),
)

func find(name string) *template.Template {
	result := rootTemplate.Lookup(name)
	if result == nil {
		panic(fmt.Errorf("template %q not found", name))
	}
	return result
}

var buffer = new(bytes.Buffer)

func runTemplate(tmpl *template.Template, data any) string {
	buffer.Reset()
	if err := tmpl.Execute(buffer, data); err != nil {
		panic(fmt.Errorf("template exec error: %w", err))
	}
	return buffer.String()
}

var (
	tmplSkyboxVertexShader   = find("skybox.vert.glsl")
	tmplSkyboxFragmentShader = find("skybox.frag.glsl")

	tmplSkycolorVertexShader   = find("skycolor.vert.glsl")
	tmplSkycolorFragmentShader = find("skycolor.frag.glsl")

	tmplDebugVertexShader   = find("debug.vert.glsl")
	tmplDebugFragmentShader = find("debug.frag.glsl")

	tmplExposureVertexShader   = find("exposure.vert.glsl")
	tmplExposureFragmentShader = find("exposure.frag.glsl")

	tmplPostprocessingVertexShader   = find("postprocess.vert.glsl")
	tmplPostprocessingFragmentShader = find("postprocess.frag.glsl")
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
)

func NewShaderCollection() graphics.ShaderCollection {
	return graphics.ShaderCollection{
		ShadowMappingSet:    newShadowMappingSet,
		PBRGeometrySet:      newPBRGeometrySet,
		DirectionalLightSet: newDirectionalLightShaderSet,
		AmbientLightSet:     newAmbientLightShaderSet,
		PointLightSet:       newPointLightShaderSet,
		SpotLightSet:        newSpotLightShaderSet,
		SkyboxSet:           newSkyboxShaderSet,
		SkycolorSet:         newSkycolorShaderSet,
		DebugSet:            newDebugShaderSet,
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
	if cfg.HasVertexColors {
		vsBuilder.AddFeature("USES_COLOR0")
		fsBuilder.AddFeature("USES_COLOR0")
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

func newSpotLightShaderSet() graphics.ShaderSet {
	vsBuilder := internal.NewShaderSourceBuilder(pointLightVertexShader)
	fsBuilder := internal.NewShaderSourceBuilder(pointLightFragmentShader)
	return graphics.ShaderSet{
		VertexShader:   vsBuilder.Build(),
		FragmentShader: fsBuilder.Build(),
	}
}

func newSkyboxShaderSet() graphics.ShaderSet {
	return graphics.ShaderSet{
		VertexShader:   runTemplate(tmplSkyboxVertexShader, struct{}{}),
		FragmentShader: runTemplate(tmplSkyboxFragmentShader, struct{}{}),
	}
}

func newSkycolorShaderSet() graphics.ShaderSet {
	return graphics.ShaderSet{
		VertexShader:   runTemplate(tmplSkycolorVertexShader, struct{}{}),
		FragmentShader: runTemplate(tmplSkycolorFragmentShader, struct{}{}),
	}
}

func newDebugShaderSet() graphics.ShaderSet {
	return graphics.ShaderSet{
		VertexShader:   runTemplate(tmplDebugVertexShader, struct{}{}),
		FragmentShader: runTemplate(tmplDebugFragmentShader, struct{}{}),
	}
}

func newExposureShaderSet() graphics.ShaderSet {
	return graphics.ShaderSet{
		VertexShader:   runTemplate(tmplExposureVertexShader, struct{}{}),
		FragmentShader: runTemplate(tmplExposureFragmentShader, struct{}{}),
	}
}

func newPostprocessingShaderSet(cfg graphics.PostprocessingShaderConfig) graphics.ShaderSet {
	var settings struct {
		UseReinhard    bool
		UseExponential bool
	}
	switch cfg.ToneMapping {
	case graphics.ReinhardToneMapping:
		settings.UseReinhard = true
	case graphics.ExponentialToneMapping:
		settings.UseExponential = true
	default:
		panic(fmt.Errorf("unknown tone mapping mode: %s", cfg.ToneMapping))
	}
	return graphics.ShaderSet{
		VertexShader:   runTemplate(tmplPostprocessingVertexShader, settings),
		FragmentShader: runTemplate(tmplPostprocessingFragmentShader, settings),
	}
}
