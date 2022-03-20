//go:build js && wasm

package internal

import (
	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking-js/webgl"
	"github.com/mokiat/wasmgl"
)

var (
	gl string // TODO: REMOVE
)

const maxVertexCount = 524288

func NewRenderer() *Renderer {
	return &Renderer{
		shape:              newShape(),
		shapeMesh:          newShapeMesh(maxVertexCount),
		shapeMaterial:      newShapeMaterial(),
		shapeBlankMaterial: newShapeBlankMaterial(),
		contour:            newContour(),
		contourMesh:        newContourMesh(maxVertexCount),
		contourMaterial:    newContourMaterial(),
		text:               newText(),
		textMesh:           newTextMesh(maxVertexCount),
		textMaterial:       newTextMaterial(),
		whiteMask:          webgl.NewTwoDTexture(),
	}
}

type Renderer struct {
	transformMatrix        sprec.Mat4
	textureTransformMatrix sprec.Mat4
	clipBounds             sprec.Vec4

	shape              *Shape
	shapeMesh          *ShapeMesh
	shapeMaterial      *Material
	shapeBlankMaterial *Material

	contour         *Contour
	contourMesh     *ContourMesh
	contourMaterial *Material

	text         *Text
	textMesh     *TextMesh
	textMaterial *Material

	whiteMask *webgl.TwoDTexture

	subMeshes []SubMesh

	target Target
}

func (r *Renderer) Init() {
	r.shapeMesh.Allocate()
	r.shapeMaterial.Allocate()
	r.shapeBlankMaterial.Allocate()
	r.contourMesh.Allocate()
	r.contourMaterial.Allocate()
	r.textMesh.Allocate()
	r.textMaterial.Allocate()
	r.whiteMask.Allocate(webgl.TwoDTextureAllocateInfo{
		Width:             1,
		Height:            1,
		WrapS:             wasmgl.CLAMP_TO_EDGE,
		WrapT:             wasmgl.CLAMP_TO_EDGE,
		MinFilter:         wasmgl.NEAREST,
		MagFilter:         wasmgl.NEAREST,
		GenerateMipmaps:   false,
		InternalFormat:    wasmgl.SRGB8_ALPHA8,
		DataFormat:        wasmgl.RGBA,
		DataComponentType: wasmgl.UNSIGNED_BYTE,
		Data:              []byte{0xFF, 0xFF, 0xFF, 0xFF},
	})
}

func (r *Renderer) Free() {
	defer r.shapeMesh.Release()
	defer r.shapeMaterial.Release()
	defer r.shapeBlankMaterial.Release()
	defer r.contourMesh.Release()
	defer r.contourMaterial.Release()
	defer r.textMesh.Release()
	defer r.textMaterial.Release()
	defer r.whiteMask.Release()
}

func (r *Renderer) Transform() sprec.Mat4 {
	return r.transformMatrix
}

func (r *Renderer) SetTransform(transform sprec.Mat4) {
	r.transformMatrix = transform
}

func (r *Renderer) TextureTransform() sprec.Mat4 {
	return r.textureTransformMatrix
}

func (r *Renderer) SetTextureTransform(textureTransform sprec.Mat4) {
	r.textureTransformMatrix = textureTransform
}

func (r *Renderer) ClipBounds() (left, right, top, bottom float32) {
	return r.clipBounds.X, r.clipBounds.Y, r.clipBounds.Z, r.clipBounds.W
}

func (r *Renderer) SetClipBounds(left, right, top, bottom float32) {
	r.clipBounds = sprec.NewVec4(left, right, top, bottom)
}

func (r *Renderer) BeginShape(fill Fill) *Shape {
	if r.shape == nil {
		panic("shape already started")
	}
	result := r.shape
	result.Init(fill)
	r.shape = nil
	return result
}

