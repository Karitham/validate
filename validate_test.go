package validate

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fail()
		}
	}()

	validators = make(map[any]any)
	err := fmt.Errorf("expected a valid string, but got empty one")

	MustRegisterValidatorFunc(func(s string) error {
		if s != "" {
			return err
		}
		return nil
	})

	type test struct {
		input string
		err   error
	}

	tests := []test{
		{input: "", err: nil},
		{input: "test", err: err},
	}

	for _, tc := range tests {
		got := Validate(tc.input)
		assert.Equal(t, got, tc.err)
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
