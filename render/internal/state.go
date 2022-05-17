package internal

type State struct {
	CullTest                   bool
	CullFace                   int
	FrontFace                  int
	DepthTest                  bool
	DepthMask                  bool
	DepthComparison            int
	StencilTest                bool
	StencilOpStencilFailFront  int
	StencilOpDepthFailFront    int
	StencilOpPassFront         int
	StencilOpStencilFailBack   int
	StencilOpDepthFailBack     int
	StencilOpPassBack          int
	StencilComparisonFuncFront int
	StencilComparisonRefFront  int
	StencilComparisonMaskFront int
	StencilComparisonFuncBack  int
	StencilComparisonRefBack   int
	StencilComparisonMaskBack  int
}
