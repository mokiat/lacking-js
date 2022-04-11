package game

import "github.com/mokiat/lacking/game/graphics/renderapi/plugin"

func NewShaderCollection() plugin.ShaderCollection {
	return plugin.ShaderCollection{
		ExposureSet:         newExposureShaderSet,
		PostprocessingSet:   newPostprocessingShaderSet,
		DirectionalLightSet: newDirectionalLightShaderSet,
		AmbientLightSet:     newAmbientLightShaderSet,
		SkyboxSet:           newSkyboxShaderSet,
		SkycolorSet:         newSkycolorShaderSet,
		PBRShaderSet:        newPBRShaderSet,
	}
}
