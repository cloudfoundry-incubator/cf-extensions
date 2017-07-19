package bot

import (
	"context"
	"fmt"
	"os"

	"encoding/json"

	"github.com/google/go-github/github"

	"github.com/maximilien/cf-extensions/models"
)

type ExtRepos struct {
	Org    string
	Topics []string
	Client *github.Client
}

func NewExtRepos(org string, topics []string, client *github.Client) *ExtRepos {
	return &ExtRepos{
		Org:    org,
		Topics: topics,
		Client: client,
	}
}

func (extRepos *ExtRepos) GetInfos() models.Infos {
	orgOpts := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 30},
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := extRepos.Client.Repositories.ListByOrg(context.Background(), extRepos.Org, orgOpts)
		if err != nil {
			fmt.Printf("err: %s", err.Error())
			os.Exit(1)
		}

		var filteredRepos []*github.Repository
		for _, r := range repos {
			if extRepos.HasTopics(r, extRepos.Topics) {
				filteredRepos = append(filteredRepos, []*github.Repository{r}...)
			}
		}

		allRepos = append(allRepos, filteredRepos...)
		if resp.NextPage == 0 {
			break
		}

		orgOpts.Page = resp.NextPage
	}

	return extRepos.FetchInfos(allRepos)
}

func (extRepos *ExtRepos) HasTopics(repo *github.Repository, topics []string) bool {
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

func (extRepos *ExtRepos) DefaultInfo(repo *github.Repository) models.Info {
	return models.Info{}
}

func (extRepos *ExtRepos) FetchInfos(repos []*github.Repository) []models.Info {
	var infos []models.Info
	for _, r := range repos {
		info, err := extRepos.FetchInfo(r)
		if err != nil {
			info = extRepos.DefaultInfo(r)
			infos = append(infos, models.Info{})
		} else {
			infos = append(infos, info)
		}
	}
	return infos
}

func (extRepos *ExtRepos) FetchInfo(repo *github.Repository) (models.Info, error) {
	fileContents, _, _, err := extRepos.Client.Repositories.GetContents(context.Background(),
		extRepos.Org, *repo.Name, ".cf-extensions", &github.RepositoryContentGetOptions{})
	if err != nil {
		return models.Info{}, err
	}

	fileBytes, err := extractFileBytes(fileContents)
	if err != nil {
		return models.Info{}, err
	}

	info := models.Info{Repo: repo}
	err = json.Unmarshal(fileBytes, &info)
	if err != nil {
		return models.Info{}, err
	}

	return info, nil
}
