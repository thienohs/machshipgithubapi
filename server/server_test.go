package server

import (
	"encoding/json"
	"fmt"
	"io"
	"machshipgithubapi/graph/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRetrieveUsers(t *testing.T) {
	tests := map[string]struct {
		Host                string
		Port                int
		GithubAPIUser       string
		Usernames           string
		ExpectedResultCount int
	}{
		"Test request with 1 username and checking the existence of username in response": {
			Host:                "",
			Port:                8777,
			GithubAPIUser:       "users",
			Usernames:           "thienohs",
			ExpectedResultCount: 1,
		},
		"Test request with 2 username (1 duplicate) and checking the existence of usernames in response": {
			Host:                "",
			Port:                8777,
			GithubAPIUser:       "users",
			Usernames:           "thienohs,apache,thienohs",
			ExpectedResultCount: 2,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			// Create an API test server
			githubAPITestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				urlParts := strings.Split(r.URL.String(), "/")
				if len(urlParts) > 0 {
					username := urlParts[len(urlParts)-1]

					responseStruct := struct {
						Login       string `json:"login"`
						Name        string `json:"name"`
						Company     string `json:"company"`
						Followers   int    `json:"followers"`
						PublicRepos int    `json:"public_repos"`
					}{
						Login:       username,
						Name:        username,
						Company:     username,
						Followers:   3,
						PublicRepos: 100,
					}

					jsonString, _ := json.Marshal(responseStruct)
					w.Write([]byte(jsonString))
				}
			}))
			defer githubAPITestServer.Close()

			// Using the API test server to mock API calling
			config := NewServerConfig(test.Host, test.Port, githubAPITestServer.URL, test.GithubAPIUser)
			s := NewServer(config)
			responseRecorder := httptest.NewRecorder()
			target := fmt.Sprintf("/retrieveUsers?usernames=%v", test.Usernames)
			request := httptest.NewRequest(http.MethodGet, target, nil)
			s.retrieveUsers(responseRecorder, request)

			response := responseRecorder.Result()
			defer response.Body.Close()
			responseData, err := io.ReadAll(response.Body)
			if err != nil {
				t.Errorf("expected no error when read data from response body, got err = %v", err)
			}

			jsonResponseData := &model.ResultRetrieveUsers{}
			err = json.Unmarshal([]byte(responseData), &jsonResponseData)
			if err != nil {
				t.Errorf("expected no error when unmarshal response data, got err = %v", err)
			} else {
				if len(jsonResponseData.Errors) > 0 {
					t.Errorf("expected no error from response, got errors = %v", jsonResponseData.Errors)
				} else if len(jsonResponseData.Users) != test.ExpectedResultCount {
					t.Errorf("expected number of result records = %v, got %v records", test.ExpectedResultCount, len(jsonResponseData.Users))
				} else {
					for _, eachUsername := range strings.Split(test.Usernames, ",") {
						foundRecord := false
						for _, eachUserInfo := range jsonResponseData.Users {
							if eachUserInfo.Login == eachUsername {
								foundRecord = true
								break
							}
						}

						if !foundRecord {
							t.Errorf("expected a record with username = %v, got no record", eachUsername)
						}
					}
				}
			}
		})
	}
}

