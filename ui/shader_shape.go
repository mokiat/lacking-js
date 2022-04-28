package ui

import (
	"github.com/mokiat/lacking-js/internal"
	"github.com/mokiat/lacking/ui"
)

func newShapeShaders() ui.ShaderSet {
	return ui.ShaderSet{
		VertexShader: func() string {
			builder := internal.NewShaderSourceBuilder(shapeMaterialVertexShaderTemplate)
			return builder.Build()
		},
		FragmentShader: func() string {
			builder := internal.NewShaderSourceBuilder(shapeMaterialFragmentShaderTemplate)
			return builder.Build()
		},
	}
}

func newShapeBlankShaders() ui.ShaderSet {
	return ui.ShaderSet{
		VertexShader: func() string {
			builder := internal.NewShaderSourceBuilder(shapeBlankMaterialVertexShaderTemplate)
			return builder.Build()
		},
		FragmentShader: func() string {
			builder := internal.NewShaderSourceBuilder(shapeBlankMaterialFragmentShaderTemplate)
			return builder.Build()
		},
	}
}

const shapeMaterialVertexShaderTemplate = `
layout(location = 0) in vec2 positionIn;

uniform mat4 transformMatrixIn;
uniform mat4 textureTransformMatrixIn;
uniform mat4 projectionMatrixIn;
// uniform vec4 clipDistancesIn;

smooth out vec2 texCoordInOut;

// out gl_PerVertex
// {
//   vec4 gl_Position;
//   float gl_ClipDistance[4];
// };

void main()
{
	vec4 screenPosition = transformMatrixIn * vec4(positionIn, 0.0, 1.0);
	texCoordInOut = (textureTransformMatrixIn * vec4(positionIn, 0.0, 1.0)).xy;
	// gl_ClipDistance[0] = screenPosition.x - clipDistancesIn.x; // left
	// gl_ClipDistance[1] = clipDistancesIn.y - screenPosition.x; // right
	// gl_ClipDistance[2] = screenPosition.y - clipDistancesIn.z; // top
	// gl_ClipDistance[3] = clipDistancesIn.w - screenPosition.y; // bottom
	gl_Position = projectionMatrixIn * screenPosition;
}
`

const shapeMaterialFragmentShaderTemplate = `
layout(location = 0) out vec4 fragmentColor;

uniform sampler2D textureIn;
uniform vec4 colorIn;

smooth in vec2 texCoordInOut;

void main()
{
	fragmentColor = texture(textureIn, texCoordInOut) * colorIn;
}
`

const shapeBlankMaterialVertexShaderTemplate = `
layout(location = 0) in vec2 positionIn;

uniform mat4 transformMatrixIn;
uniform mat4 projectionMatrixIn;
// uniform vec4 clipDistancesIn;

// out gl_PerVertex
// {
//   vec4 gl_Position;
//   float gl_ClipDistance[4];
// };

void main()
{
	vec4 screenPosition = transformMatrixIn * vec4(positionIn, 0.0, 1.0);
	// gl_ClipDistance[0] = screenPosition.x - clipDistancesIn.x; // left
	// gl_ClipDistance[1] = clipDistancesIn.y - screenPosition.x; // right
	// gl_ClipDistance[2] = screenPosition.y - clipDistancesIn.z; // top
	// gl_ClipDistance[3] = clipDistancesIn.w - screenPosition.y; // bottom
	gl_Position = projectionMatrixIn * screenPosition;
}
`

const shapeBlankMaterialFragmentShaderTemplate = `
layout(location = 0) out vec4 fragmentColor;

void main()
{
	fragmentColor = vec4(1.0, 1.0, 1.0, 1.0);
}
`
