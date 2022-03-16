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
		shouldStop:   false,
	}
}

var _ app.Window = (*loop)(nil)

type loop struct {
	htmlDocument js.Value
	htmlCanvas   js.Value
	controller   app.Controller
	tasks        chan func() error
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
	// var joystick glfw.Joystick
	// switch index {
	// case 0:
	// 	joystick = glfw.Joystick1
	// case 1:
	// 	joystick = glfw.Joystick2
	// case 2:
	// 	joystick = glfw.Joystick3
	// case 3:
	// 	joystick = glfw.Joystick4
	// default:
	// 	return app.GamepadState{}, false
	// }
	// if !joystick.Present() || !joystick.IsGamepad() {
	// 	return app.GamepadState{}, false
	// }

	// state := joystick.GetGamepadState()
	// return app.GamepadState{
	// 	LeftStickX:     state.Axes[glfw.AxisLeftX],
	// 	LeftStickY:     -state.Axes[glfw.AxisLeftY],
	// 	RightStickX:    state.Axes[glfw.AxisRightX],
	// 	RightStickY:    -state.Axes[glfw.AxisRightY],
	// 	LeftTrigger:    (state.Axes[glfw.AxisLeftTrigger] + 1.0) / 2.0,
	// 	RightTrigger:   (state.Axes[glfw.AxisRightTrigger] + 1.0) / 2.0,
	// 	LeftBumper:     state.Buttons[glfw.ButtonLeftBumper] == glfw.Press,
	// 	RightBumper:    state.Buttons[glfw.ButtonRightBumper] == glfw.Press,
	// 	SquareButton:   state.Buttons[glfw.ButtonSquare] == glfw.Press,
	// 	CircleButton:   state.Buttons[glfw.ButtonCircle] == glfw.Press,
	// 	TriangleButton: state.Buttons[glfw.ButtonTriangle] == glfw.Press,
	// 	CrossButton:    state.Buttons[glfw.ButtonCross] == glfw.Press,
	// }, true
	panic("TODO")
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
