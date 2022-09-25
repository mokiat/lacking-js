layout(location = 0) out vec4 fbColor0Out;

uniform sampler2D fbColor0TextureIn;
uniform float exposureIn;

smooth in vec2 texCoordInOut;

void main()
{
	vec3 hdr = texture(fbColor0TextureIn, texCoordInOut).xyz;
	vec3 exposedHDR = hdr * exposureIn;
	#if defined(MODE_REINHARD)
	vec3 ldr = exposedHDR / (exposedHDR + vec3(1.0));
	#endif
	#if defined(MODE_EXPONENTIAL)
	vec3 ldr = vec3(1.0) - exp2(-exposedHDR);
	#endif
	fbColor0Out = vec4(ldr, 1.0);

	fbColor0Out.rgb = pow(fbColor0Out.rgb, vec3(1.0/2.2));
}