//go:build js && wasm

package app

import (
	"fmt"
	"strings"
	"syscall/js"
	"time"

	"github.com/mokiat/lacking/app"
)

const (
	taskQueueSize         = 1024
	taskProcessingTimeout = 30 * time.Millisecond
)

func newLoop(htmlDocument, htmlCanvas js.Value, controller app.Controller) *loop {
	return &loop{
		htmlDocument: htmlDocument,
		htmlCanvas:   htmlCanvas,
		controller:   controller,
		tasks:        make(chan func() error, taskQueueSize),
		gamepads:     make(map[int]struct{}),
		shouldStop:   false,
	}
}

var _ app.Window = (*loop)(nil)

type loop struct {
	htmlDocument js.Value
	htmlCanvas   js.Value
	controller   app.Controller
	tasks        chan func() error
	gamepads     map[int]struct{}
	shouldStop   bool
}

func (l *loop) Run() error {
	l.controller.OnCreate(l)
	defer l.controller.OnDestroy(l)

	keydownCallback := js.FuncOf(l.onJSKeyDown)
	defer keydownCallback.Release()
	l.htmlCanvas.Call("addEventListener", "keydown", keydownCallback)
	defer l.htmlCanvas.Call("removeEventListener", "keydown", keydownCallback)

	keyupCallback := js.FuncOf(l.onJSKeyUp)
	defer keyupCallback.Release()
	l.htmlCanvas.Call("addEventListener", "keyup", keyupCallback)
	defer l.htmlCanvas.Call("removeEventListener", "keyup", keyupCallback)

	mouseEnterCallback := js.FuncOf(l.onJSMouseEnter)
	defer mouseEnterCallback.Release()
	l.htmlCanvas.Call("addEventListener", "mouseenter", mouseEnterCallback)
	defer l.htmlCanvas.Call("removeEventListener", "mouseenter", mouseEnterCallback)

	mouseLeaveCallback := js.FuncOf(l.onJSMouseLeave)
	defer mouseLeaveCallback.Release()
	l.htmlCanvas.Call("addEventListener", "mouseleave", mouseLeaveCallback)
	defer l.htmlCanvas.Call("removeEventListener", "mouseleave", mouseLeaveCallback)

	mouseMoveCallback := js.FuncOf(l.onJSMouseMove)
	defer mouseMoveCallback.Release()
	l.htmlCanvas.Call("addEventListener", "mousemove", mouseMoveCallback)
	defer l.htmlCanvas.Call("removeEventListener", "mousemove", mouseMoveCallback)

	mouseDownCallback := js.FuncOf(l.onJSMouseDown)
	defer mouseDownCallback.Release()
	l.htmlCanvas.Call("addEventListener", "mousedown", mouseDownCallback)
	defer l.htmlCanvas.Call("removeEventListener", "mousedown", mouseDownCallback)

	mouseUpCallback := js.FuncOf(l.onJSMouseUp)
	defer mouseUpCallback.Release()
	l.htmlCanvas.Call("addEventListener", "mouseup", mouseUpCallback)
	defer l.htmlCanvas.Call("removeEventListener", "mouseup", mouseUpCallback)

	mouseScrollCallback := js.FuncOf(l.onJSMouseWheel)
	defer mouseScrollCallback.Release()
	l.htmlCanvas.Call("addEventListener", "wheel", mouseScrollCallback)
	defer l.htmlCanvas.Call("removeEventListener", "wheel", mouseScrollCallback)

	gamepadConnectedCallback := js.FuncOf(l.onJSGamepadConnected)
	defer gamepadConnectedCallback.Release()
	js.Global().Call("addEventListener", "gamepadconnected", gamepadConnectedCallback)
	defer js.Global().Call("removeEventListener", "gamepadconnected", gamepadConnectedCallback)

	gamepadDisconnectedCallback := js.FuncOf(l.onJSGamepadDisconnected)
	defer gamepadDisconnectedCallback.Release()
	js.Global().Call("addEventListener", "gamepaddisconnected", gamepadDisconnectedCallback)
	defer js.Global().Call("removeEventListener", "gamepaddisconnected", gamepadDisconnectedCallback)

	w, h := l.Size()
	l.controller.OnResize(l, w, h)
	l.controller.OnFramebufferResize(l, w, h)

	done := make(chan error, 1)
	var loopFunc js.Func
	loopFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if l.shouldStop {
			if l.processTasks(5 * time.Second) {
				done <- nil
			} else {
				done <- fmt.Errorf("failed to cleanup within timeout")
			}
			return true
		}

		l.processTasks(taskProcessingTimeout)
		l.controller.OnRender(l)

		js.Global().Call("requestAnimationFrame", loopFunc)
		return true
	})
	js.Global().Call("requestAnimationFrame", loopFunc)
	defer loopFunc.Release()
	return <-done
}

