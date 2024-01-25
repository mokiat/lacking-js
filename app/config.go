package app

import "github.com/mokiat/lacking/app"

// NewConfig creates a new Config object that contains the minimum
// required settings.
func NewConfig(canvasID string) *Config {
	return &Config{
		canvasID:     canvasID,
		audioEnabled: true,
	}
}

// Config represents an application window configuration.
type Config struct {
	canvasID      string
	title         *string
	width         *int
	height        *int
	fullscreen    bool
	cursorVisible bool
	cursor        *app.CursorDefinition
	glExtensions  []string
	audioEnabled  bool
}

// Title returns the title of the application window.
func (c *Config) Title() string {
	if c.title == nil {
		return ""
	}
	return *c.title
}

// SetTitle sets the title of the application window.
func (c *Config) SetTitle(title string) {
	c.title = &title
}

// Width returns the width of the application window.
func (c *Config) Width() int {
	if c.width == nil {
		return 0
	}
	return *c.width
}

// SetWidth sets the width of the application window.
func (c *Config) SetWidth(width int) {
	c.width = &width
}

// Height returns the height of the application window.
func (c *Config) Height() int {
	if c.height == nil {
		return 0
	}
	return *c.height
}

// SetHeight sets the height of the application window.
func (c *Config) SetHeight(height int) {
	c.height = &height
}

// Fullscreen returns whether the application window will be
// displayed in fullscreen mode.
func (c *Config) Fullscreen() bool {
	return c.fullscreen
}

// SetFullscreen specifies whether the application window should
// be displayed in fullscreen mode.
func (c *Config) SetFullscreen(fullscreen bool) {
	c.fullscreen = fullscreen
}

// CursorVisible returns whether the cursor will be shown
// when hovering over the window.
func (c *Config) CursorVisible() bool {
	return c.cursorVisible
}

// SetCursorVisible specifies whether the cursor should be
// displayed when moved over the window.
func (c *Config) SetCursorVisible(visible bool) {
	c.cursorVisible = visible
}

// Cursor returns the cursor configuration for this application.
func (c *Config) Cursor() *app.CursorDefinition {
	return c.cursor
}

// SetCursor configures a custom cursor to be used.
// Specifying nil disables the custom cursor.
func (c *Config) SetCursor(definition *app.CursorDefinition) {
	c.cursor = definition
}

// AddGLExtension adds a new OpenGL extension that should be
// enabled when creating the OpenGL context.
func (c *Config) AddGLExtension(name string) {
	c.glExtensions = append(c.glExtensions, name)
}

// AudioEnabled returns whether audio is enabled for this application.
func (c *Config) AudioEnabled() bool {
	return c.audioEnabled
}

// SetAudioEnabled specifies whether audio should be enabled for this
// application.
func (c *Config) SetAudioEnabled(enabled bool) {
	c.audioEnabled = enabled
}
