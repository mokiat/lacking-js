package internal

import (
	"fmt"
	"unsafe"

	"github.com/mokiat/lacking/render"
	"github.com/mokiat/wasmgl"
)

func NewCommandQueue() *CommandQueue {
	return &CommandQueue{
		data: make([]byte, 1024*1024), // TODO: Start smaller and allow for growth
	}
}

type CommandQueue struct {
	data        []byte
	writeOffset uintptr
	readOffset  uintptr
}

func (q *CommandQueue) Reset() {
	q.readOffset = 0
	q.writeOffset = 0
}

func (q *CommandQueue) BindPipeline(pipeline render.Pipeline) {
	PushCommand(q, CommandHeader{
		Kind: CommandKindBindPipeline,
	})
	intPipeline := pipeline.(*Pipeline)
	PushCommand(q, CommandBindPipeline{
		ProgramID:        intPipeline.ProgramID,
		Topology:         intPipeline.Topology,
		CullTest:         intPipeline.CullTest,
		FrontFace:        intPipeline.FrontFace,
		DepthTest:        intPipeline.DepthTest,
		DepthWrite:       intPipeline.DepthWrite,
		DepthComparison:  intPipeline.DepthComparison,
		StencilTest:      intPipeline.StencilTest,
		StencilOpFront:   intPipeline.StencilOpFront,
		StencilOpBack:    intPipeline.StencilOpBack,
		StencilFuncFront: intPipeline.StencilFuncFront,
		StencilFuncBack:  intPipeline.StencilFuncBack,
		StencilMaskFront: intPipeline.StencilMaskFront,
		StencilMaskBack:  intPipeline.StencilMaskBack,
		ColorWrite:       intPipeline.ColorWrite,
		BlendEnabled:     intPipeline.BlendEnabled,
		BlendColor:       intPipeline.BlendColor,
		BlendEquation:    intPipeline.BlendEquation,
		BlendFunc:        intPipeline.BlendFunc,
		VertexArray:      intPipeline.VertexArray,
	})
}

func (q *CommandQueue) Uniform1f(location render.UniformLocation, value float32) {
	PushCommand(q, CommandHeader{
		Kind: CommandKindUniform1f,
	})
	intLocation := location.(*UniformLocation)
	PushCommand(q, CommandUniform1f{
		Location: intLocation.id,
		Value:    value,
	})
}

func (q *CommandQueue) Uniform1i(location render.UniformLocation, value int) {
	PushCommand(q, CommandHeader{
		Kind: CommandKindUniform1i,
	})
	intLocation := location.(*UniformLocation)
	PushCommand(q, CommandUniform1i{
		Location: intLocation.id,
		Value:    int32(value),
	})
}

func (q *CommandQueue) Uniform3f(location render.UniformLocation, values [3]float32) {
	PushCommand(q, CommandHeader{
		Kind: CommandKindUniform3f,
	})
	intLocation := location.(*UniformLocation)
	PushCommand(q, CommandUniform3f{
		Location: intLocation.id,
		Values:   values,
	})
}

func (q *CommandQueue) Uniform4f(location render.UniformLocation, values [4]float32) {
	PushCommand(q, CommandHeader{
		Kind: CommandKindUniform4f,
	})
	intLocation := location.(*UniformLocation)
	PushCommand(q, CommandUniform4f{
		Location: intLocation.id,
		Values:   values,
	})
}

func (q *CommandQueue) UniformMatrix4f(location render.UniformLocation, values [16]float32) {
	PushCommand(q, CommandHeader{
		Kind: CommandKindUniformMatrix4f,
	})
	intLocation := location.(*UniformLocation)
	PushCommand(q, CommandUniformMatrix4f{
		Location: intLocation.id,
		Values:   values,
	})
}

func (q *CommandQueue) TextureUnit(index int, texture render.Texture) {
	PushCommand(q, CommandHeader{
		Kind: CommandKindTextureUnit,
	})
	PushCommand(q, CommandTextureUnit{
		Index:     uint32(index),
		TextureID: texture.(*Texture).id,
	})
}

func (q *CommandQueue) Draw(vertexOffset, vertexCount, instanceCount int) {
	PushCommand(q, CommandHeader{
		Kind: CommandKindDraw,
	})
	PushCommand(q, CommandDraw{
		VertexOffset:  int32(vertexOffset),
		VertexCount:   int32(vertexCount),
		InstanceCount: int32(instanceCount),
	})
}

func (q *CommandQueue) DrawIndexed(indexOffset, indexCount, instanceCount int) {
	PushCommand(q, CommandHeader{
		Kind: CommandKindDrawIndexed,
	})
	PushCommand(q, CommandDrawIndexed{
		IndexOffset:   int32(indexOffset),
		IndexCount:    int32(indexCount),
		InstanceCount: int32(instanceCount),
	})
}

