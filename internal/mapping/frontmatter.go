package mapping

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gosimple/slug"

	"github.com/lakrizz/logsync/internal/config"
)

func (l *LogseqPage) addFrontMatter(mapping *config.Mapping) error {
	sb := strings.Builder{}

	// input filename without ext
	fileWithoutExtension := l.InputFilename[:strings.LastIndex(l.InputFilename, filepath.Ext(l.InputFilename))]

	// static frontmatter (e.g., date)
	mapping.Frontmatter["date"] = time.Now().Format(time.RFC3339)
	mapping.Frontmatter["title"] = slug.Make(fileWithoutExtension)

	for k, v := range mapping.Frontmatter {
		_, err := sb.WriteString(fmt.Sprintf("%s = '%v'\n", k, v))
		if err != nil {
			return err
		}
	}

	l.ParsedContent = fmt.Sprintf("+++\n%v+++\n%v", sb.String(), l.ParsedContent)
	return nil
}
