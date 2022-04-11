package game

import (
	"github.com/mokiat/lacking-js/webgl"
	"github.com/mokiat/lacking/game/graphics/renderapi/plugin"
)

func newSkycolorShaderSet() plugin.ShaderSet {
	vsBuilder := webgl.NewShaderSourceBuilder(colorSkyboxVertexShader)
	fsBuilder := webgl.NewShaderSourceBuilder(colorSkyboxFragmentShader)
	return plugin.ShaderSet{
		VertexShader:   vsBuilder.Build,
		FragmentShader: fsBuilder.Build,
	}
}

const colorSkyboxVertexShader = `
layout(location = 0) in vec3 coordIn;

uniform mat4 projectionMatrixIn;
uniform mat4 viewMatrixIn;

void main()
{
	// ensure that translations are ignored by setting w to 0.0
	vec4 viewPosition = viewMatrixIn * vec4(coordIn, 0.0);

	// restore w to 1.0 so that projection works
	vec4 position = projectionMatrixIn * vec4(viewPosition.xyz, 1.0);

	// set z to w so that it has maximum depth (1.0) after projection division
	gl_Position = vec4(position.xy, position.w, position.w);
}`

const colorSkyboxFragmentShader = `
layout(location = 0) out vec4 fbColor0Out;

uniform vec4 albedoColorIn;

void main()
{
	fbColor0Out = albedoColorIn;
}`
