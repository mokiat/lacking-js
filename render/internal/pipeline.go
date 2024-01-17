package internal

import (
	"fmt"

	"github.com/mokiat/lacking/render"
	"github.com/mokiat/wasmgl"
)

func NewPipeline(info render.PipelineInfo) *Pipeline {
	intProgram := info.Program.(*Program)
	intVertexArray := info.VertexArray.(*VertexArray)

	pipeline := &Pipeline{
		ProgramID: intProgram.id,
		VertexArray: CommandBindVertexArray{
			VertexArrayID: intVertexArray.id,
			IndexFormat:   uint32(intVertexArray.indexFormat),
		},
	}

	switch info.Topology {
	case render.TopologyPoints:
		pipeline.Topology.Topology = wasmgl.POINTS
	case render.TopologyLineStrip:
		pipeline.Topology.Topology = wasmgl.LINE_STRIP
	case render.TopologyLineList:
		pipeline.Topology.Topology = wasmgl.LINES
	case render.TopologyTriangleList:
		pipeline.Topology.Topology = wasmgl.TRIANGLES
	case render.TopologyTriangleStrip:
		pipeline.Topology.Topology = wasmgl.TRIANGLE_STRIP
	case render.TopologyTriangleFan:
		pipeline.Topology.Topology = wasmgl.TRIANGLE_FAN
	}

	switch info.Culling {
	case render.CullModeNone:
		pipeline.CullTest.Enabled = false
	case render.CullModeBack:
		pipeline.CullTest.Enabled = true
		pipeline.CullTest.Face = wasmgl.BACK
	case render.CullModeFront:
		pipeline.CullTest.Enabled = true
		pipeline.CullTest.Face = wasmgl.FRONT
	case render.CullModeFrontAndBack:
		pipeline.CullTest.Enabled = true
		pipeline.CullTest.Face = wasmgl.FRONT_AND_BACK
	}

	switch info.FrontFace {
	case render.FaceOrientationCCW:
		pipeline.FrontFace.Orientation = wasmgl.CCW
	case render.FaceOrientationCW:
		pipeline.FrontFace.Orientation = wasmgl.CW
	}

	pipeline.DepthTest.Enabled = info.DepthTest
	pipeline.DepthWrite.Enabled = info.DepthWrite
	pipeline.DepthComparison.Mode = glEnumFromComparison(info.DepthComparison)

	pipeline.StencilTest.Enabled = info.StencilTest

	pipeline.StencilOpFront.Face = wasmgl.FRONT
	pipeline.StencilOpFront.StencilFail = glEnumFromStencilOp(info.StencilFront.StencilFailOp)
	pipeline.StencilOpFront.DepthFail = glEnumFromStencilOp(info.StencilFront.DepthFailOp)
	pipeline.StencilOpFront.Pass = glEnumFromStencilOp(info.StencilFront.PassOp)

	pipeline.StencilOpBack.Face = wasmgl.BACK
	pipeline.StencilOpBack.StencilFail = glEnumFromStencilOp(info.StencilBack.StencilFailOp)
	pipeline.StencilOpBack.DepthFail = glEnumFromStencilOp(info.StencilBack.DepthFailOp)
	pipeline.StencilOpBack.Pass = glEnumFromStencilOp(info.StencilBack.PassOp)

	pipeline.StencilFuncFront.Face = wasmgl.FRONT
	pipeline.StencilFuncFront.Func = glEnumFromComparison(info.StencilFront.Comparison)
	pipeline.StencilFuncFront.Ref = info.StencilFront.Reference
	pipeline.StencilFuncFront.Mask = info.StencilFront.ComparisonMask

	pipeline.StencilFuncBack.Face = wasmgl.BACK
	pipeline.StencilFuncBack.Func = glEnumFromComparison(info.StencilBack.Comparison)
	pipeline.StencilFuncBack.Ref = info.StencilBack.Reference
	pipeline.StencilFuncBack.Mask = info.StencilBack.ComparisonMask

	pipeline.StencilMaskFront.Face = wasmgl.FRONT
	pipeline.StencilMaskFront.Mask = info.StencilFront.WriteMask

	pipeline.StencilMaskBack.Face = wasmgl.BACK
	pipeline.StencilMaskBack.Mask = info.StencilBack.WriteMask

	pipeline.ColorWrite.Mask = info.ColorWrite

	pipeline.BlendEnabled = info.BlendEnabled
	pipeline.BlendColor.Color = info.BlendColor

	pipeline.BlendEquation.ModeRGB = glEnumFromBlendOp(info.BlendOpColor)
	pipeline.BlendEquation.ModeAlpha = glEnumFromBlendOp(info.BlendOpAlpha)

	pipeline.BlendFunc.SourceFactorRGB = glEnumFromBlendFactor(info.BlendSourceColorFactor)
	pipeline.BlendFunc.DestinationFactorRGB = glEnumFromBlendFactor(info.BlendDestinationColorFactor)
	pipeline.BlendFunc.SourceFactorAlpha = glEnumFromBlendFactor(info.BlendSourceAlphaFactor)
	pipeline.BlendFunc.DestinationFactorAlpha = glEnumFromBlendFactor(info.BlendDestinationAlphaFactor)

	return pipeline
}

