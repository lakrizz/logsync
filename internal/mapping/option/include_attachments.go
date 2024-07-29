package option

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lakrizz/logsync/internal/config"
)

var ()

type IncludeAttachments struct {
	LogseqRepositoryPath string
	HugoRepositoryPath   string
}

func (r *IncludeAttachments) IsEnabled(opts *config.Options) (bool, error) {
	r.HugoRepositoryPath = opts.HugoRepositoryPath
	r.LogseqRepositoryPath = opts.LogseqRepositoryPath

	return opts.IncludeAttachments, nil
}

func (r *IncludeAttachments) Apply(input string) (string, error) {
	slog.Info("applying attachment copier")
	// now iterate through all links, check if they're internal
	// redirect them, etc.
	// and hit next recursion level
	// Regular expression to match image URLs in Markdown
	re := regexp.MustCompile(`!\[[^\]]*\]\(([^)]+)\)`)

	// Find all matches
	matches := re.FindAllStringSubmatch(input, -1)
	targetDirectory := filepath.Join(r.HugoRepositoryPath, "static")

	err := r.createFolder(targetDirectory)
	if err != nil {
		return "", err
	}

	// Extract URLs from matches
	for _, match := range matches {
		if len(match) > 1 {
			// check whether this is an external link
			if parsedURL, err := url.Parse(match[1]); err == nil {
				if parsedURL.Scheme != "" && (parsedURL.Host != "" || strings.Contains(parsedURL.Scheme, "file")) {
					slog.Info("this attachment is external, skipping", "url", match[1])
					continue // skip this match if it's a parseable url

				}
			}

			// copy this file, which consists of the logseq repo path and match[1]
			_, pureFilename := filepath.Split(match[1])
			srcFile := filepath.Join(r.LogseqRepositoryPath, "assets", pureFilename)
			targetFile := filepath.Join(targetDirectory, pureFilename)

			// figure out of the file actually exists
			if _, err := os.Stat(srcFile); os.IsNotExist(err) {
				slog.Info("[include attachments option] source file not found", "source_file", srcFile)
				return input, err
			}

			srcData, err := os.ReadFile(srcFile)
			if err != nil {
				return "", err
			}

			err = os.WriteFile(targetFile, srcData, 0777)
			if err != nil {
				return "", err
			}

			// now we need to rewrite the source url for hugo, a.k.a. removing the subfolder prefix
			input = strings.ReplaceAll(input, match[1], fmt.Sprintf("/%v", pureFilename))
		}
	}

	// this option does not modify the page, thus we'll just return the input
	return input, nil
}

func (r *IncludeAttachments) createFolder(folder string) error {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		slog.Info("target folder does not exist")
		return os.MkdirAll(folder, 0777)
	}
	return nil
}
