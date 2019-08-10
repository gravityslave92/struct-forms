package fields

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFields(t *testing.T) {
	var nilStructPtr *struct {
		Name string
		Age  int
	}

	testCases := map[string]struct {
		stract interface{}
		want   []field
	}{
		"Simple Test": {
			stract: struct {
				Name string
			}{},
			want: []field{
				{
					Label:       "Name",
					Name:        "Name",
					Type:        "text",
					Placeholder: "Name",
					Value:       "",
				},
			},
		},
		"Field names should be determined": {
			stract: struct {
				FullName string
				Email    string
				Age      int
			}{},
			want: []field{
				{
					Label:       "FullName",
					Name:        "FullName",
					Type:        "text",
					Placeholder: "FullName",
					Value:       "",
				},
				{
					Label:       "Email",
					Name:        "Email",
					Type:        "text",
					Placeholder: "Email",
					Value:       "",
				},
				{
					Label:       "Age",
					Name:        "Age",
					Type:        "text",
					Placeholder: "Age",
					Value:       0,
				},
			},
		},
		"Values should be parsed": {
			stract: struct {
				FullName string
				Email    string
				Age      int
			}{
				FullName: "Jon Doe",
				Email:    "jondoe@gmail.com",
				Age:      123,
			},
			want: []field{
				{
					Label:       "FullName",
					Name:        "FullName",
					Type:        "text",
					Placeholder: "FullName",
					Value:       "Jon Doe",
				},
				{
					Label:       "Email",
					Name:        "Email",
					Type:        "text",
					Placeholder: "Email",
					Value:       "jondoe@gmail.com",
				},
				{
					Label:       "Age",
					Name:        "Age",
					Type:        "text",
					Placeholder: "Age",
					Value:       123,
				},
			},
		},
		"Pointers to structs should be supported": {
			stract: &struct {
				Name string
				Age  int
			}{
				"Jon Doe",
				321,
			},
			want: []field{
				{
					Label:       "Name",
					Name:        "Name",
					Type:        "text",
					Placeholder: "Name",
					Value:       "Jon Doe",
				},
				{
					Label:       "Age",
					Name:        "Age",
					Type:        "text",
					Placeholder: "Age",
					Value:       321,
				},
			},
		},
		"nilPointerstruct must be supported": {
			stract: nilStructPtr,
			want: []field{
				{
					Label:       "Name",
					Name:        "Name",
					Type:        "text",
					Placeholder: "Name",
					Value:       "",
				},
				{
					Label:       "Age",
					Name:        "Age",
					Type:        "text",
					Placeholder: "Age",
					Value:       0,
				},
			},
		},
		"Pointer fields should be supported": {
			stract: struct {
				Name *string
				Age  *int
			}{},
			want: []field{
				{
					Label:       "Name",
					Name:        "Name",
					Type:        "text",
					Placeholder: "Name",
					Value:       "",
				},
				{
					Label:       "Age",
					Name:        "Age",
					Type:        "text",
					Placeholder: "Age",
					Value:       0,
				},
			},
		},
		"Nested structs should be supported": {
			stract: struct {
				Name    string
				Address struct {
					Street string
					Zip    int
				}
			}{
				Name: "Jon Doe",
				Address: struct {
					Street string
					Zip    int
				}{
					"random street 456",
					987654,
				},
			},
			want: []field{
				{
					Label:       "Name",
					Name:        "Name",
					Type:        "text",
					Placeholder: "Name",
					Value:       "Jon Doe",
				},
				{
					Label:       "Street",
					Name:        "Address.Street",
					Type:        "text",
					Placeholder: "Street",
					Value:       "random street 456",
				},
				{
					Label:       "Zip",
					Name:        "Address.Zip",
					Type:        "text",
					Placeholder: "Zip",
					Value:       987654,
				},
			},
		},
		"Dobule Nested structs should be supported": {
			stract: struct {
				A struct {
					B struct {
						C1 string
						C2 int
					}
				}
			}{
				A: struct {
					B struct {
						C1 string
						C2 int
					}
				}{
					B: struct {
						C1 string
						C2 int
					}{
						C1: "A string",
						C2: 123,
					},
				},
			},
			want: []field{
				{
					Label:       "C1",
					Name:        "A.B.C1",
					Type:        "text",
					Placeholder: "C1",
					Value:       "A string",
				},
				{
					Label:       "C2",
					Name:        "A.B.C2",
					Type:        "text",
					Placeholder: "C2",
					Value:       123,
				},
			},
		},
		"Nested pointer structs should be supported": {
			stract: struct {
				Name    string
				Address *struct {
					Street string
					Zip    int
				}
				ContactInfo *struct {
					Phone string
				}
			}{
				Name: "Jon Doe",
				Address: &struct {
					Street string
					Zip    int
				}{
					Street: "123 Elm street",
					Zip:    123456,
				},
			},
			want: []field{
				{
					Label:       "Name",
					Name:        "Name",
					Type:        "text",
					Placeholder: "Name",
					Value:       "Jon Doe",
				},
				{
					Label:       "Street",
					Name:        "Address.Street",
					Type:        "text",
					Placeholder: "Street",
					Value:       "123 Elm street",
				},
				{
					Label:       "Zip",
					Name:        "Address.Zip",
					Type:        "text",
					Placeholder: "Zip",
					Value:       123456,
				},
				{
					Label:       "Phone",
					Name:        "ContactInfo.Phone",
					Type:        "text",
					Placeholder: "Phone",
					Value:       "",
				},
			},
		},
		"Tag fields": struct {
			stract interface{}
			want   []field
		}{
			stract: struct {
				LabelTest       string `form:"label=This is custom"`
				NameTest        string `form:"name=age"`
				TypeTest        int    `form:"type=number"`
				PlaceholderTest string `form:"placeholder=put your value here..."`
				Nested          struct {
					MultiTest string `form:"name=NestedMultiTest;label=This is a nested label;type=email;placeholder=example@gmail.com"`
				}
			}{
				PlaceholderTest: "placeholding",
			},
			want: []field{
				{
					Label:       "This is custom",
					Name:        "LabelTest",
					Type:        "text",
					Placeholder: "LabelTest",
					Value:       "",
				},
				{
					Label:       "NameTest",
					Name:        "age",
					Type:        "text",
					Placeholder: "NameTest",
					Value:       "",
				},
				{
					Label:       "TypeTest",
					Name:        "TypeTest",
					Type:        "number",
					Placeholder: "TypeTest",
					Value:       0,
				},
				{
					Label:       "PlaceholderTest",
					Name:        "PlaceholderTest",
					Type:        "text",
					Placeholder: "put your value here...",
					Value:       "placeholding",
				},
				{
					Label:       "This is a nested label",
					Name:        "NestedMultiTest",
					Type:        "email",
					Placeholder: "example@gmail.com",
					Value:       "",
				},
			},
		},
	}

	for key, testCase := range testCases {
		t.Run(key, func(t *testing.T) {
			got := Fields(testCase.stract)
			if reflect.DeepEqual(got, testCase.want) {
				return
			}

			if len(got) != len(testCase.want) {
				t.Errorf("Fields() len = %d; want %d", len(got), len(testCase.want))
			}

			for index, gotField := range got {
				if index > len(testCase.want) {
					break
				}

				wantField := testCase.want[index]
				if reflect.DeepEqual(gotField, wantField) {
					continue
				}

				t.Errorf("Fields()[%d]", index)
				if gotField.Label != wantField.Label {
					t.Errorf(" .Label = %v; want %v", gotField.Label, wantField.Label)
				}

				if gotField.Name != wantField.Name {
					t.Errorf(" .Name = %v; want %v", gotField.Name, wantField.Name)
				}

				if gotField.Type != wantField.Type {
					t.Errorf(" .Type = %v; want %v", gotField.Type, wantField.Type)
				}

				if gotField.Placeholder != wantField.Placeholder {
					t.Errorf(" .Placeholder = %v; want %v", gotField.Placeholder, wantField.Placeholder)
				}

				if gotField.Value != wantField.Value {
					t.Errorf(" .Value = %v; want %v", gotField.Value, wantField.Value)
				}
			}
		})
	}
}

