package option

import (
	"log/slog"
	"regexp"

	"github.com/lakrizz/logsync/internal/config"
)

var ()

type InternalLinkRemover struct {
}

func (r *InternalLinkRemover) IsEnabled(opts *config.Options) (bool, error) {
	return opts.RemoveInternalLinks, nil
}

func (r *InternalLinkRemover) Apply(input string) (string, error) {
	slog.Info("applying internal link remover")
	// now iterate through all links, check if they're internal
	// redirect them, etc.
	// and hit next recursion level
	regex := regexp.MustCompile(`\[\[([^\]]+)\]\]`)
	output := regex.ReplaceAllString(input, "$1")

	return output, nil
}
