//go:build js && wasm

package app

import (
	"fmt"
	"strings"
	"syscall/js"
	"time"

	jsrender "github.com/mokiat/lacking-js/render"
	"github.com/mokiat/lacking/app"
	"github.com/mokiat/lacking/audio"
	"github.com/mokiat/lacking/render"
)

const (
	taskQueueSize         = 1024
	taskProcessingTimeout = 30 * time.Millisecond
)

func newLoop(htmlDocument, htmlCanvas js.Value, controller app.Controller) *loop {
	return &loop{
		platform:     newPlatform(),
		htmlDocument: htmlDocument,
		htmlCanvas:   htmlCanvas,
		controller:   controller,
		renderAPI:    jsrender.NewAPI(),
		tasks:        make(chan func(), taskQueueSize),
		gamepads: [4]*Gamepad{
			newGamepad(0),
			newGamepad(1),
			newGamepad(2),
			newGamepad(3),
		},
		shouldStop: false,
	}
}

var _ app.Window = (*loop)(nil)

type loop struct {
	platform     *platform
	htmlDocument js.Value
	htmlCanvas   js.Value
	controller   app.Controller
	renderAPI    render.API
	cursor       *Cursor
	tasks        chan func()
	gamepads     [4]*Gamepad
	shouldStop   bool

	knownFramebufferWidth  int
	knownFramebufferHeight int
	knownWidth             int
	knownHeight            int

	clipboardCallback js.Func
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
	l.htmlCanvas.Call("addEventListener", "pointerenter", mouseEnterCallback)
	defer l.htmlCanvas.Call("removeEventListener", "pointerenter", mouseEnterCallback)

	mouseLeaveCallback := js.FuncOf(l.onJSMouseLeave)
	defer mouseLeaveCallback.Release()
	l.htmlCanvas.Call("addEventListener", "pointerleave", mouseLeaveCallback)
	defer l.htmlCanvas.Call("removeEventListener", "pointerleave", mouseLeaveCallback)

	mouseMoveCallback := js.FuncOf(l.onJSMouseMove)
	defer mouseMoveCallback.Release()
	l.htmlCanvas.Call("addEventListener", "pointermove", mouseMoveCallback)
	defer l.htmlCanvas.Call("removeEventListener", "pointermove", mouseMoveCallback)

	mouseDownCallback := js.FuncOf(l.onJSMouseDown)
	defer mouseDownCallback.Release()
	l.htmlCanvas.Call("addEventListener", "pointerdown", mouseDownCallback)
	defer l.htmlCanvas.Call("removeEventListener", "pointerdown", mouseDownCallback)

	mouseUpCallback := js.FuncOf(l.onJSMouseUp)
	defer mouseUpCallback.Release()
	l.htmlCanvas.Call("addEventListener", "pointerup", mouseUpCallback)
	defer l.htmlCanvas.Call("removeEventListener", "pointerup", mouseUpCallback)

	mouseScrollCallback := js.FuncOf(l.onJSMouseWheel)
	defer mouseScrollCallback.Release()
	l.htmlCanvas.Call("addEventListener", "wheel", mouseScrollCallback)
	defer l.htmlCanvas.Call("removeEventListener", "wheel", mouseScrollCallback)

	closeCallback := js.FuncOf(l.onCloseRequested)
	defer closeCallback.Release()
	js.Global().Set("onbeforeunload", closeCallback)

	l.clipboardCallback = js.FuncOf(l.onClipboardReadText)
	defer l.clipboardCallback.Release()

	l.knownFramebufferWidth, l.knownFramebufferHeight = l.FramebufferSize()
	l.controller.OnFramebufferResize(l, l.knownFramebufferWidth, l.knownFramebufferHeight)

	l.knownWidth, l.knownHeight = l.Size()
	l.controller.OnResize(l, l.knownWidth, l.knownHeight)

	done := make(chan error, 1)
	var loopFunc js.Func
	loopFunc = js.FuncOf(func(this js.Value, args []js.Value) any {
		l.checkResized()

		if l.shouldStop {
			if l.processTasks(5 * time.Second) {
				done <- nil
			} else {
				done <- fmt.Errorf("failed to cleanup within timeout")
			}
			return true
		}

		for _, gamepad := range l.gamepads {
			gamepad.markDirty()
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

func (l *loop) Platform() app.Platform {
	return l.platform
}

func (l *loop) Title() string {
	return l.htmlDocument.Get("title").String()
}

func (l *loop) SetTitle(title string) {
	l.htmlDocument.Set("title", title)
}

func (l *loop) SetSize(width, height int) {
	l.htmlCanvas.Set("clientWidth", width)
	l.htmlCanvas.Set("clientHeight", height)
}

func (l *loop) Size() (int, int) {
	width := l.htmlCanvas.Get("clientWidth").Int()
	height := l.htmlCanvas.Get("clientHeight").Int()
	return width, height
}

func (l *loop) FramebufferSize() (int, int) {
	width := l.htmlCanvas.Get("width").Int()
	height := l.htmlCanvas.Get("height").Int()
	return width, height
}

func (l *loop) Gamepads() [4]app.Gamepad {
	var result [4]app.Gamepad
	for i := range result {
		result[i] = l.gamepads[i]
	}
	return result
}

func (l *loop) Schedule(fn func()) {
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
	return &Cursor{
		path:     definition.Path,
		hotspotX: definition.HotspotX,
		hotspotY: definition.HotspotY,
	}
}

func (l *loop) UseCursor(cursor app.Cursor) {
	if appCursor, ok := cursor.(*Cursor); ok {
		l.cursor = appCursor
	} else {
		l.cursor = nil
	}
	l.SetCursorVisible(l.CursorVisible()) // force refresh
}

func (l *loop) CursorVisible() bool {
	cursor := l.htmlCanvas.Get("style").Get("cursor").String()
	return !strings.EqualFold(cursor, "none")
}

func (l *loop) SetCursorVisible(visible bool) {
	if visible {
		if l.cursor != nil {
			cursorStyle := fmt.Sprintf("url(%s) %d %d, auto",
				l.cursor.path, l.cursor.hotspotX, l.cursor.hotspotY)
			l.htmlCanvas.Get("style").Set("cursor", cursorStyle)
		} else {
			l.htmlCanvas.Get("style").Set("cursor", "auto")
		}
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

func (l *loop) RequestCopy(text string) {
	jsNavigator := js.Global().Get("navigator")
	if jsNavigator.IsUndefined() || jsNavigator.IsNull() {
		appLogger.Warn("JavaScript navigator not found!")
		return
	}
	jsClipboard := jsNavigator.Get("clipboard")
	if jsClipboard.IsUndefined() || jsClipboard.IsNull() {
		appLogger.Warn("JavaScript clipboard not found!")
		return
	}
	jsClipboard.Call("writeText", text)
}

func (l *loop) RequestPaste() {
	jsNavigator := js.Global().Get("navigator")
	if jsNavigator.IsUndefined() || jsNavigator.IsNull() {
		appLogger.Warn("JavaScript navigator not found!")
		return
	}
	jsClipboard := jsNavigator.Get("clipboard")
	if jsClipboard.IsUndefined() || jsClipboard.IsNull() {
		appLogger.Warn("JavaScript clipboard not found!")
		return
	}
	jsPromise := jsClipboard.Call("readText")
	if jsPromise.IsUndefined() || jsPromise.IsNull() {
		appLogger.Warn("JavaScript clipboard.readText promise missing!")
		return
	}
	jsPromise.Call("then", l.clipboardCallback)
}

func (l *loop) RenderAPI() render.API {
	return l.renderAPI
}

func (l *loop) AudioAPI() audio.API {
	return nil
}

func (l *loop) Close() {
	l.shouldStop = true
}

func (l *loop) checkResized() {
	framebufferWidth, framebufferHeight := l.FramebufferSize()
	if framebufferWidth != l.knownFramebufferWidth || framebufferHeight != l.knownFramebufferHeight {
		l.knownFramebufferWidth = framebufferWidth
		l.knownFramebufferHeight = framebufferHeight
		l.controller.OnFramebufferResize(l, framebufferWidth, framebufferHeight)
	}

	width, height := l.Size()
	if width != l.knownWidth || height != l.knownHeight {
		l.knownWidth = width
		l.knownHeight = height
		l.controller.OnResize(l, width, height)
	}
}

func (l *loop) processTasks(limit time.Duration) bool {
	startTime := time.Now()
	for time.Since(startTime) < limit {
		select {
		case task := <-l.tasks:
			task()
		default:
			// No more tasks, we have consumed everything there
			// is for now.
			return true
		}
	}
	// We did not consume all available tasks within our time window.
	return false
}

func (l *loop) onJSKeyDown(this js.Value, args []js.Value) any {
	event := args[0]
	event.Call("preventDefault")

	var downConsumed bool
	code := event.Get("code").String()
	if keyCode, ok := keyboardCodeMapping[code]; ok {
		action := app.KeyboardActionDown
		if event.Get("repeat").Bool() {
			action = app.KeyboardActionRepeat
		}
		downConsumed = l.controller.OnKeyboardEvent(l, app.KeyboardEvent{
			Action: action,
			Code:   keyCode,
		})
	}

	// NOTE: JS has the keypress callback deprecated so we fake it here.
	var pressConsumed bool
	key := []rune(event.Get("key").String())
	if len(key) == 1 {
		pressConsumed = l.controller.OnKeyboardEvent(l, app.KeyboardEvent{
			Action:    app.KeyboardActionType,
			Character: key[0],
		})
	}

	return downConsumed || pressConsumed
}

func (l *loop) onJSKeyUp(this js.Value, args []js.Value) any {
	event := args[0]
	event.Call("preventDefault")

	code := event.Get("code").String()
	keyCode, ok := keyboardCodeMapping[code]
	if !ok {
		return false
	}

	return l.controller.OnKeyboardEvent(l, app.KeyboardEvent{
		Action: app.KeyboardActionUp,
		Code:   keyCode,
	})
}

func (l *loop) onJSMouseEnter(this js.Value, args []js.Value) any {
	event := args[0]
	return l.controller.OnMouseEvent(l, app.MouseEvent{
		Index:  0,
		Action: app.MouseActionEnter,
		X:      int(event.Get("offsetX").Float()),
		Y:      int(event.Get("offsetY").Float()),
	})
}

func (l *loop) onJSMouseLeave(this js.Value, args []js.Value) any {
	event := args[0]
	return l.controller.OnMouseEvent(l, app.MouseEvent{
		Index:  0,
		Action: app.MouseActionLeave,
		X:      int(event.Get("offsetX").Float()),
		Y:      int(event.Get("offsetY").Float()),
	})
}

func (l *loop) onJSMouseMove(this js.Value, args []js.Value) any {
	event := args[0]
	return l.controller.OnMouseEvent(l, app.MouseEvent{
		Index:  0,
		Action: app.MouseActionMove,
		X:      int(event.Get("offsetX").Float()),
		Y:      int(event.Get("offsetY").Float()),
	})
}

func (l *loop) onJSMouseDown(this js.Value, args []js.Value) any {
	event := args[0]
	l.htmlCanvas.Call("setPointerCapture", event.Get("pointerId"))

	// NOTE: Don't prevent this event or the user will never be able
	// to select the canvas for keyboard events.
	buttonIndex := event.Get("button").Int()
	return l.controller.OnMouseEvent(l, app.MouseEvent{
		Index:  0,
		Action: app.MouseActionDown,
		X:      int(event.Get("offsetX").Float()),
		Y:      int(event.Get("offsetY").Float()),
		Button: mouseButtonMapping[buttonIndex],
	})
}

func (l *loop) onJSMouseUp(this js.Value, args []js.Value) any {
	event := args[0]
	event.Call("preventDefault")
	buttonIndex := event.Get("button").Int()
	return l.controller.OnMouseEvent(l, app.MouseEvent{
		Index:  0,
		Action: app.MouseActionUp,
		X:      int(event.Get("offsetX").Float()),
		Y:      int(event.Get("offsetY").Float()),
		Button: mouseButtonMapping[buttonIndex],
	})
}

func (l *loop) onJSMouseWheel(this js.Value, args []js.Value) any {
	event := args[0]
	event.Call("preventDefault")
	return l.controller.OnMouseEvent(l, app.MouseEvent{
		Index:   0,
		Action:  app.MouseActionScroll,
		X:       int(event.Get("offsetX").Float()),
		Y:       int(event.Get("offsetY").Float()),
		ScrollX: event.Get("deltaX").Float() / 100.0,
		ScrollY: event.Get("deltaY").Float() / 100.0,
	})
}

func (l *loop) onCloseRequested(this js.Value, args []js.Value) any {
	if !l.controller.OnCloseRequested(l) {
		return "reject"
	}
	return js.Null()
}

func (l *loop) onClipboardReadText(this js.Value, args []js.Value) any {
	data := args[0]
	text := data.String()
	l.Schedule(func() {
		l.controller.OnClipboardEvent(l, app.ClipboardEvent{
			Text: text,
		})
	})
	return nil
}
