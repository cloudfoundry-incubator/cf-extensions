package models

import (
	"time"

	"github.com/google/go-github/github"
)

const PROPOSAL_DEFAULT_URL = "https://docs.google.com/document/d/1cpyBmds7WYNLKO1qkjhCdS8bNSJjWH5MqTE-h1UCQkQ/edit?usp=sharing"
const LOGO_DEFAULT_URL = "https://github.com/cloudfoundry-incubator/cf-extensions/blob/master/docs/images/cf-extensions-proposal-logo.png"
const ICON_DEFAULT_URL = "https://github.com/cloudfoundry-incubator/cf-extensions/blob/master/docs/images/cf-extensions-proposal-icon.png"

type Status struct {
	Status      string `json:"status"`
	ChangedDate string `json:"status_changed_date"`
}

type Statistics struct {
	ForksCount      int `json:"forks_count,omitempty"`
	OpenIssuesCount int `json:"open_issues_count,omitempty"`
	StargazersCount int `json:"stargazers_count,omitempty"`
	WatchersCount   int `json:"watchers_count,omitempty"`
}

type Info struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	GitUrl      string `json:"git_url"`
	TrackerUrl  string `json:"tracker_url"`
	ProposalUrl string `json:"proposal_url"`

	LogoUrl string `json:"logo_url"`
	IconUrl string `json:"icon_url"`

	OwnerCompany string `json:"owner_company"`
	ContactEmail string `json:"contact_email"`
	ProposedDate string `json:"proposed_date"`

	Status

	Stats Statistics `json:"-"`

	Repo              *github.Repository        `json:"-"`
	LatestRepoRelease *github.RepositoryRelease `json:"-"`
}

// Infos methods

type Infos []Info

func (infos Infos) Len() int {
	return len(infos)
}

func (infos Infos) Swap(i, j int) {
	infos[i], infos[j] = infos[j], infos[i]
}
func (infos Infos) Less(i, j int) bool {
	return infos[i].Name < infos[j].Name
}

// Info methods

func CreateInfo(repo *github.Repository) *Info {
	info := Info{
		Repo: repo,
	}

	info.AddDefaults()
	info.UpdateFromRepo()

	return &info
}

func (info *Info) AddDefaults() {
	if info.ProposalUrl == "" {
		info.ProposalUrl = PROPOSAL_DEFAULT_URL
	}

	if info.LogoUrl == "" {
		info.LogoUrl = LOGO_DEFAULT_URL
	}

	if info.IconUrl == "" {
		info.IconUrl = ICON_DEFAULT_URL
	}

	if info.ProposedDate == "" {
		info.ProposedDate = time.Now().String()
	}
}

func (info *Info) UpdateFromRepo() {
	if info.Repo == nil {
		return
	}

	info.Stats = Statistics{
		ForksCount:      *(info.Repo).ForksCount,
		OpenIssuesCount: *(info.Repo).OpenIssuesCount,
		StargazersCount: *(info.Repo).StargazersCount,
		WatchersCount:   *(info.Repo).WatchersCount,
	}
}
