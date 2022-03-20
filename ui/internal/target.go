package internal

import "github.com/mokiat/lacking-js/webgl"

type Target struct {
	Framebuffer *webgl.Framebuffer
	Width       int
	Height      int
}
