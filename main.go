package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"

	version "github.com/mcuadros/go-version"
	"github.com/shurcooL/githubv4"
)

type node struct {
	IsDraft      bool // githubv4.Boolean
	IsPrerelease bool // githubv4.Boolean
	Tag          struct {
		Name string // githubv4.String
	}
}

type edge struct {
	Tag struct {
		Name string // githubv4.String
	} `graphql:"tag:node"`
}

var q struct {
	Repository struct {
		Releases struct {
			Nodes []node
		} `graphql:"releases(last: 100)"`
		Tags struct {
			Edges []edge
		} `graphql:"tags:refs(refPrefix:\"refs/tags/\", last:100)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func main() {
	flag.Parse()
	tailArgs := flag.Args()
	if len(tailArgs) != 3 {
		log.Fatalln("Invoke with [repoOwner] [repoName] [versionInUse].")
	}
	owner, repoName, versionInUse := tailArgs[0], tailArgs[1], tailArgs[2]

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(repoName),
	}

	err := client.Query(context.Background(), &q, variables)
	if err != nil {
		log.Fatalln(err)
	}

	if len(q.Repository.Releases.Nodes) > 0 {
		// use releases if they exist
		// example: roundcube roundcubemail 1.3.6
		for _, release := range q.Repository.Releases.Nodes {
			if release.IsDraft || release.IsPrerelease {
				continue
			}
			if version.Compare(release.Tag.Name, versionInUse, ">") {
				fmt.Println("Found newer version:", release.Tag.Name)
			}
		}
	} else {
		// fall back to tags if no releases exist
		// example: pytest-dev pytest 3.9.0
		for _, tagEdge := range q.Repository.Tags.Edges {
			if version.Compare(tagEdge.Tag.Name, versionInUse, ">") {
				fmt.Println("Found newer version:", tagEdge.Tag.Name)
			}
		}
	}
}
