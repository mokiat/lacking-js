layout(location = 0) out vec4 fbColor0Out;
layout(location = 1) out vec4 fbColor1Out;

#if defined(USES_ALBEDO_TEXTURE)
uniform sampler2D albedoTwoDTextureIn;
#endif

layout (std140) uniform Material
{
	vec4 albedoColorIn;
	float alphaThresholdIn;
	float normalScaleIn;
	float metallicIn;
	float roughnessIn;
};

smooth in vec3 normalInOut;
#if defined(USES_TEX_COORD0)
smooth in vec2 texCoordInOut;
#endif
#if defined(USES_COLOR0)
smooth in vec4 colorInOut;
#endif

void main()
{
#if defined(USES_ALBEDO_TEXTURE) && defined(USES_TEX_COORD0)
	vec4 color = texture(albedoTwoDTextureIn, texCoordInOut);
#elif defined(USES_COLOR0)
	vec4 color = colorInOut;
#else
	vec4 color = albedoColorIn;
#endif

#if defined(USES_ALPHA_TEST)
	if (color.a < alphaThresholdIn) {
		discard;
	}
#endif

	fbColor0Out = vec4(color.xyz, metallicIn);
	fbColor1Out = vec4(normalize(normalInOut), roughnessIn);
}
