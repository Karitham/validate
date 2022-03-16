package validate

import (
	"fmt"
)

type Validator[T any] interface {
	Validate(t T) bool
}

type validatorFunc[T any] func(t T) bool

func (vf validatorFunc[T]) Validate(t T) bool {
	return vf(t)
}

// We have to give up types here so this map can map any T to any Validator[T]
// As this isnt exposed we can be more certain that this map only contains expected values
var validators map[any]any = make(map[any]any)

func RegisterValidator[T any](validator Validator[T]) error {
	var _t T
	if _, ok := validators[_t]; ok {
		return fmt.Errorf("validator for type %T already registered", _t)
	}

	validators[_t] = validator
	return nil
}

func RegisterValidatorFunc[T any](validator func(t T) bool) error {
	return RegisterValidator[T](validatorFunc[T](validator))
}

func MustRegisterValidator[T any](validator Validator[T]) {
	err := RegisterValidator(validator)
	if err != nil {
		panic(err)
	}
}

func MustRegisterValidatorFunc[T any](validator func(t T) bool) {
	MustRegisterValidator[T](validatorFunc[T](validator))
}

func Validate[T any](t T) bool {
	var _t T

	found, ok := validators[_t]
	if !ok {
		panic(fmt.Errorf("no validator for type %T is registered", _t))
	}

	validator, ok := found.(Validator[T])

	if !ok {
		panic(fmt.Errorf("validator registered for type %T is not a validator of %T", _t, _t))
	}

	return validator.Validate(t)
}
