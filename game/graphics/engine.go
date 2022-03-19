//go:build js && wasm

package graphics

import (
	"fmt"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking-js/game/graphics/internal"
	"github.com/mokiat/lacking-js/webgl"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/wasmgl"
)

func NewEngine() *Engine {
	return &Engine{
		renderer: newRenderer(),
	}
}

var _ graphics.Engine = (*Engine)(nil)

type Engine struct {
	renderer *Renderer
}

func (e *Engine) Create() {
	e.renderer.Allocate()
}

func (e *Engine) CreateScene() graphics.Scene {
	return newScene(e.renderer)
}

func (e *Engine) CreateTwoDTexture(definition graphics.TwoDTextureDefinition) graphics.TwoDTexture {
	result := newTwoDTexture()
	result.TwoDTexture.Allocate(webgl.TwoDTextureAllocateInfo{
		Width:             definition.Width,
		Height:            definition.Height,
		WrapS:             e.convertWrap(definition.WrapS),
		WrapT:             e.convertWrap(definition.WrapT),
		MinFilter:         e.convertMinFilter(definition.MinFilter),
		MagFilter:         e.convertMagFilter(definition.MagFilter),
		GenerateMipmaps:   e.needsMipmaps(definition.MinFilter),
		DataFormat:        e.convertDataFormat(definition.DataFormat),
		DataComponentType: e.convertDataComponentType(definition.DataFormat),
		InternalFormat:    e.convertInternalFormat(definition.InternalFormat),
		Data:              definition.Data,
	})
	return result
}

func (e *Engine) CreateCubeTexture(definition graphics.CubeTextureDefinition) graphics.CubeTexture {
	result := newCubeTexture()
	result.CubeTexture.Allocate(webgl.CubeTextureAllocateInfo{
		Dimension:         definition.Dimension,
		WrapS:             wasmgl.CLAMP_TO_EDGE, // pointless
		WrapT:             wasmgl.CLAMP_TO_EDGE, // pointless
		MinFilter:         e.convertMinFilter(definition.MinFilter),
		MagFilter:         e.convertMagFilter(definition.MagFilter),
		GenerateMipmaps:   e.needsMipmaps(definition.MinFilter),
		DataFormat:        e.convertDataFormat(definition.DataFormat),
		DataComponentType: e.convertDataComponentType(definition.DataFormat),
		InternalFormat:    e.convertInternalFormat(definition.InternalFormat),
		FrontSideData:     definition.FrontSideData,
		BackSideData:      definition.BackSideData,
		LeftSideData:      definition.LeftSideData,
		RightSideData:     definition.RightSideData,
		TopSideData:       definition.TopSideData,
		BottomSideData:    definition.BottomSideData,
	})
	return result
}

