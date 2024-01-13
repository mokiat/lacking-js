/*template "version.glsl"*/

layout(location = 0) in vec2 positionIn;
layout(location = 1) in vec2 texCoordIn;

/*template "ubo_camera.glsl"*/

/*template "ubo_model.glsl"*/

smooth out vec4 clipDistancesInOut;
smooth out vec2 texCoordInOut;

void main()
{
	texCoordInOut = texCoordIn;
	vec4 screenPosition = transformMatrixIn * vec4(positionIn, 0.0, 1.0);

	clipDistancesInOut = clipMatrixIn * screenPosition;

	gl_Position = projectionMatrixIn * screenPosition;
}
