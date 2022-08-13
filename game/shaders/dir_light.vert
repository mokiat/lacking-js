layout(location = 0) in vec3 coordIn;

smooth out vec2 texCoordInOut;

void main()
{
	texCoordInOut = (coordIn.xy + 1.0) / 2.0;
	gl_Position = vec4(coordIn.xy, 0.0, 1.0);
}