func (l *loop) Title() string {
	return l.htmlDocument.Get("title").String()
}

func (l *loop) SetTitle(title string) {
	l.htmlDocument.Set("title", title)
}

func (l *loop) SetSize(width, height int) {
	l.htmlCanvas.Set("width", width)
	l.htmlCanvas.Set("height", height)
}

func (l *loop) Size() (int, int) {
	width := l.htmlCanvas.Get("width").Int()
	height := l.htmlCanvas.Get("height").Int()
	return width, height
}

func (l *loop) GamepadState(index int) (app.GamepadState, bool) {
	if _, ok := l.gamepads[index]; !ok {
		return app.GamepadState{}, false
	}

	gamepads := js.Global().Get("navigator").Call("getGamepads")
	if gamepads.IsNull() {
		return app.GamepadState{}, false
	}

	gamepad := gamepads.Index(index)
	if gamepad.IsNull() {
		return app.GamepadState{}, false
	}

	if gamepad.Get("mapping").String() != "standard" {
		return app.GamepadState{}, false
	}

	buttons := gamepad.Get("buttons")
	axes := gamepad.Get("axes")
	return app.GamepadState{
		LeftStickX:     float32(axes.Index(0).Float()),
		LeftStickY:     -float32(axes.Index(1).Float()),
		RightStickX:    float32(axes.Index(2).Float()),
		RightStickY:    -float32(axes.Index(3).Float()),
		LeftTrigger:    float32(buttons.Index(6).Get("value").Float()),
		RightTrigger:   float32(buttons.Index(7).Get("value").Float()),
		LeftBumper:     buttons.Index(4).Get("pressed").Bool(),
		RightBumper:    buttons.Index(5).Get("pressed").Bool(),
		SquareButton:   buttons.Index(2).Get("pressed").Bool(),
		CircleButton:   buttons.Index(1).Get("pressed").Bool(),
		TriangleButton: buttons.Index(3).Get("pressed").Bool(),
		CrossButton:    buttons.Index(0).Get("pressed").Bool(),
	}, true
}

func (l *loop) Schedule(fn func() error) {
	select {
	case l.tasks <- fn:
	default:
		panic(fmt.Errorf("failed to queue task; queue is full"))
	}
}

func (l *loop) Invalidate() {
	// Nothing to do here. Current implementation always invalidates.
}

func (l *loop) CreateCursor(definition app.CursorDefinition) app.Cursor {
	panic("TODO")
}

func (l *loop) UseCursor(cursor app.Cursor) {
	panic("TODO")
}

func (l *loop) CursorVisible() bool {
	cursor := l.htmlCanvas.Get("style").Get("cursor").String()
	return !strings.EqualFold(cursor, "none")
}

func (l *loop) SetCursorVisible(visible bool) {
	if visible {
		l.htmlCanvas.Get("style").Set("cursor", "auto")
	} else {
		l.htmlCanvas.Get("style").Set("cursor", "none")
	}
}

func (l *loop) SetCursorLocked(locked bool) {
	// FIXME: This should be recorded in a variable
	// and the actual call should be made as part of a user
	// gesture (e.g. click)
	if locked {
		l.htmlCanvas.Call("requestPointerLock")
	} else {
		l.htmlCanvas.Call("exitPointerLock")
	}
}

func (l *loop) Close() {
	l.shouldStop = true
}

func (l *loop) processTasks(limit time.Duration) bool {
	startTime := time.Now()
	for time.Since(startTime) < limit {
		select {
		case task := <-l.tasks:
			if err := task(); err != nil {
				panic(fmt.Errorf("task error: %w", err))
			}
		default:
			// No more tasks, we have consumed everything there
			// is for now.
			return true
		}
	}
	// We did not consume all available tasks within our time window.
	return false
}

func (l *loop) onJSKeyDown(this js.Value, args []js.Value) interface{} {
	event := args[0]
	event.Call("preventDefault")

	var downConsumed bool
	code := event.Get("code").String()
	if keyCode, ok := keyboardCodeMapping[code]; ok {
		var modifiers app.KeyModifierSet
		if event.Get("ctrlKey").Bool() {
			modifiers = modifiers | app.KeyModifierSet(app.KeyModifierControl)
		}
		if event.Get("shiftKey").Bool() {
			modifiers = modifiers | app.KeyModifierSet(app.KeyModifierShift)
		}
		if event.Get("altKey").Bool() {
			modifiers = modifiers | app.KeyModifierSet(app.KeyModifierAlt)
		}
		downConsumed = l.controller.OnKeyboardEvent(l, app.KeyboardEvent{
			Type:      app.KeyboardEventTypeKeyDown,
			Code:      keyCode,
			Modifiers: modifiers,
		})
	}

	// NOTE: JS has the keypress callback deprecated so we fake it here
	var pressConsumed bool
	key := event.Get("key").String()
	if len(key) == 1 {
		pressConsumed = l.controller.OnKeyboardEvent(l, app.KeyboardEvent{
			Type: app.KeyboardEventTypeType,
			Rune: ([]rune(key))[0],
		})
	}

	return downConsumed || pressConsumed
}

