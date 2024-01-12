package game

import (
	"fmt"

	"github.com/mokiat/lacking-js/render"
	"github.com/mokiat/lacking-js/shader"
	"github.com/mokiat/lacking/game/graphics"
	renderapi "github.com/mokiat/lacking/render"
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

func newShadowMappingSet(cfg graphics.ShadowMappingShaderConfig) renderapi.ProgramCode {
	var settings struct {
		UseArmature bool
	}
	if cfg.HasArmature {
		settings.UseArmature = true
	}
	return render.ProgramCode{
		VertexCode:   shader.RunTemplate("shadow.vert.glsl", settings),
		FragmentCode: shader.RunTemplate("shadow.frag.glsl", settings),
	}
}

func newPBRGeometrySet(cfg graphics.PBRGeometryShaderConfig) renderapi.ProgramCode {
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
	return render.ProgramCode{
		VertexCode:   shader.RunTemplate("pbr_geometry.vert.glsl", settings),
		FragmentCode: shader.RunTemplate("pbr_geometry.frag.glsl", settings),
	}
}

func newAmbientLightShaderSet() renderapi.ProgramCode {
	return render.ProgramCode{
		VertexCode:   shader.RunTemplate("ambient_light.vert.glsl", struct{}{}),
		FragmentCode: shader.RunTemplate("ambient_light.frag.glsl", struct{}{}),
	}
}

func newPointLightShaderSet() renderapi.ProgramCode {
	return render.ProgramCode{
		VertexCode:   shader.RunTemplate("point_light.vert.glsl", struct{}{}),
		FragmentCode: shader.RunTemplate("point_light.frag.glsl", struct{}{}),
	}
}

func newSpotLightShaderSet() renderapi.ProgramCode {
	return render.ProgramCode{
		VertexCode:   shader.RunTemplate("spot_light.vert.glsl", struct{}{}),
		FragmentCode: shader.RunTemplate("spot_light.frag.glsl", struct{}{}),
	}
}

func newDirectionalLightShaderSet() renderapi.ProgramCode {
	var settings struct {
		UseShadowMapping bool
	}
	settings.UseShadowMapping = true // TODO
	return render.ProgramCode{
		VertexCode:   shader.RunTemplate("directional_light.vert.glsl", settings),
		FragmentCode: shader.RunTemplate("directional_light.frag.glsl", settings),
	}
}

func newSkyboxShaderSet() renderapi.ProgramCode {
	return render.ProgramCode{
		VertexCode:   shader.RunTemplate("skybox.vert.glsl", struct{}{}),
		FragmentCode: shader.RunTemplate("skybox.frag.glsl", struct{}{}),
	}
}

func newSkycolorShaderSet() renderapi.ProgramCode {
	return render.ProgramCode{
		VertexCode:   shader.RunTemplate("skycolor.vert.glsl", struct{}{}),
		FragmentCode: shader.RunTemplate("skycolor.frag.glsl", struct{}{}),
	}
}

func newDebugShaderSet() renderapi.ProgramCode {
	return render.ProgramCode{
		VertexCode:   shader.RunTemplate("debug.vert.glsl", struct{}{}),
		FragmentCode: shader.RunTemplate("debug.frag.glsl", struct{}{}),
	}
}

func newExposureShaderSet() renderapi.ProgramCode {
	return render.ProgramCode{
		VertexCode:   shader.RunTemplate("exposure.vert.glsl", struct{}{}),
		FragmentCode: shader.RunTemplate("exposure.frag.glsl", struct{}{}),
	}
}

func newPostprocessingShaderSet(cfg graphics.PostprocessingShaderConfig) renderapi.ProgramCode {
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
	return render.ProgramCode{
		VertexCode:   shader.RunTemplate("postprocess.vert.glsl", settings),
		FragmentCode: shader.RunTemplate("postprocess.frag.glsl", settings),
	}
}
