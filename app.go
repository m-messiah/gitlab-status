package main

import (
	"net/http"
	"os"
)

var (
	RedirectURL     = os.Getenv("REDIRECT_URL")
	GitlabURL       = os.Getenv("GITLAB_URL")
	GitlabAppKey    = os.Getenv("GITLAB_APP_KEY")
	GitlabAppSecret = os.Getenv("GITLAB_APP_SECRET")

	projects map[int]ApiResponse
)

func main() {
	gitlab := *NewGitLab(GitlabURL, GitlabAppKey, GitlabAppSecret)
	http.HandleFunc("/", gitlab.Index)
	http.HandleFunc("/status/", gitlab.GetStatus)
	http.HandleFunc("/api/", gitlab.IndexAPI)
	http.HandleFunc("/api/status/", gitlab.GetStatusAPI)
	http.HandleFunc("/login", gitlab.Authorize)
	http.HandleFunc("/oauth-authorized", gitlab.AuthorizedResponse)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return
	}
}
