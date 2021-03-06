// SPDX-License-Identifier: MIT

package repo

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GetLatestTagFromRepository - Get the latest Tag reference from the repo
func GetLatestTagFromRepository(repository *git.Repository) (*plumbing.Reference, *plumbing.Reference, error) {
	tagRefs, err := repository.Tags()
	if err != nil {
		return nil, nil, err
	}

	var latestTagCommit *object.Commit
	var latestTagName *plumbing.Reference
	var previousTag *plumbing.Reference
	var previousTagReturn *plumbing.Reference

	err = tagRefs.ForEach(func(tagRef *plumbing.Reference) error {
		revision := plumbing.Revision(tagRef.Name().String())

		tagCommitHash, err := repository.ResolveRevision(revision)
		if err != nil {
			return err
		}

		commit, err := repository.CommitObject(*tagCommitHash)
		if err != nil {
			return err
		}

		if latestTagCommit == nil {
			latestTagCommit = commit
			latestTagName = tagRef
			previousTagReturn = previousTag
		}

		if commit.Committer.When.After(latestTagCommit.Committer.When) {
			latestTagCommit = commit
			latestTagName = tagRef
			previousTagReturn = previousTag
		}

		previousTag = tagRef

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return latestTagName, previousTagReturn, nil
}
