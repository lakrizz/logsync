package config

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/adrg/xdg"
)

var (
	errConfigNotFound = errors.New("could not find config")
)

type Config struct {
	HugoExecutable string `json:"-"`
	Git            *struct {
		Token              string `json:"token"`
		LogseqRepoURL      string `json:"logseq_repo_url"`
		HugoRepoPath       string `json:"hugo_repo_path"`
		PrivateKeyPath     string `json:"private_key_path"`
		PrivateKeyPassword string `json:"private_key_password"`
		Username           string `json:"username"`
		Email              string `json:"email"`
	} `json:"git"`

	Ngrok *struct {
		AuthToken string `json:"auth_token"`
	} `json:"ngrok"`
	Mappings []*Mapping `json:"mappings"`
}

type Mapping struct {
	Options     *Options       `json:"options"`
	Frontmatter map[string]any `json:"frontmatter"`
	Source      string         `json:"source"`
	Target      string         `json:"target"`
}

type Options struct {
	RecursionTarget     string `json:"recursion_target,omitempty"`
	RecursionDepth      int    `json:"recursion_depth,omitempty"`
	Recursive           bool   `json:"recursive,omitempty"`
	RemoveInternalLinks bool   `json:"remove_internal_links,omitempty"`
	RecursionSkipSource bool   `json:"recursion_skip_source,omitempty"`
	RemoveEmptyTrails   bool   `json:"remove_empty_trails,omitempty"`
	UnindentFirstLevel  bool   `json:"unindent_first_level,omitempty"`
	IncludeAttachments  bool   `json:"include_attachments,omitempty"`

	// these values should be available to all options but need no manual work
	LogseqRepositoryPath string `json:"-"`
	HugoRepositoryPath   string `json:"-"`
}

func Load() (*Config, error) {
	filename := filepath.Join(xdg.ConfigHome, "logsync", "config.json")

	dat, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errConfigNotFound
		}
		return nil, err
	}
	cfg := &Config{}

	err = json.Unmarshal(dat, cfg)
	if err != nil {
		return nil, err
	}

	err = cfg.trySetHugoExecutable()
	if err != nil {
		return nil, err
	}

	err = cfg.SetStaticValuesForAllOptions()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) trySetHugoExecutable() error {
	path, err := exec.LookPath("hugo")
	if err != nil {
		return err
	}

	slog.Info("found hugo executable", "path", path)

	c.HugoExecutable = path
	return nil
}

func (c *Config) IsValid() (bool, error) {
	errs := make([]error, 8) // we have 8 required config fields, might as well reserve the mem now

	if c.HugoExecutable == "" {
		errs = append(errs, errors.New("hugo executable not found"))
	}

	if c.Git.Token == "" {
		errs = append(errs, errors.New("github token not"))
	}

	if c.Git.HugoRepoPath == "" {
		errs = append(errs, errors.New("hugo repository path not given"))
	}

	if c.Git.Username == "" {
		errs = append(errs, errors.New("git username not set"))
	}

	if c.Git.Email == "" {
		errs = append(errs, errors.New("git email not set"))
	}

	if c.Git.LogseqRepoURL == "" {
		errs = append(errs, errors.New("logseq repository url not set"))
	}

	if c.Git.PrivateKeyPath == "" {
		errs = append(errs, errors.New("private key path not set"))
	}

	if len(errs) == 0 {
		return false, errors.Join(errs...)
	}
	return true, nil
}

func (c *Config) SetStaticValuesForAllOptions() error {
	for _, mapping := range c.Mappings {
		mapping.Options.HugoRepositoryPath = c.Git.HugoRepoPath
	}

	return nil
}

func (c *Config) SetLogseqRepositoryPath(path string) {
	for _, mapping := range c.Mappings {
		mapping.Options.LogseqRepositoryPath = path
	}
}
