package graphics

import (
	"fmt"

	"github.com/mokiat/gomath/sprec"
	"github.com/mokiat/lacking-js/game/graphics/internal"
	"github.com/mokiat/lacking-js/webgl"
	"github.com/mokiat/lacking/game/graphics"
	"github.com/mokiat/wasmgl"
)

const (
	framebufferWidth  = 1920
	framebufferHeight = 1080

	coordAttributeIndex    = 0
	normalAttributeIndex   = 1
	tangentAttributeIndex  = 2
	texCoordAttributeIndex = 3
	colorAttributeIndex    = 4
)

func newRenderer() *Renderer {
	return &Renderer{
		framebufferWidth:  framebufferWidth,
		framebufferHeight: framebufferHeight,

		geometryAlbedoTexture: webgl.NewTwoDTexture(),
		geometryNormalTexture: webgl.NewTwoDTexture(),
		geometryDepthTexture:  webgl.NewTwoDTexture(),
		geometryFramebuffer:   webgl.NewFramebuffer(),

		lightingAlbedoTexture: webgl.NewTwoDTexture(),
		lightingDepthTexture:  webgl.NewTwoDTexture(),
		lightingFramebuffer:   webgl.NewFramebuffer(),

		exposureAlbedoTexture: webgl.NewTwoDTexture(),
		exposureFramebuffer:   webgl.NewFramebuffer(),
		exposureBuffer:        webgl.NewBuffer(),
		exposureSync:          wasmgl.NilSync,
		exposureTarget:        1.0,

		screenFramebuffer: webgl.DefaultFramebuffer(),

		quadMesh: newQuadMesh(),

		skyboxMesh: newSkyboxMesh(),
	}
}

type Renderer struct {
	framebufferWidth  int
	framebufferHeight int

	geometryAlbedoTexture *webgl.TwoDTexture
	geometryNormalTexture *webgl.TwoDTexture
	geometryDepthTexture  *webgl.TwoDTexture
	geometryFramebuffer   *webgl.Framebuffer

	lightingAlbedoTexture *webgl.TwoDTexture
	lightingDepthTexture  *webgl.TwoDTexture
	lightingFramebuffer   *webgl.Framebuffer

	exposureAlbedoTexture *webgl.TwoDTexture
	exposureFramebuffer   *webgl.Framebuffer
	exposurePresentation  *internal.LightingPresentation
	exposureBuffer        *webgl.Buffer
	exposureSync          wasmgl.Sync
	exposureTarget        float32

	screenFramebuffer *webgl.Framebuffer

	postprocessingPresentation *internal.PostprocessingPresentation

	directionalLightPresentation *internal.LightingPresentation
	ambientLightPresentation     *internal.LightingPresentation

	quadMesh *QuadMesh

	skyboxPresentation   *internal.SkyboxPresentation
	skycolorPresentation *internal.SkyboxPresentation
	skyboxMesh           *SkyboxMesh
}

