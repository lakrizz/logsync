package mapping

import (
	"github.com/lakrizz/logsync/internal/config"
	"github.com/lakrizz/logsync/internal/mapping/option"
)

type opt interface {
	IsEnabled(*config.Options) (bool, error)
	Apply(string) (string, error)
}

func (l *LogseqPage) getOptions() []opt {
	return []opt{
		&option.Recursion{Depth: 1, SkipSource: false}, // default values
		&option.IncludeAttachments{},
		&option.RemoveEmptyTrails{},
		&option.UnindentFirstLevel{},
		&option.IncludeAttachments{},
	}
}

func (l *LogseqPage) parseOptions(mapping *config.Options) error {
	output := l.InputContent
	for _, v := range l.getOptions() {
		enabled, err := v.IsEnabled(mapping)
		if !enabled {
			continue
		}
		if err != nil {
			return err
		}

		optionContent, err := v.Apply(output)
		if err != nil {
			return err
		}

		output = optionContent
	}

	l.ParsedContent = output
	return nil
}
