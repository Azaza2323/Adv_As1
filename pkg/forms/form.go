package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"unicode/utf8"
)

type Form struct {
	Values url.Values
	Errors map[string]string
}

var EmailRX = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func (f *Form) MinLength(field string, minLen int) {
	value := f.Values.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < minLen {
		f.Errors[field] = fmt.Sprintf("Field '%s' is too short (minimum is %d characters)", field, minLen)
	}
}

func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp, errMsg string) {
	value := f.Values.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors[field] = errMsg
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
