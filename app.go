package main

import (
	"net/http"
	"os"
)

var REDIRECT_URL = os.Getenv("REDIRECT_URL")
var GITLAB_URL = os.Getenv("GITLAB_URL")
var GITLAB_APP_KEY = os.Getenv("GITLAB_APP_KEY")
var GITLAB_APP_SECRET = os.Getenv("GITLAB_APP_SECRET")
var projects map[int]ApiResponse

func main() {
	gitlab := *NewGitLab(GITLAB_URL, GITLAB_APP_KEY, GITLAB_APP_SECRET)
	http.HandleFunc("/", gitlab.Index)
	http.HandleFunc("/status/", gitlab.Get_status)
	http.HandleFunc("/api/", gitlab.IndexAPI)
	http.HandleFunc("/api/status/", gitlab.Get_statusAPI)
	http.HandleFunc("/login", gitlab.Authorize)
	http.HandleFunc("/oauth-authorized", gitlab.Authorized_response)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	http.ListenAndServe(":"+port, nil)
}
