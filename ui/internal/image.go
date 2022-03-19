//go:build js && wasm

package internal

import (
	"image"
	"image/draw"

	"github.com/mokiat/lacking-js/webgl"
	"github.com/mokiat/lacking/ui"
	"github.com/mokiat/wasmgl"
)

func NewImage() *Image {
	return &Image{
		texture: webgl.NewTwoDTexture(),
		size:    ui.NewSize(0, 0),
	}
}

type Image struct {
	texture *webgl.TwoDTexture
	size    ui.Size
}

func (i *Image) Allocate(img image.Image) {
	bounds := img.Bounds()
	var rgbaImg *image.NRGBA
	switch img := img.(type) {
	case *image.NRGBA:
		rgbaImg = img
	default:
		rgbaImg = image.NewNRGBA(bounds)
		draw.Draw(rgbaImg, bounds, img, bounds.Min, draw.Src)
	}
	i.size = ui.NewSize(bounds.Dx(), bounds.Dy())
	i.texture.Allocate(webgl.TwoDTextureAllocateInfo{
		Width:             bounds.Dx(),
		Height:            bounds.Dy(),
		WrapS:             wasmgl.CLAMP_TO_EDGE,
		WrapT:             wasmgl.CLAMP_TO_EDGE,
		MinFilter:         wasmgl.LINEAR,
		MagFilter:         wasmgl.LINEAR,
		InternalFormat:    wasmgl.SRGB8_ALPHA8,
		DataFormat:        wasmgl.RGBA,
		DataComponentType: wasmgl.UNSIGNED_BYTE,
		Data:              rgbaImg.Pix,
	})
}

func (i *Image) Release() {
	i.size = ui.NewSize(0, 0)
	i.texture.Release()
}

func (i *Image) Size() ui.Size {
	return i.size
}

func (i *Image) Destroy() {
	i.Release()
}
