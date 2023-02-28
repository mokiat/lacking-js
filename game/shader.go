package game

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"

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
	tmplShadowMappingVertexShader   = find("shadow.vert.glsl")
	tmplShadowMappingFragmentShader = find("shadow.frag.glsl")

	tmplPBRGeometryVertexShader   = find("pbr_geometry.vert.glsl")
	tmplPBRGeometryFragmentShader = find("pbr_geometry.frag.glsl")

	tmplAmbientLightVertexShader   = find("ambient_light.vert.glsl")
	tmplAmbientLightFragmentShader = find("ambient_light.frag.glsl")

	tmplPointLightVertexShader   = find("point_light.vert.glsl")
	tmplPointLightFragmentShader = find("point_light.frag.glsl")

	tmplSpotLightVertexShader   = find("spot_light.vert.glsl")
	tmplSpotLightFragmentShader = find("spot_light.frag.glsl")

	tmplDirectionalLightVertexShader   = find("directional_light.vert.glsl")
	tmplDirectionalLightFragmentShader = find("directional_light.frag.glsl")

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
	var settings struct {
		UseArmature bool
	}
	if cfg.HasArmature {
		settings.UseArmature = true
	}
	return graphics.ShaderSet{
		VertexShader:   runTemplate(tmplShadowMappingVertexShader, settings),
		FragmentShader: runTemplate(tmplShadowMappingFragmentShader, settings),
	}
}

func newPBRGeometrySet(cfg graphics.PBRGeometryShaderConfig) graphics.ShaderSet {
	var settings struct {
		UseArmature       bool
		UseAlphaTest      bool
		UseVertexColoring bool
		UseTexturing      bool
		UseAlbedoTexture  bool
	}
	if cfg.HasArmature {
		settings.UseArmature = true
	}
	if cfg.HasAlphaTesting {
		settings.UseAlphaTest = true
	}
	if cfg.HasVertexColors {
		settings.UseVertexColoring = true
	}
	if cfg.HasAlbedoTexture {
		settings.UseTexturing = true
		settings.UseAlbedoTexture = true
	}
	return graphics.ShaderSet{
		VertexShader:   runTemplate(tmplPBRGeometryVertexShader, settings),
		FragmentShader: runTemplate(tmplPBRGeometryFragmentShader, settings),
	}
}

func newDirectionalLightShaderSet() graphics.ShaderSet {
	var settings struct {
		UseShadowMapping bool
	}
	settings.UseShadowMapping = true // TODO
	return graphics.ShaderSet{
		VertexShader:   runTemplate(tmplDirectionalLightVertexShader, settings),
		FragmentShader: runTemplate(tmplDirectionalLightFragmentShader, settings),
	}
}

func newAmbientLightShaderSet() graphics.ShaderSet {
	return graphics.ShaderSet{
		VertexShader:   runTemplate(tmplAmbientLightVertexShader, struct{}{}),
		FragmentShader: runTemplate(tmplAmbientLightFragmentShader, struct{}{}),
	}
}

func newPointLightShaderSet() graphics.ShaderSet {
	return graphics.ShaderSet{
		VertexShader:   runTemplate(tmplPointLightVertexShader, struct{}{}),
		FragmentShader: runTemplate(tmplPointLightFragmentShader, struct{}{}),
	}
}

func newSpotLightShaderSet() graphics.ShaderSet {
	return graphics.ShaderSet{
		VertexShader:   runTemplate(tmplSpotLightVertexShader, struct{}{}),
		FragmentShader: runTemplate(tmplSpotLightFragmentShader, struct{}{}),
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
