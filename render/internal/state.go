package internal

type State struct {
	CullTest        bool
	CullFace        int
	FrontFace       int
	DepthTest       bool
	DepthMask       bool
	DepthComparison int
	StencilTest     bool
}
