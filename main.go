package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/src-d/go-git.v4"
)

func main() {
	// Tempdir to clone the repository
	dir, err := ioutil.TempDir("", "clone-example")
	log.Println(dir)
	if err != nil {
		log.Fatal(err)
	}
	dir = "/var/folders/bk/88q3ywjn3m1d4sxcfxt3nx7r0000gp/T/clone-example420022365"
	// defer os.RemoveAll(dir) // clean up

	// Clones the repository into the given dir, just as a normal git clone does
	start := time.Now()
	var gitRepo *git.Repository
	gitRepo, err = git.PlainCloneContext(context.Background(), dir, false, &git.CloneOptions{
		URL:           "git@github.com:tcnksm/misc.git",
		ReferenceName: "refs/heads/test",
		SingleBranch:  true,
		Depth:         1,
		Tags:          git.NoTags,
	})
	if err != nil {
		if err == git.ErrRepositoryAlreadyExists {
			gitRepo, err = git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{})
			if err != nil {
				log.Fatal("Open: ", err)
			}
		} else {
			log.Fatal(err)
		}
	}

	workTree, err := gitRepo.Worktree()
	if err != nil {
		log.Fatal("WorkTree: ", err)
	}

	err = workTree.Pull(&git.PullOptions{
		ReferenceName: "refs/heads/test",
		SingleBranch:  true,
		Depth:         1,
	})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
		} else {
			log.Fatal("Pull: ", err)
		}
	}

	fmt.Println(time.Now().Sub(start))
}
