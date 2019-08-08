package form

import (
	"html/template"
	"reflect"
)

type field struct {
	Label       string
	Name        string
	Type        string
	Placeholder string
	Value       interface{}
}

func HTML(stract interface{}, templ *template.Template) template.HTML {
	return template.HTML("")
}

func valueOf(val interface{}) reflect.Value {
	var reflectValue reflect.Value
	switch value := val.(type) {
	case reflect.Value:
		reflectValue = value
	default:
		reflectValue = reflect.ValueOf(val)
	}

	if reflectValue.Kind() == reflect.Ptr {
		if reflectValue.IsNil() {
			reflectValue = reflect.New(reflectValue.Type().Elem())
		}

		reflectValue = reflectValue.Elem()
	}

	return reflectValue
}

func fields(stract interface{}) []field {
	reflectValue := valueOf(stract)

	if reflectValue.Kind() != reflect.Struct {
		panic("invalid value: only structs are applicable for forms")
	}
	reflectType := reflectValue.Type()

	var ret []field
	for i := 0; i < reflectType.NumField(); i++ {
		typeField := reflectType.Field(i)
		reflectValueField := valueOf(reflectValue.Field(i))
		if !reflectValueField.CanInterface() {
			continue
		}

		if reflectValueField.Kind() == reflect.Struct {
			nestedFields := fields(reflectValueField.Interface())
			for index, nestedField := range nestedFields {
				nestedFields[index].Name = typeField.Name + "." + nestedField.Name
			}

			ret = append(ret, nestedFields...)
			continue
		}

		ret = append(ret, field{
			Label:       typeField.Name,
			Name:        typeField.Name,
			Type:        "text",
			Placeholder: typeField.Name,
			Value:       reflectValueField.Interface(),
		})
	}

	return ret
}
