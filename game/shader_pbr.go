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
	if definition.Armature {
		vsBuilder.AddFeature("USES_BONES")
		fsBuilder.AddFeature("USES_BONES")
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
#if defined(USES_BONES)
layout(location = 5) in vec4 weightsIn;
layout(location = 6) in uvec4 jointsIn;
#endif

layout (std140) uniform Camera
{
	mat4 projectionMatrixIn;
	mat4 viewMatrixIn;
	mat4 cameraMatrixIn;
};

#if defined(USES_BONES)
layout (std140) uniform Model
{
	mat4 modelMatrixIn;
	mat4 boneMatrixIn[255];
};
#else
layout (std140) uniform Model
{
	mat4 modelMatrixIn[256];
};
#endif

smooth out vec3 normalInOut;
#if defined(USES_TEX_COORD0)
smooth out vec2 texCoordInOut;
#endif

void main()
{
#if defined(USES_TEX_COORD0)
	texCoordInOut = texCoordIn;
#endif
#if defined(USES_BONES)
	mat4 modelMatrixA = modelMatrixIn * boneMatrixIn[jointsIn.x];
	mat4 modelMatrixB = modelMatrixIn * boneMatrixIn[jointsIn.y];
	mat4 modelMatrixC = modelMatrixIn * boneMatrixIn[jointsIn.z];
	mat4 modelMatrixD = modelMatrixIn * boneMatrixIn[jointsIn.w];
	vec4 worldPosition =
		modelMatrixA * (weightsIn.x * coordIn) +
		modelMatrixB * (weightsIn.y * coordIn) +
		modelMatrixC * (weightsIn.z * coordIn) +
		modelMatrixD * (weightsIn.w * coordIn);
	vec3 worldNormal =
		inverse(transpose(mat3(modelMatrixA))) * (weightsIn.x * normalIn) +
		inverse(transpose(mat3(modelMatrixB))) * (weightsIn.y * normalIn) +
		inverse(transpose(mat3(modelMatrixC))) * (weightsIn.z * normalIn) +
		inverse(transpose(mat3(modelMatrixD))) * (weightsIn.w * normalIn);
#else
	mat4 modelMatrix = modelMatrixIn[gl_InstanceID];
	vec4 worldPosition = modelMatrix * coordIn;
	vec3 worldNormal = inverse(transpose(mat3(modelMatrix))) * normalIn;
#endif
	normalInOut = worldNormal;
	gl_Position = projectionMatrixIn * (viewMatrixIn * worldPosition);
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
