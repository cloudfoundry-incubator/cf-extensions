package bot

import (
	"context"
	"os"
	"path"
	"sort"
	"time"

	"encoding/json"
	"html/template"
	"io/ioutil"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"

	"github.com/cloudfoundry-incubator/cf-extensions/models"
)

type App struct {
	accessToken string
	Username    string
	Email       string
	ExtRepos    *ExtRepos
	Client      *github.Client
}

func NewApp(accessToken, username, email string) *App {
	app := &App{
		accessToken: accessToken,
		Username:    username,
		Email:       email,
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: app.accessToken})
	app.Client = github.NewClient(oauth2.NewClient(context.Background(), ts))

	return app
}

func (app *App) Run(org string, topics []string) {
	Printf("Finding CF-Extensions projects in org: `%s` using topics: `%s`\n", org, topics)
	Printf("Current time: `%s`\n", time.Now().String())

	app.ExtRepos = NewExtRepos(app.Username, org, topics, app.Client)
	trackedInfos, untrackedInfos := app.ExtRepos.GetInfos()

	sort.Sort(models.Infos(trackedInfos))
	sort.Sort(models.Infos(untrackedInfos))

	projectsPath := path.Join("data", "projects.json")
	projects := models.Projects{Org: app.ExtRepos.Org, Infos: trackedInfos}
	err := app.PushProjectsJsonDb(projects, projectsPath)
	if err != nil {
		Printf("ERROR: saving / pushing `%s` file: %s\n", projectsPath, err.Error())
	}

	untrackedProjectsPath := path.Join("data", "untracked_projects.json")
	untrackedProjects := models.Projects{Org: app.ExtRepos.Org, Infos: untrackedInfos}
	err = app.PushProjectsJsonDb(untrackedProjects, untrackedProjectsPath)
	if err != nil {
		Printf("ERROR: saving / pushing `%s` file: %s\n", untrackedProjectsPath, err.Error())
	}

	err = app.GenerateMarkdowns(projects)
	if err != nil {
		Printf("ERROR: generating markdown file for projects: %s\n", err.Error())
	}

	print(app.ExtRepos.Org, trackedInfos)

	Println("Done.")
}

func (app *App) GenerateMarkdowns(projects models.Projects) error {
	err := app.GenerateProjectsMarkdown(projects)
	if err != nil {
		return err
	}

	err = app.GenerateIndexMarkdown(projects)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) GenerateProjectsMarkdown(projects models.Projects) error {
	fileContents, _, _, err := app.Client.Repositories.GetContents(
		context.Background(),
		"cloudfoundry-incubator",
		"cf-extensions",
		"data/projects.json",
		&github.RepositoryContentGetOptions{})
	if err != nil {
		return err
	}

	if !hasProjectsChanged(projects, fileContents) {
		Printf("Commited projects.md has not changed, last commit SHA: %s\n", *fileContents.SHA)
		return nil
	}

	funcMap := template.FuncMap{
		"Length":           Length,
		"CurrentTime":      CurrentTime,
		"FormatAsDate":     FormatAsDate,
		"FormatAsDateTime": FormatAsDateTime,
		"ParseAsDate":      ParseAsDate,
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	templatePath := path.Join(wd, "templates", "projects.md.tmpl")
	t := template.Must(template.New("projects.md.tmpl").Funcs(funcMap).ParseFiles(templatePath))

	tmpFile, err := ioutil.TempFile(os.TempDir(), "cf-extensions")
	defer os.Remove(tmpFile.Name())
	if err != nil {
		return err
	}

	err = t.Execute(tmpFile, projects)
	if err != nil {
		return err
	}

	contents, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		return err
	}

	//docs/projects.md
	projectsMdFileContents, _, _, err := app.Client.Repositories.GetContents(
		context.Background(),
		"cloudfoundry-incubator",
		"cf-extensions",
		"docs/projects.md",
		&github.RepositoryContentGetOptions{})
	if err != nil {
		return err
	}

	message := "Updating cf-extensions projects.md file"
	repositoryContentsOptions := &github.RepositoryContentFileOptions{
		Message:   &message,
		Content:   contents,
		SHA:       projectsMdFileContents.SHA,
		Committer: &github.CommitAuthor{Name: github.String(app.Username), Email: github.String(app.Email)},
	}

	updateResponse, _, err := app.Client.Repositories.UpdateFile(
		context.Background(),
		"cloudfoundry-incubator",
		"cf-extensions",
		"docs/projects.md",
		repositoryContentsOptions)
	if err != nil {
		Printf("Repositories.UpdateFile returned error: %v", err)
		return err
	}

	Printf("Commited projects.md %s\n", *updateResponse.Commit.SHA)

	return nil
}