func TestAverageFollowersPerPublicRepo(t *testing.T) {
	tests := map[string]struct {
		Host          string
		Port          int
		GithubAPIUser string
		Username      string
		ExpectedAvg   float32
	}{
		"Test average followers per public report, followers = 10, public_repos = 300": {
			Host:          "",
			Port:          8777,
			GithubAPIUser: "users",
			Username:      "abc",
			ExpectedAvg:   0.03,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			// Create an API test server
			githubAPITestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				responseStruct := struct {
					Login       string `json:"login"`
					Name        string `json:"name"`
					Company     string `json:"company"`
					Followers   int    `json:"followers"`
					PublicRepos int    `json:"public_repos"`
				}{
					Login:       "abc",
					Name:        "abc",
					Company:     "abc",
					Followers:   3,
					PublicRepos: 100,
				}

				jsonString, _ := json.Marshal(responseStruct)
				w.Write([]byte(jsonString))
			}))
			defer githubAPITestServer.Close()

			// Using the API test server to mock API calling
			config := NewServerConfig(test.Host, test.Port, githubAPITestServer.URL, test.GithubAPIUser)
			s := NewServer(config)
			responseRecorder := httptest.NewRecorder()
			target := fmt.Sprintf("/retrieveUsers?usernames=%v", test.Username)
			request := httptest.NewRequest(http.MethodGet, target, nil)
			s.retrieveUsers(responseRecorder, request)

			response := responseRecorder.Result()
			defer response.Body.Close()
			responseData, err := io.ReadAll(response.Body)
			if err == nil {
				jsonResponseData := &model.ResultRetrieveUsers{}
				err = json.Unmarshal([]byte(responseData), &jsonResponseData)
				if err == nil {
					if len(jsonResponseData.Errors) > 0 {
						t.Errorf("expected no error from response, got errors = %v", jsonResponseData.Errors)
					} else {
						foundRecord := false
						for _, eachUserInfo := range jsonResponseData.Users {
							if eachUserInfo.Login == test.Username {
								foundRecord = true

								if eachUserInfo.AvgFollowersPerPublicRepo != test.ExpectedAvg {
									t.Errorf("expected AvgFollowersPerPublicRepo = %v, got %v", test.ExpectedAvg, eachUserInfo.AvgFollowersPerPublicRepo)
								}

								break
							}
						}

						if !foundRecord {
							t.Errorf("expected a record with username = %v, got no record", test.Username)
						}
					}
				} else {
					t.Fail()
				}
			} else {
				t.Fail()
			}
		})
	}
}

func TestNotFoundUser(t *testing.T) {
	tests := map[string]struct {
		Host          string
		Port          int
		GithubAPIUser string
		Username      string
	}{
		"Test user not found": {
			Host:          "",
			Port:          8777,
			GithubAPIUser: "users",
			Username:      "notfound",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			// Create an API test server
			githubAPITestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				responseStruct := struct {
					Message string `json:"message"`
				}{
					Message: "Not Found",
				}

				jsonString, _ := json.Marshal(responseStruct)
				w.Write([]byte(jsonString))
			}))
			defer githubAPITestServer.Close()

			// Using the API test server to mock API calling
			config := NewServerConfig(test.Host, test.Port, githubAPITestServer.URL, test.GithubAPIUser)
			s := NewServer(config)
			responseRecorder := httptest.NewRecorder()
			target := fmt.Sprintf("/retrieveUsers?usernames=%v", test.Username)
			request := httptest.NewRequest(http.MethodGet, target, nil)
			s.retrieveUsers(responseRecorder, request)

			response := responseRecorder.Result()
			defer response.Body.Close()
			responseData, err := io.ReadAll(response.Body)
			if err == nil {
				jsonResponseData := &model.ResultRetrieveUsers{}
				err = json.Unmarshal([]byte(responseData), &jsonResponseData)
				if err == nil {
					// There should be an error
					if len(jsonResponseData.Errors) == 0 {
						t.Errorf("expected an error from response, got no errors = %v", jsonResponseData.Errors)
					}

					// There should be no user record
					if len(jsonResponseData.Users) > 0 {
						t.Errorf("expected no user records, got = %v", jsonResponseData.Users)
					}
				} else {
					t.Fail()
				}
			} else {
				t.Fail()
			}
		})
	}
}

