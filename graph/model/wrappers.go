package model

import "encoding/json"

// GithubUserInfo github user info wrapper
type GithubUserInfo struct {
	Message                   string  `json:"message,omitempty"`
	Name                      string  `json:"name"`
	Login                     string  `json:"login"`
	Company                   string  `json:"company"`
	Followers                 int     `json:"followers"`
	PublicRepos               int     `json:"public_repos"`
	AvgFollowersPerPublicRepo float32 `json:"avg_followers_per_public_repo"`
}

// String GithubUserInfo should comply with server.ICacheable which required String() implementation
func (i GithubUserInfo) String() string {
	result, err := json.Marshal(i)
	if err != nil {
		return ""
	}
	return string(result)
}

// ResultError error to include in result object
type ResultError struct {
	Message string `json:"message"`
}

// String return text representation of the struct
func (re ResultError) String() string {
	bytes, err := json.Marshal(re)
	if err == nil {
		return string(bytes)
	}
	return ""
}

// ResultRetrieveUsers result struct when calling retrieveUsers
type ResultRetrieveUsers struct {
	Users  []*GithubUserInfo `json:"users"`
	Errors []*ResultError    `json:"errors"`
}

// String return text representation of the struct
func (rru ResultRetrieveUsers) String() string {
	bytes, err := json.Marshal(rru)
	if err == nil {
		return string(bytes)
	}
	return ""
}
