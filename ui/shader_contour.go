package ui

import (
	"github.com/mokiat/lacking-js/internal"
	"github.com/mokiat/lacking/ui"
)

func newContourShaders() ui.ShaderSet {
	return ui.ShaderSet{
		VertexShader: func() string {
			builder := internal.NewShaderSourceBuilder(contourMaterialVertexShaderTemplate)
			return builder.Build()
		},
		FragmentShader: func() string {
			builder := internal.NewShaderSourceBuilder(contourMaterialFragmentShaderTemplate)
			return builder.Build()
		},
	}
}

const contourMaterialVertexShaderTemplate = `
layout(location = 0) in vec2 positionIn;
layout(location = 2) in vec4 colorIn;

uniform mat4 projectionMatrixIn;
uniform mat4 transformMatrixIn;
uniform mat4 clipMatrixIn;

smooth out vec4 clipDistancesInOut;
smooth out vec4 colorInOut;

void main()
{
	colorInOut = colorIn;
	vec4 screenPosition = transformMatrixIn * vec4(positionIn, 0.0, 1.0);

	clipDistancesInOut = clipMatrixIn * screenPosition;

	gl_Position = projectionMatrixIn * screenPosition;
}
`

const contourMaterialFragmentShaderTemplate = `
layout(location = 0) out vec4 fragmentColor;

smooth in vec4 clipDistancesInOut;
smooth in vec4 colorInOut;

void main()
{
	float dist = min(min(clipDistancesInOut.x, clipDistancesInOut.y), min(clipDistancesInOut.z, clipDistancesInOut.w));
	if (dist < 0.0) {
		discard;
	}

	fragmentColor = colorInOut;
}
`
