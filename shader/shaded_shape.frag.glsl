/*template "version.glsl"*/
/*template "precision.glsl"*/

layout(location = 0) out vec4 fragmentColor;

uniform sampler2D textureIn;
uniform vec4 colorIn;

smooth in vec4 clipDistancesInOut;
smooth in vec2 texCoordInOut;

void main()
{
	float dist = min(min(clipDistancesInOut.x, clipDistancesInOut.y), min(clipDistancesInOut.z, clipDistancesInOut.w));
	if (dist < 0.0) {
		discard;
	}

	fragmentColor = texture(textureIn, texCoordInOut) * colorIn;
}
