package mapping

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/lakrizz/logsync/internal/config"
)

type LogseqPage struct {
	logger        *slog.Logger
	InputContent  string
	ParsedContent string
	InputFilename string
}

func ParsePage(log *slog.Logger, filename string, mapping *config.Mapping) (*LogseqPage, error) {
	_, fn := filepath.Split(filename)
	l := &LogseqPage{InputContent: "", ParsedContent: "", logger: log, InputFilename: fn}
	err := l.readFile(filename)
	if err != nil {
		return nil, err
	}

	err = l.parseOptions(mapping.Options)
	if err != nil {
		return nil, err
	}

	err = l.addFrontMatter(mapping)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (l *LogseqPage) readFile(path string) error {
	dat, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	l.InputContent = string(dat)
	return nil
}

func (l *LogseqPage) Save(path string) error {
	return os.WriteFile(path, []byte(l.ParsedContent), 0777)
}
