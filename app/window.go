//go:build js && wasm

package app

import (
	"fmt"
	"syscall/js"

	"github.com/mokiat/lacking/app"
	"github.com/mokiat/wasmgl"
)

// Run starts a new application by attaching to the specified in the config
// HTML canvas element. The configuration is used to determine how the
// canvas is further initialized.
//
// The specified controller will be used to send notifications
// on window state changes.
func Run(cfg *Config, controller app.Controller) error {
	htmlDocument := js.Global().Get("document")
	if htmlDocument.IsUndefined() {
		return fmt.Errorf("could not locate document element")
	}

	htmlCanvas := htmlDocument.Call("getElementById", cfg.canvasID)
	if htmlCanvas.IsNull() {
		return fmt.Errorf("could not locate canvas element")
	}

	if cfg.title != nil {
		htmlDocument.Set("title", *cfg.title)
	}

	if cfg.fullscreen {
		htmlBody := htmlDocument.Get("body")
		bodyWidth := htmlBody.Get("clientWidth").Int()
		bodyHeight := htmlBody.Get("clientHeight").Int()
		htmlCanvas.Set("width", bodyWidth)
		htmlCanvas.Set("height", bodyHeight)
	} else {
		if cfg.width != nil {
			htmlCanvas.Set("width", *cfg.width)
		}
		if cfg.height != nil {
			htmlCanvas.Set("height", *cfg.height)
		}
	}

	// TODO: Make graphics library configurable
	err := wasmgl.InitFromCanvas(htmlCanvas,
		wasmgl.WithOptionPowerPreference(wasmgl.PowerPreferenceHighPerformance),
	)
	if err != nil {
		return fmt.Errorf("error initializing webgl: %w", err)
	}
	for _, ext := range cfg.glExtensions {
		if wasmgl.GetExtension(ext) == nil {
			appLogger.Warn("[app] Extension %q might not be supported", ext)
		}
	}

	l := newLoop(htmlDocument, htmlCanvas, controller, cfg.audioEnabled)
	if cfg.cursor != nil {
		cursor := l.CreateCursor(*cfg.cursor)
		defer cursor.Destroy()
		l.UseCursor(cursor)
		defer l.UseCursor(nil)
	}
	if !cfg.cursorVisible {
		l.SetCursorVisible(false)
	}
	return l.Run()
}
