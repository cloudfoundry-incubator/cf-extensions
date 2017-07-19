package main

import (
	"fmt"

	"github.com/maximilien/cf-extensions/bot"
)

const VERSION = "0.1.0"

func main() {
	fmt.Printf("CF-Extensions github bot v%s\n", VERSION)

	app := bot.NewApp("30bb3196bd7d24eeba37b0e6def3e556b6ed49f1")
	app.Run("cloudfoundry-incubator", []string{"cf-extensions"})
}
