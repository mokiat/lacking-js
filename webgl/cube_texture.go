package webgl

import (
	"github.com/mokiat/wasmgl"
)

func NewCubeTexture() *CubeTexture {
	return &CubeTexture{}
}

type CubeTexture struct {
	Texture
}

func (t *CubeTexture) Allocate(info CubeTextureAllocateInfo) {
	t.raw = wasmgl.CreateTexture()
	wasmgl.BindTexture(wasmgl.TEXTURE_CUBE_MAP, t.raw)

	wasmgl.TexParameteri(wasmgl.TEXTURE_CUBE_MAP, wasmgl.TEXTURE_WRAP_S, info.wrapS())
	wasmgl.TexParameteri(wasmgl.TEXTURE_CUBE_MAP, wasmgl.TEXTURE_WRAP_T, info.wrapT())
	wasmgl.TexParameteri(wasmgl.TEXTURE_CUBE_MAP, wasmgl.TEXTURE_MIN_FILTER, info.minFilter())
	wasmgl.TexParameteri(wasmgl.TEXTURE_CUBE_MAP, wasmgl.TEXTURE_MAG_FILTER, info.magFilter())

	// Note: Top and Bottom are flipped due to OpenGL's renderman issue
	wasmgl.TexStorage2D(wasmgl.TEXTURE_CUBE_MAP, 1, info.internalFormat(), info.Dimension, info.Dimension)

	wasmgl.TexSubImage2D(wasmgl.TEXTURE_CUBE_MAP_POSITIVE_X, 0, 0, 0, info.Dimension, info.Dimension, info.dataFormat(), info.dataComponentType(), info.RightSideData)
	wasmgl.TexSubImage2D(wasmgl.TEXTURE_CUBE_MAP_NEGATIVE_X, 0, 0, 0, info.Dimension, info.Dimension, info.dataFormat(), info.dataComponentType(), info.LeftSideData)
	wasmgl.TexSubImage2D(wasmgl.TEXTURE_CUBE_MAP_POSITIVE_Y, 0, 0, 0, info.Dimension, info.Dimension, info.dataFormat(), info.dataComponentType(), info.BottomSideData)
	wasmgl.TexSubImage2D(wasmgl.TEXTURE_CUBE_MAP_NEGATIVE_Y, 0, 0, 0, info.Dimension, info.Dimension, info.dataFormat(), info.dataComponentType(), info.TopSideData)
	wasmgl.TexSubImage2D(wasmgl.TEXTURE_CUBE_MAP_POSITIVE_Z, 0, 0, 0, info.Dimension, info.Dimension, info.dataFormat(), info.dataComponentType(), info.FrontSideData)
	wasmgl.TexSubImage2D(wasmgl.TEXTURE_CUBE_MAP_NEGATIVE_Z, 0, 0, 0, info.Dimension, info.Dimension, info.dataFormat(), info.dataComponentType(), info.BackSideData)
	// if info.GenerateMipmaps {
	// 	wasmgl.GenerateMipmap(wasmgl.TEXTURE_CUBE_MAP)
	// }
}

func (t *CubeTexture) Use() {
	wasmgl.BindTexture(wasmgl.TEXTURE_CUBE_MAP, t.raw)
}

func (t *CubeTexture) Release() {
	wasmgl.DeleteTexture(t.raw)
	t.raw = wasmgl.NilTexture
}

type CubeTextureAllocateInfo struct {
	Dimension         int
	WrapS             int
	WrapT             int
	MinFilter         int
	MagFilter         int
	GenerateMipmaps   bool
	InternalFormat    int
	DataFormat        int
	DataComponentType int
	FrontSideData     []byte
	BackSideData      []byte
	LeftSideData      []byte
	RightSideData     []byte
	TopSideData       []byte
	BottomSideData    []byte
}

func (i CubeTextureAllocateInfo) wrapS() int {
	if i.WrapS == 0 {
		return wasmgl.CLAMP_TO_EDGE
	}
	return i.WrapS
}

func (i CubeTextureAllocateInfo) wrapT() int {
	if i.WrapT == 0 {
		return wasmgl.CLAMP_TO_EDGE
	}
	return i.WrapT
}

func (i CubeTextureAllocateInfo) minFilter() int {
	if i.MinFilter == 0 {
		return wasmgl.LINEAR_MIPMAP_LINEAR
	}
	return i.MinFilter
}

func (i CubeTextureAllocateInfo) magFilter() int {
	if i.MagFilter == 0 {
		return wasmgl.LINEAR
	}
	return i.MagFilter
}

func (i CubeTextureAllocateInfo) internalFormat() int {
	if i.InternalFormat == 0 {
		return wasmgl.SRGB8_ALPHA8
	}
	return i.InternalFormat
}

func (i CubeTextureAllocateInfo) dataFormat() int {
	if i.DataFormat == 0 {
		return wasmgl.RGBA
	}
	return i.DataFormat
}

func (i CubeTextureAllocateInfo) dataComponentType() int {
	if i.DataComponentType == 0 {
		return wasmgl.UNSIGNED_BYTE
	}
	return i.DataComponentType
}