func TestEmptyUsername(t *testing.T) {
	tests := map[string]struct {
		Host          string
		Port          int
		GithubAPIUser string
		Usernames     string
	}{
		"Test single empty username": {
			Host:          "",
			Port:          8777,
			GithubAPIUser: "users",
			Usernames:     "",
		},
		"Test multiple empty usernames": {
			Host:          "",
			Port:          8777,
			GithubAPIUser: "users",
			Usernames:     ",,,,,",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			// Create an API test server
			githubAPITestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				responseStruct := struct {
					Message string `json:"message"`
				}{
					Message: "Not Found",
				}

				jsonString, _ := json.Marshal(responseStruct)
				w.Write([]byte(jsonString))
			}))
			defer githubAPITestServer.Close()

			// Using the API test server to mock API calling
			config := NewServerConfig(test.Host, test.Port, githubAPITestServer.URL, test.GithubAPIUser)
			s := NewServer(config)
			responseRecorder := httptest.NewRecorder()
			target := fmt.Sprintf("/retrieveUsers?usernames=%v", test.Usernames)
			request := httptest.NewRequest(http.MethodGet, target, nil)
			s.retrieveUsers(responseRecorder, request)

			response := responseRecorder.Result()
			defer response.Body.Close()
			responseData, err := io.ReadAll(response.Body)
			if err == nil {
				jsonResponseData := &model.ResultRetrieveUsers{}
				err = json.Unmarshal([]byte(responseData), &jsonResponseData)

				if err == nil {
					// There should be no error
					if len(jsonResponseData.Errors) > 0 {
						t.Errorf("expected no error from response, got = %v", jsonResponseData.Errors)
					}

					// There should be no user record
					if len(jsonResponseData.Users) > 0 {
						t.Errorf("expected no user records, got = %v", jsonResponseData.Users)
					}
				} else {
					t.Fail()
				}
			} else {
				t.Fail()
			}
		})
	}
}

func TestEmptyUsernameTogetherWithNotFound(t *testing.T) {
	tests := map[string]struct {
		Host          string
		Port          int
		GithubAPIUser string
		Usernames     string
	}{
		"Test multiple empty usernames together with not found": {
			Host:          "",
			Port:          8777,
			GithubAPIUser: "users",
			Usernames:     ",,notfound,,,",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			// Create an API test server
			githubAPITestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				responseStruct := struct {
					Message string `json:"message"`
				}{
					Message: "Not Found",
				}

				jsonString, _ := json.Marshal(responseStruct)
				w.Write([]byte(jsonString))
			}))
			defer githubAPITestServer.Close()

			// Using the API test server to mock API calling
			config := NewServerConfig(test.Host, test.Port, githubAPITestServer.URL, test.GithubAPIUser)
			s := NewServer(config)
			responseRecorder := httptest.NewRecorder()
			target := fmt.Sprintf("/retrieveUsers?usernames=%v", test.Usernames)
			request := httptest.NewRequest(http.MethodGet, target, nil)
			s.retrieveUsers(responseRecorder, request)

			response := responseRecorder.Result()
			defer response.Body.Close()
			responseData, err := io.ReadAll(response.Body)
			if err == nil {
				jsonResponseData := &model.ResultRetrieveUsers{}
				err = json.Unmarshal([]byte(responseData), &jsonResponseData)

				if err == nil {
					// There should be one error
					if len(jsonResponseData.Errors) == 0 {
						t.Errorf("expected an error from response, got no errors = %v", jsonResponseData.Errors)
					}

					// There should be no user record
					if len(jsonResponseData.Users) > 0 {
						t.Errorf("expected no user records, got = %v", jsonResponseData.Users)
					}
				} else {
					t.Fail()
				}
			} else {
				t.Fail()
			}
		})
	}
}

