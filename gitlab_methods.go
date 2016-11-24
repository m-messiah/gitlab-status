package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

func NewGitLab(gitlab_url, access_key, secret_key string) *GitLab {
	return &GitLab{
		url:			  gitlab_url,
		api_url:          gitlab_url + "/api/v3/projects/",
		authorize_url:    gitlab_url + "/oauth/authorize",
		access_token_url: gitlab_url + "/oauth/token",
		access_key:       access_key,
		secret_key:       secret_key,
	}
}

func (self GitLab) Authorize(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(self.authorize_url)
	u.RawQuery = url.Values{
		"client_id":     {self.access_key},
		"response_type": {"code"},
		"redirect_uri":  {REDIRECT_URL},
	}.Encode()
	http.Redirect(w, r, u.String(), http.StatusFound)
}

func (self GitLab) Authorized_response(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query()["code"][0]
	resp, err := http.PostForm(
		self.access_token_url,
		url.Values{
			"client_id":     {self.access_key},
			"client_secret": {self.secret_key},
			"code":          {code},
			"grant_type":    {"authorization_code"},
			"redirect_uri":  {REDIRECT_URL},
		})
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	var gitlab_response AuthResponse
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&gitlab_response)
	http.SetCookie(w, &http.Cookie{Name: "gitlab_token", Value: gitlab_response.Access_token})
	http.Redirect(w, r, "/", http.StatusFound)
}

func (self GitLab) Get_API_raw(endpoint, archived, per_page, token string) []byte {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", self.api_url+endpoint+"?per_page="+per_page+"&archived="+archived, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body
}

func (self GitLab) Get_API(endpoint, archived, per_page, token string, output interface{}) {
	body := self.Get_API_raw(endpoint, archived, per_page, token)
	json.Unmarshal(body, output)
}

func (self GitLab) Get_projects(token string) map[int]ApiResponse {
	var projects []ApiResponse
	self.Get_API("", "0", "100", token, &projects)
	projects_map := make(map[int]ApiResponse)
	for _, project := range projects {
		if project.Builds_enabled {
			projects_map[project.Id] = project
		}
	}
	return projects_map
}

func (self GitLab) Get_commit(project_id, branch, token string) *Commit {
	var response ApiResponse
	self.Get_API(project_id+"/repository/branches/"+branch, "0", "1", token, &response)
	return &response.Commit
}

func (self GitLab) Get_builds(project_id, commit_id, token string) []Build {
	var builds []Build
	url := project_id
	if commit_id != "" {
		url += "/repository/commits/" + commit_id
	}
	self.Get_API(url+"/builds", "0", "100", token, &builds)
	return builds
}
