layout(location = 0) in vec4 coordIn;
layout(location = 1) in vec3 normalIn;
#if defined(USES_TEX_COORD0)
layout(location = 3) in vec2 texCoordIn;
#endif
#if defined(USES_COLOR0)
layout(location = 4) in vec4 colorIn;
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
	mat4 boneMatrixIn[256];
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
#if defined(USES_COLOR0)
smooth out vec4 colorInOut;
#endif

void main()
{
#if defined(USES_TEX_COORD0)
	texCoordInOut = texCoordIn;
#endif
#if defined(USES_COLOR0)
	colorInOut = colorIn;
#endif
#if defined(USES_BONES)
	mat4 modelMatrixA = boneMatrixIn[jointsIn.x];
	mat4 modelMatrixB = boneMatrixIn[jointsIn.y];
	mat4 modelMatrixC = boneMatrixIn[jointsIn.z];
	mat4 modelMatrixD = boneMatrixIn[jointsIn.w];
	vec4 worldPosition =
		modelMatrixA * (coordIn * weightsIn.x) +
		modelMatrixB * (coordIn * weightsIn.y) +
		modelMatrixC * (coordIn * weightsIn.z) +
		modelMatrixD * (coordIn * weightsIn.w);
	vec3 worldNormal =
		inverse(transpose(mat3(modelMatrixA))) * (normalIn * weightsIn.x) +
		inverse(transpose(mat3(modelMatrixB))) * (normalIn * weightsIn.y) +
		inverse(transpose(mat3(modelMatrixC))) * (normalIn * weightsIn.z) +
		inverse(transpose(mat3(modelMatrixD))) * (normalIn * weightsIn.w);
#else
	mat4 modelMatrix = modelMatrixIn[gl_InstanceID];
	vec4 worldPosition = modelMatrix * coordIn;
	vec3 worldNormal = inverse(transpose(mat3(modelMatrix))) * normalIn;
#endif
	normalInOut = worldNormal;
	gl_Position = projectionMatrixIn * (viewMatrixIn * worldPosition);
}
