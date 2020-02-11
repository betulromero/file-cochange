package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type config struct {
	help         bool
	countOfFiles int
}

func getConfig() config {
	cfg := config{}
	flag.BoolVar(&cfg.help, "help", false, "show help for command")
	flag.IntVar(&cfg.countOfFiles, "count", 10, "total number of files to be listed")
	flag.Parse()
	return cfg
}

func trimAndRemoveBlanks(items []string, trimStr string) []string {
	var retItems []string
	for _, s := range items {
		trimmed := strings.Trim(s, trimStr)
		if trimmed == "" {
			continue
		}
		retItems = append(retItems, trimmed)
	}
	return retItems
}

// gitChangedFiles returns all changed files in a commit
func gitChangedFiles() ([]string, error) {
	out, err := exec.Command("git", "-c", "status.relativePaths=false", "status", "-s").Output()
	if err != nil {
		return nil, fmt.Errorf("could not run git: %v", err)
	}

	lines := strings.Split(string(out), "\n")

	var files []string
	lines = trimAndRemoveBlanks(lines, " ")

	for _, l := range lines {
		parts := strings.Split(l, " ")
		if len(parts) != 2 || parts[0] != "M" {
			continue
		}
		files = append(files, parts[1])
	}
	return files, nil
}

// processAllCommits goes through all commits in a repo,
// saving all the data.
func processAllCommits() ([]string, error) {
	out, err := exec.Command("git", "--no-pager", "log", "--pretty=\"%H\"").Output()
	if err != nil {
		return nil, fmt.Errorf("couldn't run git: %v", err)
	}
	lines := trimAndRemoveBlanks(strings.Split(string(out), "\n"), "\" ")
	return lines, nil
}

func getAllFilesChangedInGivenCommit(sha string) ([]string, error) {
	out, err := exec.Command("git", "diff-tree", "--no-commit-id", "--name-only", "-r", sha).Output()
	if err != nil {
		return nil, fmt.Errorf("couldn't run git: %v", err)
	}
	lines := trimAndRemoveBlanks(strings.Split(string(out), "\n"), "\" ")
	return lines, nil
}

// FilePairings contains a map of files to the files that changed with it in a commit with
// its assocciated co-change counts.
type FilePairings map[string]map[string]int

func getAllChangedFilesInAllCommits() (FilePairings, error) {
	filePairings := make(FilePairings)

	commitshas, err := processAllCommits()
	if err != nil {
		return nil, err
	}
	for _, s := range commitshas {
		changedfiles, err := getAllFilesChangedInGivenCommit(s)
		if err != nil {
			return nil, err
		}
		for _, baseFile := range changedfiles {
			// we have a new file that we are going to add to
			// the map.
			if filePairings[baseFile] == nil { // Is the map instantiated for the file?
				filePairings[baseFile] = make(map[string]int)
			}
			// OK. Now that that is initialized, lets continue as normal...cd ..
			for _, friend := range changedfiles {
				filePairings[baseFile][friend]++
			}
		}
	}

	return filePairings, nil
}

func (fp FilePairings) getCochangeRatiOfTwoFiles(basefile, friend string) float64 {
	ratio := float64(fp[basefile][friend]) / float64(fp[basefile][basefile])
	return ratio
}

type fileScore struct {
	fname string
	count float64
}

func (fp FilePairings) generateReport() ([]fileScore, error) {
	score := make(map[string]float64)

	modified, err := gitChangedFiles()
	if err != nil {
		return nil, err
	}

otherFileLoop:
	for otherFile := range fp {
		for _, basefile := range modified {
			if basefile == otherFile {
				continue otherFileLoop
			}
		}

		var total float64
		for _, file := range modified {
			total += fp.getCochangeRatiOfTwoFiles(otherFile, file)

		}
		score[otherFile] = total
	}

	tmp := []fileScore{}
	for fname, count := range score {
		if count > 0 {
			tmp = append(tmp, fileScore{fname: fname, count: count})
		}
	}

	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].count > tmp[j].count
	})

	return tmp, nil
}

func main() {
	cfg := getConfig()
	if cfg.help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	pairings, err := getAllChangedFilesInAllCommits()
	if err != nil {
		fmt.Printf("unable to find file shas: %v", err)
		return
	}

	scores, err := pairings.generateReport()
	if err != nil {
		fmt.Printf("unable to find file shas: %v", err)
		return
	}

	var count int
	for _, score := range scores {
		if len(score.fname) > 30 {
			showName := "..." + score.fname[len(score.fname)-25:]
			fmt.Printf("%-40v%5.4f\n", showName, score.count)
		} else {
			fmt.Printf("%-40v%5.4f\n", score.fname, score.count)
		}
		count++
		if cfg.countOfFiles == count {
			break
		}
	}
}
