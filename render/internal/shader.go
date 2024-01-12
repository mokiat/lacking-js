package internal

import (
	"errors"

	"github.com/mokiat/wasmgl"
)

func newVertexShader(programLabel, sourceCode string) *Shader {
	shader := &Shader{
		raw: wasmgl.CreateShader(wasmgl.VERTEX_SHADER),
	}
	shader.setSourceCode(sourceCode)
	if err := shader.compile(); err != nil {
		logger.Error("Vertex Shader (for program %q) compilation error: %v!", programLabel, err)
	}
	return shader
}

func newFragmentShader(programLabel, sourceCode string) *Shader {
	shader := &Shader{
		raw: wasmgl.CreateShader(wasmgl.FRAGMENT_SHADER),
	}
	shader.setSourceCode(sourceCode)
	if err := shader.compile(); err != nil {
		logger.Error("Fragment Shader (for program %q) compilation error: %v!", programLabel, err)
	}
	return shader
}

type Shader struct {
	raw wasmgl.Shader
}

func (s *Shader) Release() {
	wasmgl.DeleteShader(s.raw)
	s.raw = wasmgl.NilShader
}

func (s *Shader) setSourceCode(code string) {
	wasmgl.ShaderSource(s.raw, code)
}

func (s *Shader) compile() error {
	wasmgl.CompileShader(s.raw)
	if !s.isCompileSuccessful() {
		return errors.New(s.getInfoLog())
	}
	return nil
}

func (s *Shader) isCompileSuccessful() bool {
	result := wasmgl.GetShaderParameter(s.raw, wasmgl.COMPILE_STATUS)
	return result.GLboolean()
}

func (s *Shader) getInfoLog() string {
	return wasmgl.GetShaderInfoLog(s.raw)
}