func (q *CommandQueue) CopyContentToBuffer(info render.CopyContentToBufferInfo) {
	PushCommand(q, CommandHeader{
		Kind: CommandKindCopyContentToBuffer,
	})
	var format, xtype uint32
	switch info.Format {
	case render.DataFormatRGBA8:
		format = wasmgl.RGBA
		xtype = wasmgl.UNSIGNED_BYTE
	case render.DataFormatRGBA16F:
		format = wasmgl.RGBA
		xtype = wasmgl.HALF_FLOAT
	case render.DataFormatRGBA32F:
		format = wasmgl.RGBA
		xtype = wasmgl.FLOAT
	default:
		panic(fmt.Errorf("unsupported data format %v", info.Format))
	}
	PushCommand(q, CommandCopyContentToBuffer{
		BufferID:     info.Buffer.(*Buffer).id,
		X:            int32(info.X),
		Y:            int32(info.Y),
		Width:        int32(info.Width),
		Height:       int32(info.Height),
		Format:       format,
		XType:        xtype,
		BufferOffset: uint32(info.Offset),
	})
}

func (q *CommandQueue) Release() {
	q.data = nil
}

func MoreCommands(queue *CommandQueue) bool {
	return queue.writeOffset > queue.readOffset
}

func PushCommand[T any](queue *CommandQueue, command T) {
	target := (*T)(unsafe.Add(unsafe.Pointer(&queue.data[0]), queue.writeOffset))
	*target = command
	queue.writeOffset += unsafe.Sizeof(command)
}

func PopCommand[T any](queue *CommandQueue) T {
	target := (*T)(unsafe.Add(unsafe.Pointer(&queue.data[0]), queue.readOffset))
	command := *target
	queue.readOffset += unsafe.Sizeof(command)
	return command
}

type CommandKind uint8

const (
	CommandKindBindPipeline CommandKind = iota
	CommandKindTopology
	CommandKindCullTest
	CommandKindFrontFace
	CommandKindDepthTest
	CommandKindDepthWrite
	CommandKindDepthComparison
	CommandKindStencilTest
	CommandKindStencilOperation
	CommandKindStencilFunc
	CommandKindStencilMask
	CommandKindColorWrite
	CommandKindBlendColor
	CommandKindBlendEquation
	CommandKindBlendFunc
	CommandBindKindVertexArray
	CommandKindUniform1f
	CommandKindUniform1i
	CommandKindUniform3f
	CommandKindUniform4f
	CommandKindUniformMatrix4f
	CommandKindTextureUnit
	CommandKindDraw
	CommandKindDrawIndexed
	CommandKindCopyContentToBuffer
)

type CommandHeader struct {
	Kind CommandKind
}

type CommandBindPipeline struct {
	ProgramID        uint32 // not dynamic
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
	BlendEnabled     bool // not dynamic
	BlendEquation    CommandBlendEquation
	BlendFunc        CommandBlendFunc
	BlendColor       CommandBlendColor
	VertexArray      CommandBindVertexArray
}

type CommandTopology struct {
	Topology uint32
}

type CommandCullTest struct {
	Enabled bool
	Face    uint32
}

type CommandFrontFace struct {
	Orientation uint32
}

type CommandDepthTest struct {
	Enabled bool
}

type CommandDepthWrite struct {
	Enabled bool
}

type CommandDepthComparison struct {
	Mode uint32
}

type CommandStencilTest struct {
	Enabled bool
}

type CommandStencilOperation struct {
	Face        uint32
	StencilFail uint32
	DepthFail   uint32
	Pass        uint32
}

type CommandStencilFunc struct {
	Face uint32
	Func uint32
	Ref  uint32
	Mask uint32
}

type CommandStencilMask struct {
	Face uint32
	Mask uint32
}

type CommandColorWrite struct {
	Mask [4]bool
}

type CommandBlendColor struct {
	Color [4]float32
}

type CommandBlendEquation struct {
	ModeRGB   uint32
	ModeAlpha uint32
}

type CommandBlendFunc struct {
	SourceFactorRGB        uint32
	DestinationFactorRGB   uint32
	SourceFactorAlpha      uint32
	DestinationFactorAlpha uint32
}

type CommandBindVertexArray struct {
	VertexArrayID uint32
	IndexFormat   uint32
}

type CommandUniform1f struct {
	Location int32
	Value    float32
}

type CommandUniform1i struct {
	Location int32
	Value    int32
}

type CommandUniform3f struct {
	Location int32
	Values   [3]float32
}

type CommandUniform4f struct {
	Location int32
	Values   [4]float32
}

type CommandUniformMatrix4f struct {
	Location int32
	Values   [16]float32
}

type CommandTextureUnit struct {
	Index     uint32
	TextureID uint32
}

type CommandDraw struct {
	VertexOffset  int32
	VertexCount   int32
	InstanceCount int32
}

type CommandDrawIndexed struct {
	IndexOffset   int32
	IndexCount    int32
	InstanceCount int32
}

type CommandCopyContentToBuffer struct {
	BufferID     uint32
	X            int32
	Y            int32
	Width        int32
	Height       int32
	Format       uint32
	XType        uint32
	BufferOffset uint32
}