type Pipeline struct {
	render.PipelineMarker
	ProgramID        uint32
	Topology         CommandTopology
	CullTest         CommandCullTest
	FrontFace        CommandFrontFace
	DepthTest        CommandDepthTest
	DepthWrite       CommandDepthWrite
	DepthComparison  CommandDepthComparison
	StencilTest      CommandStencilTest
	StencilOpFront   CommandStencilOperation
	StencilOpBack    CommandStencilOperation
	StencilFuncFront CommandStencilFunc
	StencilFuncBack  CommandStencilFunc
	StencilMaskFront CommandStencilMask
	StencilMaskBack  CommandStencilMask
	ColorWrite       CommandColorWrite
	BlendEnabled     bool
	BlendColor       CommandBlendColor
	BlendEquation    CommandBlendEquation
	BlendFunc        CommandBlendFunc
	VertexArray      CommandBindVertexArray
}

func (p *Pipeline) Release() {
}

func glEnumFromComparison(comparison render.Comparison) uint32 {
	switch comparison {
	case render.ComparisonNever:
		return wasmgl.NEVER
	case render.ComparisonLess:
		return wasmgl.LESS
	case render.ComparisonEqual:
		return wasmgl.EQUAL
	case render.ComparisonLessOrEqual:
		return wasmgl.LEQUAL
	case render.ComparisonGreater:
		return wasmgl.GREATER
	case render.ComparisonNotEqual:
		return wasmgl.NOTEQUAL
	case render.ComparisonGreaterOrEqual:
		return wasmgl.GEQUAL
	case render.ComparisonAlways:
		return wasmgl.ALWAYS
	default:
		panic(fmt.Errorf("unknown comparison: %d", comparison))
	}
}

func glEnumFromStencilOp(op render.StencilOperation) uint32 {
	switch op {
	case render.StencilOperationKeep:
		return wasmgl.KEEP
	case render.StencilOperationZero:
		return wasmgl.ZERO
	case render.StencilOperationReplace:
		return wasmgl.REPLACE
	case render.StencilOperationIncrease:
		return wasmgl.INCR
	case render.StencilOperationIncreaseWrap:
		return wasmgl.INCR_WRAP
	case render.StencilOperationDecrease:
		return wasmgl.DECR
	case render.StencilOperationDecreaseWrap:
		return wasmgl.DECR_WRAP
	case render.StencilOperationInvert:
		return wasmgl.INVERT
	default:
		panic(fmt.Errorf("unknown op: %d", op))
	}
}

func glEnumFromBlendOp(op render.BlendOperation) uint32 {
	switch op {
	case render.BlendOperationAdd:
		return wasmgl.FUNC_ADD
	case render.BlendOperationSubtract:
		return wasmgl.FUNC_SUBTRACT
	case render.BlendOperationReverseSubtract:
		return wasmgl.FUNC_REVERSE_SUBTRACT
	case render.BlendOperationMin:
		return wasmgl.MIN
	case render.BlendOperationMax:
		return wasmgl.MAX
	default:
		panic(fmt.Errorf("unknown op: %d", op))
	}
}

func glEnumFromBlendFactor(factor render.BlendFactor) uint32 {
	switch factor {
	case render.BlendFactorZero:
		return wasmgl.ZERO
	case render.BlendFactorOne:
		return wasmgl.ONE
	case render.BlendFactorSourceColor:
		return wasmgl.SRC_COLOR
	case render.BlendFactorOneMinusSourceColor:
		return wasmgl.ONE_MINUS_SRC_COLOR
	case render.BlendFactorDestinationColor:
		return wasmgl.DST_COLOR
	case render.BlendFactorOneMinusDestinationColor:
		return wasmgl.ONE_MINUS_DST_COLOR
	case render.BlendFactorSourceAlpha:
		return wasmgl.SRC_ALPHA
	case render.BlendFactorOneMinusSourceAlpha:
		return wasmgl.ONE_MINUS_SRC_ALPHA
	case render.BlendFactorDestinationAlpha:
		return wasmgl.DST_ALPHA
	case render.BlendFactorOneMinusDestinationAlpha:
		return wasmgl.ONE_MINUS_DST_ALPHA
	case render.BlendFactorConstantColor:
		return wasmgl.CONSTANT_COLOR
	case render.BlendFactorOneMinusConstantColor:
		return wasmgl.ONE_MINUS_CONSTANT_COLOR
	case render.BlendFactorConstantAlpha:
		return wasmgl.CONSTANT_ALPHA
	case render.BlendFactorOneMinusConstantAlpha:
		return wasmgl.ONE_MINUS_CONSTANT_ALPHA
	case render.BlendFactorSourceAlphaSaturate:
		return wasmgl.SRC_ALPHA_SATURATE
	default:
		panic(fmt.Errorf("unknown factor: %d", factor))
	}
}
