package models

import (
	"time"

	"github.com/google/go-github/github"
)

const PROPOSAL_DEFAULT_URL = "https://docs.google.com/document/d/1cpyBmds7WYNLKO1qkjhCdS8bNSJjWH5MqTE-h1UCQkQ/edit?usp=sharing"
const TRACKER_DEFAULT_URL = "https://www.pivotaltracker.com"

const LOGO_DEFAULT_URL = "https://github.com/cloudfoundry-incubator/cf-extensions/blob/master/docs/images/cf-extensions-proposal-logo.png"
const ICON_DEFAULT_URL = "https://github.com/cloudfoundry-incubator/cf-extensions/blob/master/docs/images/cf-extensions-proposal-icon.png"

const CATEGORY_DEFAULT = "unknown"
const COMMIT_STYLE = "pairing"

var CATEGORIES = map[string]string{"bosh": "BOSH", "runtime": "Runtime", "apis": "APIs", "cli": "CLI", "tool": "Tool", "unkown": "Unknown"}
var COMMIT_STYLES = map[string]string{"pairing": "Pairing", "distributed": "Distributed"}

type Status struct {
	Status      string `json:"status,omitempty" yaml:"status,omitempty"`
	ChangedDate string `json:"status_changed_date,omitempty" yaml:"status_changed_date,omitempty"`
	Category    string `json:"category,omitempty" yaml:"category,omitempty"`
	CommitStyle string `json:"commit_style,omitempty" yaml:"commit_style,omitempty"`
}

type Statistics struct {
	ForksCount      int `json:"-" yaml:"-"`
	OpenIssuesCount int `json:"-" yaml:"-"`
	StargazersCount int `json:"-" yaml:"-"`
	WatchersCount   int `json:"-" yaml:"-"`
}

type Info struct {
	// Optionally provided by owner
	OwnerCompany string `json:"owner_company" yaml:"owner_company"`
	ContactEmail string `json:"contact_email" yaml:"contact_email"`

	Description string `json:"description" yaml:"description"`

	TrackerUrl   string `json:"tracker_url" yaml:"tracker_url"`
	ProposalUrl  string `json:"proposal_url" yaml:"proposal_url"`
	ProposedDate string `json:"proposed_date" yaml:"proposed_date"`

	LogoUrl string `json:"logo_url,omitempty" yaml:"logo_url,omitempty"`
	IconUrl string `json:"icon_url,omitempty" yaml:"icon_url,omitempty"`

	// Computed fields
	Name   string `json:"name,omitempty" yaml:"name,omitempty"`
	GitUrl string `json:"git_url,omitempty" yaml:"git_url,omitempty"`

	Stats Statistics `json:"-" yaml:"-"`

	Repo              *github.Repository        `json:"-" yaml:"-"`
	LatestRepoRelease *github.RepositoryRelease `json:"-" yaml:"-"`

	// Protected fields

	Status `json:"-" yaml:"-"`
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
