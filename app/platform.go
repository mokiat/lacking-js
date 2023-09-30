package app

import (
	"strings"
	"syscall/js"

	"github.com/mokiat/lacking/app"
)

func newPlatform() *platform {
	userAgent := js.Global().Get("navigator").Get("platform")
	return &platform{
		os: determineOS(userAgent.String()),
	}
}

type platform struct {
	os app.OS
}

var _ app.Platform = (*platform)(nil)

func (p *platform) OS() app.OS {
	return p.os
}

func determineOS(userAgent string) app.OS {
	userAgent = strings.ToLower(userAgent)
	switch {
	case strings.Contains(userAgent, "linux"):
		return app.OSLinux
	case strings.Contains(userAgent, "windows"):
		return app.OSWindows
	case strings.Contains(userAgent, "darwin") || strings.Contains(userAgent, "mac"):
		return app.OSDarwin
	default:
		return app.OSUnknown
	}
}
