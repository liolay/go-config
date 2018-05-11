package main

import (
	"gopkg.in/src-d/go-git.v4"
	"log"
	"os"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"fmt"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func Clone(home string, sshKey []byte, repoUrl string) *git.Repository {
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
		log.Fatal("can't fetch repo from git server", err)
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


func main() {

//	idRsa := []byte(`-----BEGIN RSA PRIVATE KEY-----
//MIIEpAIBAAKCAQEAyB+gq1xjj4HCe4hPLfA3a9y4pXVhrJfQQheHR1Mj9Y2eBJTm
//4jKtHwYWOxVbUahmim+vUuVia5nH5LYA9c/GZxUKrhHwND0hF7q7kKRMboT2w/2J
//JRPkVC3T2I02ptADTyFCZLNUViviFm6JoVEnEZAwlBDakvBQBqYgQstrSE2mH8wN
//niM7T7U6gPChpgPVpuTOBNbsgfag7bgAUVl56i3IFvht+N0LXKuNU0TTuMUq2+lR
//i+iha86OtkSYZqA34SXpoL2cyQJUikMkPsjXsxfpq1f35AbPC4FUfNnnN/NRBD/2
//jOHHQSHcMOq0Y/SapCXgxa4urlLkCOra1lC5hQIDAQABAoIBAQCeK5E7nzv5cp+a
//L3QVZOUI1V0DOTFHzn2Fnz8GeonTTGj2ShHp+g+mk5MCg7C3a5gQFpHFvRL65IJ/
//G/LKVbwEQTc9uWPWhfIf5TDV82WNfH3lDgBVU9GFTus/Hu1xDrtu0WS+XpZrvSdm
//f1s8Kv3r/cDHZkK7HEDD4I1i/Y//hinEOFbZglgmgGCwF3aZUdfj0wq1t1CnVJsj
//P3c+rbxtvBLPblizHunqMiNHM1DTYqK5sW92yyMvy6ZbN3fSntcEdUNLWl61uWh1
//E0eeg+GSuQoY/STFNqeTCqPju+gW8YR+Qix1R2n/AMo/D3HXcKxZCxWSD5zY8bH0
//C/5oZerBAoGBAOdLjhQb2L1xU0V7vHWIHsWeOOTaDu1IhratY+detY/D0VCbnTXx
//CZQmCgwT88/iPLVbEg3tsBBlVRYdTSOwtGk4Dvis83Izs1kU4RV8AI3+pZ74qA6O
//8V3xrvaIrnJWHX0M078EAOFnj9zvkMFQnovnSyli+GSzPAreHVAmOMANAoGBAN1/
//uojTaFg+cCfD6mQACo8DMN2CbWxJjXo4cpGv2HSYcjrVK+knjCdFhEnquLa0lQ+u
//TxvNTFDhPDnEXFUegZwD+ge2oDI7wP/qAJUpgU53W12OvaRs5gMHKyU8AHWrMgYk
//PvEwZIOlEc7ceyWzKCyt7nZwPHiteRsv+QgSS4lZAoGAUH27sQXL1ImWmAyqliBL
//zSv10raMEUl3ECWhKciM2L4lnq649Cew1Ky0PGXJKGQsClTqIIzCA8Kv7KU/zhbV
//gfRvSV0uz2RsmqioeAiSTNf8nSkdmwtltfLAl60TQFj1pCoNmmDzSX3308RPFOdQ
//dZGFV57IoIq7b3DCtLzIbRUCgYB8ZFAYqUlPTXllC6SlllRXrn4R2D6lcsUuX2cQ
//JEYWbMqx+aeYX+pY37SEYnpruQyBau3oeioives5sen8r44wVRdkn45lx6MC1aKQ
//ImgI7gT0jMY6AiJGjw8O8Rx8+LC2PELQ5tF8EQboOnA6YtvsA54JC80aJKn/t7hO
//bR/YuQKBgQDik5Es/Y2NM0JSKQwy7myGP4g5SadZd5Pv4r89OY/8mEOhKvHlaxki
//M2o2is6RiFTKn/aFWXl0/QKXyupq35jsesq78y29zoUuph7J5PYcLgFF2t49q3B4
//Pn/53pMethW8zhphp3fM77vBkUu/NcusttbQjtdLEI7dkR/kqEc4fQ==
//-----END RSA PRIVATE KEY-----`)
//
//Clone("/Users/liolay/aaa",idRsa,"git@github.com:liolay/go-config.git")
	repo, _ := git.PlainOpen("/Users/liolay/aaa")
	//reference, _ := repo.Reference("refs/remotes/origin/master",true)
	//fmt.Print(reference.Hash())
	//fmt.Print(":"+reference.Name())
	//referenceIter, _ := repo.References()
	//referenceIter.ForEach(func(reference *plumbing.Reference) error {
	//	fmt.Print(reference.Target().IsNote())
	//	fmt.Print(":"+reference.Name())
	//	fmt.Println("："+reference.Hash().String())
	//	return nil
	//})

	
	//iter, _ := repo.Branches()
	//_ = iter.ForEach(func(reference *plumbing.Reference) error {
	//	fmt.Print(reference.Name())
	//	fmt.Println("："+reference.Hash().String())
	//	return nil
	//})

	//branch, err := repo.Branch("test")

	//CheckOut(repo,"547bda7cfd54748cbb77e8517e420e126e82524b")

	//ref, _ := repo.Head()
	commit, _ := repo.CommitObject(plumbing.NewHash("1645c1431b46e7d5d6838cbea0af386b3b5c7f88"))

	tree, _ := commit.Tree()
	tree.Files().ForEach(func(f *object.File) error {

		fmt.Printf("100644 blob %s    %s\n", f.Hash, f.Name)
		return nil
	})

	//ref, _ := r.Head()
	//commit, _ := r.CommitObject(ref.Hash())
	//tree, _ := commit.Tree()
	//tree.Files().ForEach(func(f *object.File) error {
	//	fmt.Printf("100644 blob %s    %s\n", f.Hash, f.Name)
	//	return nil
	//})
	//commit.Files().ForEach(func(f *object.File) error {
	//	fmt.Printf("100644 blob %s    %s\n", f.Hash, f.Name)
	//	return nil
	//})
	//println(reference.Hash().String())
	//CheckOut(repo,"db90ee80e0684f87b75666f1651e37caa8685aff")

}
