//go:build js && wasm

package app

import "github.com/mokiat/lacking/app"

var (
	mouseButtonMapping map[int]app.MouseButton
)

func init() {
	mouseButtonMapping = make(map[int]app.MouseButton)
	mouseButtonMapping[0] = app.MouseButtonLeft
	mouseButtonMapping[1] = app.MouseButtonMiddle
	mouseButtonMapping[2] = app.MouseButtonRight
}
