package game

import "github.com/mokiat/lacking/game/graphics"

func NewShaderCollection() graphics.ShaderCollection {
	return graphics.ShaderCollection{
		ExposureSet:         newExposureShaderSet,
		PostprocessingSet:   newPostprocessingShaderSet,
		DirectionalLightSet: newDirectionalLightShaderSet,
		AmbientLightSet:     newAmbientLightShaderSet,
		SkyboxSet:           newSkyboxShaderSet,
		SkycolorSet:         newSkycolorShaderSet,
		PBRShaderSet:        newPBRShaderSet,
	}
}
