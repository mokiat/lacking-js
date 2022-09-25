layout(location = 0) in vec3 coordIn;

layout (std140) uniform Camera
{
	mat4 projectionMatrixIn;
	mat4 viewMatrixIn;
	mat4 cameraMatrixIn;
};

void main()
{
	// ensure that translations are ignored by setting w to 0.0
	vec4 viewPosition = viewMatrixIn * vec4(coordIn, 0.0);

	// restore w to 1.0 so that projection works
	vec4 position = projectionMatrixIn * vec4(viewPosition.xyz, 1.0);

	// set z to w so that it has maximum depth (1.0) after projection division
	gl_Position = vec4(position.xy, position.w, position.w);
}
