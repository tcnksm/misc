/*
Command 'sync-github-release' syncs github release pages with one with
the other. For example, this is used to sync fork and upstream release.

  $ sync-github-release REPO DIST_OWNER SRC_OWNER

For example, if you want to sync releases on https://github.com/tcnksm/xxxx with
https://github.com/deeeet/xxxx

  $ sync-github-release xxxx deeeet tcnksm

It doensn't sync uploaded artifacts.

To use this command, you need to prepare GitHub API Token (with repo priviledge).
You can set it via TOKEN env var.

*/
package main

import (
	"log"
	"os"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const EnvToken = "TOKEN"

func main() {

	if len(os.Args) != 4 {
		log.Fatal("[Usage] sync-github-release REPO DIST_OWNER SRC_OWNER")
	}

	token := os.Getenv(EnvToken)
	if len(token) == 0 {
		log.Fatal("You need GitHub API token (repo priviledge) via TOKEN env var")
	}

	repo, distOwner, srcOwner := os.Args[1], os.Args[2], os.Args[3]
	log.Printf("[INFO] Import GitHub release from %s/%s to %s/%s",
		srcOwner, repo, distOwner, repo)

	// Construct github HTTP client
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	// List all releases on upstream repository
	releases, _, err := client.Repositories.ListReleases(srcOwner, repo, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create release which is on upsteram on fork repository
	log.Printf("[INFO] Found %d releases", len(releases))
	var success int
	for _, release := range releases {
		log.Println("Sync", *release.TagName)

		_, _, err = client.Repositories.CreateRelease(distOwner, repo, &github.RepositoryRelease{
			Name:       release.TagName,
			TagName:    release.TagName,
			Body:       release.Body,
			Draft:      github.Bool(false),
			Prerelease: github.Bool(false),
		})

		if err != nil {
			log.Println("[ERROR] Failed to create release:", err)
			continue
		}

		// Prevent to DDos to GitHub
		time.Sleep(5 * time.Second)
		success++
	}

	log.Printf("[INFO] Successfully sync %d releases", success)
}
