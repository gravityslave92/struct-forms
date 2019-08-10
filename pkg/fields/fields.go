package fields

import (
	"reflect"
	"strings"
)

type field struct {
	Label       string
	Name        string
	Type        string
	Placeholder string
	Value       interface{}
	Errors      []string
}

type ErrorField struct {
	Field string
	Error string
}

func (f *field) applyTags(tags map[string]string) {
	if value, ok := tags["name"]; ok {
		f.Name = value
	}

	if value, ok := tags["label"]; ok {
		f.Label = value
	}

	if value, ok := tags["type"]; ok {
		f.Type = value
	}

	if value, ok := tags["placeholder"]; ok {
		f.Placeholder = value
	}

	if value, ok := tags["value"]; ok {
		f.Value = value
	}
}

func (f *field) setErrors(errors []ErrorField)  {
	for _, fieldErr := range errors {
		if fieldErr.Field == f.Name {
			f.Errors = append(f.Errors, fieldErr.Error)
		}
	}
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

func Fields(stract interface{}, parentNames ...string) []field {
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
			nestedParentNames := append(parentNames, typeField.Name)
			nestedFields := Fields(reflectValueField.Interface(), nestedParentNames...)

			ret = append(ret, nestedFields...)
			continue
		}
		names := append(parentNames, typeField.Name)
		name := strings.Join(names, ".")
		fld := field{
			Label:       typeField.Name,
			Name:        name,
			Type:        "text",
			Placeholder: typeField.Name,
			Value:       reflectValueField.Interface(),
		}
		// parse struct field tag and apply if exists
		fld.applyTags(parseTags(typeField))
		ret = append(ret, fld)
	}

	return ret
}

func parseTags(stractField reflect.StructField) map[string]string {
	stractTag := stractField.Tag.Get("form")
	if len(stractTag) == 0 {
		return nil
	}
	ret := make(map[string]string)

	tags := strings.Split(stractTag, ";")
	for _, tag := range tags {
		keyValuePair := strings.Split(tag, "=")
		if len(keyValuePair) != 2 {
			panic("invalid struct field tag")
		}

		key, value := keyValuePair[0], keyValuePair[1]
		ret[key] = value
	}
	return ret
}
