package bot

import (
	"fmt"
	"io"
	"os"
	"time"

	"io/ioutil"
	"net/http"

	"github.com/google/go-github/github"

	"github.com/maximilien/cf-extensions/models"
)

func extractFileBytes(fileContent *github.RepositoryContent) ([]byte, error) {
	response, err := http.Get(*fileContent.DownloadURL)
	if err != nil {
		return []byte{}, err
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "ExtRepos")
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

func length(infos []models.Info) int {
	return len(infos)
}

func currentTime() time.Time {
	return time.Now()
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()

	return fmt.Sprintf("%d/%d/%d", month, day, year)
}

func formatAsDateTime(t time.Time) string {
	year, month, day := t.Date()

	return fmt.Sprintf("%d/%d/%d @ %d:%d:%d", month, day, year, t.Hour(), t.Minute(), t.Second())
}

func parseAsDate(timeString string) string {
	stringTime, err := time.Parse("2017-02-03T12:00:00Z07:00", timeString)
	if err != nil {
		fmt.Printf("ERROR parsing time: %s, message: %s\n", timeString, err.Error())
		return formatAsDate(time.Now())
	}

	return formatAsDate(stringTime)
}
