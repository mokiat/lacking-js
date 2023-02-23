layout(location = 0) out vec4 fbColor0Out;

smooth in vec3 colorInOut;

void main()
{
	fbColor0Out = vec4(colorInOut, 1.0);
	gl_FragDepth = gl_FragCoord.z - 0.001;
}
