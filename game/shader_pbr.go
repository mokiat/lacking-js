package game

import (
	"github.com/mokiat/lacking-js/internal"
	"github.com/mokiat/lacking/game/graphics"
)

func newPBRShaderSet(definition graphics.PBRMaterialDefinition) graphics.ShaderSet {
	vsBuilder := internal.NewShaderSourceBuilder(pbrGeometryVertexShader)
	fsBuilder := internal.NewShaderSourceBuilder(pbrGeometryFragmentShader)
	if definition.AlbedoTexture != nil {
		vsBuilder.AddFeature("USES_ALBEDO_TEXTURE")
		fsBuilder.AddFeature("USES_ALBEDO_TEXTURE")
		vsBuilder.AddFeature("USES_TEX_COORD0")
		fsBuilder.AddFeature("USES_TEX_COORD0")
	}
	return graphics.ShaderSet{
		VertexShader:   vsBuilder.Build,
		FragmentShader: fsBuilder.Build,
	}
}

const pbrGeometryVertexShader = `
layout(location = 0) in vec4 coordIn;
layout(location = 1) in vec3 normalIn;
#if defined(USES_TEX_COORD0)
layout(location = 3) in vec2 texCoordIn;
#endif

uniform mat4 projectionMatrixIn;
uniform mat4 modelMatrixIn;
uniform mat4 viewMatrixIn;

smooth out vec3 normalInOut;
#if defined(USES_TEX_COORD0)
smooth out vec2 texCoordInOut;
#endif

void main()
{
#if defined(USES_TEX_COORD0)
	texCoordInOut = texCoordIn;
#endif
	normalInOut = inverse(transpose(mat3(modelMatrixIn))) * normalIn;
	gl_Position = projectionMatrixIn * (viewMatrixIn * (modelMatrixIn * coordIn));
}
`

const pbrGeometryFragmentShader = `
layout(location = 0) out vec4 fbColor0Out;
layout(location = 1) out vec4 fbColor1Out;

#if defined(USES_ALBEDO_TEXTURE)
uniform sampler2D albedoTwoDTextureIn;
#endif
uniform vec4 albedoColorIn;

uniform float metalnessIn;
uniform float roughnessIn;
uniform float alphaThresholdIn;

smooth in vec3 normalInOut;
#if defined(USES_TEX_COORD0)
smooth in vec2 texCoordInOut;
#endif

void main()
{
#if defined(USES_ALBEDO_TEXTURE) && defined(USES_TEX_COORD0)
	vec4 color = texture(albedoTwoDTextureIn, texCoordInOut);
	if (color.a < 0.5) { // FIXME: USE alphaThresholdIn
		discard;
	}
#else
	vec4 color = albedoColorIn;
#endif

	fbColor0Out = vec4(color.xyz, metalnessIn);
	fbColor1Out = vec4(normalize(normalInOut), roughnessIn);
}
`
