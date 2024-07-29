package github

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"

	"github.com/google/go-github/v60/github"
)

// SetWebhook function deletes all webhooks and adds the given url as the sole webhook
func SetWebhook(ctx context.Context, github_token, target_url, logseq_repository_url string) error {
	client := github.NewClient(nil).WithAuthToken(github_token)

	github_user, github_repository, err := extractUsernameAndRepo(logseq_repository_url)
	if err != nil {
		return err
	}

	repos, _, err := client.Repositories.ListHooks(ctx, github_user, github_repository, &github.ListOptions{})
	if err != nil {
		return err
	}

	for _, v := range repos {
		_, err := client.Repositories.DeleteHook(ctx, github_user, github_repository, *v.ID)
		if err != nil {
			return err
		}
	}

	hook_config := &github.Hook{
		Name: github.String("web"),
		Config: &github.HookConfig{
			ContentType: github.String("json"),
			URL:         github.String(target_url),
		},
		Events: []string{"push"},
		Active: github.Bool(true),
	}

	hook, _, err := client.Repositories.CreateHook(ctx, github_user, github_repository, hook_config)
	if err != nil {
		return err
	}

	slog.Debug("new hook created", "id", hook.ID)

	return nil

}

// ExtractUsernameAndRepo extracts the username and repository name from a given URL.
// TODO: refactor
func extractUsernameAndRepo(url string) (string, string, error) {
	// define the regex pattern
	pattern := `(?:git@github\.com:|https://github\.com/)([^/]+)/([^/]+)\.git`

	re := regexp.MustCompile(pattern)
	// find matches
	matches := re.FindStringSubmatch(url)
	if len(matches) < 3 { 
		return "", "", fmt.Errorf("no match found")
	}

	// return the username and repository name
	return matches[1], matches[2], nil
}