func (app *App) GenerateIndexMarkdown(projects models.Projects) error {
	fileContents, _, _, err := app.Client.Repositories.GetContents(
		context.Background(),
		"cloudfoundry-incubator",
		"cf-extensions",
		"data/projects.json",
		&github.RepositoryContentGetOptions{})
	if err != nil {
		return err
	}

	if !hasProjectsChanged(projects, fileContents) {
		Printf("Commited projects.md has not changed, last commit SHA: %s\n", *fileContents.SHA)
		return nil
	}

	funcMap := template.FuncMap{
		"Length":           Length,
		"CurrentTime":      CurrentTime,
		"FormatAsDate":     FormatAsDate,
		"FormatAsDateTime": FormatAsDateTime,
		"ParseAsDate":      ParseAsDate,
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	templatePath := path.Join(wd, "templates", "index.md.tmpl")
	t := template.Must(template.New("index.md.tmpl").Funcs(funcMap).ParseFiles(templatePath))

	tmpFile, err := ioutil.TempFile(os.TempDir(), "cf-extensions")
	defer os.Remove(tmpFile.Name())
	if err != nil {
		return err
	}

	err = t.Execute(tmpFile, projects)
	if err != nil {
		return err
	}

	contents, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		return err
	}

	indexMdFileContents, _, _, err := app.Client.Repositories.GetContents(
		context.Background(),
		"cloudfoundry-incubator",
		"cf-extensions",
		"docs/index.md",
		&github.RepositoryContentGetOptions{})
	if err != nil {
		return err
	}

	message := "Updating cf-extensions index.md file"
	repositoryContentsOptions := &github.RepositoryContentFileOptions{
		Message:   &message,
		Content:   contents,
		SHA:       indexMdFileContents.SHA,
		Committer: &github.CommitAuthor{Name: github.String(app.Username), Email: github.String(app.Email)},
	}

	updateResponse, _, err := app.Client.Repositories.UpdateFile(
		context.Background(),
		"cloudfoundry-incubator",
		"cf-extensions",
		"docs/index.md",
		repositoryContentsOptions)
	if err != nil {
		Printf("Repositories.UpdateFile returned error: %v", err)
		return err
	}

	Printf("Commited index.md %s\n", *updateResponse.Commit.SHA)

	return nil
}

func (app *App) PushProjectsJsonDb(projects models.Projects, filePath string) error {
	fileContents, _, _, err := app.Client.Repositories.GetContents(
		context.Background(),
		"cloudfoundry-incubator",
		"cf-extensions",
		filePath,
		&github.RepositoryContentGetOptions{})
	if err != nil {
		return err
	}

	if !hasProjectsChanged(projects, fileContents) {
		Printf("Commited projects.json has not changed, last commit SHA: %s\n", *fileContents.SHA)
		return nil
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "cf-extensions")
	defer os.Remove(tmpFile.Name())
	if err != nil {
		return err
	}

	contents, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return err
	}

	message := "Updating cf-extensions repos info"
	repositoryContentsOptions := &github.RepositoryContentFileOptions{
		Message:   &message,
		Content:   contents,
		SHA:       fileContents.SHA,
		Committer: &github.CommitAuthor{Name: github.String(app.Username), Email: github.String(app.Email)},
	}

	updateResponse, _, err := app.Client.Repositories.UpdateFile(
		context.Background(),
		"cloudfoundry-incubator",
		"cf-extensions",
		filePath,
		repositoryContentsOptions)
	if err != nil {
		Printf("Repositories.UpdateFile returned error: %v", err)
		return err
	}

	Printf("Commited `%s` %s\n", filePath, *updateResponse.Commit.SHA)

	return nil
}

// Private utility functions

func print(org string, infos []models.Info) {
	sort.Sort(models.Infos(infos))

	Println()
	Printf("Repos for %s, total: %d\n", org, len(infos))
	Println("-----------------")
	for _, info := range infos {
		Printf("Repo name: %s, URL: %s\n", *info.Repo.Name, *info.Repo.GitURL)
		Printf("Topics:     %s\n", *info.Repo.Topics)
		Println()
	}
	Println("-----------------")
	Printf("Total repos: %d\n", len(infos))
}

func hasProjectsChanged(projects models.Projects, fileContent *github.RepositoryContent) bool {
	fileBytes, err := extractFileBytes(fileContent)
	if err != nil {
		return true
	}

	downloadedProjects := models.Projects{}
	err = json.Unmarshal(fileBytes, &downloadedProjects)
	if err != nil {
		return true
	}

	return projects.Equals(downloadedProjects) != true
}
