package main

import (
	"fmt"
	"os"

	"encoding/json"
	"github.com/cloudfoundry-incubator/cf-extensions/bot"
	"io/ioutil"
)

const VERSION = "0.9.0"

type Credentials struct {
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	AccessToken  string   `json:"accessToken"`
	TopicFilters []string `json:"topicFilters"`
	Orgs         []string `json:"orgs"`
}

const credentialsExample = `{
	"username": "john",
	"email": "john@smith.com",
	"accessToken": "ADD-GITHUB-ACCESS-TOKEN-HERE",
	"topicFilters": ["test", "topic"],
	"orgs": ["org0", "org1"]
}
`

func main() {
	fmt.Printf("CF-Extensions github bot v%s\n", VERSION)

	if len(os.Args) < 2 {
		usage()
		os.Exit(0)
	}

	credentials, err := parseCredentials(os.Args[1])
	if err != nil {
		fmt.Printf("Error parsing credentials from file: %s, message: %s\n", os.Args[1], err.Error())
		usage()
		os.Exit(1)
	}

	app := bot.NewApp(credentials.AccessToken, credentials.Username, credentials.Email)
	app.Run(credentials.Orgs[0], credentials.TopicFilters)
}

func usage() {
	fmt.Printf("Usage:\n")
	fmt.Printf("  $cf-extensions credentials.json\n")
	fmt.Printf("    where credentials.json contains GitHub credentials and info, e.g., \n%s\n", credentialsExample)
}

func parseCredentials(filePath string) (Credentials, error) {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Credentials{}, err
	}

	credentials := Credentials{}
	err = json.Unmarshal(fileBytes, &credentials)
	if err != nil {
		return Credentials{}, err
	}

	return credentials, nil
}
