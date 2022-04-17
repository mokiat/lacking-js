package game

import (
	"fmt"

	"github.com/mokiat/lacking-js/webgl"
	"github.com/mokiat/lacking/game/graphics"
)

func newPostprocessingShaderSet(mapping graphics.ToneMapping) graphics.ShaderSet {
	vsBuilder := webgl.NewShaderSourceBuilder(tonePostprocessingVertexShader)
	fsBuilder := webgl.NewShaderSourceBuilder(tonePostprocessingFragmentShader)
	switch mapping {
	case graphics.ReinhardToneMapping:
		fsBuilder.AddFeature("MODE_REINHARD")
	case graphics.ExponentialToneMapping:
		fsBuilder.AddFeature("MODE_EXPONENTIAL")
	default:
		panic(fmt.Errorf("unknown tone mapping mode: %s", mapping))
	}
	return graphics.ShaderSet{
		VertexShader:   vsBuilder.Build,
		FragmentShader: fsBuilder.Build,
	}
}

const tonePostprocessingVertexShader = `
layout(location = 0) in vec2 coordIn;

smooth out vec2 texCoordInOut;

void main()
{
	texCoordInOut = (coordIn + 1.0) / 2.0;
	gl_Position = vec4(coordIn, 0.0, 1.0);
}
`

const tonePostprocessingFragmentShader = `
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
`
