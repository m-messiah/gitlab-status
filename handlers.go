package main

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func (self GitLab) Index(w http.ResponseWriter, r *http.Request) {
	gitlab_token, err := r.Cookie("gitlab_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	// Global variable projects
	projects = self.Get_projects(gitlab_token.Value)
	t := template.Must(template.New("index").Parse(TEMPLATE_INDEX))
	t.Execute(w, projects)
}

func (self GitLab) Get_status(w http.ResponseWriter, r *http.Request) {
	gitlab_token, err := r.Cookie("gitlab_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if projects == nil {
		http.Error(w, "Projects list not built", http.StatusNotFound)
		return
	}
	project_ids, ok := r.URL.Query()["id"]
	if !ok {
		http.Error(w, "Project ID not found in request", http.StatusBadRequest)
		return
	}
	project_id := project_ids[0]
	project_id_for_map, _ := strconv.Atoi(project_id)
	if _, exists := projects[project_id_for_map]; !exists {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}
	commit := self.Get_commit(project_id, projects[project_id_for_map].Default_branch, gitlab_token.Value)
	all_builds := self.Get_builds(project_id, commit.Id, projects[project_id_for_map].Default_branch, gitlab_token.Value)
	builds := make(map[string]Build)
	var coverage *float32
	for _, b := range all_builds {
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
				return string(s[:8])
			} else {
				return s
			}
		},
		"Title": func(s string) string {
			return strings.Split(s, "\n")[0]
		},
	}
	t := template.Must(template.New("status").Funcs(funcMap).Parse(TEMPLATE_STATUS))
	status := &Status{Commit: *commit, Project: projects[project_id_for_map], Coverage: coverage, Builds: builds, Url: self.url}
	t.Execute(w, status)
}

// API functions

func (self GitLab) IndexAPI(w http.ResponseWriter, r *http.Request) {
	gitlab_token, err := r.Cookie("gitlab_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	// Global variable projects
	body := self.Get_API_raw("", "0", "100", gitlab_token.Value)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (self GitLab) Get_statusAPI(w http.ResponseWriter, r *http.Request) {
	gitlab_token, err := r.Cookie("gitlab_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	project_ids, ok := r.URL.Query()["id"]
	if !ok {
		http.Error(w, "Project ID not found in request", http.StatusBadRequest)
		return
	}
	project_id := project_ids[0]
	body := self.Get_API_raw(project_id+"/builds", "0", "100", gitlab_token.Value)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
