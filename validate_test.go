package validate

import (
	"errors"
	"reflect"
	"testing"
)

func TestValidate(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fail()
		}
	}()

	validator := ValidatorFunc[string](func(s string) bool {
		return s == ""
	})

	validators = make(map[any]any)

	MustRegisterValidator[string](validator)

	type test struct {
		input string
		want  bool
	}

	tests := []test{
		{input: "", want: true},
		{input: "test", want: false},
	}

	for _, tc := range tests {
		got := Validate(tc.input)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}

func TestValidate_PanicsWhenNoValidatorRegistered(t *testing.T) {
	expected := errors.New("no validator for type string is registered")
	defer func() {
		r := recover()

		got, ok := r.(error)

		if !ok {
			t.Fatalf("recovered value was not an error, got: %v", r)
		}

		if got.Error() != expected.Error() {
			t.Fatalf("expected: %v, got: %v", expected, got)
		}
	}()

	validators = make(map[any]any)

	Validate("")

	t.Fail()
}
