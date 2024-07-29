package main

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"

	"github.com/lakrizz/logsync/internal/config"
	"github.com/lakrizz/logsync/internal/git"
	"github.com/lakrizz/logsync/internal/github"
	"github.com/lakrizz/logsync/internal/hugo"
	"github.com/lakrizz/logsync/internal/ngrok"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}))
	cfg, err := config.Load()
	if err != nil {
		log.Error("error loading config", "error", "error", err)
		return
	}

	if valid, errs := cfg.IsValid(); !valid {
		log.Error("invalid config", "hints", errs)
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Error("error getting working directory", "error", err)
		return
	}

	pwd := path.Join(wd, "git", "logseq")
	logseqRepo, err := git.CloneOrOpen(pwd, cfg.Git.LogseqRepoURL, cfg.Git.PrivateKeyPath, cfg.Git.PrivateKeyPassword)
	if err != nil {
		log.Error("error opening/cloning logseq repository", "error", err)
		return
	}
	cfg.SetLogseqRepositoryPath(pwd)
	log.Info("logseq repository opened", "path", pwd)

	hugoRepo, err := git.Open(cfg.Git.HugoRepoPath)
	if err != nil {
		log.Error("error opening hugo repository", "error", err)
		return
	}
	log.Info("hugo repository opened")

	// now we cloned the repository
	// we want to open the reverse proxy (in this case ngrok)
	ctx := context.Background()
	targetURL, errChan, err := ngrok.Start(ctx, cfg, func(w http.ResponseWriter, r *http.Request) {
		log.Info("received webhook call")
		if r == nil {
			log.Error("request is nil, this should not happen")
			return
		}
		// 1. decode request body
		if r.Body == nil {
			log.Info("github webhook request has no body, skipping")
			return
		}
		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error("error while reading github webhook request body", "error", err)
			return
		}

		obj := &github.PushPayload{}
		err = json.Unmarshal(data, &obj)
		if err != nil {
			// handle the initial ping from github (see: https://docs.github.com/en/webhooks/webhook-events-and-payloads#ping)
			pingObject := &github.PingPayload{}
			if err = json.Unmarshal(data, &pingObject); err == nil {
				log.Info("✓ github webhook ping successfully handled")
				w.WriteHeader(http.StatusOK)
				return
			}

			log.Error("error while unmarshalling github webhook payload into push struct", "error", err)
			return
		}

		log.Info("received push in logseq repository")

		// 2. check if any of the changes fit to any mapping
		changedFiles := make([]string, 0)
		for _, commit := range obj.Commits {
			for _, file := range commit.Modified {
				for _, mapping := range cfg.Mappings {
					if mapping.Source == file {
						changedFiles = append(changedFiles, file)
					}
				}
			}
		}

		if len(changedFiles) == 0 {
			log.Info("no mapped files are part of this push")
			return
		}

		// 3. refresh repository and fetch all mapped pages (with depth n)
		err = git.Pull(logseqRepo, cfg.Git.PrivateKeyPath, cfg.Git.PrivateKeyPassword)
		if err != nil {
			log.Error("error while pulling the logseq repo", "error", err)
			return
		}

		// 4. convert pages to hugo
		_, err = logseqRepo.Worktree()
		if err != nil {
			log.Error("error fetching the git worktree", "error", err)
			return
		}

		// pull current hugo-repo state to prevent non-fast-foward updates
		err = git.Pull(hugoRepo, cfg.Git.PrivateKeyPath, cfg.Git.PrivateKeyPassword)
		if err != nil {
			log.Error("error pulling hugo repo", "error", err)
			return
		}
		log.Info("successfully pushed changes from hugo repository")

		// send all new and changed files to the hugo function
		err = hugo.HandleModifiedPages(changedFiles, cfg, logseqRepo, hugoRepo)
		if err != nil {
			log.Error("error handling modified pages", "error", err)
			git.Reset(hugoRepo)
			return
		}
		log.Info("successfully created updated pages for hugo")

		// before pushing, we need to pull all changes to the hugo repo
		// since this tool might not be the only thing that changes
		// hugo :D

		err = git.Push(hugoRepo, cfg.Git.PrivateKeyPath, cfg.Git.PrivateKeyPassword)
		if err != nil {
			log.Error("error pushing changeset", "error", err)
			return
		}
		log.Info("successfully pushed changes to hugo repository")

	})

	log.Info("started ngrok", "ngrok_url", targetURL)

	if err != nil {
		return
	}

	// now add this url as the webhook url in the repo
	err = github.SetWebhook(ctx, cfg.Git.Token, targetURL, cfg.Git.LogseqRepoURL)
	if err != nil {
		log.Error("error setting webhook", "error", err)
		return
	}
	log.Info("github webhook added")

	<-errChan
}
