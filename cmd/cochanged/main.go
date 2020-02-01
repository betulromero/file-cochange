package main

import (
	"flag"
	"fmt"
	"os"
)

type config struct {
	help bool
	file string
}

func getConfig() config {
	cfg := config{}
	flag.BoolVar(&cfg.help, "help", false, "show help for command")
	flag.StringVar(&cfg.file, "f", "", "specific file to check change probability")
	flag.Parse()
	return cfg
}

// gitChangedFiles returns all changed files in a commit
func gitChangedFiles() []string {
	return []string{"cmd/cochanged/main.go", "README.md"}
}

// processAllCommits goes through all commits in a repo,
// saving all the data.
func processAllCommits() {

}

func main() {
	cfg := getConfig()
	if cfg.help {
		fmt.Println("Usage: cochanged ")
		os.Exit(0)
	}
	changedFiles := gitChangedFiles()

	fmt.Println("hello")
}
