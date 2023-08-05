package app

import "github.com/mokiat/lacking/app"

var _ app.Cursor = (*Cursor)(nil)

type Cursor struct {
	path     string
	hotspotX int
	hotspotY int
}

func (c *Cursor) Destroy() {
	c.path = ""
}
