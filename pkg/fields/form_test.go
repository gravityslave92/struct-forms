package fields

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

var (
	testTemplate = template.Must(template.New("").Parse(
		`<input type="{{.Type }}" name="{{.Name}}"{{with .Value}} value="{{.}}"{{end}}>`))
	testFullTemplate = template.Must(template.New("").Parse(`
	<label>{{.Label}}</label>
	<input
		type="{{.Type}}"
		name="{{.Name}}"
		placeholder="{{.Placeholder}}"
		{{with .Value}} value="{{.}}"{{end}}>`))
	testErrorTemplate = template.Must(template.New("").Parse(`
<label>{{.Label}}</label>
	<input
		class="{{with .Errors}}border-red{{end}}"
		type="{{.Type}}"
		name="{{.Name}}"
		placeholder="{{.Placeholder}}"
		{{with .Value}}value="{{.}}"{{end}}>
{{range .Errors}}<p class="text-red text-xs italic">{{.}}</p>
{{end}}`))
)

var updateFlag bool

func init() {
	flag.BoolVar(&updateFlag, "update", false,
		"set the update flag in order to refine expected output of golden file tests")
}

func TestHTML(t *testing.T) {
	testCases := map[string]struct {
		templ   *template.Template
		stract  interface{}
		errors  []ErrorField
		want    string
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
			want: "TestHTML_basic.golden",
		},
		"form with custom tags": {
			templ: testFullTemplate,
			stract: struct {
				LabelTest       string `form:"label=this is custom label"`
				NameTest        string `form:"name=full_name"`
				TypeTest        string `form:"type=number"`
				PlaceholderTest string `form:"placeholder=this is placeholder..."`
				Nested          struct {
					MultiTest string `form:"name=NestedMultiTest;label=This is a nested label;type=email;placeholder=example@gmail.com"`
				}
			}{
				PlaceholderTest: "placeholding",
			},
			want: "TestHTML_custom_tags.golden",
		},
		"form with errors": {
			templ: testErrorTemplate,
			stract: struct {
				Email    string `form:"label=Email Address;placeholder=your@domain.com;
												 type=email;name=EmailAddress"`
				Password string `form:'type=password'`
			}{
				Email:    "jondoe@gmail.com",
				Password: "admin",
			},
			errors: []ErrorField{
				{
					Field: "Email",
					Error: "Email address has already been taken",
				},
				{
					Field: "Password",
					Error: "Password must be between 10 and 200 characters long",
				},
				{
					Field: "Password",
					Error: "Password must contain a latin letter",
				},
				{
					Field: "Password",
					Error: "Password must be a palindrome",
				},
				{
					Field: "Password",
					Error: "Password must contain an emoji",
				},
			},
			want: "TestHTML_errors.golden",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			got, gotErr := HTML(testCase.templ, testCase.stract, testCase.errors...)
			if gotErr != testCase.wantErr {
				t.Fatalf("HTML() err = %v; want %v;", gotErr, testCase.wantErr)
			}
			gotFileName := strings.Replace(testCase.want, ".golden", ".got", 1)
			// remove .got file if exists
			os.Remove(gotFileName)

			if updateFlag {
				writeFile(t, testCase.want, string(got))
				t.Logf("updated golden file %s", testCase.want)
			}

			want := template.HTML(readFile(t, testCase.want))
			if got != want {
				t.Error("HTML() results does not match golden file.")
				writeFile(t, gotFileName, string(got))
				t.Errorf(" To compare run: diff %s %s", gotFileName, testCase.want)
			}
		})
	}

}

func readFile(t *testing.T, filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("Error opening file %s: %v", filename, err)
	}
	defer file.Close()

	baites, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("error reading file %s: %v", filename, err)
	}

	return baites
}

func writeFile(t *testing.T, filename, contents string) {
	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Error creating file %s: %v", filename, err)
	}
	defer file.Close()

	fmt.Fprint(file, contents)
}
