package internal

import "github.com/mokiat/wasmal"

type Node interface {
	Input() wasmal.AudioNode
	Output() wasmal.AudioNode
}
