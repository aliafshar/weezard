// Package weezard asks the user questions on standard IO and stores answers.
//
// Quesiton metadata is provided as struct field tags, similar to how they are
// used by `json`, etc.
package weezard

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/template"
)

// Template contains the template used when displaying question prompts.
//
// The variable that is passed on execution is a *Question
var Template = "{{.Usage|.Bold}} (default={{.Default|.Blue}}) > "

// Question is a single promptable unit.
type Question struct {

  // Usage is the content of the question.
	Usage   string

  // Default is the default answer.
	Default string

  // Set is called with an answer.
  Set     func(string)
}

// Bold generates ansi-escaped bold text.
func (q *Question) Bold(str string) string {
	return "\033[1m" + str + "\033[0m"
}

// Blue generates ansi-escaped blue text.
func (q *Question) Blue(str string) string {
	return q.Bold("\033[34m" + str + "\033[0m")
}

// parseTag extracts required metadata from a field's struct tag
func parseTag(tag string) (*Question, error) {
	var ts []string
	if tag == "" {
		ts = []string{"", ""}
	} else {
		ts = strings.SplitN(tag, ",", 2)
		if len(ts) < 2 {
			return nil, errors.New("Must provide <default>,<question>")
		}
	}
	q := &Question{Usage: ts[1], Default: ts[0]}
	return q, nil
}

// newQuestion creates a question from a struct/value field pair.
func newQuestion(tfield reflect.StructField, setter func(string)) (*Question, error) {
	tag := tfield.Tag.Get("question")
	q, err := parseTag(tag)
	if err != nil {
		return nil, err
	}
	q.Set = setter
	// For fields with no tag
	if q.Usage == "" {
		q.Usage = tfield.Name + "?"
	}
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

// QuestionsFor builds a list of questions from a non-nil struct pointer.
func QuestionsFor(s interface{}) ([]*Question, error) {
	qs := []*Question{}
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
		q, err := newQuestion(tfield, vfield.SetString)
		if err != nil {
			return qs, err
		}
		qs = append(qs, q)
	}
	return qs, nil
}

// AskQuestion prompts the user for a single question, and calls the provided
// setter, and returns the answer.
func AskQuestion(q *Question) (string, error) {
	var s string
	for s == "" {
		tmpl, err := template.New("question").Parse(Template)
		if err != nil {
			return "", err
		}
		tmpl.Execute(os.Stdout, q)
		_, err = fmt.Scanln(&s)
		if err != nil && err.Error() != "unexpected newline" {
			return "", err
		}
		if s == "" {
			s = q.Default
		}
    if q.Set != nil {
      q.Set(s)
    }
	}
	return s, nil
}

// AskQuestions prompts the user for all passed questions, and sets the
// appropriate values from the answers.
func AskQuestions(qs []*Question) error {
	for _, q := range qs {
		_, err := AskQuestion(q)
		if err != nil {
			return err
		}
	}
  return nil
}

// Ask prompts the user for all answers in the given struct, and sets the
// appropriate values from the answers.
//
// v must be a non-nil pointer to a struct.
func Ask(v interface{}) error {
	qs, err := QuestionsFor(v)
	if err != nil {
		return err
	}
  return AskQuestions(qs)
}