func TestEmptyUsernameTogetherWithNotFoundAndValid(t *testing.T) {
	tests := map[string]struct {
		Host              string
		Port              int
		GithubAPIUser     string
		Usernames         string
		ExpectedUsernames string
	}{
		"Test multiple empty usernames together with not found and single valid": {
			Host:              "",
			Port:              8777,
			GithubAPIUser:     "users",
			Usernames:         ",,notfound,,abc,,,,",
			ExpectedUsernames: "abc",
		},
		"Test multiple empty usernames together with multiple not found and multiple valid": {
			Host:              "",
			Port:              8777,
			GithubAPIUser:     "users",
			Usernames:         ",,notfound1,,abc,,,cde,,notfound2,",
			ExpectedUsernames: "abc,cde",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			// Create an API test server
			githubAPITestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				urlParts := strings.Split(r.URL.String(), "/")
				if len(urlParts) > 0 {
					username := urlParts[len(urlParts)-1]
					if !strings.Contains(username, "notfound") {
						responseStruct := struct {
							Login       string `json:"login"`
							Name        string `json:"name"`
							Company     string `json:"company"`
							Followers   int    `json:"followers"`
							PublicRepos int    `json:"public_repos"`
						}{
							Login:       username,
							Name:        username,
							Company:     username,
							Followers:   3,
							PublicRepos: 100,
						}

						jsonString, _ := json.Marshal(responseStruct)
						w.Write([]byte(jsonString))
					} else {
						responseStruct := struct {
							Message string `json:"message"`
						}{
							Message: "Not Found",
						}

						jsonString, _ := json.Marshal(responseStruct)
						w.Write([]byte(jsonString))
					}
				}
			}))
			defer githubAPITestServer.Close()

			// Using the API test server to mock API calling
			config := NewServerConfig(test.Host, test.Port, githubAPITestServer.URL, test.GithubAPIUser)
			s := NewServer(config)
			responseRecorder := httptest.NewRecorder()
			target := fmt.Sprintf("/retrieveUsers?usernames=%v", test.Usernames)
			request := httptest.NewRequest(http.MethodGet, target, nil)
			s.retrieveUsers(responseRecorder, request)

			response := responseRecorder.Result()
			defer response.Body.Close()
			responseData, err := io.ReadAll(response.Body)
			if err == nil {
				jsonResponseData := &model.ResultRetrieveUsers{}
				err = json.Unmarshal([]byte(responseData), &jsonResponseData)

				if err == nil {
					// There should be an error
					if len(jsonResponseData.Errors) == 0 {
						t.Errorf("expected an error from response, got no errors = %v", jsonResponseData.Errors)
					}

					// There should be a user record
					expectedUsernames := strings.Split(test.ExpectedUsernames, ",")
					if len(jsonResponseData.Users) != len(expectedUsernames) {
						t.Errorf("expected %v user record(s), got = %v", len(expectedUsernames), len(jsonResponseData.Users))
					} else {
						for _, eachExpectedUsername := range expectedUsernames {
							foundRecord := false
							for _, eachUserRecord := range jsonResponseData.Users {
								if eachUserRecord.Login == eachExpectedUsername {
									foundRecord = true
									break
								}
							}

							if !foundRecord {
								t.Errorf("expected record for username %v, got none", eachExpectedUsername)
							}
						}
					}
				} else {
					t.Fail()
				}
			} else {
				t.Fail()
			}
		})
	}
}