func (r *Renderer) EndShape(shape *Shape) {
	if r.shape != nil {
		panic("shape already ended")
	}
	r.shape = shape

	vertexOffset := r.shapeMesh.Offset()
	for _, point := range shape.points {
		r.shapeMesh.Append(ShapeVertex{
			position: point.coords,
		})
	}

	// draw stencil mask for all sub-shapes
	if shape.fill.mode != StencilModeNone {
		for _, subShape := range shape.subShapes {
			r.subMeshes = append(r.subMeshes, SubMesh{
				clipBounds:      r.clipBounds,
				material:        r.shapeBlankMaterial,
				vertexArray:     r.shapeMesh.vertexArray,
				transformMatrix: r.transformMatrix,
				vertexOffset:    vertexOffset + subShape.pointOffset,
				vertexCount:     subShape.pointCount,
				primitive:       wasmgl.TRIANGLE_FAN,
				skipColor:       true,
				stencil:         true,
				stencilCfg: stencilConfig{
					stencilFuncFront: stencilFunc{
						fn:   wasmgl.ALWAYS,
						ref:  0,
						mask: 0xFF,
					},
					stencilFuncBack: stencilFunc{
						fn:   wasmgl.ALWAYS,
						ref:  0,
						mask: 0xFF,
					},
					stencilOpFront: stencilOp{
						sfail:  wasmgl.KEEP,
						dpfail: wasmgl.KEEP,
						dppass: wasmgl.INCR_WRAP,
					},
					stencilOpBack: stencilOp{
						sfail:  wasmgl.KEEP,
						dpfail: wasmgl.KEEP,
						dppass: wasmgl.DECR_WRAP,
					},
				},
			})
		}
	}

	// render color shader shape and clear stencil buffer
	for _, subShape := range shape.subShapes {
		texture := r.whiteMask
		if shape.fill.image != nil {
			texture = shape.fill.image.texture
		}

		stencilMask := int(0xFF)
		if shape.fill.mode == StencilModeOdd {
			stencilMask = int(0x01)
		}

		r.subMeshes = append(r.subMeshes, SubMesh{
			clipBounds:             r.clipBounds,
			material:               r.shapeMaterial,
			vertexArray:            r.shapeMesh.vertexArray,
			transformMatrix:        r.transformMatrix,
			textureTransformMatrix: r.textureTransformMatrix,
			texture:                texture,
			color:                  shape.fill.color,
			vertexOffset:           vertexOffset + subShape.pointOffset,
			vertexCount:            subShape.pointCount,
			primitive:              wasmgl.TRIANGLE_FAN,
			stencil:                shape.fill.mode != StencilModeNone,
			stencilCfg: stencilConfig{
				stencilFuncFront: stencilFunc{
					fn:   wasmgl.NOTEQUAL,
					ref:  0,
					mask: stencilMask,
				},
				stencilFuncBack: stencilFunc{
					fn:   wasmgl.NOTEQUAL,
					ref:  0,
					mask: stencilMask,
				},
				stencilOpFront: stencilOp{
					sfail:  wasmgl.KEEP,
					dpfail: wasmgl.KEEP,
					dppass: wasmgl.REPLACE,
				},
				stencilOpBack: stencilOp{
					sfail:  wasmgl.KEEP,
					dpfail: wasmgl.KEEP,
					dppass: wasmgl.REPLACE,
				},
			},
		})
	}
}

func (r *Renderer) BeginContour() *Contour {
	if r.contour == nil {
		panic("contour already started")
	}
	result := r.contour
	result.Init()
	r.contour = nil
	return result
}

