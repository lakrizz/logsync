package option

import (
	"errors"

	"github.com/lakrizz/logsync/internal/config"
)

var (
	errNoRecursionTarget = errors.New("[recursion] no recursion path given, please add 'recursion_target' to your mapping options")
)

type Recursion struct {
	Target     string
	Depth      int
	SkipSource bool
}

func (r *Recursion) IsEnabled(opts *config.Options) (bool, error) {
	if !opts.Recursive {
		return false, nil
	}

	if opts.RecursionTarget == "" {
		return false, errNoRecursionTarget
	}

	return true, nil
}

func (r *Recursion) Apply(input string) (string, error) {
	// now iterate through all links, check if they're internal
	// redirect them, etc.
	// and hit next recursion level
	return input, nil
}
