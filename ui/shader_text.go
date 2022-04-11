package ui

import (
	"github.com/mokiat/lacking-js/webgl"
	"github.com/mokiat/lacking/ui/renderapi/plugin"
)

func newTextShaders() plugin.ShaderSet {
	return plugin.ShaderSet{
		VertexShader: func() string {
			builder := webgl.NewShaderSourceBuilder(textMaterialVertexShaderTemplate)
			return builder.Build()
		},
		FragmentShader: func() string {
			builder := webgl.NewShaderSourceBuilder(textMaterialFragmentShaderTemplate)
			return builder.Build()
		},
	}
}

const textMaterialVertexShaderTemplate = `
layout(location = 0) in vec2 positionIn;
layout(location = 1) in vec2 texCoordIn;

uniform mat4 transformMatrixIn;
uniform mat4 projectionMatrixIn;
uniform vec4 clipDistancesIn;

// out gl_PerVertex
// {
//   vec4 gl_Position;
//   float gl_ClipDistance[4];
// };

smooth out vec2 texCoordInOut;

void main()
{
	texCoordInOut = texCoordIn;
	vec4 screenPosition = transformMatrixIn * vec4(positionIn, 0.0, 1.0);
	// gl_ClipDistance[0] = screenPosition.x - clipDistancesIn.x; // left
	// gl_ClipDistance[1] = clipDistancesIn.y - screenPosition.x; // right
	// gl_ClipDistance[2] = screenPosition.y - clipDistancesIn.z; // top
	// gl_ClipDistance[3] = clipDistancesIn.w - screenPosition.y; // bottom
	gl_Position = projectionMatrixIn * screenPosition;
}
`

const textMaterialFragmentShaderTemplate = `
layout(location = 0) out vec4 fragmentColor;

uniform sampler2D textureIn;
uniform vec4 colorIn;

smooth in vec2 texCoordInOut;

void main()
{
	float amount = pow(clamp(texture(textureIn, texCoordInOut).x, 0.0, 1.0), 0.5);
	if (amount > 0.8) {
		amount = 1.0;
	}
	if (amount < 0.2) {
		amount = 0.0;
	}
	fragmentColor = vec4(amount) * colorIn;
}
`