func (l *loop) onJSKeyUp(this js.Value, args []js.Value) interface{} {
	event := args[0]
	event.Call("preventDefault")

	code := event.Get("code").String()
	keyCode, ok := keyboardCodeMapping[code]
	if !ok {
		return false
	}
	var modifiers app.KeyModifierSet
	if event.Get("ctrlKey").Bool() {
		modifiers = modifiers | app.KeyModifierSet(app.KeyModifierControl)
	}
	if event.Get("shiftKey").Bool() {
		modifiers = modifiers | app.KeyModifierSet(app.KeyModifierShift)
	}
	if event.Get("altKey").Bool() {
		modifiers = modifiers | app.KeyModifierSet(app.KeyModifierAlt)
	}
	return l.controller.OnKeyboardEvent(l, app.KeyboardEvent{
		Type:      app.KeyboardEventTypeKeyUp,
		Code:      keyCode,
		Modifiers: modifiers,
	})
}

func (l *loop) onJSMouseEnter(this js.Value, args []js.Value) interface{} {
	event := args[0]
	return l.controller.OnMouseEvent(l, app.MouseEvent{
		Index: 0,
		X:     int(event.Get("offsetX").Float()),
		Y:     int(event.Get("offsetY").Float()),
		Type:  app.MouseEventTypeEnter,
	})
}

func (l *loop) onJSMouseLeave(this js.Value, args []js.Value) interface{} {
	event := args[0]
	return l.controller.OnMouseEvent(l, app.MouseEvent{
		Index: 0,
		X:     int(event.Get("offsetX").Float()),
		Y:     int(event.Get("offsetY").Float()),
		Type:  app.MouseEventTypeLeave,
	})
}

func (l *loop) onJSMouseMove(this js.Value, args []js.Value) interface{} {
	event := args[0]
	return l.controller.OnMouseEvent(l, app.MouseEvent{
		Index: 0,
		X:     int(event.Get("offsetX").Float()),
		Y:     int(event.Get("offsetY").Float()),
		Type:  app.MouseEventTypeMove,
	})
}

func (l *loop) onJSMouseDown(this js.Value, args []js.Value) interface{} {
	event := args[0]
	// NOTE: Don't prevent this event or the user will never be able
	// to select canvas for keyboard events.
	buttonIndex := event.Get("button").Int()
	return l.controller.OnMouseEvent(l, app.MouseEvent{
		Index:  0,
		X:      int(event.Get("offsetX").Float()),
		Y:      int(event.Get("offsetY").Float()),
		Type:   app.MouseEventTypeDown,
		Button: mouseButtonMapping[buttonIndex],
	})
}

func (l *loop) onJSMouseUp(this js.Value, args []js.Value) interface{} {
	event := args[0]
	event.Call("preventDefault")
	buttonIndex := event.Get("button").Int()
	return l.controller.OnMouseEvent(l, app.MouseEvent{
		Index:  0,
		X:      int(event.Get("offsetX").Float()),
		Y:      int(event.Get("offsetY").Float()),
		Type:   app.MouseEventTypeUp,
		Button: mouseButtonMapping[buttonIndex],
	})
}

func (l *loop) onJSMouseWheel(this js.Value, args []js.Value) interface{} {
	event := args[0]
	event.Call("preventDefault")
	return l.controller.OnMouseEvent(l, app.MouseEvent{
		Index:   0,
		X:       int(event.Get("offsetX").Float()),
		Y:       int(event.Get("offsetY").Float()),
		Type:    app.MouseEventTypeScroll,
		ScrollX: event.Get("deltaX").Float(),
		ScrollY: event.Get("deltaY").Float(),
	})
}

func (l *loop) onJSGamepadConnected(this js.Value, args []js.Value) interface{} {
	event := args[0]
	index := event.Get("gamepad").Get("index").Int()
	l.gamepads[index] = struct{}{}
	return true
}

func (l *loop) onJSGamepadDisconnected(this js.Value, args []js.Value) interface{} {
	event := args[0]
	index := event.Get("gamepad").Get("index").Int()
	delete(l.gamepads, index)
	return true
}
