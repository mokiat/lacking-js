/*template "version.glsl"*/

layout(location = 0) in vec2 positionIn;
layout(location = 2) in vec4 colorIn;

uniform mat4 projectionMatrixIn;
uniform mat4 transformMatrixIn;
uniform mat4 clipMatrixIn;

smooth out vec4 clipDistancesInOut;
smooth out vec4 colorInOut;

void main()
{
	colorInOut = colorIn;
	vec4 screenPosition = transformMatrixIn * vec4(positionIn, 0.0, 1.0);

	clipDistancesInOut = clipMatrixIn * screenPosition;

	gl_Position = projectionMatrixIn * screenPosition;
}
