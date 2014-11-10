// Package weezard asks the user questions on the command line and takes
// answers.
package weezard

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"text/template"
)

var Template = "{{.Usage}} {{.Name}} [default={{.Default}}]"

// parseTag extracts required metadata from a field's struct tag
func parseTag(tag string) (*question, error) {
	var ts []string
	if tag == "" {
		ts = []string{"", ""}
	} else {
		ts = strings.SplitN(tag, ",", 2)
		if len(ts) < 2 {
			return nil, errors.New("Must provide <default>,<question>")
		}
	}
	q := &question{Usage: ts[1], Default: ts[0]}
	return q, nil
}

// newQuestion creates a question from a struct/value field pair.
func newQuestion(tfield reflect.StructField, vfield reflect.Value) (*question, error) {
	tag := tfield.Tag.Get("question")
	q, err := parseTag(tag)
	if err != nil {
		return nil, err
	}
	q.field = vfield
	q.Name = tfield.Name
	return q, nil
}

// reflectIfPointer validates a correct non-nil pointer type and reflects it.
func reflectIfPointer(s interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return v, errors.New("Non-nil pointer required.")
	}
	return v, nil
}

// newQuestions builds a list of questions from a non-nil struct pointer.
func newQuestions(s interface{}) ([]*question, error) {
	qs := []*question{}
	pv, err := reflectIfPointer(s)
	if err != nil {
		return nil, err
	}
	v := pv.Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		vfield := v.Field(i)
		tfield := t.Field(i)
		if tfield.Name == "" {
			// Private fields
			continue
		}
		q, err := newQuestion(tfield, vfield)
		if err != nil {
			return qs, err
		}
		qs = append(qs, q)
	}
	return qs, nil
}

type question struct {
	Name    string
	Usage   string
	Default string
	field   reflect.Value
}

func (q *question) set(v string) {
	q.field.SetString(v)
}

// bold generates ansi-escaped bold text.
func bold(str string) string {
	return "\033[1m" + str + "\033[0m"
}

// blue generates ansi-escaped blue text.
func blue(str string) string {
	return bold("\033[34m" + str + "\033[0m")
}

func askQuestion(q *question) string {
	var s string
	for s == "" {
		// fmt.Printf("%v %v [default=%v]>> ", q.Usage, bold(q.Name), blue(q.Default))
		tmpl, err := template.New("question").Parse(Template)
		if err != nil {
			log.Fatalln("Bad template")
		}
		tmpl.Execute(os.Stdout, q)
		_, err = fmt.Scanln(&s)
		if err != nil && err.Error() != "unexpected newline" {
			log.Fatalln("Bad scan", err)
		}
		if s == "" {
			s = q.Default
		}
	}
	return s
}

func Ask(s interface{}) error {
	qs, err := newQuestions(s)
	if err != nil {
		return err
	}
	for _, q := range qs {
		a := askQuestion(q)
		q.set(a)
	}
	return nil
}
