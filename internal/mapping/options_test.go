package mapping_test

import (
	"log"
	"log/slog"
	"testing"

	"github.com/lakrizz/logsync/internal/mapping/option"
)

func TestRemoveEmptyTrails(t *testing.T) {
	t.Helper()

	s := `- irgendwie nen "text ist gehighlighted" effekt einbauen
- ![image.png](./input.md)
- so aber anders, also so in gut :D
-`
	ret := &option.IncludeAttachments{LogseqRepositoryPath: "/home/krizz/src/krizz.org/logsync/test", HugoRepositoryPath: "/home/krizz/src/krizz.org/logsync/test"}
	res, err := ret.Apply(s)

	if err != nil {
		slog.Error("error with option", "error", err)
	}
	log.Println(res, err)
}
