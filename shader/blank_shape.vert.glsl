/*template "version.glsl"*/

layout(location = 0) in vec2 positionIn;

uniform mat4 projectionMatrixIn;
uniform mat4 transformMatrixIn;
uniform mat4 clipMatrixIn;

smooth out vec4 clipDistancesInOut;

void main()
{
	vec4 screenPosition = transformMatrixIn * vec4(positionIn, 0.0, 1.0);

	clipDistancesInOut = clipMatrixIn * screenPosition;

	gl_Position = projectionMatrixIn * screenPosition;
}