func (r *Renderer) Allocate() {
	if wasmgl.GetExtension("EXT_color_buffer_float") == nil {
		panic(fmt.Errorf("EXT_color_buffer_float not supported"))
	}

	r.geometryAlbedoTexture.Allocate(webgl.TwoDTextureAllocateInfo{
		Width:             framebufferWidth,
		Height:            framebufferHeight,
		MinFilter:         wasmgl.NEAREST,
		MagFilter:         wasmgl.NEAREST,
		InternalFormat:    wasmgl.RGBA8,
		DataFormat:        wasmgl.RGBA,
		DataComponentType: wasmgl.UNSIGNED_BYTE,
	})

	r.geometryNormalTexture.Allocate(webgl.TwoDTextureAllocateInfo{
		Width:             framebufferWidth,
		Height:            framebufferHeight,
		MinFilter:         wasmgl.NEAREST,
		MagFilter:         wasmgl.NEAREST,
		InternalFormat:    wasmgl.RGBA32F,
		DataFormat:        wasmgl.RGBA,
		DataComponentType: wasmgl.FLOAT,
	})

	r.geometryDepthTexture.Allocate(webgl.TwoDTextureAllocateInfo{
		Width:             framebufferWidth,
		Height:            framebufferHeight,
		MinFilter:         wasmgl.NEAREST,
		MagFilter:         wasmgl.NEAREST,
		InternalFormat:    wasmgl.DEPTH_COMPONENT32F,
		DataFormat:        wasmgl.DEPTH_COMPONENT,
		DataComponentType: wasmgl.FLOAT,
	})

	r.geometryFramebuffer.Allocate(webgl.FramebufferAllocateInfo{
		ColorAttachments: []*webgl.Texture{
			&r.geometryAlbedoTexture.Texture,
			&r.geometryNormalTexture.Texture,
		},
		DepthAttachment: &r.geometryDepthTexture.Texture,
	})

	r.lightingAlbedoTexture.Allocate(webgl.TwoDTextureAllocateInfo{
		Width:             framebufferWidth,
		Height:            framebufferHeight,
		MinFilter:         wasmgl.NEAREST,
		MagFilter:         wasmgl.NEAREST,
		InternalFormat:    wasmgl.RGBA32F,
		DataFormat:        wasmgl.RGBA,
		DataComponentType: wasmgl.FLOAT,
	})

	r.lightingDepthTexture.Allocate(webgl.TwoDTextureAllocateInfo{
		Width:             framebufferWidth,
		Height:            framebufferHeight,
		MinFilter:         wasmgl.NEAREST,
		MagFilter:         wasmgl.NEAREST,
		InternalFormat:    wasmgl.DEPTH_COMPONENT32F,
		DataFormat:        wasmgl.DEPTH_COMPONENT,
		DataComponentType: wasmgl.FLOAT,
	})

	r.lightingFramebuffer.Allocate(webgl.FramebufferAllocateInfo{
		ColorAttachments: []*webgl.Texture{
			&r.lightingAlbedoTexture.Texture,
		},
		DepthAttachment: &r.lightingDepthTexture.Texture,
	})

	r.exposureAlbedoTexture.Allocate(webgl.TwoDTextureAllocateInfo{
		Width:             1,
		Height:            1,
		MinFilter:         wasmgl.NEAREST,
		MagFilter:         wasmgl.NEAREST,
		InternalFormat:    wasmgl.RGBA32F,
		DataFormat:        wasmgl.RGBA,
		DataComponentType: wasmgl.FLOAT,
	})

	r.exposureFramebuffer.Allocate(webgl.FramebufferAllocateInfo{
		ColorAttachments: []*webgl.Texture{
			&r.exposureAlbedoTexture.Texture,
		},
	})
	r.exposurePresentation = internal.NewExposurePresentation()

	r.exposureBuffer.Allocate(webgl.BufferAllocateInfo{
		BufferType: wasmgl.PIXEL_PACK_BUFFER,
		Dynamic:    true,
		Data:       make([]byte, 4*4),
	})

	r.postprocessingPresentation = internal.NewTonePostprocessingPresentation(internal.ExponentialToneMapping)

	r.directionalLightPresentation = internal.NewDirectionalLightPresentation()
	r.ambientLightPresentation = internal.NewAmbientLightPresentation()

	r.quadMesh.Allocate()

	r.skyboxPresentation = internal.NewCubeSkyboxPresentation()
	r.skycolorPresentation = internal.NewColorSkyboxPresentation()
	r.skyboxMesh.Allocate()
}

func (r *Renderer) Release() {
	r.skyboxPresentation.Delete()
	r.skycolorPresentation.Delete()
	r.skyboxMesh.Release()

	r.ambientLightPresentation.Delete()
	r.directionalLightPresentation.Delete()

	r.quadMesh.Release()

	r.postprocessingPresentation.Delete()

	r.exposureBuffer.Release()
	r.exposurePresentation.Delete()
	r.exposureFramebuffer.Release()
	r.exposureAlbedoTexture.Release()

	r.lightingFramebuffer.Release()
	r.lightingAlbedoTexture.Release()
	r.lightingDepthTexture.Release()

	r.geometryFramebuffer.Release()
	r.geometryDepthTexture.Release()
	r.geometryNormalTexture.Release()
	r.geometryAlbedoTexture.Release()
}

type renderCtx struct {
	scene            *Scene
	x                int
	y                int
	width            int
	height           int
	projectionMatrix [16]float32
	cameraMatrix     [16]float32
	viewMatrix       [16]float32
	camera           *Camera
}

