package ui

import "github.com/mokiat/lacking/ui/renderapi/plugin"

func NewShaderCollection() plugin.ShaderCollection {
	return plugin.ShaderCollection{
		ShapeMaterial:      newShapeShaders(),
		ShapeBlankMaterial: newShapeBlankShaders(),
		ContourMaterial:    newContourShaders(),
		TextMaterial:       newTextShaders(),
	}
}
