package main

import (
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
	"context"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "30bb3196bd7d24eeba37b0e6def3e556b6ed49f1"})
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	//listUserRepos("maximilien", client, ctx)
	listReposByOrg("cloudfoundry-incubator", []string{ "cf-extensions" }, client, ctx)
	//listReposByOrg("cloudfoundry-incubator", []string{}, client, ctx)
}

func listUserRepos(user string, client *github.Client, ctx context.Context) {
	opts :=  &github.RepositoryListOptions{Visibility: "public"}
	repos, _, err := client.Repositories.List(ctx, user, opts)
	if err != nil {
		fmt.Printf("err: %s", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Repos for user: %s, total: %d\n", "maximilien", len(repos))
	for _, r := range repos {
		fmt.Printf("Repo name:  %s\n", *r.Name)
		fmt.Printf("Git URL:    %s\n", *r.GitURL)
		fmt.Printf("Labels URL: %s\n", *r.LabelsURL)
		fmt.Printf("Tags URL:   %s\n", *r.TagsURL)
		fmt.Printf("Topics:     %s\n", *r.Topics)
		fmt.Println("-----------------\n")
	}
}

func listReposByOrg(org string, topicsFilter []string, client *github.Client, ctx context.Context) {
	orgOpts :=  &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 30},
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, org, orgOpts)
		if err != nil {
			fmt.Printf("err: %s", err.Error())
			os.Exit(1)
		}

		var filteredRepos []*github.Repository
		for _, r := range repos {
			if repoHasTopics(r, topicsFilter) {
				filteredRepos = append(allRepos, []*github.Repository{r}...)
			}
		}

		allRepos = append(allRepos, filteredRepos...)
		if resp.NextPage == 0 {
			break
		}

		orgOpts.Page = resp.NextPage
	}

	fmt.Printf("Repo s for org: %s, total: %d\n", org, len(allRepos))
	fmt.Println("-----------------\n")
	for _, r := range allRepos {
		fmt.Printf("Repo name: %s, URL: %s\n", *r.Name, *r.GitURL)
		fmt.Printf("Topics:     %s\n", *r.Topics)
	}
	fmt.Println("-----------------\n")
	fmt.Printf("Total repos: %d\n", len(allRepos))
}

func repoHasTopics(repo *github.Repository , topics []string) bool {
	for _, topic := range topics {
		found := false
		for _, repoTopic := range *repo.Topics {
			if topic == repoTopic {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}