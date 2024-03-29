package actions

import (
	"html/template"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gofrs/uuid"
	"github.com/hyeoncheon/skel/public"
	"github.com/hyeoncheon/skel/templates"
)

var r *render.Engine

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.html",

		// Box containing all of the templates:
		TemplatesFS: templates.FS(),
		AssetsFS:    public.FS(),

		// Add template helpers here:
		Helpers: render.Helpers{
			// uncomment for non-Bootstrap form helpers:
			// "form":     plush.FormHelper,
			// "form_for": plush.FormForHelper,
			"iconize": func(s string) template.HTML {
				switch s {
				case "admin":
					return template.HTML(`<i class="fab fa-empire"></i>`)
				case "doctor":
					return template.HTML(`<i class="fas fa-pencil-square"></i>`)
				default:
					return template.HTML(`<i class="fas fa-` + s + `"></i>`)
				}
			},
			"trunc": func(t interface{}, args ...int) string {
				length := 20
				var s string
				switch t.(type) {
				case string:
					s = t.(string)
				case uuid.UUID:
					s = t.(uuid.UUID).String()
					length = 14
				}

				if len(args) > 0 {
					length = args[0]
				}
				if length > len(s)-4 {
					return s
				}
				return s[0:length] + "..."
			},
		},
	})
}
