package internal

import (
	"time"

	"github.com/mokiat/wasmal"
)

type Media struct {
	buffer wasmal.AudioBuffer
}

func (m *Media) Length() time.Duration {
	return time.Duration(m.buffer.Duration() * float64(time.Second))
}

func (m *Media) Delete() {
}
