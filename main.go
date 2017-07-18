package main

func main() {
	app := NewApp("30bb3196bd7d24eeba37b0e6def3e556b6ed49f1")
	app.Run("cloudfoundry-incubator", []string{"cf-extensions"})
}
