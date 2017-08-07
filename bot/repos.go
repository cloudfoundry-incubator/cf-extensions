package bot

import (
	"context"
	"fmt"
	"os"
	"time"

	"encoding/json"
	"io/ioutil"
	"text/template"

	"github.com/google/go-github/github"

	"github.com/maximilien/cf-extensions/models"
)

type ExtRepos struct {
	Username string
	Org      string
	Topics   []string
	Client   *github.Client
}

const ISSUE_TITLE = "Add .cf-extensions to your repo to be listed in cloudfoundry-incubator.cf-extensions"

const INFO_ISSUE_BODY = `Add {{.Filename}} file to your repo so that it shows correctly in the CF-Extensions catalog.

{{.InfoJson}}

This is a JSON formatted file. The default values in the file are for you to get started. You should edit to match your project's data.

For example, the field {{.TrackerUrl}} should contain your project's tracker URL, and so on.
`

func NewExtRepos(username, org string, topics []string, client *github.Client) *ExtRepos {
	return &ExtRepos{
		Username: username,
		Org:      org,
		Topics:   topics,
		Client:   client,
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
	info := models.Info{
		Name:   *repo.Name,
		GitUrl: *repo.GitURL,

		Description: "ADD DESCRIPTION HERE",

		ProposalUrl: models.PROPOSAL_DEFAULT_URL,

		LogoUrl: models.LOGO_DEFAULT_URL,
		IconUrl: models.ICON_DEFAULT_URL,

		OwnerCompany: "ADD OWNER COMPANY HERE",
		ContactEmail: "contact@owner-company.com",

		ProposedDate: time.Now().String(),

		Repo: repo,
	}

	info.UpdateFromRepo()

	return info
}

func (extRepos *ExtRepos) FetchInfos(repos []*github.Repository) []models.Info {
	var infos []models.Info
	for _, r := range repos {
		info, err := extRepos.FetchInfo(r)
		if err != nil {
			info = extRepos.DefaultInfo(r)
			if !extRepos.InfoIssueExists(info) {
				issue, err := extRepos.CreateInfoIssue(info, r)
				if err != nil {
					fmt.Printf("ERROR creating default info issue to: %s, message: %s\n", info.Name, err.Error())
				}
				fmt.Printf("Created default info issue #%d to: %s\n", *issue.Number, info.Name)
			} else {
				fmt.Printf("Info issue already exists in %s\n", info.Name)
			}
		} else {
			latestRepoRelease, err := extRepos.FetchLatestRepoRelease(r)
			if err != nil {
				fmt.Printf("Error getting latest release for repo: %s\n", info.Name)
			} else {
				info.LatestRepoRelease = latestRepoRelease
			}
			info.AddDefaults()
			infos = append(infos, info)
		}
	}
	return infos
}

func (extRepos *ExtRepos) FetchLatestRepoRelease(repo *github.Repository) (*github.RepositoryRelease, error) {
	latestRepoRelease, _, err := extRepos.Client.Repositories.GetLatestRelease(context.Background(), extRepos.Org, *repo.Name)
	if err != nil {
		return nil, err
	}

	return latestRepoRelease, nil
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
	info.UpdateFromRepo()

	return info, nil
}

func (extRepos *ExtRepos) CreateInfoIssue(info models.Info, repo *github.Repository) (*github.Issue, error) {
	infoBytes, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		fmt.Printf("Could not marshall info info string error: %v\n", err)
		return nil, err
	}

	type IssueInfo struct {
		Filename   string
		InfoJson   string
		TrackerUrl string
	}
	issueInfo := IssueInfo{
		"`.cf-extensions`",
		fmt.Sprintf("```json\n%s\n```", string(infoBytes)),
		"`tracker_url`",
	}
	issueInfoTemplate, err := template.New("issue-info").Parse(INFO_ISSUE_BODY)
	if err != nil {
		fmt.Printf("Could not create issue info error: %v\n", err)
		return nil, err
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "cf-extensions")
	defer os.Remove(tmpFile.Name())
	if err != nil {
		return nil, err
	}

	err = issueInfoTemplate.Execute(tmpFile, issueInfo)
	if err != nil {
		return nil, err
	}

	issueInfoContents, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, err
	}

	issueRequest := github.IssueRequest{
		Title: github.String(ISSUE_TITLE),
		Body:  github.String(string(issueInfoContents)),
	}

	issue, _, err := extRepos.Client.Issues.Create(context.Background(), extRepos.Org, info.Name, &issueRequest)
	if err != nil {
		fmt.Printf("Issues.Create returned error: %v\n", err)
		return nil, err
	}

	return issue, nil
}

func (extRepos *ExtRepos) InfoIssueExists(info models.Info) bool {
	issueListByRepoOpts := github.IssueListByRepoOptions{
		State:   "open",
		Creator: extRepos.Username,
	}

	issues, _, err := extRepos.Client.Issues.ListByRepo(context.Background(), extRepos.Org, info.Name, &issueListByRepoOpts)
	if err != nil {
		fmt.Printf("Issues.List returned error: %v\n", err)
		return false
	}

	for _, issue := range issues {
		if *issue.Title == ISSUE_TITLE {
			return true
		}
	}

	return false
}
