/*template "version.glsl"*/

layout(location = 0) in vec2 positionIn;

/*template "ubo_camera.glsl"*/

/*template "ubo_model.glsl"*/

/*template "ubo_material.glsl"*/

smooth out vec4 clipDistancesInOut;
smooth out vec2 texCoordInOut;

void main()
{
	vec4 screenPosition = transformMatrixIn * vec4(positionIn, 0.0, 1.0);
	texCoordInOut = (textureTransformMatrixIn * vec4(positionIn, 0.0, 1.0)).xy;

	clipDistancesInOut = clipMatrixIn * screenPosition;

	gl_Position = projectionMatrixIn * screenPosition;
}
