package grifts

import (
	"github.com/gobuffalo/buffalo"

	"github.com/hyeoncheon/skel/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
