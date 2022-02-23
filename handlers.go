package main

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func (c GitLab) Index(w http.ResponseWriter, r *http.Request) {
	gitlabToken, err := r.Cookie("gitlab_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	// Global variable projects
	projects = c.GetProjects(gitlabToken.Value)
	t := template.Must(template.New("index").Parse(TemplateIndex))
	if err = t.Execute(w, projects); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c GitLab) GetStatus(w http.ResponseWriter, r *http.Request) {
	gitlabToken, err := r.Cookie("gitlab_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if projects == nil {
		http.Error(w, "Projects list not built", http.StatusNotFound)
		return
	}
	projectIds, ok := r.URL.Query()["id"]
	if !ok {
		http.Error(w, "Project ID not found in request", http.StatusBadRequest)
		return
	}
	projectId := projectIds[0]
	projectIdForMap, _ := strconv.Atoi(projectId)
	if _, exists := projects[projectIdForMap]; !exists {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}
	commit := c.GetCommit(projectId, projects[projectIdForMap].DefaultBranch, gitlabToken.Value)
	allBuilds := c.GetBuilds(projectId, commit.Id, projects[projectIdForMap].DefaultBranch, gitlabToken.Value)
	builds := make(map[string]Build)
	var coverage *float32
	for _, b := range allBuilds {
		if _, ok := builds[b.Name]; !ok {
			builds[b.Name] = b
			if b.Coverage != nil {
				coverage = b.Coverage
			}
		}
	}
	funcMap := template.FuncMap{
		"Short": func(s string) string {
			if len(s) > 8 {
				return s[:8]
			} else {
				return s
			}
		},
		"Title": func(s string) string {
			return strings.Split(s, "\n")[0]
		},
	}
	t := template.Must(template.New("status").Funcs(funcMap).Parse(TemplateStatus))
	status := &Status{Commit: *commit, Project: projects[projectIdForMap], Coverage: coverage, Builds: builds, Url: c.url}
	if err = t.Execute(w, status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// API functions

func (c GitLab) IndexAPI(w http.ResponseWriter, r *http.Request) {
	gitlabToken, err := r.Cookie("gitlab_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	// Global variable projects
	body := c.GetApiRaw("", "0", "100", gitlabToken.Value)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func (c GitLab) GetStatusAPI(w http.ResponseWriter, r *http.Request) {
	gitlabToken, err := r.Cookie("gitlab_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	projectIds, ok := r.URL.Query()["id"]
	if !ok {
		http.Error(w, "Project ID not found in request", http.StatusBadRequest)
		return
	}
	projectId := projectIds[0]
	body := c.GetApiRaw(projectId+"/builds", "0", "100", gitlabToken.Value)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}
