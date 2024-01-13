/*template "version.glsl"*/
/*template "precision.glsl"*/

layout(location = 0) out vec4 fragmentColor;

/*template "ubo_material.glsl"*/

uniform sampler2D fontTextureIn;

smooth in vec4 clipDistancesInOut;
smooth in vec2 texCoordInOut;

void main()
{
	float dist = min(min(clipDistancesInOut.x, clipDistancesInOut.y), min(clipDistancesInOut.z, clipDistancesInOut.w));
	if (dist < 0.0) {
		discard;
	}

	float amount = texture(fontTextureIn, texCoordInOut).x;
	fragmentColor = vec4(amount) * colorIn;
}
