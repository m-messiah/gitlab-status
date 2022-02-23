package main

type GitLab struct {
	url            string
	apiURL         string
	authorizeURL   string
	accessTokenURL string
	accessKey      string
	secretKey      string
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type Commit struct {
	Id         string
	Message    string
	AuthorName string `json:"author_name"`
}

type Build struct {
	Id       uint64
	Status   string
	Name     string
	Ref      string
	Coverage *float32
}

type ApiResponse struct {
	Id            int
	DefaultBranch string `json:"default_branch"`
	Name          string `json:"path_with_namespace"`
	BuildsEnabled bool   `json:"builds_enabled"`
	Commit        Commit
}

type Status struct {
	Commit   Commit
	Project  ApiResponse
	Builds   map[string]Build
	Coverage *float32
	Url      string
}
