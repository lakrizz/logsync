package option

import (
	"strings"

	"github.com/lakrizz/logsync/internal/config"
)

// RemoveEmptyTrails removes the first level of the input file, e.g., to convert bullet points to paragraphs
type UnindentFirstLevel struct {
}

func (r *UnindentFirstLevel) IsEnabled(opts *config.Options) (bool, error) {
	return opts.UnindentFirstLevel, nil
}

func (r *UnindentFirstLevel) Apply(input string) (string, error) {
	// here we shall start at the end of the newline splitted string and move up each line until
	// it is not empty (e.g., "- " or "-" or "" or " "), then return the result
	lines := strings.Split(input, "\n")
	for i := len(lines) - 1; i > 0; i-- {
		// TODO: pretty sure there's nicer ways to archive this
		if lines[i] != "- " && lines[i] != "-" && lines[i] != " " && lines[i] != "" {
			return strings.Join(lines[:i+1], "\n"), nil
		}
	}

	return input, nil
}
