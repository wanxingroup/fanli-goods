package validator

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/shomali11/util/xstrings"
)

func init() {
	//修改提示信息，not blank留给NotBlank Rule
	validation.Required = validation.Required.Error("cannot be empty")
	validation.NilOrNotEmpty = validation.NilOrNotEmpty.Error("cannot be empty")
}

//组合AND
func And(rules ...validation.Rule) AndRule {
	return AndRule{rules: rules}
}

type AndRule struct {
	rules []validation.Rule
}

func (r AndRule) Validate(value interface{}) error {
	for _, rule := range r.rules {
		if err := validation.Validate(value, rule); err != nil {
			return err
		}
	}
	return nil
}

//组合OR
func Or(rules ...validation.Rule) OrRule {
	return OrRule{rules: rules}
}

type OrRule struct {
	rules []validation.Rule
}

func (r OrRule) Validate(value interface{}) error {
	var err error
	for _, rule := range r.rules {
		e := validation.Validate(value, rule)
		if e == nil {
			return nil
		}
		if err == nil {
			err = e
		}
	}
	return err
}

//字符串not blank
var NotBlank = validation.NewStringRule(func(s string) bool {
	return xstrings.IsNotBlank(s)
}, "cannot be blank")