func (r *Renderer) EndContour(contour *Contour) {
	if r.contour != nil {
		panic("contour already ended")
	}
	r.contour = contour

	for _, subContour := range contour.subContours {
		pointIndex := subContour.pointOffset
		current := contour.points[pointIndex]
		next := contour.points[pointIndex+1]
		currentNormal := endPointNormal(
			current.coords,
			next.coords,
		)
		pointIndex++

		vertexOffset := r.contourMesh.Offset()
		for pointIndex < subContour.pointOffset+subContour.pointCount {
			prev := current
			prevNormal := currentNormal

			current = contour.points[pointIndex]
			if pointIndex != subContour.pointOffset+subContour.pointCount-1 {
				next := contour.points[pointIndex+1]
				currentNormal = midPointNormal(
					prev.coords,
					current.coords,
					next.coords,
				)
			} else {
				currentNormal = endPointNormal(
					prev.coords,
					current.coords,
				)
			}

			prevLeft := ContourVertex{
				position: sprec.Vec2Sum(prev.coords, sprec.Vec2Prod(prevNormal, prev.stroke.innerSize)),
				color:    prev.stroke.color,
			}
			prevRight := ContourVertex{
				position: sprec.Vec2Diff(prev.coords, sprec.Vec2Prod(prevNormal, prev.stroke.outerSize)),
				color:    prev.stroke.color,
			}
			currentLeft := ContourVertex{
				position: sprec.Vec2Sum(current.coords, sprec.Vec2Prod(currentNormal, current.stroke.innerSize)),
				color:    prev.stroke.color,
			}
			currentRight := ContourVertex{
				position: sprec.Vec2Diff(current.coords, sprec.Vec2Prod(currentNormal, current.stroke.outerSize)),
				color:    prev.stroke.color,
			}

			r.contourMesh.Append(prevLeft)
			r.contourMesh.Append(prevRight)
			r.contourMesh.Append(currentLeft)

			r.contourMesh.Append(prevRight)
			r.contourMesh.Append(currentRight)
			r.contourMesh.Append(currentLeft)

			pointIndex++
		}
		vertexCount := r.contourMesh.Offset() - vertexOffset

		r.subMeshes = append(r.subMeshes, SubMesh{
			clipBounds:      r.clipBounds,
			material:        r.contourMaterial,
			vertexArray:     r.contourMesh.vertexArray,
			transformMatrix: r.transformMatrix,
			vertexOffset:    vertexOffset,
			vertexCount:     vertexCount,
			primitive:       wasmgl.TRIANGLES,
		})
	}
}

func (r *Renderer) BeginText(typography Typography) *Text {
	if r.text == nil {
		panic("text already started")
	}
	result := r.text
	result.Init(typography)
	r.text = nil
	return result
}

func (r *Renderer) EndText(text *Text) {
	if r.text != nil {
		panic("text already ended")
	}
	r.text = text

	vertexOffset := r.textMesh.Offset()
	for _, paragraph := range text.paragraphs {
		offset := paragraph.position
		lastGlyph := (*fontGlyph)(nil)

		paragraphChars := text.characters[paragraph.charOffset : paragraph.charOffset+paragraph.charCount]
		for _, ch := range paragraphChars {
			lineHeight := text.font.lineHeight * text.fontSize
			lineAscent := text.font.lineAscent * text.fontSize
			if ch == '\r' {
				offset.X = paragraph.position.X
				lastGlyph = nil
				continue
			}
			if ch == '\n' {
				offset.X = paragraph.position.X
				offset.Y += lineHeight
				lastGlyph = nil
				continue
			}

			if glyph, ok := text.font.glyphs[ch]; ok {
				advance := glyph.advance * text.fontSize
				leftBearing := glyph.leftBearing * text.fontSize
				rightBearing := glyph.rightBearing * text.fontSize
				ascent := glyph.ascent * text.fontSize
				descent := glyph.descent * text.fontSize

				vertTopLeft := TextVertex{
					position: sprec.Vec2Sum(
						sprec.NewVec2(
							leftBearing,
							lineAscent-ascent,
						),
						offset,
					),
					texCoord: sprec.NewVec2(glyph.leftU, glyph.topV),
				}
				vertTopRight := TextVertex{
					position: sprec.Vec2Sum(
						sprec.NewVec2(
							advance-rightBearing,
							lineAscent-ascent,
						),
						offset,
					),
					texCoord: sprec.NewVec2(glyph.rightU, glyph.topV),
				}
				vertBottomLeft := TextVertex{
					position: sprec.Vec2Sum(
						sprec.NewVec2(
							leftBearing,
							lineAscent+descent,
						),
						offset,
					),
					texCoord: sprec.NewVec2(glyph.leftU, glyph.bottomV),
				}
				vertBottomRight := TextVertex{
					position: sprec.Vec2Sum(
						sprec.NewVec2(
							advance-rightBearing,
							lineAscent+descent,
						),
						offset,
					),
					texCoord: sprec.NewVec2(glyph.rightU, glyph.bottomV),
				}

				r.textMesh.Append(vertTopLeft)
				r.textMesh.Append(vertBottomLeft)
				r.textMesh.Append(vertBottomRight)

				r.textMesh.Append(vertTopLeft)
				r.textMesh.Append(vertBottomRight)
				r.textMesh.Append(vertTopRight)

				offset.X += advance
				if lastGlyph != nil {
					offset.X += lastGlyph.kerns[ch] * text.fontSize
				}
				lastGlyph = glyph
			}
		}
	}
	vertexCount := r.textMesh.Offset() - vertexOffset

	r.subMeshes = append(r.subMeshes, SubMesh{
		clipBounds:      r.clipBounds,
		material:        r.textMaterial,
		vertexArray:     r.textMesh.vertexArray,
		transformMatrix: r.transformMatrix,
		texture:         text.font.texture,
		color:           text.color,
		vertexOffset:    vertexOffset,
		vertexCount:     vertexCount,
		primitive:       wasmgl.TRIANGLES,
	})
}

