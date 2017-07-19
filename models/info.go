package models

import "github.com/google/go-github/github"

type Status struct {
	Status      string `json:"status"`
	ChangedDate string `json:"status_changed_date"`
}

type Info struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	GitUrl     string `json:"git_url"`
	TrackerUrl string `json:"tracker_url"`

	OwnerCompany string `json:"owner_company"`
	ContactEmail string `json:"contact_email"`
	ProposedDate string `json:"proposed_date"`

	Status

	Repo *github.Repository
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