func (e *Engine) CreateMeshTemplate(definition graphics.MeshTemplateDefinition) graphics.MeshTemplate {
	vertexBuffer := webgl.NewBuffer()
	vertexBuffer.Allocate(webgl.BufferAllocateInfo{
		BufferType: wasmgl.ARRAY_BUFFER,
		Dynamic:    false,
		Data:       definition.VertexData,
	})

	indexBuffer := webgl.NewBuffer()
	indexBuffer.Allocate(webgl.BufferAllocateInfo{
		BufferType: wasmgl.ELEMENT_ARRAY_BUFFER,
		Dynamic:    false,
		Data:       definition.IndexData,
	})

	var attributes []webgl.VertexArrayAttribute
	if definition.VertexFormat.HasCoord {
		attributes = append(attributes, webgl.VertexArrayAttribute{
			Buffer:         vertexBuffer,
			Index:          coordAttributeIndex,
			ComponentCount: 3,
			ComponentType:  wasmgl.FLOAT,
			Normalized:     false,
			StrideBytes:    definition.VertexFormat.CoordStrideBytes, // FIXME: Not accurate
			OffsetBytes:    definition.VertexFormat.CoordOffsetBytes,
		})
	}
	if definition.VertexFormat.HasNormal {
		attributes = append(attributes, webgl.VertexArrayAttribute{
			Buffer:         vertexBuffer,
			Index:          normalAttributeIndex,
			ComponentCount: 3,
			ComponentType:  wasmgl.FLOAT,
			Normalized:     false,
			StrideBytes:    definition.VertexFormat.CoordStrideBytes, // FIXME: Not accurate
			OffsetBytes:    definition.VertexFormat.NormalOffsetBytes,
		})
	}
	if definition.VertexFormat.HasTangent {
		attributes = append(attributes, webgl.VertexArrayAttribute{
			Buffer:         vertexBuffer,
			Index:          tangentAttributeIndex,
			ComponentCount: 3,
			ComponentType:  wasmgl.FLOAT,
			Normalized:     false,
			StrideBytes:    definition.VertexFormat.CoordStrideBytes, // FIXME: Not accurate
			OffsetBytes:    definition.VertexFormat.TangentOffsetBytes,
		})
	}
	if definition.VertexFormat.HasTexCoord {
		attributes = append(attributes, webgl.VertexArrayAttribute{
			Buffer:         vertexBuffer,
			Index:          texCoordAttributeIndex,
			ComponentCount: 2,
			ComponentType:  wasmgl.FLOAT,
			Normalized:     false,
			StrideBytes:    definition.VertexFormat.CoordStrideBytes, // FIXME: Not accurate
			OffsetBytes:    definition.VertexFormat.TexCoordOffsetBytes,
		})
	}
	if definition.VertexFormat.HasColor {
		attributes = append(attributes, webgl.VertexArrayAttribute{
			Buffer:         vertexBuffer,
			Index:          colorAttributeIndex,
			ComponentCount: 4,
			ComponentType:  wasmgl.FLOAT,
			Normalized:     false,
			StrideBytes:    definition.VertexFormat.CoordStrideBytes, // FIXME: Not accurate
			OffsetBytes:    definition.VertexFormat.ColorOffsetBytes,
		})
	}

	vertexArray := webgl.NewVertexArray()
	vertexArray.Allocate(webgl.VertexArrayAllocateInfo{
		Attributes:  attributes,
		IndexBuffer: indexBuffer,
	})

	result := &MeshTemplate{
		vertexBuffer: vertexBuffer,
		indexBuffer:  indexBuffer,
		vertexArray:  vertexArray,
		subMeshes:    make([]SubMeshTemplate, len(definition.SubMeshes)),
	}
	for i, subMesh := range definition.SubMeshes {
		result.subMeshes[i] = SubMeshTemplate{
			material:         subMesh.Material.(*Material),
			primitive:        e.convertPrimitive(subMesh.Primitive),
			indexCount:       subMesh.IndexCount,
			indexOffsetBytes: subMesh.IndexOffset,
			indexType:        e.convertIndexType(definition.IndexFormat),
		}
	}
	return result
}

func (e *Engine) CreatePBRMaterial(definition graphics.PBRMaterialDefinition) graphics.Material {
	extractTwoDTexture := func(src graphics.TwoDTexture) *webgl.TwoDTexture {
		if src == nil {
			return nil
		}
		return src.(*TwoDTexture).TwoDTexture
	}
	return &Material{
		backfaceCulling: definition.BackfaceCulling,
		alphaBlending:   definition.AlphaBlending,
		alphaTesting:    definition.AlphaTesting,
		alphaThreshold:  definition.AlphaThreshold,
		twoDTextures: []*webgl.TwoDTexture{
			extractTwoDTexture(definition.AlbedoTexture),
			extractTwoDTexture(definition.NormalTexture),
			extractTwoDTexture(definition.MetalnessTexture),
			extractTwoDTexture(definition.RoughnessTexture),
		},
		cubeTextures: []*webgl.CubeTexture{},
		vectors: []sprec.Vec4{
			definition.AlbedoColor,
			sprec.NewVec4(definition.NormalScale, definition.Metalness, definition.Roughness, 0.0),
		},
		geometryPresentation: internal.NewPBRGeometryPresentation(definition),
		shadowPresentation:   nil, // TODO
	}
}

func (e *Engine) Destroy() {
	e.renderer.Release()
}

