package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// NewGitLab creates GitLab client instance
func NewGitLab(gitlabURL, accessKey, secretKey string) *GitLab {
	return &GitLab{
		url:            gitlabURL,
		apiURL:         gitlabURL + "/api/v3/projects/",
		authorizeURL:   gitlabURL + "/oauth/authorize",
		accessTokenURL: gitlabURL + "/oauth/token",
		accessKey:      accessKey,
		secretKey:      secretKey,
	}
}

func (c GitLab) Authorize(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(c.authorizeURL)
	u.RawQuery = url.Values{
		"client_id":     {c.accessKey},
		"response_type": {"code"},
		"redirect_uri":  {RedirectURL},
	}.Encode()
	http.Redirect(w, r, u.String(), http.StatusFound)
}

func (c GitLab) AuthorizedResponse(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query()["code"][0]
	resp, err := http.PostForm(
		c.accessTokenURL,
		url.Values{
			"client_id":     {c.accessKey},
			"client_secret": {c.secretKey},
			"code":          {code},
			"grant_type":    {"authorization_code"},
			"redirect_uri":  {RedirectURL},
		})
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	var gitlabResponse AuthResponse
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&gitlabResponse); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "gitlab_token", Value: gitlabResponse.AccessToken})
	http.Redirect(w, r, "/", http.StatusFound)
}

func (c GitLab) GetApiRaw(endpoint, archived, perPage, token string) []byte {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", c.apiURL+endpoint+"?per_page="+perPage+"&archived="+archived, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	return body
}

func (c GitLab) GetApi(endpoint, archived, perPage, token string, output interface{}) {
	body := c.GetApiRaw(endpoint, archived, perPage, token)
	_ = json.Unmarshal(body, output)
}

func (c GitLab) GetProjects(token string) map[int]ApiResponse {
	var projects []ApiResponse
	c.GetApi("", "0", "100", token, &projects)
	projectsMap := make(map[int]ApiResponse)
	for _, project := range projects {
		if project.BuildsEnabled {
			projectsMap[project.Id] = project
		}
	}
	return projectsMap
}

func (c GitLab) GetCommit(projectId, branch, token string) *Commit {
	var response ApiResponse
	c.GetApi(projectId+"/repository/branches/"+branch, "0", "1", token, &response)
	return &response.Commit
}

func (c GitLab) GetBuilds(projectId, commitId, branch, token string) []Build {
	var allBuilds []Build
	builds := make([]Build, 0)
	buildsUrl := projectId
	if commitId != "" {
		buildsUrl += "/repository/commits/" + commitId
	}
	c.GetApi(buildsUrl+"/builds", "0", "100", token, &allBuilds)
	for _, build := range allBuilds {
		if build.Ref == branch {
			builds = append(builds, build)
		}
	}
	return builds
}
