package hugo

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/lakrizz/logsync/internal/config"
	"github.com/lakrizz/logsync/internal/mapping"
)

func HandleModifiedPages(files []string, cfg *config.Config, logseq_repo, hugo_repo *git.Repository) error {
	logseqWorktree, err := logseq_repo.Worktree()
	if err != nil {
		return err
	}

	hugoWorktree, err := hugo_repo.Worktree()
	if err != nil {
		return err
	}

	for _, file := range files {
		i := slices.IndexFunc(cfg.Mappings, func(s *config.Mapping) bool {
			return s.Source == file
		})
		if i == -1 {
			continue
		}

		target := cfg.Mappings[i]

		sourceFS, err := logseqWorktree.Filesystem.Open(file)
		if err != nil {
			return err
		}

		targetFile, err := hugoWorktree.Filesystem.Create(target.Target)
		if err != nil {
			return err
		}

		slog.Info("parsing logseq file...", "filename", logseqWorktree.Filesystem.Root())
		// here we need to convert the logseq pages to hugo pages
		// by adding frontmatter, etc.
		parsedPage, err := mapping.ParsePage(slog.Default(), filepath.Join(logseqWorktree.Filesystem.Root(), sourceFS.Name()), target)

		if err != nil {
			return err
		}

		slog.Info("copying new hugo file")
		err = parsedPage.Save(filepath.Join(hugoWorktree.Filesystem.Root(), targetFile.Name()))
		if err != nil {
			return fmt.Errorf("cannot copy file: %w", err)
		}

		slog.Info("adding file to index")
		_, err = hugoWorktree.Add(targetFile.Name())
		if err != nil {
			return fmt.Errorf("cannot add to worktree: %w", err)
		}
	}

	// now we need to create a commit and push import
	commitMessage := fmt.Sprintf("logsync autocommit %v / files: %v", time.Now().Format(time.RFC3339), strings.Join(files, ","))
	_, err = hugoWorktree.Commit(commitMessage, &git.CommitOptions{Author: &object.Signature{Name: cfg.Git.Username, Email: cfg.Git.Email, When: time.Now()}})
	if err != nil {
		return errors.Join(errors.New("cannot commit"), err)
	}

	// now we copied the files and created and pushed the commits
	// problem here is that we cloned the hugo_repo ourselves
	os.Chdir(hugoWorktree.Filesystem.Root())
	cmd := exec.Command(cfg.HugoExecutable)

	_, err = cmd.CombinedOutput()
	if err != nil {
		return errors.Join(errors.New("cannot execute hugo"), err)
	}

	return nil
}
