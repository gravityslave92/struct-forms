package fields

import (
	"html/template"
	"strings"
)



func HTML(t *template.Template, stract interface{}, errors ...ErrorField) (template.HTML, error) {
	var inputs []string
	for _, field := range Fields(stract) {
		field.setErrors(errors)
		var sb strings.Builder
		err := t.Execute(&sb, field)
		if err != nil {
			return "", err
		}
		inputs = append(inputs, sb.String())
	}

	return template.HTML(strings.Join(inputs, "")), nil
}
