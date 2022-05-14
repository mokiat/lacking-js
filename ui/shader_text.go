package ui

import (
	"github.com/mokiat/lacking-js/internal"
	"github.com/mokiat/lacking/ui"
)

func newTextShaders() ui.ShaderSet {
	return ui.ShaderSet{
		VertexShader: func() string {
			builder := internal.NewShaderSourceBuilder(textMaterialVertexShaderTemplate)
			return builder.Build()
		},
		FragmentShader: func() string {
			builder := internal.NewShaderSourceBuilder(textMaterialFragmentShaderTemplate)
			return builder.Build()
		},
	}
}

const textMaterialVertexShaderTemplate = `
layout(location = 0) in vec2 positionIn;
layout(location = 1) in vec2 texCoordIn;

uniform mat4 projectionMatrixIn;
uniform mat4 transformMatrixIn;
uniform mat4 clipMatrixIn;

smooth out vec4 clipDistancesInOut;
smooth out vec2 texCoordInOut;

void main()
{
	texCoordInOut = texCoordIn;
	vec4 screenPosition = transformMatrixIn * vec4(positionIn, 0.0, 1.0);

	clipDistancesInOut = clipMatrixIn * screenPosition;

	gl_Position = projectionMatrixIn * screenPosition;
}
`

const textMaterialFragmentShaderTemplate = `
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

	float amount = texture(textureIn, texCoordInOut).x;
	fragmentColor = vec4(amount) * colorIn;
}
`
