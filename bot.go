package main

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v35/github"
	"github.com/jspc-bots/bottom"
	"github.com/lrstanley/girc"
)

type Bot struct {
	bottom bottom.Bottom
	github *github.Client
}

func New(user, password, server string, verify bool, gh *github.Client) (b Bot, err error) {
	b.github = gh

	b.bottom, err = bottom.New(user, password, server, verify)
	if err != nil {
		return
	}

	b.bottom.Client.Handlers.Add(girc.CONNECTED, func(c *girc.Client, e girc.Event) {
		c.Cmd.Join(Chan)
	})

	router := bottom.NewRouter()
	router.AddRoute(`bump\s+(major|minor|patch)\s+version\s+for\s+([\w\-\_]+)\/([\w\-\_\.]+)`, b.newRelease)

	b.bottom.Middlewares.Push(router)

	return
}

func (b Bot) newRelease(sender, channel string, groups []string) (err error) {
	ctx := context.Background()
	owner := groups[2]
	repo := groups[3]

	// validate repo exists (if it does but GetLatestRelease fails, then we're creating the first
	// release)
	_, _, err = b.github.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return
	}

	// Get latest release
	rel, resp, err := b.github.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		// A status of 404 on a repo which exists means there is no latest release
		// and so this release will be the first
		if resp.StatusCode != 404 {
			return
		}

		err = nil

		// Create an empty, dummy release with a version of 0.0.0
		rel = &github.RepositoryRelease{
			TagName: github.String("v0.0.0"),
		}
	}

	// If latest does not follow semver2, or is a pre-release then return error
	v, err := semver.NewVersion(*rel.TagName)
	if err != nil {
		return
	}

	// Increment relevant field
	var nv semver.Version
	switch groups[1] {
	case "major":
		nv = v.IncMajor()

	case "minor":
		nv = v.IncMinor()

	case "patch":
		nv = v.IncPatch()
	}

	// Push to github
	rel.TagName = github.String(fmt.Sprintf("v%s", nv.String()))
	rel.TargetCommitish = github.String("") // let it pick up the latest commit on the default branch

	_, _, err = b.github.Repositories.CreateRelease(context.Background(), owner, repo, rel)
	if err != nil {
		return
	}

	b.bottom.Client.Cmd.Messagef(channel, "%s/%s %s released", owner, repo, nv.String())

	// return
	return
}
