layout(location = 0) in vec4 coordIn;
layout(location = 4) in vec3 colorIn;

layout (std140) uniform Camera
{
	mat4 projectionMatrixIn;
	mat4 viewMatrixIn;
	mat4 cameraMatrixIn;
};

smooth out vec3 colorInOut;

void main()
{
	colorInOut = colorIn;
	gl_Position = projectionMatrixIn * (viewMatrixIn * coordIn);
}
