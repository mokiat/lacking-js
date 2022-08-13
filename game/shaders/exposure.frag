layout(location = 0) out vec4 fbColor0Out;

uniform sampler2D fbColor0TextureIn;

void main()
{
	vec3 mixture = vec3(0.0, 0.0, 0.0);
	float count = 0.0;
	for (float u = 0.0; u <= 1.0; u += 0.05) {
		for (float v = 0.0; v <= 1.0; v += 0.05) {
			mixture += clamp(texture(fbColor0TextureIn, vec2(u, v)).xyz, 0.0, 100.0);
			count += 1.0;
		}
	}
	fbColor0Out = vec4(mixture / count, 1.0);
}