func (e *Engine) convertWrap(wrap graphics.Wrap) int {
	switch wrap {
	case graphics.WrapClampToEdge:
		return wasmgl.CLAMP_TO_EDGE
	case graphics.WrapMirroredClampToEdge:
		// WebGL does not support mirrored clamp so fallback to default clamp
		return wasmgl.CLAMP_TO_EDGE
	case graphics.WrapRepeat:
		return wasmgl.REPEAT
	case graphics.WrapMirroredRepat:
		return wasmgl.MIRRORED_REPEAT
	default:
		panic(fmt.Errorf("unknown wrap mode: %d", wrap))
	}
}

func (e *Engine) needsMipmaps(filter graphics.Filter) bool {
	switch filter {
	case graphics.FilterNearestMipmapNearest:
		fallthrough
	case graphics.FilterNearestMipmapLinear:
		fallthrough
	case graphics.FilterLinearMipmapNearest:
		fallthrough
	case graphics.FilterLinearMipmapLinear:
		return true
	default:
		return false
	}
}

func (e *Engine) convertMinFilter(filter graphics.Filter) int {
	switch filter {
	case graphics.FilterNearest:
		return wasmgl.NEAREST
	case graphics.FilterLinear:
		return wasmgl.LINEAR
	case graphics.FilterNearestMipmapNearest:
		return wasmgl.NEAREST_MIPMAP_NEAREST
	case graphics.FilterNearestMipmapLinear:
		return wasmgl.NEAREST_MIPMAP_LINEAR
	case graphics.FilterLinearMipmapNearest:
		return wasmgl.LINEAR_MIPMAP_NEAREST
	case graphics.FilterLinearMipmapLinear:
		return wasmgl.LINEAR_MIPMAP_LINEAR
	default:
		panic(fmt.Errorf("unknown min filter mode: %d", filter))
	}
}

func (e *Engine) convertMagFilter(filter graphics.Filter) int {
	switch filter {
	case graphics.FilterNearest:
		return wasmgl.NEAREST
	case graphics.FilterLinear:
		return wasmgl.LINEAR
	default:
		panic(fmt.Errorf("unknown mag filter mode: %d", filter))
	}
}

func (e *Engine) convertDataFormat(format graphics.DataFormat) int {
	switch format {
	case graphics.DataFormatRGBA8:
		return wasmgl.RGBA
	case graphics.DataFormatRGBA32F:
		return wasmgl.RGBA
	default:
		panic(fmt.Errorf("unknown data format: %d", format))
	}
}

func (e *Engine) convertDataComponentType(format graphics.DataFormat) int {
	switch format {
	case graphics.DataFormatRGBA8:
		return wasmgl.UNSIGNED_BYTE
	case graphics.DataFormatRGBA32F:
		return wasmgl.FLOAT
	default:
		panic(fmt.Errorf("unknown data format: %d", format))
	}
}

func (e *Engine) convertInternalFormat(format graphics.InternalFormat) int {
	switch format {
	case graphics.InternalFormatRGBA8:
		return wasmgl.SRGB8_ALPHA8
	case graphics.InternalFormatRGBA32F:
		return wasmgl.RGBA32F
	default:
		panic(fmt.Errorf("unknown internal format: %d", format))
	}
}

func (e *Engine) convertPrimitive(primitive graphics.Primitive) int {
	switch primitive {
	case graphics.PrimitivePoints:
		return wasmgl.POINTS
	case graphics.PrimitiveLines:
		return wasmgl.LINES
	case graphics.PrimitiveLineStrip:
		return wasmgl.LINE_STRIP
	case graphics.PrimitiveLineLoop:
		return wasmgl.LINE_LOOP
	case graphics.PrimitiveTriangles:
		return wasmgl.TRIANGLES
	case graphics.PrimitiveTriangleStrip:
		return wasmgl.TRIANGLE_STRIP
	case graphics.PrimitiveTriangleFan:
		return wasmgl.TRIANGLE_FAN
	default:
		panic(fmt.Errorf("unknown primitive: %d", primitive))
	}
}

func (e *Engine) convertIndexType(indexFormat graphics.IndexFormat) int {
	switch indexFormat {
	case graphics.IndexFormatU16:
		return wasmgl.UNSIGNED_SHORT
	default:
		panic(fmt.Errorf("unknown index format: %d", indexFormat))
	}
}
