//go:build js && wasm

package app

import "github.com/mokiat/lacking/app"

var (
	keyboardCodeMapping map[string]app.KeyCode
)

func init() {
	keyboardCodeMapping = make(map[string]app.KeyCode)
	keyboardCodeMapping["Escape"] = app.KeyCodeEscape
	keyboardCodeMapping["Enter"] = app.KeyCodeEnter
	keyboardCodeMapping["Space"] = app.KeyCodeSpace
	keyboardCodeMapping["Tab"] = app.KeyCodeTab
	keyboardCodeMapping["CapsLock"] = app.KeyCodeCaps
	keyboardCodeMapping["ShiftLeft"] = app.KeyCodeLeftShift
	keyboardCodeMapping["ShiftRight"] = app.KeyCodeRightShift
	keyboardCodeMapping["ControlLeft"] = app.KeyCodeLeftControl
	keyboardCodeMapping["ControlRight"] = app.KeyCodeRightControl
	keyboardCodeMapping["AltLeft"] = app.KeyCodeLeftAlt
	keyboardCodeMapping["AltRight"] = app.KeyCodeRightAlt
	keyboardCodeMapping["MetaLeft"] = app.KeyCodeLeftSuper
	keyboardCodeMapping["MetaRight"] = app.KeyCodeRightSuper
	keyboardCodeMapping["Backspace"] = app.KeyCodeBackspace
	keyboardCodeMapping["Insert"] = app.KeyCodeInsert
	keyboardCodeMapping["Delete"] = app.KeyCodeDelete
	keyboardCodeMapping["Home"] = app.KeyCodeHome
	keyboardCodeMapping["End"] = app.KeyCodeEnd
	keyboardCodeMapping["PageUp"] = app.KeyCodePageUp
	keyboardCodeMapping["PageDown"] = app.KeyCodePageDown
	keyboardCodeMapping["ArrowLeft"] = app.KeyCodeArrowLeft
	keyboardCodeMapping["ArrowRight"] = app.KeyCodeArrowRight
	keyboardCodeMapping["ArrowUp"] = app.KeyCodeArrowUp
	keyboardCodeMapping["ArrowDown"] = app.KeyCodeArrowDown
	keyboardCodeMapping["Minus"] = app.KeyCodeMinus
	keyboardCodeMapping["Equal"] = app.KeyCodeEqual
	keyboardCodeMapping["BracketLeft"] = app.KeyCodeLeftBracket
	keyboardCodeMapping["BracketRight"] = app.KeyCodeRightBracket
	keyboardCodeMapping["Semicolon"] = app.KeyCodeSemicolon
	keyboardCodeMapping["Comma"] = app.KeyCodeComma
	keyboardCodeMapping["Period"] = app.KeyCodePeriod
	keyboardCodeMapping["Slash"] = app.KeyCodeSlash
	keyboardCodeMapping["Backslash"] = app.KeyCodeBackslash
	keyboardCodeMapping["Quote"] = app.KeyCodeApostrophe
	keyboardCodeMapping["Backquote"] = app.KeyCodeGraveAccent
	keyboardCodeMapping["KeyA"] = app.KeyCodeA
	keyboardCodeMapping["KeyB"] = app.KeyCodeB
	keyboardCodeMapping["KeyC"] = app.KeyCodeC
	keyboardCodeMapping["KeyD"] = app.KeyCodeD
	keyboardCodeMapping["KeyE"] = app.KeyCodeE
	keyboardCodeMapping["KeyF"] = app.KeyCodeF
	keyboardCodeMapping["KeyG"] = app.KeyCodeG
	keyboardCodeMapping["KeyH"] = app.KeyCodeH
	keyboardCodeMapping["KeyI"] = app.KeyCodeI
	keyboardCodeMapping["KeyJ"] = app.KeyCodeJ
	keyboardCodeMapping["KeyK"] = app.KeyCodeK
	keyboardCodeMapping["KeyL"] = app.KeyCodeL
	keyboardCodeMapping["KeyM"] = app.KeyCodeM
	keyboardCodeMapping["KeyN"] = app.KeyCodeN
	keyboardCodeMapping["KeyO"] = app.KeyCodeO
	keyboardCodeMapping["KeyP"] = app.KeyCodeP
	keyboardCodeMapping["KeyQ"] = app.KeyCodeQ
	keyboardCodeMapping["KeyR"] = app.KeyCodeR
	keyboardCodeMapping["KeyS"] = app.KeyCodeS
	keyboardCodeMapping["KeyT"] = app.KeyCodeT
	keyboardCodeMapping["KeyU"] = app.KeyCodeU
	keyboardCodeMapping["KeyV"] = app.KeyCodeV
	keyboardCodeMapping["KeyW"] = app.KeyCodeW
	keyboardCodeMapping["KeyX"] = app.KeyCodeX
	keyboardCodeMapping["KeyY"] = app.KeyCodeY
	keyboardCodeMapping["KeyZ"] = app.KeyCodeZ
	keyboardCodeMapping["Digit0"] = app.KeyCode0
	keyboardCodeMapping["Digit1"] = app.KeyCode1
	keyboardCodeMapping["Digit2"] = app.KeyCode2
	keyboardCodeMapping["Digit3"] = app.KeyCode3
	keyboardCodeMapping["Digit4"] = app.KeyCode4
	keyboardCodeMapping["Digit5"] = app.KeyCode5
	keyboardCodeMapping["Digit6"] = app.KeyCode6
	keyboardCodeMapping["Digit7"] = app.KeyCode7
	keyboardCodeMapping["Digit8"] = app.KeyCode8
	keyboardCodeMapping["Digit9"] = app.KeyCode9
	keyboardCodeMapping["F1"] = app.KeyCodeF1
	keyboardCodeMapping["F2"] = app.KeyCodeF2
	keyboardCodeMapping["F3"] = app.KeyCodeF3
	keyboardCodeMapping["F4"] = app.KeyCodeF4
	keyboardCodeMapping["F5"] = app.KeyCodeF5
	keyboardCodeMapping["F6"] = app.KeyCodeF6
	keyboardCodeMapping["F7"] = app.KeyCodeF7
	keyboardCodeMapping["F8"] = app.KeyCodeF8
	keyboardCodeMapping["F9"] = app.KeyCodeF9
	keyboardCodeMapping["F10"] = app.KeyCodeF10
	keyboardCodeMapping["F11"] = app.KeyCodeF11
	keyboardCodeMapping["F12"] = app.KeyCodeF12
}
