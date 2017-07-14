package main

import (
	"fmt"
	"os"
	"io"
	"net/http"
	"context"
	"io/ioutil"
	"encoding/json"

	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
)

type CFExtensionsInfo struct {
	Name string `json:"name"`
	GitUrl string  `json:"git_url"`
	TrackerUrl string  `json:"tracker_url"`
	Description string `json:"description"`
	OwnerCompany string `json:"owner_company"`
	ContactEmail string `json:"contact_email"`
	Status string `json:"status"`
	ProposedDate string `json:"proposed_date"`
	StatusChangedDate string `json:"status_changed_date"`
}

type Projects struct {
	Projects []CFExtensionsInfo `json:"projects"`
}

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "30bb3196bd7d24eeba37b0e6def3e556b6ed49f1"})
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	listReposByOrg("cloudfoundry-incubator", []string{ "cf-extensions" }, client)
}

func listReposByOrg(org string, topicsFilter []string, client *github.Client) {
	orgOpts :=  &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 30},
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(context.Background(), org, orgOpts)
		if err != nil {
			fmt.Printf("err: %s", err.Error())
			os.Exit(1)
		}

		var filteredRepos []*github.Repository
		for _, r := range repos {
			if repoHasTopics(r, topicsFilter) {
				filteredRepos = append(filteredRepos, []*github.Repository{r}...)
			}
		}

		allRepos = append(allRepos, filteredRepos...)
		if resp.NextPage == 0 {
			break
		}

		orgOpts.Page = resp.NextPage
	}

	cfExtensionsInfos := fetchCFExtensionsInfos(allRepos, client)
	printRepos(org, allRepos, cfExtensionsInfos)
}

func fetchCFExtensionsInfos(repos []*github.Repository, client *github.Client) []CFExtensionsInfo {
	var cfExtensionsInfos []CFExtensionsInfo
	for _, r := range repos {
		cfExtensionsInfo, err := fetchCFExtensionsInfo(r, client)
		if err != nil {
			cfExtensionsInfos = append(cfExtensionsInfos, CFExtensionsInfo{})
		} else {
			cfExtensionsInfos = append(cfExtensionsInfos, cfExtensionsInfo)
		}
	}
	return cfExtensionsInfos
}

func printRepos(org string, repos []*github.Repository, infos []CFExtensionsInfo) {
	fmt.Printf("Repo s for org: %s, total: %d\n", org, len(repos))
	fmt.Println("-----------------\n")
	for i, r := range repos {
		fmt.Printf("Repo name: %s, URL: %s\n", *r.Name, *r.GitURL)
		fmt.Printf("Topics:     %s\n", *r.Topics)
		fmt.Printf(".cf-extensions: %v\n", infos[i])
		fmt.Println()
	}
	fmt.Println("-----------------\n")
	fmt.Printf("Total repos: %d\n", len(repos))
}

func repoHasTopics(repo *github.Repository, topics []string) bool {
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

func fetchCFExtensionsInfo(repo *github.Repository, client *github.Client) (CFExtensionsInfo, error) {
	fileContents, _, _, err := client.Repositories.GetContents(context.Background(),
		"cloudfoundry-incubator", *repo.Name, ".cf-extensions", &github.RepositoryContentGetOptions{})
	if err != nil {
		return CFExtensionsInfo{}, err
	}

	response, err := http.Get(*fileContents.DownloadURL)
	if err != nil {
		return CFExtensionsInfo{}, err
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "cf-extensions")
	defer os.Remove(tmpFile.Name())

	defer response.Body.Close()
	_, err = io.Copy(tmpFile, response.Body)
	if err != nil {
		return CFExtensionsInfo{}, err
	}

	fileBytes, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		return CFExtensionsInfo{}, err
	}

	cfExtensionsInfo := CFExtensionsInfo{}
	err = json.Unmarshal(fileBytes, &cfExtensionsInfo)
	if err != nil {
		return CFExtensionsInfo{}, err
	}

	return cfExtensionsInfo, nil
}