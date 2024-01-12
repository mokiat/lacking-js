/*template "version.glsl"*/
/*template "precision.glsl"*/

layout(location = 0) out vec4 fragmentColor;

smooth in vec4 clipDistancesInOut;
smooth in vec4 colorInOut;

void main()
{
	float dist = min(min(clipDistancesInOut.x, clipDistancesInOut.y), min(clipDistancesInOut.z, clipDistancesInOut.w));
	if (dist < 0.0) {
		discard;
	}

	fragmentColor = colorInOut;
}
