layout(location = 0) in vec4 coordIn;
#if defined(USES_BONES)
layout(location = 5) in vec4 weightsIn;
layout(location = 6) in uvec4 jointsIn;
#endif

layout (std140) uniform Light
{
	mat4 projectionMatrixIn;
	mat4 viewMatrixIn;
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

void main()
{
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
#else
	mat4 modelMatrix = modelMatrixIn[gl_InstanceID];
	vec4 worldPosition = modelMatrix * coordIn;
#endif
	gl_Position = projectionMatrixIn * (viewMatrixIn * worldPosition);
}
