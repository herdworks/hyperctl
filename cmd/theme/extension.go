package theme

import (
	httheme "github.com/herdworks/httt/theme"
)

type Themes []httheme.Theme

func (themes Themes) ColumnValues(field string) []string {
	values := []string{}
	for i := range themes {
		switch field {
		case "Path", "path":
			values = append(values, themes[i].URI.Path)
		case "Name", "name":
			values = append(values, themes[i].Name)
		case "Version", "version":
			values = append(values, themes[i].Version)
		}
	}
	return values
}
