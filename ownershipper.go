package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"gopkg.in/src-d/go-git.v4"
)

func main() {

	generateCmd := flag.NewFlagSet("", flag.ExitOnError)
	numOwners := generateCmd.Int("numowners", 1, "max number of owners to assign")
	outFmt := generateCmd.String("out", "stdout", "output format is file or stdout")
	overwright := generateCmd.Bool("overwrite", false, "whether to overwrite an existing CODEOWNERS file")
	generateCmd.Parse(os.Args[2:]) // Args[0] = "ownershipper", Args[1] = <path to git project>

	owners := owners{generateCodeOwners(*numOwners)}
	if len(owners.ownerList) == 0 {
		log.Fatal("no owners were written")
	}
	owners.write(*outFmt, *overwright)
}

// genCodeOwners returns code owners sorted by descreasing commits up to the numowners flag
func generateCodeOwners(numOwners int) []owner {
	owners := []owner{}
	directory := os.Args[1]
	r, err := git.PlainOpen(directory)
	if err != nil {
		log.Fatal("could not open directory ", directory, err)
	}

	// CommitIter is a generic closable interface for iterating over commits.
	commitIter, err := r.Log(&git.LogOptions{Order: git.LogOrderCommitterTime})
	if err != nil {
		log.Fatal("could not create objectCommitIter", err)
	}
	authorMap := map[string]*owner{}
	for {
		commit, err := commitIter.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		if _, ok := authorMap[commit.Author.Email]; !ok {
			authorMap[commit.Author.Email] = &owner{
				email:   commit.Author.Email,
				name:    commit.Author.Name,
				commits: 1,
			}
			continue
		}
		authorMap[commit.Author.Email].commits++
	}
	commitIter.Close()

	for k := range authorMap {
		owner := *authorMap[k]
		owners = append(owners, owner)
	}

	// sort owners by number of commits
	sort.Slice(owners, func(i, j int) bool {
		return owners[i].commits > owners[j].commits
	})
	if numOwners > len(owners) {
		return owners
	}

	return owners[:numOwners]
}

// write owners to CODEOWNERS file
func (owners owners) write(outFmt string, overwright bool) {
	var ownersStr string
	for _, owner := range owners.ownerList {
		ownersStr += owner.email + " "
	}
	ownersStr = strings.TrimRight(ownersStr, " ")

	text := fmt.Sprintf("# Example of CODEOWNER file syntax:\n"+
		"# https://help.github.com/en/github/creating-cloning-and-archiving-repositories/about-code-owners\n"+
		"*\t%s", ownersStr)

	if outFmt == "file" {
		if !fileExists("CODEOWNERS") || overwright {
			f, err := os.Create("CODEOWNERS")
			if err != nil {
				log.Fatal("could not create file", err)
			}
			f.Write([]byte(text))
			f.Close()
		} else {
			log.Println("Did not overwrite existing CODEOWNERS file")
		}
	} else {
		fmt.Println(text)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

type owners struct {
	ownerList []owner
}

type owner struct {
	email    string
	name     string
	handle   string
	language string
	commits  int
}
