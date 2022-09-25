package internal

import "github.com/mokiat/wasmgl"

type State struct {
	CullTest                    bool
	CullFace                    wasmgl.GLenum
	FrontFace                   wasmgl.GLenum
	DepthTest                   bool
	DepthMask                   bool
	DepthComparison             wasmgl.GLenum
	StencilTest                 bool
	StencilOpStencilFailFront   wasmgl.GLenum
	StencilOpDepthFailFront     wasmgl.GLenum
	StencilOpPassFront          wasmgl.GLenum
	StencilOpStencilFailBack    wasmgl.GLenum
	StencilOpDepthFailBack      wasmgl.GLenum
	StencilOpPassBack           wasmgl.GLenum
	StencilComparisonFuncFront  wasmgl.GLenum
	StencilComparisonRefFront   wasmgl.GLint
	StencilComparisonMaskFront  wasmgl.GLuint
	StencilComparisonFuncBack   wasmgl.GLenum
	StencilComparisonRefBack    wasmgl.GLint
	StencilComparisonMaskBack   wasmgl.GLuint
	StencilMaskFront            wasmgl.GLuint
	StencilMaskBack             wasmgl.GLuint
	ColorMask                   [4]bool
	Blending                    bool
	BlendColor                  [4]float32
	BlendModeRGB                wasmgl.GLenum
	BlendModeAlpha              wasmgl.GLenum
	BlendSourceFactorRGB        wasmgl.GLenum
	BlendDestinationFactorRGB   wasmgl.GLenum
	BlendSourceFactorAlpha      wasmgl.GLenum
	BlendDestinationFactorAlpha wasmgl.GLenum
}
