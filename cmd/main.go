package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/src-d/go-git.v4"
)

func main() {
	generalCodeBaseOwners := getCodeBaseOwners()
	// languageOwners := getLanguageBasedOwners()
	write(generalCodeBaseOwners)
}

func getCodeBaseOwners() []owner {
	directory := os.Args[1]
	r, err := git.PlainOpen(directory)
	if err != nil {
		log.Fatal("could not open directory ", directory, err)
	}
	// w, err := r.Worktree()
	// if err != nil {
	// 	log.Fatal("could not create worktree", err)
	// }
	fmt.Println(r.Log(&git.LogOptions{}))
	//git shortlog -s -n --all --some-option 3
	return nil
}

func write(owners []owner) {
	f, err := os.Create("CODEOWNERS")
	if err != nil {
		log.Println("could not create file", err)
	}
	f.Write([]byte("# Example of CODEOWNER file syntax: \n" +
		"# https://help.github.com/en/github/creating-cloning-and-archiving-repositories/about-code-owners\n"))
	f.Close()
}

type owner struct {
	name     string
	email    string
	handle   string
	language string
	commits  int
}

// type language struct {
// 	Language string
// 	Percent  float64
// }
