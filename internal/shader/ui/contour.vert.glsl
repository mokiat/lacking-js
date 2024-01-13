/*template "version.glsl"*/

layout(location = 0) in vec2 positionIn;
layout(location = 2) in vec4 colorIn;

/*template "ubo_camera.glsl"*/

/*template "ubo_model.glsl"*/

smooth out vec4 clipDistancesInOut;
smooth out vec4 colorInOut;

void main()
{
	colorInOut = colorIn;
	vec4 screenPosition = transformMatrixIn * vec4(positionIn, 0.0, 1.0);

	clipDistancesInOut = clipMatrixIn * screenPosition;

	gl_Position = projectionMatrixIn * screenPosition;
}
