package main

import (
	"fmt"
	"os"

	"encoding/json"
	"io/ioutil"

	"github.com/jessevdk/go-flags"

	"github.com/cloudfoundry-incubator/cf-extensions/bot"
)

const VERSION = "0.9.1"

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

type Options struct {
	Verbose    bool   `short:"v" long:"verbose" description:"Show verbose debug information"`
	File       string `short:"f" long:"file" description:"The credential file"`
	Credential string `short:"c" long:"credentials" description:"The credentials string in JSON"`
}

func main() {
	fmt.Printf("CF-Extensions github bot v%s\n", VERSION)

	opts := Options{}

	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	credentials := Credentials{}
	if opts.File != "" {
		credentials, err = parseCredentialsFile(opts.File)
		if err != nil {
			fmt.Printf("Error parsing credentials from file: %s, message: %s\n", opts.File, err.Error())
			os.Exit(1)
		}
	} else if opts.Credential != "" {
		credentials, err = parseCredentials(opts.Credential)
		if err != nil {
			fmt.Printf("Error parsing credentials from JSON:\n `%s`\n   Message: %s\n", opts.Credential, err.Error())
			os.Exit(1)
		}
	} else {
		fmt.Println("No credentials passed")
		os.Exit(1)
	}

	app := bot.NewApp(credentials.AccessToken, credentials.Username, credentials.Email)
	app.Run(credentials.Orgs[0], credentials.TopicFilters)
}

func parseCredentialsFile(filePath string) (Credentials, error) {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Credentials{}, err
	}

	return parseCredentials(string(fileBytes))
}

func parseCredentials(credentialJson string) (Credentials, error) {
	credentials := Credentials{}
	err := json.Unmarshal([]byte(credentialJson), &credentials)
	if err != nil {
		return Credentials{}, err
	}

	return credentials, nil
}
