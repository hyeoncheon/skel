package actions

import (
	"strings"

	"github.com/gobuffalo/buffalo"
)

const defaultLanguage = "en"

func currentLanguage(c buffalo.Context) string {
	if languages := c.Value("languages"); languages != nil {
		if langs, ok := languages.([]string); ok && len(langs) > 0 {
			return strings.Split(langs[0], "-")[0]
		}
	}
	return defaultLanguage
}
