package main

import (
	"log"
	"machshipgithubapi/server"
	"os"
	"strconv"
)

const (
	defaultPort          = 8777
	defaultGithubAPIURL  = "https://api.github.com"
	defaultGithubAPIUser = "users"
)

func main() {
	portEnv := os.Getenv("PORT")
	port := defaultPort
	if portEnv != "" {
		port, _ = strconv.Atoi(portEnv)
	}

	githubAPIURL := os.Getenv("GITHUB_API_URL")
	if githubAPIURL == "" {
		githubAPIURL = defaultGithubAPIURL
	}

	githubAPIUser := os.Getenv("GITHUB_API_USER")
	if githubAPIUser == "" {
		githubAPIUser = defaultGithubAPIUser
	}

	config := server.NewServerConfig("", port, githubAPIURL, githubAPIUser)
	s := server.NewServer(config)
	err := s.Serve()
	if err != nil {
		log.Fatalln("Server.Serve encounter error", err)
	}
}
