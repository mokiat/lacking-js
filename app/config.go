package app

// NewConfig creates a new Config object that contains the minimum
// required settings.
func NewConfig(canvasID string) *Config {
	return &Config{
		canvasID: canvasID,
	}
}

// Config represents an application window configuration.
type Config struct {
	canvasID   string
	title      *string
	width      *int
	height     *int
	fullscreen bool
}

func (c *Config) Title() string {
	if c.title == nil {
		return ""
	}
	return *c.title
}

func (c *Config) SetTitle(title string) {
	c.title = &title
}

func (c *Config) Width() int {
	if c.width == nil {
		return 0
	}
	return *c.width
}

func (c *Config) SetWidth(width int) {
	c.width = &width
}

func (c *Config) Height() int {
	if c.height == nil {
		return 0
	}
	return *c.height
}

func (c *Config) SetHeight(height int) {
	c.height = &height
}

func (c *Config) Fullscreen() bool {
	return c.fullscreen
}

func (c *Config) SetFullscreen(fullscreen bool) {
	c.fullscreen = fullscreen
}