func (r *Renderer) DrawSurface(surface Surface) {
	r.subMeshes = append(r.subMeshes, SubMesh{
		surface:    surface,
		clipBounds: r.clipBounds,
	})
}

func (r *Renderer) Begin(target Target) {
	r.target = target
	r.transformMatrix = sprec.IdentityMat4()
	r.textureTransformMatrix = sprec.IdentityMat4()
	r.clipBounds = sprec.NewVec4(
		0.0, float32(target.Width),
		0.0, float32(target.Height),
	)
	r.shapeMesh.Reset()
	r.contourMesh.Reset()
	r.textMesh.Reset()
	r.subMeshes = r.subMeshes[:0]
}

func (r *Renderer) End() {
	r.shapeMesh.Update()
	r.contourMesh.Update()
	r.textMesh.Update()

	r.enableOptions()

	projectionMatrix := sprec.OrthoMat4(
		0.0, float32(r.target.Width),
		0.0, float32(r.target.Height),
		0.0, 1.0,
	).ColumnMajorArray()

	// TODO: Maybe optimize by accumulating draw commands
	// if they are similar.
	for _, subMesh := range r.subMeshes {
		if subMesh.surface != nil {
			r.disableOptions()

			x := int(subMesh.clipBounds.X)
			y := int(subMesh.clipBounds.Z)
			width := int(subMesh.clipBounds.Y - subMesh.clipBounds.X)
			height := int(subMesh.clipBounds.W - subMesh.clipBounds.Z)
			subMesh.surface.Render(
				x,
				r.target.Height-height-y,
				width,
				height,
			)

			r.enableOptions()
			continue
		}

		material := subMesh.material
		transformMatrix := subMesh.transformMatrix.ColumnMajorArray()
		textureTransformMatrix := subMesh.textureTransformMatrix.ColumnMajorArray()

		if subMesh.skipColor {
			wasmgl.ColorMask(false, false, false, false)
		} else {
			wasmgl.ColorMask(true, true, true, true)
		}
		if subMesh.stencil {
			wasmgl.Enable(wasmgl.STENCIL_TEST)

			cfg := subMesh.stencilCfg
			wasmgl.StencilFuncSeparate(wasmgl.FRONT, cfg.stencilFuncFront.fn, cfg.stencilFuncFront.ref, cfg.stencilFuncFront.mask)
			wasmgl.StencilFuncSeparate(wasmgl.BACK, cfg.stencilFuncBack.fn, cfg.stencilFuncBack.ref, cfg.stencilFuncBack.mask)
			wasmgl.StencilOpSeparate(wasmgl.FRONT, cfg.stencilOpFront.sfail, cfg.stencilOpFront.dpfail, cfg.stencilOpFront.dppass)
			wasmgl.StencilOpSeparate(wasmgl.BACK, cfg.stencilOpBack.sfail, cfg.stencilOpBack.dpfail, cfg.stencilOpBack.dppass)
		} else {
			wasmgl.Disable(wasmgl.STENCIL_TEST)
		}
		if subMesh.culling {
			wasmgl.Enable(wasmgl.CULL_FACE)
			wasmgl.CullFace(subMesh.cullFace)
		} else {
			wasmgl.Disable(wasmgl.CULL_FACE)
		}
		material.program.Use()
		wasmgl.Uniform4f(material.clipDistancesLocation, subMesh.clipBounds.X, subMesh.clipBounds.Y, subMesh.clipBounds.Z, subMesh.clipBounds.W)
		wasmgl.UniformMatrix4fv(material.transformMatrixLocation, false, transformMatrix[:])
		if material.textureTransformMatrixLocation.IsValid() {
			wasmgl.UniformMatrix4fv(material.textureTransformMatrixLocation, false, textureTransformMatrix[:])
		}
		wasmgl.UniformMatrix4fv(material.projectionMatrixLocation, false, projectionMatrix[:])
		if material.colorLocation.IsValid() {
			wasmgl.Uniform4f(material.colorLocation, subMesh.color.X, subMesh.color.Y, subMesh.color.Z, subMesh.color.W)
		}
		if material.textureLocation.IsValid() {
			wasmgl.ActiveTexture(wasmgl.TEXTURE0)
			subMesh.texture.Use()
			wasmgl.Uniform1i(material.textureLocation, 0)
		}
		subMesh.vertexArray.Use()
		wasmgl.DrawArrays(subMesh.primitive, subMesh.vertexOffset, subMesh.vertexCount)
	}

	r.disableOptions()
}

