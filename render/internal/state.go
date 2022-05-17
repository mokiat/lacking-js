package internal

type State struct {
	CullTest                    bool
	CullFace                    int
	FrontFace                   int
	DepthTest                   bool
	DepthMask                   bool
	DepthComparison             int
	StencilTest                 bool
	StencilOpStencilFailFront   int
	StencilOpDepthFailFront     int
	StencilOpPassFront          int
	StencilOpStencilFailBack    int
	StencilOpDepthFailBack      int
	StencilOpPassBack           int
	StencilComparisonFuncFront  int
	StencilComparisonRefFront   int
	StencilComparisonMaskFront  int
	StencilComparisonFuncBack   int
	StencilComparisonRefBack    int
	StencilComparisonMaskBack   int
	StencilMaskFront            int
	StencilMaskBack             int
	ColorMask                   [4]bool
	Blending                    bool
	BlendColor                  [4]float32
	BlendModeRGB                int
	BlendModeAlpha              int
	BlendSourceFactorRGB        int
	BlendDestinationFactorRGB   int
	BlendSourceFactorAlpha      int
	BlendDestinationFactorAlpha int
}
