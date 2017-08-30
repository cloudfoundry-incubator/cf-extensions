package main

import (
	"fmt"
	"os"
	"sync"

	"encoding/json"
	"io/ioutil"
	"os/signal"

	"github.com/jessevdk/go-flags"
	"github.com/robfig/cron"

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
	Schedule   string `short:"s" long:"schedule" description:"The schedule for when to run app in cron format, e.g., '@every 12h'"`
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
	if opts.Schedule == "" {
		app.Run(credentials.Orgs[0], credentials.TopicFilters)
	} else {
		cronJob := cron.New()
		cronJob.AddFunc(opts.Schedule, func() { app.Run(credentials.Orgs[0], credentials.TopicFilters) })

		fmt.Printf("Running as per cron schedule: `%s`\n", opts.Schedule)
		cronJob.Start()
	}

	fmt.Printf("Press Ctrl+C to end\n")
	waitForCtrlC()
	fmt.Println()
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

func waitForCtrlC() {
	var end_waiter sync.WaitGroup

	end_waiter.Add(1)

	var signal_channel chan os.Signal

	signal_channel = make(chan os.Signal, 1)
	signal.Notify(signal_channel, os.Interrupt)

	go func() {
		<-signal_channel
		end_waiter.Done()
	}()

	end_waiter.Wait()
}
