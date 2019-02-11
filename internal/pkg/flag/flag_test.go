package flag

import (
	"testing"

	"github.com/pkg/errors"
)

func TestParse(t *testing.T) {
	var fl FlagSet

	err := fl.Parse([]string{"-unknown"})

	if err == nil {
		t.Error("expected error; got none", err, errors.Cause(err))
	}

	if err != nil && errors.Cause(err) != ErrFlagNotDefined {
		t.Error("expected error; got none", err, errors.Cause(err))
	}
}
