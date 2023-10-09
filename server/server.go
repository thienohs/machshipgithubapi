package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"machshipgithubapi/graph"
	"machshipgithubapi/graph/model"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

type Server struct {
	httpServer          *http.Server
	serverMux           *http.ServeMux
	config              *ServerConfig
	githubUserInfoCache *ServerCache[model.GithubUserInfo]
}

const (
	// GITHUB_API_MESSAGE_USER_NOT_FOUND this message return from Github API when github user is not found
	GITHUB_API_MESSAGE_USER_NOT_FOUND = "Not Found"
)

// NewServer return a new server instance
func NewServer(config *ServerConfig) *Server {
	serverMux := http.NewServeMux()
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.host, config.port),
		Handler: serverMux,
	}
	return &Server{
		httpServer:          httpServer,
		serverMux:           serverMux,
		githubUserInfoCache: NewServerCache[model.GithubUserInfo](config.defaultCachePartitions),
		config:              config,
	}
}

// retrieveUsers handling retrieving users
func (s *Server) retrieveUsers(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	usernamesFormValue := r.FormValue("usernames")
	processedUserMap := make(map[string]bool)
	resultObj := &model.ResultRetrieveUsers{
		Users:  make([]*model.GithubUserInfo, 0),
		Errors: make([]*model.ResultError, 0),
	}

	if len(usernamesFormValue) > 0 {
		// Split the usernames by separator ,
		usernames := strings.Split(usernamesFormValue, ",")

		if len(usernames) > 0 {
			for _, eachUsername := range usernames {
				// Validate length of username (trimmed)
				if len(strings.TrimSpace(eachUsername)) == 0 {
					continue
				}

				// Check if this login has been processed before
				_, processed := processedUserMap[eachUsername]
				if processed {
					// Skip if it has been processed before
					continue
				}

				// Mark this username has been processed
				processedUserMap[eachUsername] = true

				var userInfo *model.GithubUserInfo
				// Get from cache (if have)
				userInfo = s.githubUserInfoCache.Get(eachUsername)

				// Can not find in cache, calling API to get data
				if userInfo == nil {
					apiURL := fmt.Sprintf("%s/%s/%s", s.config.githubAPIURL, s.config.githubAPIUser, eachUsername)
					req, err := http.NewRequest(http.MethodGet, apiURL, nil)
					if err != nil {
						resultObj.Errors = append(resultObj.Errors, &model.ResultError{
							Message: fmt.Sprintf("encounter err for username %q: %v", eachUsername, err),
						})
					} else {
						req.Header.Add("Accept", "application/vnd.github+json")
						req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
						resp, err := client.Do(req)
						if err != nil {
							resultObj.Errors = append(resultObj.Errors, &model.ResultError{
								Message: fmt.Sprintf("encounter err for username %q: %v", eachUsername, err),
							})
						} else {
							responseData, err := io.ReadAll(resp.Body)
							if err != nil {
								resultObj.Errors = append(resultObj.Errors, &model.ResultError{
									Message: fmt.Sprintf("encounter err for username %q: %v", eachUsername, err),
								})
							} else {
								userInfo = &model.GithubUserInfo{}
								err = json.Unmarshal(responseData, &userInfo)
								if err != nil {
									resultObj.Errors = append(resultObj.Errors, &model.ResultError{
										Message: fmt.Sprintf("encounter err for username %q: %v", eachUsername, err),
									})
								} else {
									// Calculate AvgFollowersPerPublicRepo
									if userInfo.PublicRepos > 0 {
										userInfo.AvgFollowersPerPublicRepo = float32(userInfo.Followers) / float32(userInfo.PublicRepos)
									}

									// Cache the data
									s.githubUserInfoCache.Set(eachUsername, userInfo)
								}
							}
						}
					}
				}

				// Process result (whether from cache or from API call)
				if userInfo != nil {
					if userInfo.Message == GITHUB_API_MESSAGE_USER_NOT_FOUND {
						resultObj.Errors = append(resultObj.Errors, &model.ResultError{
							Message: fmt.Sprintf("username %q not found", eachUsername),
						})
					} else {
						// Add the user to result object's user list
						resultObj.Users = append(resultObj.Users, userInfo)
					}
				}
			}
		}
	}

	// Sort users data
	sort.SliceStable(resultObj.Users, func(i, j int) bool {
		return strings.Compare(resultObj.Users[i].Login, resultObj.Users[j].Login) < 0
	})

	// Write response (pretty JSON format)
	w.Header().Set("Content-Type", "application/json")
	responseData, err := json.MarshalIndent(resultObj, "", "    ")
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(responseData)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unexpected error: %v", err)))
	}
	// json.NewEncoder(w).Encode(resultObj)
}

// Serve server will use this function to register and serve handlers, this function will block and listen to connections
func (s *Server) Serve() error {
	// Register handler
	s.serverMux.HandleFunc("/retrieveUsers", s.retrieveUsers)

	// Register graphql
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		RetrieveUsersHandler: s.retrieveUsers,
	}}))
	s.serverMux.Handle("/graphql/playground", playground.Handler("GraphQL playground", "/graphql/query"))
	s.serverMux.Handle("/graphql/query", srv)

	// Listen and serve
	log.Println("Server is listening on", fmt.Sprintf("%s:%d", s.config.host, s.config.port))
	log.Println("GraphQL playground is available on", fmt.Sprintf("%s:%d/graphql/playground", s.config.host, s.config.port))
	return s.httpServer.ListenAndServe()
}

// Shutdown shutdown the server
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // Gracefully shutdown
	defer cancel()
	log.Println("Server is stopping...")
	return s.httpServer.Shutdown(ctx)
}
