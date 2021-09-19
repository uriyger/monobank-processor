package main

import (
	"monobank-processor/app"
	"os"
)

func main() {
	os.Exit((&app.App{}).Run())
}
