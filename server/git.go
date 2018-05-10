package main

import (
	"log"
	"gopkg.in/src-d/go-git.v4"
	"os"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func clone(home string, sshKey []byte, repoUrl string) *git.Repository {
	keys, err := gitssh.NewPublicKeys("git", sshKey, "")
	if err != nil {
		log.Println("can't create PublicKeys", err)
		panic(err)
	}

	repo, err := git.PlainClone(home, false, &git.CloneOptions{
		URL:               repoUrl,
		Auth:              keys,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Progress:          os.Stdout,
	})

	if err != nil {
		log.Fatal("can't fetch repo from git server")
		return nil
	}

	return repo

}

func openLocalRepo(repoPath string) *git.Repository {
	if repository, err := git.PlainOpen(repoPath); err == nil {
		return repository
	}
	return nil
}

func checkOut(repo *git.Repository, label string) {
	if reference, err := repo.Head(); err == nil {
		log.Printf("git show-ref --head HEAD : %s",reference.Hash())
	}

	if worktree, err := repo.Worktree(); err == nil {
		log.Printf("git checkout %s", label)
		worktree.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(label),
		})
	}

	if reference, err := repo.Head(); err == nil {
		log.Printf("git show-ref --head HEAD : %s",reference.Hash())
	}
}
