package util

import (
	"log"
	"gopkg.in/src-d/go-git.v4"
	"os"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"strings"
)

func Clone(home string, sshKey []byte, repoUrl string) *git.Repository {
	print(home,":"+repoUrl)
	keys, err := gitssh.NewPublicKeys("git", sshKey, "")
	if err != nil {
		log.Println("can't create PublicKeys", err)
		return nil
	}

	repo, err := git.PlainClone(home, false, &git.CloneOptions{
		URL:               repoUrl,
		Auth:              keys,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Progress:          os.Stdout,
	})

	if err != nil {
		log.Println("can't fetch repo from git server:",err)
		return nil
	}

	return repo

}

func OpenLocalRepo(repoPath string) *git.Repository {
	if repository, err := git.PlainOpen(repoPath); err == nil {
		return repository
	}
	return nil
}

func CheckOut(repo *git.Repository, label string) {
	if reference, err := repo.Head(); err == nil {
		log.Printf("git show-ref --head HEAD : %s", reference.Hash())
	}

	if worktree, err := repo.Worktree(); err == nil {
		log.Printf("git checkout %s", label)
		worktree.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(label),
		})
	}

	if reference, err := repo.Head(); err == nil {
		log.Printf("git show-ref --head HEAD : %s", reference.Hash())
	}
}

func FileIterator(repo *git.Repository, label string) *object.FileIter {
	if referenceIter, err := repo.References(); err == nil {
		var hash plumbing.Hash
		referenceIter.ForEach(func(reference *plumbing.Reference) error {
			if strings.HasSuffix(reference.Name().Short(), "/"+label) || reference.Name().Short() == label {
				hash = reference.Hash()
			}
			return nil
		})

		if len(hash) == 0 {
			hash = plumbing.NewHash(label)
		}

		commit, err := repo.CommitObject(hash);
		if err != nil {
			log.Fatal(err)
		}

		tree, err := commit.Tree();
		if err != nil {
			log.Fatal(err)
		}
		return tree.Files()
	}
	return nil
}
