package ui

import "github.com/mokiat/lacking/ui"

func NewShaderCollection() ui.ShaderCollection {
	return ui.ShaderCollection{
		ShapeMaterial:      newShapeShaders(),
		ShapeBlankMaterial: newShapeBlankShaders(),
		ContourMaterial:    newContourShaders(),
		TextMaterial:       newTextShaders(),
	}
}
