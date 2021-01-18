package repo

import (
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

// Clone performs a git clone on the upstream url repository.
func Clone(url string) (err error) {
	location := filepath.Join(".binoc/sources/", filepath.Base(url))
	_, err = git.PlainClone(location, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	return err
}
