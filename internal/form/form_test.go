package form

import (
	"html/template"
	"testing"
)

var (
	testTemplate = template.Must(template.New("").Parse(
		`<input type="{{.Type }}" name="{{.Name}}"{{with .Value}} value="{{.}}"{{end}}>`))
)

func TestHTML(t *testing.T) {
	testCases := map[string]struct {
		templ   *template.Template
		stract  interface{}
		want    template.HTML
		wantErr error
	}{
		"simple form with values": {
			templ: testTemplate,
			stract: struct {
				Name  string
				Email string
			}{
				Name:  "Jon Doe",
				Email: "jondoe@gmail.com",
			},
			want: `<input type="text" name="Name" value="Jon Doe">` +
				`<input type="text" name="Email" value="jondoe@gmail.com">`,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			got, gotErr := HTML(testCase.templ, testCase.stract)
			if gotErr != testCase.wantErr {
				t.Fatalf("HTML() err = %v; want %v;", gotErr, testCase.wantErr)
			}

			if got != testCase.want {
				t.Errorf("HTML() = %q; want %q", got, testCase.want)
			}
		})
	}

}