func TestFields_invalidValues(t *testing.T) {
	testCases := []struct {
		notAStruct interface{}
	}{
		{"this is a string"},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%t", testCase.notAStruct), func(t *testing.T) {
			defer func() {
				if err := recover(); err == nil {
					t.Errorf("fields(%v) did not panic", testCase.notAStruct)
				}
			}()
			Fields(testCase.notAStruct)
		})
	}
}

func TestParseTags(t *testing.T) {
	testCases := map[string]struct {
		arg  reflect.StructField
		want map[string]string
	}{
		"empty tag": struct {
			arg  reflect.StructField
			want map[string]string
		}{
			arg:  reflect.StructField{},
			want: nil,
		},
		"label tag": struct {
			arg  reflect.StructField
			want map[string]string
		}{
			arg: reflect.StructField{Tag: `form:"label=Full Name"`},
			want: map[string]string{
				"label": "Full Name",
			},
		},
		"multiple tags": struct {
			arg  reflect.StructField
			want map[string]string
		}{
			arg: reflect.StructField{Tag: `form:"label=Full Name;name=full_name"`},
			want: map[string]string{
				"label": "Full Name",
				"name":  "full_name",
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			got := parseTags(testCase.arg)
			if len(got) != len(testCase.want) {
				t.Errorf("parseTags() len =%d, want %d", len(got), len(testCase.want))
			}

			for key, value := range testCase.want {
				gotVal, ok := got[key]
				if !ok {
					t.Errorf("parseTags() miising %s", key)
					continue
				}

				if gotVal != value {
					t.Errorf("parseTags()[%q] = %q, want %q", key, gotVal, value)
				}
				delete(got, key)
			}

			for gotKey, gotValue := range got {
				t.Errorf("parseTags() exta key %q, value = %q", gotKey, gotValue)
			}
		})
	}
}

func TestParseTags_Invalid(t *testing.T) {
	testCases := []struct {
		arg reflect.StructField
	}{
		{reflect.StructField{Tag: `form:"invalid-value"`}},
	}

	for _, testCase := range testCases {
		t.Run(string(testCase.arg.Tag), func(t *testing.T) {
			defer func() {
				if err := recover(); err == nil {
					t.Error("parseTags() did not panic")
				}
			}()

			parseTags(testCase.arg)
		})
	}
}
