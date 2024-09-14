package internal

import (
	"github.com/mokiat/lacking/render"
	"github.com/mokiat/wasmgl"
)

func NewColorTexture2D(info render.ColorTexture2DInfo) *Texture {
	if glLogger.IsDebugEnabled() {
		defer trackError("Error creating color texture 2D (%v)", info.Label)()
	}

	raw := wasmgl.CreateTexture()
	wasmgl.BindTexture(wasmgl.TEXTURE_2D, raw)
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D, wasmgl.TEXTURE_WRAP_S, wasmgl.CLAMP_TO_EDGE)
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D, wasmgl.TEXTURE_WRAP_T, wasmgl.CLAMP_TO_EDGE)
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D, wasmgl.TEXTURE_MIN_FILTER, wasmgl.NEAREST)
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D, wasmgl.TEXTURE_MAG_FILTER, wasmgl.NEAREST)

	levels := glMipmapLevels(info.Width, info.Height, info.GenerateMipmaps)
	internalFormat := glInternalFormat(info.Format, info.GammaCorrection)
	wasmgl.TexStorage2D(wasmgl.TEXTURE_2D, levels, internalFormat, wasmgl.GLsizei(info.Width), wasmgl.GLsizei(info.Height))

	if info.Data != nil {
		dataFormat := glDataFormat(info.Format)
		componentType := glDataComponentType(info.Format)
		wasmgl.TexSubImage2D(wasmgl.TEXTURE_2D, 0, 0, 0, wasmgl.GLsizei(info.Width), wasmgl.GLsizei(info.Height), dataFormat, componentType, info.Data)
	}

	if info.GenerateMipmaps {
		// TODO: Move as separate command
		wasmgl.GenerateMipmap(wasmgl.TEXTURE_2D)
	}

	result := &Texture{
		label:  info.Label,
		raw:    raw,
		kind:   wasmgl.TEXTURE_2D,
		width:  info.Width,
		height: info.Height,
	}
	result.id = textures.Allocate(result)
	return result
}

func NewDepthTexture2D(info render.DepthTexture2DInfo) *Texture {
	if glLogger.IsDebugEnabled() {
		defer trackError("Error creating depth texture 2D (%v)", info.Label)()
	}

	raw := wasmgl.CreateTexture()
	wasmgl.BindTexture(wasmgl.TEXTURE_2D, raw)
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D, wasmgl.TEXTURE_WRAP_S, wasmgl.CLAMP_TO_EDGE)
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D, wasmgl.TEXTURE_WRAP_T, wasmgl.CLAMP_TO_EDGE)
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D, wasmgl.TEXTURE_MIN_FILTER, wasmgl.NEAREST)
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D, wasmgl.TEXTURE_MAG_FILTER, wasmgl.NEAREST)
	if info.Comparable {
		wasmgl.TexParameteri(wasmgl.TEXTURE_2D, wasmgl.TEXTURE_COMPARE_MODE, wasmgl.COMPARE_REF_TO_TEXTURE)
		wasmgl.TexStorage2D(wasmgl.TEXTURE_2D, 1, wasmgl.DEPTH_COMPONENT32F, wasmgl.GLsizei(info.Width), wasmgl.GLsizei(info.Height))
	} else {
		wasmgl.TexStorage2D(wasmgl.TEXTURE_2D, 1, wasmgl.DEPTH_COMPONENT24, wasmgl.GLsizei(info.Width), wasmgl.GLsizei(info.Height))
	}

	result := &Texture{
		label:  info.Label,
		raw:    raw,
		kind:   wasmgl.TEXTURE_2D,
		width:  info.Width,
		height: info.Height,
	}
	result.id = textures.Allocate(result)
	return result
}

