package graphics

import (
	"github.com/mokiat/lacking-js/webgl"
	"github.com/mokiat/lacking/game/graphics"
)

func newTwoDTexture() *TwoDTexture {
	return &TwoDTexture{
		TwoDTexture: webgl.NewTwoDTexture(),
	}
}

var _ graphics.TwoDTexture = (*TwoDTexture)(nil)

type TwoDTexture struct {
	*webgl.TwoDTexture
}

func (t *TwoDTexture) Delete() {
	t.Release()
}

func newCubeTexture() *CubeTexture {
	return &CubeTexture{
		CubeTexture: webgl.NewCubeTexture(),
	}
}

var _ graphics.CubeTexture = (*CubeTexture)(nil)

type CubeTexture struct {
	*webgl.CubeTexture
}

func (t *CubeTexture) Delete() {
	t.Release()
}