func (r *Renderer) Render(viewport graphics.Viewport, scene *Scene, camera *Camera) {
	projectionMatrix := r.evaluateProjectionMatrix(camera, viewport.Width, viewport.Height)
	cameraMatrix := camera.ModelMatrix()
	viewMatrix := sprec.InverseMat4(cameraMatrix)

	ctx := renderCtx{
		scene:            scene,
		x:                viewport.X,
		y:                viewport.Y,
		width:            viewport.Width,
		height:           viewport.Height,
		projectionMatrix: projectionMatrix.ColumnMajorArray(),
		cameraMatrix:     cameraMatrix.ColumnMajorArray(),
		viewMatrix:       viewMatrix.ColumnMajorArray(),
		camera:           camera,
	}
	r.renderGeometryPass(ctx)
	// gl.TextureBarrier()
	r.renderLightingPass(ctx)
	r.renderForwardPass(ctx)
	if camera.autoExposureEnabled {
		// gl.TextureBarrier()
		r.renderExposureProbePass(ctx)
	}
	r.renderPostprocessingPass(ctx)
}

func (r *Renderer) evaluateProjectionMatrix(camera *Camera, width, height int) sprec.Mat4 {
	const (
		near = float32(0.5)
		far  = float32(900.0)
	)
	var (
		fWidth  = sprec.Max(1.0, float32(width))
		fHeight = sprec.Max(1.0, float32(height))
	)

	switch camera.fovMode {
	case graphics.FoVModeHorizontalPlus:
		halfHeight := near * sprec.Tan(camera.fov/2.0)
		halfWidth := halfHeight * (fWidth / fHeight)
		return sprec.PerspectiveMat4(
			-halfWidth, halfWidth, -halfHeight, halfHeight, near, far,
		)

	case graphics.FoVModeVertialMinus:
		halfWidth := near * sprec.Tan(camera.fov/2.0)
		halfHeight := halfWidth * (fHeight / fWidth)
		return sprec.PerspectiveMat4(
			-halfWidth, halfWidth, -halfHeight, halfHeight, near, far,
		)

	case graphics.FoVModePixelBased:
		halfWidth := fWidth / 2.0
		halfHeight := fHeight / 2.0
		return sprec.OrthoMat4(
			-halfWidth, halfWidth, halfHeight, -halfHeight, near, far,
		)

	default:
		panic(fmt.Errorf("unsupported fov mode: %s", camera.fovMode))
	}
}

func (r *Renderer) renderGeometryPass(ctx renderCtx) {
	r.geometryFramebuffer.Use()

	wasmgl.Viewport(0, 0, r.framebufferWidth, r.framebufferHeight)
	wasmgl.Enable(wasmgl.DEPTH_TEST)
	wasmgl.DepthMask(true)
	wasmgl.DepthFunc(wasmgl.LEQUAL)

	wasmgl.ClearColor(
		ctx.scene.sky.backgroundColor.X,
		ctx.scene.sky.backgroundColor.Y,
		ctx.scene.sky.backgroundColor.Z,
		1.0,
	)
	wasmgl.ClearDepth(1.0)
	wasmgl.Clear(wasmgl.COLOR_BUFFER_BIT | wasmgl.DEPTH_BUFFER_BIT)

	// TODO: Traverse octree
	for mesh := ctx.scene.firstMesh; mesh != nil; mesh = mesh.next {
		r.renderMesh(ctx, mesh.ModelMatrix().ColumnMajorArray(), mesh.template)
	}
}

func (r *Renderer) renderMesh(ctx renderCtx, modelMatrix [16]float32, template *MeshTemplate) {
	for _, subMesh := range template.subMeshes {
		if subMesh.material.backfaceCulling {
			wasmgl.Enable(wasmgl.CULL_FACE)
		} else {
			wasmgl.Disable(wasmgl.CULL_FACE)
		}

		material := subMesh.material
		presentation := material.geometryPresentation
		presentation.Program.Use()

		wasmgl.UniformMatrix4fv(presentation.ProjectionMatrixLocation, false, ctx.projectionMatrix[:])
		wasmgl.UniformMatrix4fv(presentation.ViewMatrixLocation, false, ctx.viewMatrix[:])
		wasmgl.UniformMatrix4fv(presentation.ModelMatrixLocation, false, modelMatrix[:])

		wasmgl.Uniform1f(presentation.MetalnessLocation, material.vectors[1].Y)
		wasmgl.Uniform1f(presentation.RoughnessLocation, material.vectors[1].Z)
		wasmgl.Uniform4f(presentation.AlbedoColorLocation, material.vectors[0].X, material.vectors[0].Y, material.vectors[0].Z, material.vectors[0].Z)

		textureUnit := 0
		if material.twoDTextures[0] != nil {
			wasmgl.ActiveTexture(wasmgl.TEXTURE0 + textureUnit)
			material.twoDTextures[0].Use()
			wasmgl.Uniform1i(presentation.AlbedoTextureLocation, textureUnit)
		}

		template.vertexArray.Use()
		wasmgl.DrawElements(subMesh.primitive, subMesh.indexCount, subMesh.indexType, subMesh.indexOffsetBytes)
	}
}