func NewDepthTexture2DArray(info render.DepthTexture2DArrayInfo) *Texture {
	if glLogger.IsDebugEnabled() {
		defer trackError("Error creating array depth texture 2D (%v)", info.Label)()
	}

	raw := wasmgl.CreateTexture()
	wasmgl.BindTexture(wasmgl.TEXTURE_2D_ARRAY, raw)
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D_ARRAY, wasmgl.TEXTURE_WRAP_S, wasmgl.CLAMP_TO_EDGE)
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D_ARRAY, wasmgl.TEXTURE_WRAP_T, wasmgl.CLAMP_TO_EDGE)
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D_ARRAY, wasmgl.TEXTURE_MIN_FILTER, wasmgl.NEAREST)
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D_ARRAY, wasmgl.TEXTURE_MAG_FILTER, wasmgl.NEAREST)
	if info.Comparable {
		wasmgl.TexParameteri(wasmgl.TEXTURE_2D_ARRAY, wasmgl.TEXTURE_COMPARE_MODE, wasmgl.COMPARE_REF_TO_TEXTURE)
		wasmgl.TexStorage3D(wasmgl.TEXTURE_2D_ARRAY, 1, wasmgl.DEPTH_COMPONENT32F, wasmgl.GLsizei(info.Width), wasmgl.GLsizei(info.Height), wasmgl.GLsizei(info.Layers))
	} else {
		wasmgl.TexStorage3D(wasmgl.TEXTURE_2D_ARRAY, 1, wasmgl.DEPTH_COMPONENT24, wasmgl.GLsizei(info.Width), wasmgl.GLsizei(info.Height), wasmgl.GLsizei(info.Layers))
	}

	result := &Texture{
		label:  info.Label,
		raw:    raw,
		kind:   wasmgl.TEXTURE_2D_ARRAY,
		width:  info.Width,
		height: info.Height,
	}
	result.id = textures.Allocate(result)
	return result
}

func NewStencilTexture2D(info render.StencilTexture2DInfo) *Texture {
	if glLogger.IsDebugEnabled() {
		defer trackError("Error creating stencil texture 2D (%v)", info.Label)()
	}

	raw := wasmgl.CreateTexture()
	wasmgl.BindTexture(wasmgl.TEXTURE_2D, raw)
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D, wasmgl.TEXTURE_WRAP_S, wasmgl.CLAMP_TO_EDGE)
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D, wasmgl.TEXTURE_WRAP_T, wasmgl.CLAMP_TO_EDGE)
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D, wasmgl.TEXTURE_MIN_FILTER, wasmgl.NEAREST)
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D, wasmgl.TEXTURE_MAG_FILTER, wasmgl.NEAREST)
	// NOTE: Firefox does not support wasmgl.STENCIL_INDEX8
	wasmgl.TexStorage2D(wasmgl.TEXTURE_2D, 1, wasmgl.DEPTH24_STENCIL8, wasmgl.GLsizei(info.Width), wasmgl.GLsizei(info.Height))
	result := &Texture{
		label:  info.Label,
		raw:    raw,
		kind:   wasmgl.TEXTURE_2D,
		width:  info.Width,
		height: info.Height,
	}
	result.id = textures.Allocate(result)
	return result
}

func NewDepthStencilTexture2D(info render.DepthStencilTexture2DInfo) *Texture {
	if glLogger.IsDebugEnabled() {
		defer trackError("Error creating depth-stencil texture 2D (%v)", info.Label)()
	}

	raw := wasmgl.CreateTexture()
	wasmgl.BindTexture(wasmgl.TEXTURE_2D, raw)
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D, wasmgl.TEXTURE_WRAP_S, wasmgl.CLAMP_TO_EDGE)
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D, wasmgl.TEXTURE_WRAP_T, wasmgl.CLAMP_TO_EDGE)
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D, wasmgl.TEXTURE_MIN_FILTER, wasmgl.NEAREST)
	wasmgl.TexParameteri(wasmgl.TEXTURE_2D, wasmgl.TEXTURE_MAG_FILTER, wasmgl.NEAREST)
	wasmgl.TexStorage2D(wasmgl.TEXTURE_2D, 1, wasmgl.DEPTH24_STENCIL8, wasmgl.GLsizei(info.Width), wasmgl.GLsizei(info.Height))
	result := &Texture{
		label:  info.Label,
		raw:    raw,
		kind:   wasmgl.TEXTURE_2D,
		width:  info.Width,
		height: info.Height,
	}
	result.id = textures.Allocate(result)
	return result
}

