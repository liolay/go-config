package main

import (
	"os"
	"flag"
	"gopkg.in/src-d/go-git.v4"
	"strings"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"log"
)

var gitHomePath = flag.String("gitHomePath", os.TempDir()+"gitHomePath", "where to store files which fetched by git client from remote")

func init() {
	flag.Parse()
	if _, err := os.Stat(*gitHomePath); os.IsNotExist(err) {
		os.MkdirAll(*gitHomePath, os.ModePerm)
	}
}

func clone(repoUrl string) {
	repoPath := strings.TrimSuffix(repoUrl, ".git")
	repoPath = *gitHomePath + repoPath[strings.LastIndex(repoPath, "/"):]
	os.RemoveAll(repoPath)
	idRsa := []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAyB+gq1xjj4HCe4hPLfA3a9y4pXVhrJfQQheHR1Mj9Y2eBJTm
4jKtHwYWOxVbUahmim+vUuVia5nH5LYA9c/GZxUKrhHwND0hF7q7kKRMboT2w/2J
JRPkVC3T2I02ptADTyFCZLNUViviFm6JoVEnEZAwlBDakvBQBqYgQstrSE2mH8wN
niM7T7U6gPChpgPVpuTOBNbsgfag7bgAUVl56i3IFvht+N0LXKuNU0TTuMUq2+lR
i+iha86OtkSYZqA34SXpoL2cyQJUikMkPsjXsxfpq1f35AbPC4FUfNnnN/NRBD/2
jOHHQSHcMOq0Y/SapCXgxa4urlLkCOra1lC5hQIDAQABAoIBAQCeK5E7nzv5cp+a
L3QVZOUI1V0DOTFHzn2Fnz8GeonTTGj2ShHp+g+mk5MCg7C3a5gQFpHFvRL65IJ/
G/LKVbwEQTc9uWPWhfIf5TDV82WNfH3lDgBVU9GFTus/Hu1xDrtu0WS+XpZrvSdm
f1s8Kv3r/cDHZkK7HEDD4I1i/Y//hinEOFbZglgmgGCwF3aZUdfj0wq1t1CnVJsj
P3c+rbxtvBLPblizHunqMiNHM1DTYqK5sW92yyMvy6ZbN3fSntcEdUNLWl61uWh1
E0eeg+GSuQoY/STFNqeTCqPju+gW8YR+Qix1R2n/AMo/D3HXcKxZCxWSD5zY8bH0
C/5oZerBAoGBAOdLjhQb2L1xU0V7vHWIHsWeOOTaDu1IhratY+detY/D0VCbnTXx
CZQmCgwT88/iPLVbEg3tsBBlVRYdTSOwtGk4Dvis83Izs1kU4RV8AI3+pZ74qA6O
8V3xrvaIrnJWHX0M078EAOFnj9zvkMFQnovnSyli+GSzPAreHVAmOMANAoGBAN1/
uojTaFg+cCfD6mQACo8DMN2CbWxJjXo4cpGv2HSYcjrVK+knjCdFhEnquLa0lQ+u
TxvNTFDhPDnEXFUegZwD+ge2oDI7wP/qAJUpgU53W12OvaRs5gMHKyU8AHWrMgYk
PvEwZIOlEc7ceyWzKCyt7nZwPHiteRsv+QgSS4lZAoGAUH27sQXL1ImWmAyqliBL
zSv10raMEUl3ECWhKciM2L4lnq649Cew1Ky0PGXJKGQsClTqIIzCA8Kv7KU/zhbV
gfRvSV0uz2RsmqioeAiSTNf8nSkdmwtltfLAl60TQFj1pCoNmmDzSX3308RPFOdQ
dZGFV57IoIq7b3DCtLzIbRUCgYB8ZFAYqUlPTXllC6SlllRXrn4R2D6lcsUuX2cQ
JEYWbMqx+aeYX+pY37SEYnpruQyBau3oeioives5sen8r44wVRdkn45lx6MC1aKQ
ImgI7gT0jMY6AiJGjw8O8Rx8+LC2PELQ5tF8EQboOnA6YtvsA54JC80aJKn/t7hO
bR/YuQKBgQDik5Es/Y2NM0JSKQwy7myGP4g5SadZd5Pv4r89OY/8mEOhKvHlaxki
M2o2is6RiFTKn/aFWXl0/QKXyupq35jsesq78y29zoUuph7J5PYcLgFF2t49q3B4
1Pn/53pMethW8zhphp3fM77vBkUu/NcusttbQjtdLEI7dkR/kqEc4fQ==
-----END RSA PRIVATE KEY-----`)

	keys, err := gitssh.NewPublicKeys("git", idRsa, "")
	if err != nil {
		log.Println("can't create PublicKeys", err)
	}

	repo, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		URL:               repoUrl,
		Auth:              keys,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Progress:          os.Stdout,
	})

	
	println(repoPath)
}

func main() {
	clone("git@github.com:liolay/go-config.git")
	//r, _ := git.PlainClone("/Users/liolay/foo", false, &git.CloneOptions{
	//	URL:      "https://github.com/liolay/go-config.git",
	//	Progress: os.Stdout,
	//})
	//
	//iter, e := r.Log(&git.LogOptions{})
	//iter.ForEach(func(commit *object.Commit) error {
	//	commit.
	//})

}
