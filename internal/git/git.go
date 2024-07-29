package git

import (
	"errors"
	"os"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

// CloneOrOpen clones or opens a repository
// the url should be in a system-accessible format (e.g., https or git)
// the function returns the folder-name
func CloneOrOpen(dir, url, ssh_key_path, ssh_key_password string) (*git.Repository, error) {
	auth, err := sshPrivateKeyAuth(ssh_key_path, ssh_key_password)
	if err != nil {
		return nil, err
	}

	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Auth:              auth,
	})

	if err == git.ErrRepositoryAlreadyExists {
		return git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{})
	}

	if err != nil {
		return nil, err
	}

	return repo, nil
}

// Open clones or opens a repository
// the url should be in a system-accessible format (e.g., https or git)
// the function returns the folder-name
func Open(directory string) (*git.Repository, error) {
	return git.PlainOpenWithOptions(directory, &git.PlainOpenOptions{})
}

// getSSHAgentAuth
func sshPrivateKeyAuth(ssh_private_key_file, ssh_private_key_password string) (*ssh.PublicKeys, error) {
	_, err := os.Stat(ssh_private_key_file)
	if err != nil {
		return nil, err
	}

	// Clone the given repository to the given directory
	public_keys, err := ssh.NewPublicKeysFromFile("git", ssh_private_key_file, ssh_private_key_password)
	if err != nil {
		return nil, err
	}

	return public_keys, nil
}

func Pull(repo *git.Repository, ssh_key_path, ssh_key_password string) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	auth, err := sshPrivateKeyAuth(ssh_key_path, ssh_key_password)
	if err != nil {
		return err
	}

	err = Reset(repo)
	if err != nil {
		return err
	}

	err = worktree.Pull(&git.PullOptions{RemoteName: "origin", Auth: auth, Force: true})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}

	return nil
}

func Reset(repo *git.Repository) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = worktree.Reset(&git.ResetOptions{Mode: git.HardReset})
	if err != nil {
		return err
	}
	return nil
}

func Push(repo *git.Repository, ssh_key_path, ssh_key_password string) error {
	auth, err := sshPrivateKeyAuth(ssh_key_path, ssh_key_password)
	if err != nil {
		return err
	}

	err = repo.Push(&git.PushOptions{RemoteName: "origin", Auth: auth})
	if err != nil {
		return err
	}

	return nil
}
