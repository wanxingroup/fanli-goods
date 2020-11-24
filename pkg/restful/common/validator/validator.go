package validator

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"
	"reflect"
)

const ValidateParamError = "validate_param_error"

type Validatable interface {
	Validate() error
}

type ValidateError struct {
	Field   string
	Code    string
	Message string
}

func (es *ValidateError) Error() string {
	return es.Message
}

type ValidatorWrapper struct {
	Validators []Validatable
}

func WrapValidators(validators ...Validatable) *ValidatorWrapper {
	return &ValidatorWrapper{
		Validators: validators,
	}
}

func (vw *ValidatorWrapper) AddValidator(validators ...Validatable) *ValidatorWrapper {
	vw.Validators = append(vw.Validators, validators...)
	return vw
}

func (vw *ValidatorWrapper) Validate() error {
	for _, item := range vw.Validators {
		if item != nil {
			if err := item.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
}

type RuleValidator struct {
	Field    string
	FieldVal interface{}
	Rules    []validation.Rule
}

func (vr *RuleValidator) Validate() error {
	err := validation.Validate(vr.FieldVal, vr.Rules...)
	if err != nil {
		errmsg := fmt.Sprintf("%v: %v", vr.Field, err.Error())
		if errobj, ok := err.(validation.ErrorObject); ok {
			return &ValidateError{vr.Field, errobj.Code(), errmsg}
		} else {
			logrus.WithError(err).WithField("value", fmt.Sprintf("%v", vr.FieldVal)).Warn("validate request param failed")
			return &ValidateError{vr.Field, ValidateParamError, errmsg}
		}
	}
	return nil
}

func NewRuleValidator(field string, value interface{}, rules ...validation.Rule) Validatable {
	return &RuleValidator{field, value, rules}
}

//for nested object
func NewConditionValidator(condition bool, validatorNewFunc func() []Validatable) Validatable {
	if condition {
		return WrapValidators(validatorNewFunc()...)
	}
	return nullValidator(0)
}

//for nested array/map
func NewEachValidator(value interface{}, validatorNewFunc func(interface{}) []Validatable) Validatable {
	if value == nil {
		return nullValidator(0)
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Map:
		validators := make([]Validatable, 0, 5)
		for _, k := range v.MapKeys() {
			element := getInterface(v.MapIndex(k))
			validators = append(validators, validatorNewFunc(element)...)
		}
		return WrapValidators(validators...)
	case reflect.Slice, reflect.Array:
		validators := make([]Validatable, 0, 5)
		for i := 0; i < v.Len(); i++ {
			element := getInterface(v.Index(i))
			validators = append(validators, validatorNewFunc(element)...)
		}
		return WrapValidators(validators...)

	}
	return nullValidator(0)
}

func getInterface(value reflect.Value) interface{} {
	switch value.Kind() {
	case reflect.Ptr, reflect.Interface:
		if value.IsNil() {
			return nil
		}
		return value.Elem().Interface()
	default:
		return value.Interface()
	}
}

type nullValidator int

func (nullValidator) Validate() error {
	return nil
}