func (r *Renderer) renderLightingPass(ctx renderCtx) {
	wasmgl.BindFramebuffer(wasmgl.READ_FRAMEBUFFER, r.geometryFramebuffer.Raw())
	wasmgl.BindFramebuffer(wasmgl.DRAW_FRAMEBUFFER, r.lightingFramebuffer.Raw())

	wasmgl.BlitFramebuffer(
		0, 0, r.framebufferWidth, r.framebufferHeight,
		0, 0, r.framebufferWidth, r.framebufferHeight,
		wasmgl.DEPTH_BUFFER_BIT,
		wasmgl.NEAREST,
	)

	r.lightingFramebuffer.Use()

	wasmgl.Viewport(0, 0, r.framebufferWidth, r.framebufferHeight)
	wasmgl.Disable(wasmgl.DEPTH_TEST)
	wasmgl.DepthMask(false)
	wasmgl.Enable(wasmgl.CULL_FACE)

	wasmgl.ClearColor(0.0, 0.0, 0.0, 1.0)
	wasmgl.Clear(wasmgl.COLOR_BUFFER_BIT) // r.lightingFramebuffer.ClearColor(0, )

	wasmgl.Enable(wasmgl.BLEND) // 	gl.Enablei(gl.BLEND, 0)

	wasmgl.BlendEquationSeparate(wasmgl.FUNC_ADD, wasmgl.FUNC_ADD)
	wasmgl.BlendFuncSeparate(wasmgl.ONE, wasmgl.ONE, wasmgl.ONE, wasmgl.ZERO)

	// TODO: Traverse octree
	for light := ctx.scene.firstLight; light != nil; light = light.next {
		switch light.mode {
		case LightModeDirectional:
			r.renderDirectionalLight(ctx, light)
		case LightModeAmbient:
			r.renderAmbientLight(ctx, light)
		}
	}

	wasmgl.Disable(wasmgl.BLEND) // 	gl.Disablei(gl.BLEND, 0)
}

func (r *Renderer) renderAmbientLight(ctx renderCtx, light *Light) {
	presentation := r.ambientLightPresentation
	presentation.Program.Use()

	wasmgl.UniformMatrix4fv(presentation.ProjectionMatrixLocation, false, ctx.projectionMatrix[:])
	wasmgl.UniformMatrix4fv(presentation.CameraMatrixLocation, false, ctx.cameraMatrix[:])
	wasmgl.UniformMatrix4fv(presentation.ViewMatrixLocation, false, ctx.viewMatrix[:])

	textureUnit := 0

	wasmgl.ActiveTexture(wasmgl.TEXTURE0 + textureUnit)
	r.geometryAlbedoTexture.Use()
	wasmgl.Uniform1i(presentation.FramebufferDraw0Location, textureUnit)
	textureUnit++

	wasmgl.ActiveTexture(wasmgl.TEXTURE0 + textureUnit)
	r.geometryNormalTexture.Use()
	wasmgl.Uniform1i(presentation.FramebufferDraw1Location, textureUnit)
	textureUnit++

	wasmgl.ActiveTexture(wasmgl.TEXTURE0 + textureUnit)
	r.geometryDepthTexture.Use()
	wasmgl.Uniform1i(presentation.FramebufferDepthLocation, textureUnit)
	textureUnit++

	wasmgl.ActiveTexture(wasmgl.TEXTURE0 + textureUnit)
	light.reflectionTexture.Use()
	wasmgl.Uniform1i(presentation.ReflectionTextureLocation, textureUnit)
	textureUnit++

	wasmgl.ActiveTexture(wasmgl.TEXTURE0 + textureUnit)
	light.refractionTexture.Use()
	wasmgl.Uniform1i(presentation.RefractionTextureLocation, textureUnit)
	textureUnit++

	r.quadMesh.VertexArray.Use()
	wasmgl.DrawElements(r.quadMesh.Primitive, r.quadMesh.IndexCount, wasmgl.UNSIGNED_SHORT, r.quadMesh.IndexOffsetBytes)
}

