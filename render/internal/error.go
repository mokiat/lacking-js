package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mokiat/wasmgl"
)

var isDebugEnabled = logger.Enabled(context.Background(), slog.LevelDebug)

func trackError(msg, label string) func() {
	if !isDebugEnabled {
		return nopFunc
	}
	clearErrors()
	return func() {
		if err := getError(); err != "" {
			logger.Error(msg,
				slog.String("label", label),
				slog.String("error", err),
			)
		}
	}
}

func nopFunc() {}

func clearErrors() {
	for wasmgl.GetError() != wasmgl.NO_ERROR {
	}
}

func getError() string {
	switch code := wasmgl.GetError(); code {
	case wasmgl.NO_ERROR:
		return ""
	case wasmgl.INVALID_ENUM:
		return "INVALID_ENUM"
	case wasmgl.INVALID_VALUE:
		return "INVALID_VALUE"
	case wasmgl.INVALID_OPERATION:
		return "INVALID_OPERATION"
	case wasmgl.INVALID_FRAMEBUFFER_OPERATION:
		return "INVALID_FRAMEBUFFER_OPERATION"
	case wasmgl.OUT_OF_MEMORY:
		return "OUT_OF_MEMORY"
	default:
		return fmt.Sprintf("UNKNOWN_ERROR(%x)", code)
	}
}
