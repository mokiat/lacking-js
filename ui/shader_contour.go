package ui

import (
	"github.com/mokiat/lacking-js/webgl"
	"github.com/mokiat/lacking/ui"
)

func newContourShaders() ui.ShaderSet {
	return ui.ShaderSet{
		VertexShader: func() string {
			builder := webgl.NewShaderSourceBuilder(contourMaterialVertexShaderTemplate)
			return builder.Build()
		},
		FragmentShader: func() string {
			builder := webgl.NewShaderSourceBuilder(contourMaterialFragmentShaderTemplate)
			return builder.Build()
		},
	}
}

const contourMaterialVertexShaderTemplate = `
layout(location = 0) in vec2 positionIn;
layout(location = 2) in vec4 colorIn;

uniform mat4 transformMatrixIn;
uniform mat4 projectionMatrixIn;
uniform vec4 clipDistancesIn;

// out gl_PerVertex
// {
//   vec4 gl_Position;
//   float gl_ClipDistance[4];
// };

smooth out vec4 colorInOut;

void main()
{
	colorInOut = colorIn;
	vec4 screenPosition = transformMatrixIn * vec4(positionIn, 0.0, 1.0);
	// gl_ClipDistance[0] = screenPosition.x - clipDistancesIn.x; // left
	// gl_ClipDistance[1] = clipDistancesIn.y - screenPosition.x; // right
	// gl_ClipDistance[2] = screenPosition.y - clipDistancesIn.z; // top
	// gl_ClipDistance[3] = clipDistancesIn.w - screenPosition.y; // bottom
	gl_Position = projectionMatrixIn * screenPosition;
}
`

const contourMaterialFragmentShaderTemplate = `
layout(location = 0) out vec4 fragmentColor;

smooth in vec4 colorInOut;

void main()
{
	fragmentColor = colorInOut;
}
`
