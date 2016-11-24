package main

type GitLab struct {
	url				 string
	api_url          string
	authorize_url    string
	access_token_url string
	access_key       string
	secret_key       string
}

type AuthResponse struct {
	Access_token string
	Token_type   string
}

type Commit struct {
	Id          string
	Message     string
	Author_name string
}

type Build struct {
	Id       uint64
	Status   string
	Name     string
	Coverage *float32
}

type ApiResponse struct {
	Id             int
	Default_branch string
	Name           string `json:"path_with_namespace"`
	Builds_enabled bool
	Commit         Commit
}

type Status struct {
	Commit   Commit
	Project  ApiResponse
	Builds   map[string]Build
	Coverage *float32
	Url		 string
}
