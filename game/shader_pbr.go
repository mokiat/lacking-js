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
	if definition.AlphaTesting {
		vsBuilder.AddFeature("USES_ALPHA_TEST")
		fsBuilder.AddFeature("USES_ALPHA_TEST")
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

layout (std140) uniform Camera
{
	mat4 projectionMatrixIn;
	mat4 viewMatrixIn;
	mat4 cameraMatrixIn;
};

layout (std140) uniform Model
{
	mat4 modelMatrixIn[256];
};

smooth out vec3 normalInOut;
#if defined(USES_TEX_COORD0)
smooth out vec2 texCoordInOut;
#endif

void main()
{
#if defined(USES_TEX_COORD0)
	texCoordInOut = texCoordIn;
#endif
	mat4 modelMatrix = modelMatrixIn[gl_InstanceID];
	normalInOut = inverse(transpose(mat3(modelMatrix))) * normalIn;
	gl_Position = projectionMatrixIn * (viewMatrixIn * (modelMatrix * coordIn));
}
`

const pbrGeometryFragmentShader = `
layout(location = 0) out vec4 fbColor0Out;
layout(location = 1) out vec4 fbColor1Out;

#if defined(USES_ALBEDO_TEXTURE)
uniform sampler2D albedoTwoDTextureIn;
#endif

layout (std140) uniform Material
{
	vec4 albedoColorIn;
	float alphaThresholdIn;
	float normalScaleIn;
	float metallicIn;
	float roughnessIn;
};

smooth in vec3 normalInOut;
#if defined(USES_TEX_COORD0)
smooth in vec2 texCoordInOut;
#endif

void main()
{
#if defined(USES_ALBEDO_TEXTURE) && defined(USES_TEX_COORD0)
	vec4 color = texture(albedoTwoDTextureIn, texCoordInOut);
#else
	vec4 color = albedoColorIn;
#endif

#if defined(USES_ALPHA_TEST)
	if (color.a < alphaThresholdIn) {
		discard;
	}
#endif

	fbColor0Out = vec4(color.xyz, metallicIn);
	fbColor1Out = vec4(normalize(normalInOut), roughnessIn);
}
`
