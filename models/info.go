package models

import (
	"time"

	"github.com/google/go-github/github"
)

const PROPOSAL_DEFAULT_URL = "https://docs.google.com/document/d/1cpyBmds7WYNLKO1qkjhCdS8bNSJjWH5MqTE-h1UCQkQ/edit?usp=sharing"
const LOGO_DEFAULT_URL = "https://github.com/cloudfoundry-incubator/cf-extensions/blob/master/docs/images/cf-extensions-proposal-logo.png"
const ICON_DEFAULT_URL = "https://github.com/cloudfoundry-incubator/cf-extensions/blob/master/docs/images/cf-extensions-proposal-icon.png"

type Status struct {
	Status      string `json:"status,omitempty"`
	ChangedDate string `json:"status_changed_date,omitempty"`
}

type Statistics struct {
	ForksCount      int `json:"-"`
	OpenIssuesCount int `json:"-"`
	StargazersCount int `json:"-"`
	WatchersCount   int `json:"-"`
}

type Info struct {
	// Optionally provided by owner
	Description string `json:"description"`

	TrackerUrl  string `json:"tracker_url"`
	ProposalUrl string `json:"proposal_url"`

	LogoUrl string `json:"logo_url"`
	IconUrl string `json:"icon_url"`

	OwnerCompany string `json:"owner_company"`
	ContactEmail string `json:"contact_email"`
	ProposedDate string `json:"proposed_date"`

	// Computed fields
	Name   string `json:"name,omitempty"`
	GitUrl string `json:"git_url,omitempty"`

	Stats Statistics `json:"-"`

	Repo              *github.Repository        `json:"-"`
	LatestRepoRelease *github.RepositoryRelease `json:"-"`

	// Protected fields

	Status
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
	info.Update()

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

func (info *Info) Update() {
	info.updateFromRepo()
}

// Private methods

func (info *Info) updateFromRepo() {
	if info.Repo == nil {
		return
	}

	info.Name = *(info.Repo).Name
	info.GitUrl = *(info.Repo).GitURL

	info.Stats = Statistics{
		ForksCount:      *(info.Repo).ForksCount,
		OpenIssuesCount: *(info.Repo).OpenIssuesCount,
		StargazersCount: *(info.Repo).StargazersCount,
		WatchersCount:   *(info.Repo).WatchersCount,
	}
}
