//go:build js && wasm

package webgl

import "github.com/mokiat/wasmgl"

func NewTwoDTexture() *TwoDTexture {
	return &TwoDTexture{}
}

type TwoDTexture struct {
	Texture
}

func (t *TwoDTexture) Allocate(info TwoDTextureAllocateInfo) {
	t.raw = wasmgl.CreateTexture()
	wasmgl.BindTexture(wasmgl.TEXTURE_2D, t.raw)
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D, wasmgl.TEXTURE_WRAP_S, info.wrapS())
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D, wasmgl.TEXTURE_WRAP_T, info.wrapT())
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D, wasmgl.TEXTURE_MIN_FILTER, info.minFilter())
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D, wasmgl.TEXTURE_MAG_FILTER, info.magFilter())
	wasmgl.TexStorage2D(wasmgl.TEXTURE_2D, info.levels(), info.internalFormat(), info.Width, info.Height)
	if info.Data != nil {
		wasmgl.TexSubImage2D(wasmgl.TEXTURE_2D, 0, 0, 0, info.Width, info.Height, info.dataFormat(), info.dataComponentType(), info.Data)
	}
	if info.GenerateMipmaps {
		wasmgl.GenerateMipmap(wasmgl.TEXTURE_2D)
	}
}

func (t *TwoDTexture) Use() {
	wasmgl.BindTexture(wasmgl.TEXTURE_2D, t.raw)
}

func (t *TwoDTexture) Release() {
	wasmgl.DeleteTexture(t.raw)
	t.raw = wasmgl.Texture{}
}

type TwoDTextureAllocateInfo struct {
	Width              int
	Height             int
	WrapS              int
	WrapT              int
	MinFilter          int
	MagFilter          int
	PlaceholderMipmaps bool
	GenerateMipmaps    bool
	InternalFormat     int
	DataFormat         int
	DataComponentType  int
	Data               []byte
}

func (i TwoDTextureAllocateInfo) wrapS() int {
	if i.WrapS == 0 {
		return wasmgl.REPEAT
	}
	return i.WrapS
}

func (i TwoDTextureAllocateInfo) wrapT() int {
	if i.WrapT == 0 {
		return wasmgl.REPEAT
	}
	return i.WrapT
}

func (i TwoDTextureAllocateInfo) minFilter() int {
	if i.MinFilter == 0 {
		return wasmgl.LINEAR_MIPMAP_LINEAR
	}
	return i.MinFilter
}

func (i TwoDTextureAllocateInfo) magFilter() int {
	if i.MagFilter == 0 {
		return wasmgl.LINEAR
	}
	return i.MagFilter
}

func (i TwoDTextureAllocateInfo) internalFormat() int {
	if i.InternalFormat == 0 {
		return wasmgl.SRGB8_ALPHA8
	}
	return i.InternalFormat
}

func (i TwoDTextureAllocateInfo) dataFormat() int {
	if i.DataFormat == 0 {
		return wasmgl.RGBA
	}
	return i.DataFormat
}

func (i TwoDTextureAllocateInfo) dataComponentType() int {
	if i.DataComponentType == 0 {
		return wasmgl.UNSIGNED_BYTE
	}
	return i.DataComponentType
}

func (i TwoDTextureAllocateInfo) levels() int {
	if !i.GenerateMipmaps && !i.PlaceholderMipmaps {
		return 1
	}
	count := int(1)
	width, height := i.Width, i.Height
	for width > 1 || height > 1 {
		width /= 2
		height /= 2
		count++
	}
	return count
}