func (r *Renderer) renderDirectionalLight(ctx renderCtx, light *Light) {
	presentation := r.directionalLightPresentation
	presentation.Program.Use()

	wasmgl.UniformMatrix4fv(presentation.ProjectionMatrixLocation, false, ctx.projectionMatrix[:])
	wasmgl.UniformMatrix4fv(presentation.CameraMatrixLocation, false, ctx.cameraMatrix[:])
	wasmgl.UniformMatrix4fv(presentation.ViewMatrixLocation, false, ctx.viewMatrix[:])

	direction := light.Rotation().OrientationZ()
	wasmgl.Uniform3f(presentation.LightDirection, direction.X, direction.Y, direction.Z)
	intensity := light.intensity
	wasmgl.Uniform3f(presentation.LightIntensity, intensity.X, intensity.Y, intensity.Z)

	textureUnit := 0

	wasmgl.ActiveTexture(wasmgl.TEXTURE0 + textureUnit)
	r.geometryAlbedoTexture.Use()
	wasmgl.Uniform1i(presentation.FramebufferDraw0Location, textureUnit)
	textureUnit++

	wasmgl.ActiveTexture(wasmgl.TEXTURE0 + textureUnit)
	r.geometryNormalTexture.Use()
	wasmgl.Uniform1i(presentation.FramebufferDraw1Location, textureUnit)
	textureUnit++

	wasmgl.ActiveTexture(wasmgl.TEXTURE0 + textureUnit)
	r.geometryDepthTexture.Use()
	wasmgl.Uniform1i(presentation.FramebufferDepthLocation, textureUnit)
	textureUnit++

	r.quadMesh.VertexArray.Use()
	wasmgl.DrawElements(r.quadMesh.Primitive, r.quadMesh.IndexCount, wasmgl.UNSIGNED_SHORT, r.quadMesh.IndexOffsetBytes)
}

func (r *Renderer) renderForwardPass(ctx renderCtx) {
	r.lightingFramebuffer.Use()

	wasmgl.Viewport(0, 0, r.framebufferWidth, r.framebufferHeight)
	wasmgl.Enable(wasmgl.DEPTH_TEST)
	wasmgl.DepthMask(false)
	wasmgl.DepthFunc(wasmgl.LEQUAL)

	sky := ctx.scene.sky
	if texture := sky.skyboxTexture; texture != nil {
		wasmgl.Enable(wasmgl.CULL_FACE)

		presentation := r.skyboxPresentation
		program := presentation.Program
		program.Use()

		wasmgl.UniformMatrix4fv(presentation.ProjectionMatrixLocation, false, ctx.projectionMatrix[:])
		wasmgl.UniformMatrix4fv(presentation.ViewMatrixLocation, false, ctx.viewMatrix[:])

		wasmgl.ActiveTexture(wasmgl.TEXTURE0)
		wasmgl.BindTexture(wasmgl.TEXTURE_CUBE_MAP, texture.Raw())
		wasmgl.Uniform1i(presentation.AlbedoCubeTextureLocation, 0)

		r.skyboxMesh.VertexArray.Use()
		wasmgl.DrawElements(r.skyboxMesh.Primitive, r.skyboxMesh.IndexCount, wasmgl.UNSIGNED_SHORT, r.skyboxMesh.IndexOffsetBytes)
	} else {
		wasmgl.Enable(wasmgl.CULL_FACE)

		presentation := r.skycolorPresentation
		program := presentation.Program
		program.Use()

		wasmgl.UniformMatrix4fv(presentation.ProjectionMatrixLocation, false, ctx.projectionMatrix[:])
		wasmgl.UniformMatrix4fv(presentation.ViewMatrixLocation, false, ctx.viewMatrix[:])

		wasmgl.Uniform4f(presentation.AlbedoColorLocation,
			sky.backgroundColor.X,
			sky.backgroundColor.Y,
			sky.backgroundColor.Z,
			1.0,
		)
		r.skyboxMesh.VertexArray.Use()
		wasmgl.DrawElements(r.skyboxMesh.Primitive, r.skyboxMesh.IndexCount, wasmgl.UNSIGNED_SHORT, r.skyboxMesh.IndexOffsetBytes)
	}
}

