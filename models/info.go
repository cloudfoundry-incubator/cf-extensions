package models

import (
	"time"

	"github.com/google/go-github/github"
)

const PROPOSAL_DEFAULT_URL = "https://docs.google.com/document/d/1cpyBmds7WYNLKO1qkjhCdS8bNSJjWH5MqTE-h1UCQkQ/edit?usp=sharing"
const LOGO_DEFAULT_URL = "https://github.com/cloudfoundry-incubator/cf-extensions/blob/master/images/cf-extensions-proposal-logo.png"
const ICON_DEFAULT_URL = "https://github.com/cloudfoundry-incubator/cf-extensions/blob/master/images/cf-extensions-proposal-icon.png"

type Status struct {
	Status      string `json:"status"`
	ChangedDate string `json:"status_changed_date"`
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

	Repo *github.Repository `json:"-"`
}

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
