package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

const (
	githubAccessToken = "YOUR_GITHUB_ACCESS_TOKEN"
	username          = "USERNAME" // Replace with the GitHub username
	cloneDirectory    = "./repositories" // The directory where repos will be cloned
)

func main() {
	ctx := context.Background()
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubAccessToken},
	)
	tokenClient := oauth2.NewClient(ctx, tokenService)

	client := github.NewClient(tokenClient)

	// Get all repositories for the given user
	opt := &github.RepositoryListOptions{Type: "owner", Sort: "updated"}
	repos, _, err := client.Repositories.List(ctx, username, opt)
	if err != nil {
		fmt.Printf("Error fetching repositories: %v\n", err)
		return
	}

	if len(repos) == 0 {
		fmt.Println("No accessible repositories found.")
		return
	}

	// Ensure the clone directory exists
	err = os.MkdirAll(cloneDirectory, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	// Clone each repository
	for _, repo := range repos {
		repoName := *repo.Name
		repoURL := *repo.CloneURL
		fmt.Printf("Cloning %s ...\n", repoName)

		// Define the local repo path
		repoPath := filepath.Join(cloneDirectory, repoName)
		if _, err := os.Stat(repoPath); !os.IsNotExist(err) {
			fmt.Printf("Repository already cloned: %s\n", repoPath)
			continue
		}

		// Clone the repo
		cmd := exec.Command("git", "clone", repoURL, repoPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error cloning repository %s: %v\n", repoName, err)
		}
	}

	fmt.Println("Finished cloning accessible repositories.")
}

