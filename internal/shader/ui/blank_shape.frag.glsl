/*template "version.glsl"*/
/*template "precision.glsl"*/

layout(location = 0) out vec4 fragmentColor;

smooth in vec4 clipDistancesInOut;

void main()
{
	float dist = min(min(clipDistancesInOut.x, clipDistancesInOut.y), min(clipDistancesInOut.z, clipDistancesInOut.w));
	if (dist < 0.0) {
		discard;
	}

	fragmentColor = vec4(1.0, 1.0, 1.0, 1.0);
}