func NewColorTextureCube(info render.ColorTextureCubeInfo) *Texture {
	if glLogger.IsDebugEnabled() {
		defer trackError("Error creating color texture cube (%v)", info.Label)()
	}

	raw := wasmgl.CreateTexture()
	wasmgl.BindTexture(wasmgl.TEXTURE_CUBE_MAP, raw)
	wasmgl.TexParameteri(wasmgl.TEXTURE_CUBE_MAP, wasmgl.TEXTURE_WRAP_S, wasmgl.CLAMP_TO_EDGE)
	wasmgl.TexParameteri(wasmgl.TEXTURE_CUBE_MAP, wasmgl.TEXTURE_WRAP_T, wasmgl.CLAMP_TO_EDGE)
	wasmgl.TexParameteri(wasmgl.TEXTURE_CUBE_MAP, wasmgl.TEXTURE_WRAP_R, wasmgl.CLAMP_TO_EDGE)
	wasmgl.TexParameteri(wasmgl.TEXTURE_CUBE_MAP, wasmgl.TEXTURE_MIN_FILTER, wasmgl.NEAREST)
	wasmgl.TexParameteri(wasmgl.TEXTURE_CUBE_MAP, wasmgl.TEXTURE_MAG_FILTER, wasmgl.NEAREST)

	levels := glMipmapLevels(info.Dimension, info.Dimension, info.GenerateMipmaps)
	internalFormat := glInternalFormat(info.Format, info.GammaCorrection)
	wasmgl.TexStorage2D(wasmgl.TEXTURE_CUBE_MAP, levels, internalFormat, wasmgl.GLsizei(info.Dimension), wasmgl.GLsizei(info.Dimension))

	dataFormat := glDataFormat(info.Format)
	componentType := glDataComponentType(info.Format)
	if info.RightSideData != nil {
		wasmgl.TexSubImage2D(wasmgl.TEXTURE_CUBE_MAP_POSITIVE_X, 0, 0, 0, wasmgl.GLsizei(info.Dimension), wasmgl.GLsizei(info.Dimension), dataFormat, componentType, info.RightSideData)
	}
	if info.LeftSideData != nil {
		wasmgl.TexSubImage2D(wasmgl.TEXTURE_CUBE_MAP_NEGATIVE_X, 0, 0, 0, wasmgl.GLsizei(info.Dimension), wasmgl.GLsizei(info.Dimension), dataFormat, componentType, info.LeftSideData)
	}
	if info.BottomSideData != nil {
		wasmgl.TexSubImage2D(wasmgl.TEXTURE_CUBE_MAP_POSITIVE_Y, 0, 0, 0, wasmgl.GLsizei(info.Dimension), wasmgl.GLsizei(info.Dimension), dataFormat, componentType, info.BottomSideData)
	}
	if info.TopSideData != nil {
		wasmgl.TexSubImage2D(wasmgl.TEXTURE_CUBE_MAP_NEGATIVE_Y, 0, 0, 0, wasmgl.GLsizei(info.Dimension), wasmgl.GLsizei(info.Dimension), dataFormat, componentType, info.TopSideData)
	}
	if info.FrontSideData != nil {
		wasmgl.TexSubImage2D(wasmgl.TEXTURE_CUBE_MAP_POSITIVE_Z, 0, 0, 0, wasmgl.GLsizei(info.Dimension), wasmgl.GLsizei(info.Dimension), dataFormat, componentType, info.FrontSideData)
	}
	if info.BackSideData != nil {
		wasmgl.TexSubImage2D(wasmgl.TEXTURE_CUBE_MAP_NEGATIVE_Z, 0, 0, 0, wasmgl.GLsizei(info.Dimension), wasmgl.GLsizei(info.Dimension), dataFormat, componentType, info.BackSideData)
	}

	if info.GenerateMipmaps {
		// TODO: Move as separate command
		wasmgl.GenerateMipmap(wasmgl.TEXTURE_CUBE_MAP)
	}

	result := &Texture{
		label:  info.Label,
		raw:    raw,
		kind:   wasmgl.TEXTURE_CUBE_MAP,
		width:  info.Dimension,
		height: info.Dimension,
		depth:  info.Dimension,
	}
	result.id = textures.Allocate(result)
	return result
}

type Texture struct {
	render.TextureMarker

	label  string
	id     uint32
	raw    wasmgl.Texture
	kind   wasmgl.GLenum
	width  uint32
	height uint32
	depth  uint32
}

func (t *Texture) Label() string {
	return t.label
}

func (t *Texture) Width() uint32 {
	return t.width
}

func (t *Texture) Height() uint32 {
	return t.height
}

func (t *Texture) Depth() uint32 {
	return t.depth
}

