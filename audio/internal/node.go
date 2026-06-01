package internal

import (
	"github.com/mokiat/lacking/audio"
	"github.com/mokiat/wasmal"
)

type Node interface {
	audio.Node
	AudioNode() wasmal.AudioNode
}
