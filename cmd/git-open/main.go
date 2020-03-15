/*
Command 'git-open' opens the

To install, use go get,

  $ go get github.com/tcnksm/misc/cmd/git-open

*/
package main

import (
	"log"
	"regexp"
	"strconv"

	"github.com/cli/cli/git"
)

func main() {
	prHeadRef, err := git.CurrentBranch()
	if err != nil {
		log.Fatal(err)
	}

	branchConfig := git.ReadBranchConfig(prHeadRef)
	prHeadRE := regexp.MustCompile(`^refs/pull/(\d+)/head$`)
	if m := prHeadRE.FindStringSubmatch(branchConfig.MergeRef); m != nil {
		prNumber, _ := strconv.Atoi(m[1])
		log.Println(prNumber)
	}
	log.Printf("%#v", branchConfig)
	log.Println(prHeadRef)
}
