package model

import (
	"encoding/json"
	"testing"
)

func TestGithubUserInfoToString(t *testing.T) {
	tests := map[string]struct {
		Input *GithubUserInfo
	}{
		"GithubUserInfo keep message": {
			Input: &GithubUserInfo{
				Message:                   "Test",
				Name:                      "A",
				Login:                     "B",
				Company:                   "C",
				Followers:                 3,
				PublicRepos:               10,
				AvgFollowersPerPublicRepo: 0.3,
			},
		},
		"GithubUserInfo omit message": {
			Input: &GithubUserInfo{
				Message:                   "",
				Name:                      "A",
				Login:                     "B",
				Company:                   "C",
				Followers:                 3,
				PublicRepos:               10,
				AvgFollowersPerPublicRepo: 0.3,
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			result := test.Input.String()
			expectedBytes, _ := json.Marshal(test.Input)
			expected := string(expectedBytes)
			if result != expected {
				t.Errorf("expected %v, got %v", expected, result)
			}
		})
	}
}

func TestResultErrorToString(t *testing.T) {
	tests := map[string]struct {
		Input *ResultError
	}{
		"ResultError": {
			Input: &ResultError{
				Message: "Test",
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			result := test.Input.String()
			expectedBytes, _ := json.Marshal(test.Input)
			expected := string(expectedBytes)
			if result != expected {
				t.Errorf("expected %v, got %v", expected, result)
			}
		})
	}
}

func TestResultRetrieveUsersToString(t *testing.T) {
	tests := map[string]struct {
		Input *ResultRetrieveUsers
	}{
		"ResultRetrieveUsers": {
			Input: &ResultRetrieveUsers{
				Users: []*GithubUserInfo{
					{
						Message:                   "A",
						Name:                      "B",
						Login:                     "C",
						Company:                   "D",
						Followers:                 3,
						PublicRepos:               10,
						AvgFollowersPerPublicRepo: 0.3,
					},
				},
				Errors: []*ResultError{
					{
						Message: "A",
					},
				},
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			result := test.Input.String()
			expectedBytes, _ := json.Marshal(test.Input)
			expected := string(expectedBytes)
			if result != expected {
				t.Errorf("expected %v, got %v", expected, result)
			}
		})
	}
}
