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

uniform mat4 projectionMatrixIn;
uniform mat4 transformMatrixIn;
uniform mat4 clipMatrixIn;
uniform mat4 textureTransformMatrixIn;

smooth out vec4 clipDistancesInOut;
smooth out vec2 texCoordInOut;

void main()
{
	vec4 screenPosition = transformMatrixIn * vec4(positionIn, 0.0, 1.0);
	texCoordInOut = (textureTransformMatrixIn * vec4(positionIn, 0.0, 1.0)).xy;

	clipDistancesInOut = clipMatrixIn * screenPosition;

	gl_Position = projectionMatrixIn * screenPosition;
}
`

const shapeMaterialFragmentShaderTemplate = `
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

	fragmentColor = texture(textureIn, texCoordInOut) * colorIn;
}
`

const shapeBlankMaterialVertexShaderTemplate = `
layout(location = 0) in vec2 positionIn;

uniform mat4 projectionMatrixIn;
uniform mat4 transformMatrixIn;
uniform mat4 clipMatrixIn;

smooth out vec4 clipDistancesInOut;

void main()
{
	vec4 screenPosition = transformMatrixIn * vec4(positionIn, 0.0, 1.0);

	clipDistancesInOut = clipMatrixIn * screenPosition;

	gl_Position = projectionMatrixIn * screenPosition;
}
`

const shapeBlankMaterialFragmentShaderTemplate = `
layout(location = 0) out vec4 fragmentColor;

smooth in vec4 clipDistancesInOut;

void main()
{
	float dist = min(min(clipDistancesInOut.x, clipDistancesInOut.y), min(clipDistancesInOut.z, clipDistancesInOut.w));
	if (dist < 0.0) {
		discard;
	}

	fragmentColor = vec4(1.0, 1.0, 1.0, 1.0);
}
`
