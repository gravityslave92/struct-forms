package form

import (
	"forms_with_go/internal/fields"
	"html/template"
	"strings"
)

func HTML(t *template.Template, stract interface{}) (template.HTML, error) {
	var inputs []string
	for _, field := range fields.Fields(stract) {
		var sb strings.Builder
		err := t.Execute(&sb, field)
		if err != nil {
			return "", err
		}
		inputs = append(inputs, sb.String())
	}

	return template.HTML(strings.Join(inputs, "")), nil
}
