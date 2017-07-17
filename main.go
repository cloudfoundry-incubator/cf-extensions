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
	"sort"
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

type CFExtensionsInfos []CFExtensionsInfo

func (infos CFExtensionsInfos) Len() int {
	return len(infos)
}

func (infos CFExtensionsInfos) Swap(i, j int) {
	infos[i], infos[j] = infos[j], infos[i]
}
func (infos CFExtensionsInfos) Less(i, j int) bool {
	return infos[i].Name < infos[j].Name
}

type Projects struct {
	Org string `json:"org"`
	CFExtensionInfos []CFExtensionsInfo `json:"projects"`
}

func (p Projects) Equal(otherProjects Projects) bool {
	if p.Org != otherProjects.Org {
		return false
	}

	for _, info := range p.CFExtensionInfos {
		found := false

		for _, otherInfo := range otherProjects.CFExtensionInfos {
			if otherInfo.Name == info.Name {
				found = true
				if info == otherInfo {
					break
				} else {
					return false
				}
			}
		}

		if !found {
			return false
		}
	}

	return true
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
	sort.Sort(CFExtensionsInfos(cfExtensionsInfos))
	err := saveAndPush(org, cfExtensionsInfos, client)
	if err != nil {
		fmt.Printf("ERROR: saving / pushing file with CF-Extensions infos: %s\n", err.Error())
	}

	print(org, allRepos, cfExtensionsInfos)
}

func createDefaultCFExtensionInfo(repo *github.Repository, client *github.Client) CFExtensionsInfo {
	return CFExtensionsInfo{}
}

func fetchCFExtensionsInfos(repos []*github.Repository, client *github.Client) []CFExtensionsInfo {
	var cfExtensionsInfos []CFExtensionsInfo
	for _, r := range repos {
		cfExtensionsInfo, err := fetchCFExtensionsInfo(r, client)
		if err != nil {
			cfExtensionsInfo = createDefaultCFExtensionInfo(r, client)
			cfExtensionsInfos = append(cfExtensionsInfos, CFExtensionsInfo{})
		} else {
			cfExtensionsInfos = append(cfExtensionsInfos, cfExtensionsInfo)
		}
	}
	return cfExtensionsInfos
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

	fileBytes, err := extractFileBytes(fileContents)
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

func extractFileBytes(fileContent *github.RepositoryContent) ([]byte, error) {
	response, err := http.Get(*fileContent.DownloadURL)
	if err != nil {
		return []byte{}, err
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "cf-extensions")
	defer os.Remove(tmpFile.Name())
	if err != nil {
		return []byte{}, err
	}

	defer response.Body.Close()
	_, err = io.Copy(tmpFile, response.Body)
	if err != nil {
		return []byte{}, err
	}

	fileBytes, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		return []byte{}, err
	}

	return fileBytes, nil
}

func print(org string, repos []*github.Repository, infos []CFExtensionsInfo) {
	sort.Sort(CFExtensionsInfos(infos))

	fmt.Printf("Repos for %s, total: %d\n", org, len(repos))
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

func saveAndPush(org string, infos []CFExtensionsInfo, client *github.Client) error {
	projects := Projects{Org: org, CFExtensionInfos: infos}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "cf-extensions")
	defer os.Remove(tmpFile.Name())
	if err != nil {
		return err
	}

	contents, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return err
	}

	fileContents, _, _, err := client.Repositories.GetContents(context.Background(),
		"cloudfoundry-incubator", "cf-extensions", "projects.json", &github.RepositoryContentGetOptions{})
	if err != nil {
		return err
	}

	if !hasProjectsChanged(projects, fileContents) {
		fmt.Printf("Commited projects.json has not changed, last commit SHA: %s\n", *fileContents.SHA)
		return nil
	}

	message := "Updating cf-extensions repos info"
	repositoryContentsOptions := &github.RepositoryContentFileOptions{
		Message:   &message,
		Content:   contents,
		SHA: fileContents.SHA,
		Committer: &github.CommitAuthor{Name: github.String("maximilien"), Email: github.String("maxim@us.ibm.com")},
	}

	updateResponse, _, err := client.Repositories.UpdateFile(context.Background(), "cloudfoundry-incubator", "cf-extensions", "projects.json", repositoryContentsOptions)
	if err != nil {
		fmt.Printf("Repositories.UpdateFile returned error: %v", err)
		return err
	}

	fmt.Printf("Commited projects.json %s\n\n", *updateResponse.Commit.SHA)

	return nil
}

func hasProjectsChanged(projects Projects, fileContent *github.RepositoryContent) bool {
	fileBytes, err := extractFileBytes(fileContent)
	if err != nil {
		return true
	}

	downloadedProjects := Projects{}
	err = json.Unmarshal(fileBytes, &downloadedProjects)
	if err != nil {
		return true
	}

	return projects.Equal(downloadedProjects) != true
}