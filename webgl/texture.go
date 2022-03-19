//go:build js && wasm

package webgl

import "github.com/mokiat/wasmgl"

type Texture struct {
	raw wasmgl.Texture
}

func (t *Texture) Raw() wasmgl.Texture {
	return t.raw
}