func (r *Renderer) enableOptions() {
	r.target.Framebuffer.Use()

	wasmgl.Viewport(0, 0, r.target.Width, r.target.Height)
	wasmgl.ClearStencil(0)
	wasmgl.Clear(wasmgl.STENCIL_BUFFER_BIT)
	wasmgl.Disable(wasmgl.DEPTH_TEST)
	wasmgl.DepthMask(false)
	wasmgl.Enable(wasmgl.BLEND)
	wasmgl.BlendFunc(wasmgl.SRC_ALPHA, wasmgl.ONE_MINUS_SRC_ALPHA)

	// gl.Enable(gl.CLIP_DISTANCE0)
	// gl.Enable(gl.CLIP_DISTANCE1)
	// gl.Enable(gl.CLIP_DISTANCE2)
	// gl.Enable(gl.CLIP_DISTANCE3)
}

func (r *Renderer) disableOptions() {
	// gl.Disable(gl.CLIP_DISTANCE0)
	// gl.Disable(gl.CLIP_DISTANCE1)
	// gl.Disable(gl.CLIP_DISTANCE2)
	// gl.Disable(gl.CLIP_DISTANCE3)

	wasmgl.ColorMask(true, true, true, true)
	wasmgl.Disable(wasmgl.STENCIL_TEST)
	wasmgl.Enable(wasmgl.CULL_FACE)
	wasmgl.CullFace(wasmgl.BACK)
	wasmgl.Disable(wasmgl.BLEND)
}

func midPointNormal(prev, middle, next sprec.Vec2) sprec.Vec2 {
	normal1 := endPointNormal(prev, middle)
	normal2 := endPointNormal(middle, next)
	normalSum := sprec.Vec2Sum(normal1, normal2)
	dot := sprec.Vec2Dot(normal1, normalSum)
	return sprec.Vec2Quot(normalSum, dot)
}

func endPointNormal(prev, next sprec.Vec2) sprec.Vec2 {
	tangent := sprec.UnitVec2(sprec.Vec2Diff(next, prev))
	return sprec.NewVec2(tangent.Y, -tangent.X)
}

type Surface interface {
	Render(x, y, width, height int)
}
