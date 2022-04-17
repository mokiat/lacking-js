package webgl

import (
	"fmt"

	"github.com/mokiat/wasmgl"
)

func NewShader() *Shader {
	return &Shader{}
}

type Shader struct {
	raw wasmgl.Shader
}

func (s *Shader) Raw() wasmgl.Shader {
	return s.raw
}

func (s *Shader) Allocate(info ShaderAllocateInfo) {
	s.raw = wasmgl.CreateShader(info.ShaderType)
	wasmgl.ShaderSource(s.raw, info.SourceCode)
	wasmgl.CompileShader(s.raw)
	if !s.isCompileSuccessful() {
		panic(fmt.Errorf("failed to compile shader: %s", s.getInfoLog()))
	}
}

func (s *Shader) Release() {
	wasmgl.DeleteShader(s.raw)
	s.raw = wasmgl.Shader{}
}

func (s *Shader) isCompileSuccessful() bool {
	result := wasmgl.GetShaderParameter(s.raw, wasmgl.COMPILE_STATUS)
	return result.Bool()
}

func (s *Shader) getInfoLog() string {
	return wasmgl.GetShaderInfoLog(s.raw)
}

type ShaderAllocateInfo struct {
	ShaderType int
	SourceCode string
}