func (r *Renderer) renderExposureProbePass(ctx renderCtx) {
	if r.exposureSync.Valid() {
		status := wasmgl.ClientWaitSync(r.exposureSync, wasmgl.SYNC_FLUSH_COMMANDS_BIT, 0)
		switch status {
		case wasmgl.ALREADY_SIGNALED, wasmgl.CONDITION_SATISFIED:
			data := make([]float32, 4)
			r.exposureBuffer.Use()
			wasmgl.GetBufferSubData(wasmgl.PIXEL_PACK_BUFFER, 0, data)
			wasmgl.BindBuffer(wasmgl.PIXEL_PACK_BUFFER, wasmgl.NilBuffer)
			brightness := 0.2126*data[0] + 0.7152*data[1] + 0.0722*data[2]
			if brightness < 0.001 {
				brightness = 0.001
			}
			r.exposureTarget = 1.0 / (3.14 * brightness)
			if r.exposureTarget > ctx.camera.maxExposure {
				r.exposureTarget = ctx.camera.maxExposure
			}
			if r.exposureTarget < ctx.camera.minExposure {
				r.exposureTarget = ctx.camera.minExposure
			}
			wasmgl.DeleteSync(r.exposureSync)
			r.exposureSync = wasmgl.NilSync
		case wasmgl.WAIT_FAILED:
			r.exposureSync = wasmgl.NilSync
		}
	}

	ctx.camera.exposure = mix(ctx.camera.exposure, r.exposureTarget, float32(0.01))

	if !r.exposureSync.Valid() {
		r.exposureFramebuffer.Use()

		wasmgl.Viewport(0, 0, r.framebufferWidth, r.framebufferHeight)
		wasmgl.Disable(wasmgl.DEPTH_TEST)
		wasmgl.DepthMask(false)
		wasmgl.Enable(wasmgl.CULL_FACE)

		wasmgl.ClearColor(0.0, 0.0, 0.0, 0.0)
		wasmgl.Clear(wasmgl.COLOR_BUFFER_BIT)

		presentation := r.exposurePresentation
		program := presentation.Program
		program.Use()

		textureUnit := 0

		wasmgl.ActiveTexture(wasmgl.TEXTURE0 + textureUnit)
		r.lightingAlbedoTexture.Use()
		wasmgl.Uniform1i(presentation.FramebufferDraw0Location, textureUnit)
		textureUnit++

		r.quadMesh.VertexArray.Use()
		wasmgl.DrawElements(r.quadMesh.Primitive, r.quadMesh.IndexCount, wasmgl.UNSIGNED_SHORT, r.quadMesh.IndexOffsetBytes)

		// 		gl.TextureBarrier()

		r.exposureBuffer.Use()
		wasmgl.ReadPixels(0, 0, 1, 1, wasmgl.RGBA, wasmgl.FLOAT, 0)
		r.exposureSync = wasmgl.FenceSync(wasmgl.SYNC_GPU_COMMANDS_COMPLETE, 0)
		wasmgl.BindBuffer(wasmgl.PIXEL_PACK_BUFFER, wasmgl.NilBuffer)

		r.screenFramebuffer.Use()
	}
}

func (r *Renderer) renderPostprocessingPass(ctx renderCtx) {
	r.screenFramebuffer.Use()
	wasmgl.Viewport(ctx.x, ctx.y, ctx.width, ctx.height)
	wasmgl.Scissor(ctx.x, ctx.y, ctx.width, ctx.height)

	wasmgl.Disable(wasmgl.DEPTH_TEST)
	wasmgl.DepthMask(false)
	wasmgl.DepthFunc(wasmgl.ALWAYS)

	wasmgl.Enable(wasmgl.CULL_FACE)
	presentation := r.postprocessingPresentation
	presentation.Program.Use()

	wasmgl.ActiveTexture(wasmgl.TEXTURE0)
	r.lightingAlbedoTexture.Use()
	wasmgl.Uniform1i(presentation.FramebufferDraw0Location, 0)

	wasmgl.Uniform1f(presentation.ExposureLocation, ctx.camera.exposure)

	r.quadMesh.VertexArray.Use()
	wasmgl.DrawElements(r.quadMesh.Primitive, r.quadMesh.IndexCount, wasmgl.UNSIGNED_SHORT, r.quadMesh.IndexOffsetBytes)
}

// TODO: Move to gomath
func mix(a, b, amount float32) float32 {
	return a*(1.0-amount) + b*amount
}