func TestSortOrderOfRetrieveUsers(t *testing.T) {
	tests := map[string]struct {
		Host              string
		Port              int
		GithubAPIUser     string
		Usernames         string
		ExpectedUsernames string
	}{
		"Test whether the returned result is sorted by username alphabetically": {
			Host:              "",
			Port:              8777,
			GithubAPIUser:     "users",
			Usernames:         "c,b,a",
			ExpectedUsernames: "a,b,c",
		},
		"Test whether the returned result is sorted by username alphabetically (long username)": {
			Host:              "",
			Port:              8777,
			GithubAPIUser:     "users",
			Usernames:         "aba,aac,aab,aaa,aad",
			ExpectedUsernames: "aaa,aab,aac,aad,aba",
		},
		"Test whether the returned result is sorted by username alphabetically (long username, with invalid and empty)": {
			Host:              "",
			Port:              8777,
			GithubAPIUser:     "users",
			Usernames:         "aba,anotfound,aac,aab,bnotfound,aaa,,aad",
			ExpectedUsernames: "aaa,aab,aac,aad,aba",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			// Create an API test server
			githubAPITestServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				urlParts := strings.Split(r.URL.String(), "/")
				if len(urlParts) > 0 {
					username := urlParts[len(urlParts)-1]
					if !strings.Contains(username, "notfound") {
						responseStruct := struct {
							Login       string `json:"login"`
							Name        string `json:"name"`
							Company     string `json:"company"`
							Followers   int    `json:"followers"`
							PublicRepos int    `json:"public_repos"`
						}{
							Login:       username,
							Name:        username,
							Company:     username,
							Followers:   3,
							PublicRepos: 100,
						}

						jsonString, _ := json.Marshal(responseStruct)
						w.Write([]byte(jsonString))
					} else {
						responseStruct := struct {
							Message string `json:"message"`
						}{
							Message: "Not Found",
						}

						jsonString, _ := json.Marshal(responseStruct)
						w.Write([]byte(jsonString))
					}
				}
			}))
			defer githubAPITestServer.Close()

			// Using the API test server to mock API calling
			config := NewServerConfig(test.Host, test.Port, githubAPITestServer.URL, test.GithubAPIUser)
			s := NewServer(config)
			responseRecorder := httptest.NewRecorder()
			target := fmt.Sprintf("/retrieveUsers?usernames=%v", test.Usernames)
			request := httptest.NewRequest(http.MethodGet, target, nil)
			s.retrieveUsers(responseRecorder, request)

			response := responseRecorder.Result()
			defer response.Body.Close()
			responseData, err := io.ReadAll(response.Body)
			if err == nil {
				jsonResponseData := &model.ResultRetrieveUsers{}
				err = json.Unmarshal([]byte(responseData), &jsonResponseData)

				if err == nil {
					// There should be a user record
					expectedUsernames := strings.Split(test.ExpectedUsernames, ",")
					if len(jsonResponseData.Users) != len(expectedUsernames) {
						t.Errorf("expected %v user record(s), got = %v", len(expectedUsernames), len(jsonResponseData.Users))
					} else {
						for i, eachExpectedUsername := range expectedUsernames {
							foundRecord := false
							foundRecordAtIndex := -1
							for j, eachUserRecord := range jsonResponseData.Users {
								if eachUserRecord.Login == eachExpectedUsername {
									foundRecord = true
									foundRecordAtIndex = j
									break
								}
							}

							if !foundRecord {
								t.Errorf("expected record for username %v at index %d, got none", eachExpectedUsername, i)
							}

							if foundRecordAtIndex != i {
								t.Errorf("expected record for username %v at index %d, got result at index %d", eachExpectedUsername, i, foundRecordAtIndex)
							}
						}
					}
				} else {
					t.Fail()
				}
			} else {
				t.Fail()
			}
		})
	}
}

func TestServe(t *testing.T) {
	tests := map[string]struct {
		Host          string
		Port          int
		ExpectedError bool
	}{
		"No error": {
			Host:          "",
			Port:          8777,
			ExpectedError: false,
		},
		"Invalid port": {
			Host:          "",
			Port:          100000, // Invalid port
			ExpectedError: true,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			errSignal := make(chan error)
			defer close(errSignal)
			go func() {
				config := NewServerConfig(test.Host, test.Port, "https://api.github.com", "users")
				s := NewServer(config)
				defer s.Shutdown()
				err := s.Serve()
				errSignal <- err
			}()

			select {
			case err := <-errSignal:
				if !test.ExpectedError {
					t.Errorf("expected no error, got %v", err)
				}
			case <-time.After(1 * time.Second):
				if test.ExpectedError {
					t.Errorf("expected error, got timeout")
				}
				break
			}
		})
	}
}