func (t *Texture) Release() {
	textures.Release(t.id)
	wasmgl.DeleteTexture(t.raw)
	t.raw = wasmgl.NilTexture
	t.id = 0
}

func NewSampler(info render.SamplerInfo) *Sampler {
	if glLogger.IsDebugEnabled() {
		defer trackError("Error creating sampler (%v)", info.Label)()
	}

	raw := wasmgl.CreateSampler()
	wasmgl.SamplerParameteri(raw, wasmgl.TEXTURE_WRAP_S, glWrap(info.Wrapping))
	wasmgl.SamplerParameteri(raw, wasmgl.TEXTURE_WRAP_T, glWrap(info.Wrapping))
	wasmgl.SamplerParameteri(raw, wasmgl.TEXTURE_WRAP_R, glWrap(info.Wrapping))
	wasmgl.SamplerParameteri(raw, wasmgl.TEXTURE_MIN_FILTER, glFilter(info.Filtering, info.Mipmapping))
	wasmgl.SamplerParameteri(raw, wasmgl.TEXTURE_MAG_FILTER, glFilter(info.Filtering, false)) // no mipmaps when magnification
	if info.Comparison.Specified {
		wasmgl.SamplerParameteri(raw, wasmgl.TEXTURE_COMPARE_MODE, wasmgl.COMPARE_REF_TO_TEXTURE)
		wasmgl.SamplerParameteri(raw, wasmgl.TEXTURE_COMPARE_FUNC, int32(glEnumFromComparison(info.Comparison.Value)))
	}

	result := &Sampler{
		label: info.Label,
		raw:   raw,
	}
	result.id = samplers.Allocate(result)
	return result
}

type Sampler struct {
	render.SamplerMarker

	label string
	id    uint32
	raw   wasmgl.Sampler
}

func (s *Sampler) Label() string {
	return s.label
}

func (s *Sampler) Release() {
	samplers.Release(s.id)
	wasmgl.DeleteSampler(s.raw)
	s.raw = wasmgl.NilSampler
	s.id = 0
}

func glWrap(wrap render.WrapMode) wasmgl.GLint {
	switch wrap {
	case render.WrapModeClamp:
		return wasmgl.CLAMP_TO_EDGE
	case render.WrapModeRepeat:
		return wasmgl.REPEAT
	case render.WrapModeMirroredRepeat:
		return wasmgl.MIRRORED_REPEAT
	default:
		return wasmgl.CLAMP_TO_EDGE
	}
}

func glFilter(filter render.FilterMode, mipmaps bool) wasmgl.GLint {
	switch filter {
	case render.FilterModeNearest:
		if mipmaps {
			return wasmgl.NEAREST_MIPMAP_NEAREST
		}
		return wasmgl.NEAREST
	case render.FilterModeLinear, render.FilterModeAnisotropic:
		if mipmaps {
			return wasmgl.LINEAR_MIPMAP_LINEAR
		}
		return wasmgl.LINEAR
	default:
		return wasmgl.NEAREST
	}
}

func glMipmapLevels(width, height uint32, mipmapping bool) wasmgl.GLsizei {
	if !mipmapping {
		return 1
	}
	count := wasmgl.GLsizei(1)
	for width > 1 || height > 1 {
		width /= 2
		height /= 2
		count++
	}
	return count
}

func glInternalFormat(format render.DataFormat, gammaCorrection bool) wasmgl.GLenum {
	switch format {
	case render.DataFormatRGBA8:
		if gammaCorrection {
			return wasmgl.SRGB8_ALPHA8
		}
		return wasmgl.RGBA8
	case render.DataFormatRGBA16F:
		return wasmgl.RGBA16F
	case render.DataFormatRGBA32F:
		return wasmgl.RGBA32F
	default:
		return wasmgl.RGBA8
	}
}

func glDataFormat(format render.DataFormat) wasmgl.GLenum {
	switch format {
	default:
		return wasmgl.RGBA
	}
}

func glDataComponentType(format render.DataFormat) wasmgl.GLenum {
	switch format {
	case render.DataFormatRGBA8:
		return wasmgl.UNSIGNED_BYTE
	case render.DataFormatRGBA16F:
		return wasmgl.HALF_FLOAT
	case render.DataFormatRGBA32F:
		return wasmgl.FLOAT
	default:
		return wasmgl.UNSIGNED_BYTE
	}
}
