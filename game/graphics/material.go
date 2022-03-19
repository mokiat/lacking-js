package graphics

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking-js/game/graphics/internal"
	"github.com/mokiat/lacking-js/webgl"
)

type Material struct {
	backfaceCulling bool
	alphaTesting    bool
	alphaBlending   bool
	alphaThreshold  float32

	geometryPresentation *internal.GeometryPresentation
	shadowPresentation   *internal.ShadowPresentation

	twoDTextures []*webgl.TwoDTexture
	cubeTextures []*webgl.CubeTexture
	vectors      []sprec.Vec4
}

func (m *Material) Delete() {
	if m.geometryPresentation != nil {
		m.geometryPresentation.Delete()
	}
	if m.shadowPresentation != nil {
		m.shadowPresentation.Delete()
	}
}
