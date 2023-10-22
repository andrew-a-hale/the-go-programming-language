package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"html/template"
)

const (
	APIURL        string = "https://api.github.com"
	ISSUEPAGE     string = "issues"
	MILESTONEPAGE string = "milestones"
	USERPAGE      string = "users"
)

type Data struct {
	issues     *Issues
	milestones *Milestones
	users      *Users
	cached     bool
}

var (
	issues     *Issues     = &Issues{}
	milestones *Milestones = &Milestones{}
	users      *Users      = &Users{}
)
var data = Data{
	issues:     issues,
	milestones: milestones,
	users:      users,
	cached:     false,
}

type Issue struct {
	Number    int       `json:"number"`
	HTMLURL   string    `json:"html_url"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	User      *User     `json:"user"`
	CreatedAt time.Time `json:"created_at"`
}

type Issues struct {
	TotalCount int
	Items      []Issue
}

type Milestone struct {
	Number    int       `json:"number"`
	HTMLURL   string    `json:"html_url"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	Creator   *User     `json:"creator"`
	CreatedAt time.Time `json:"created_at"`
}

type Milestones struct {
	TotalCount int
	Items      []Milestone
}

type User struct {
	Login   string `json:"login"`
	HTMLURL string `json:"html_url"`
}

type Users struct {
	TotalCount int
	Items      []User
}

func main() {
	http.HandleFunc("/", router)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func router(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "favicon") {
		homeView(w, r)
		return
	}

	path := strings.Split(r.URL.Path, "/")[1:]
	if len(path) < 3 {
		homeView(w, r)
		return
	}

	owner, repo := string(path[0]), string(path[1])
	cacheData(owner, repo)

	switch path[2] {
	case ISSUEPAGE:
		issues.issueView(w, r)
	case MILESTONEPAGE:
		milestones.milestoneView(w, r)
	case USERPAGE:
		users.userView(w, r)
	default:
		homeView(w, r)
	}
}

func homeView(w http.ResponseWriter, r *http.Request) {
	templ := `<h1>Home Page!</h1>
	<h2><a href='{{.Path}}/issues'>Issues</a></h2>
	<h2><a href='{{.Path}}/milestones'>Milestones</a></h2>
	<h2><a href='{{.Path}}/users'>Users</a></h2>`

	dat := struct {
		Path string
	}{
		r.URL.String(),
	}
	template.Must(template.New("home").Parse(templ)).Execute(w, dat)
}

func cacheData(owner, repo string) {
	if data.cached {
		return
	}

	cacheIssues(owner, repo)
	cacheMilestones(owner, repo)
	cacheUsers()
}

func cacheIssues(owner, repo string) error {
	url := fmt.Sprintf("%s/repos/%s/%s/issues", APIURL, owner, repo)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error calling %s: %s", url, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading resp %s: %s", url, err)
	}

	var issueList []Issue
	if err := json.Unmarshal(body, &issueList); err != nil {
		return fmt.Errorf("error parsing resp %s: %s", url, err)
	}
	issues.TotalCount = len(issueList)
	issues.Items = issueList

	return nil
}

func (issues *Issues) issueView(w http.ResponseWriter, r *http.Request) {
	const templ = `<h1>Issues Page</h1>
<h2>Count: {{.TotalCount}}</h2>
<hr>
<h2>Issues</h2>
<table>
<tr>
	<th>#</th>
	<th>State</th>
	<th>User</th>
	<th>Title</th>
</tr>
{{range .Items}}
<tr>
	<td><a href='{{.HTMLURL}}'>{{.Number}}</a></td>
	<td>{{.State}}</td>
	<td><a href='{{.User.HTMLURL}}'>{{.User.Login}}</a></td>
	<td><a href='{{.HTMLURL}}'>{{.Title}}</a></td>
</tr>
{{end}}
</table>`

	err := template.Must(template.New("issues").
		Parse(templ)).
		Execute(w, issues)
	if err != nil {
		fmt.Println("failed to build issue view")
	}
}

func cacheMilestones(owner, repo string) error {
	url := fmt.Sprintf("%s/repos/%s/%s/milestones", APIURL, owner, repo)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error calling %s: %s", url, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading resp %s: %s", url, err)
	}

	var milestoneList []Milestone
	if err := json.Unmarshal(body, &milestoneList); err != nil {
		return fmt.Errorf("error parsing resp %s: %s", url, err)
	}
	milestones.TotalCount = len(milestoneList)
	milestones.Items = milestoneList

	return nil
}

func (milestones *Milestones) milestoneView(w http.ResponseWriter, r *http.Request) {
	const templ = `<h1>Milestone Page</h1>
<h2>Count: {{.TotalCount}}</h2>
<hr>
<h2>Milestones</h2>
<table>
<tr>
	<th>#</th>
	<th>State</th>
	<th>Creator</th>
	<th>Title</th>
</tr>
{{range .Items}}
<tr>
	<td><a href='{{.HTMLURL}}'>{{.Number}}</a></td>
	<td>{{.State}}</td>
	<td><a href='{{.Creator.HTMLURL}}'>{{.Creator.Login}}</a></td>
	<td><a href='{{.HTMLURL}}'>{{.Title}}</a></td>
</tr>
{{end}}
</table>`

	err := template.Must(template.New("milestones").
		Parse(templ)).
		Execute(w, milestones)
	if err != nil {
		fmt.Println("failed to build milestones view")
	}
}

func cacheUsers() {
	var usersList []User
	for _, issue := range issues.Items {
		usersList = append(usersList, *issue.User)
	}
	for _, milestone := range milestones.Items {
		usersList = append(usersList, *milestone.Creator)
	}
	users.TotalCount = len(usersList)
	users.Items = usersList
}

func (users *Users) userView(w http.ResponseWriter, r *http.Request) {
	const templ = `<h1>Users Page</h1>
<h2>Count: {{.TotalCount}}</h2>
<hr>
<h2>Users</h2>
<table>
<tr>
	<th>User</th>
</tr>
{{range .Items}}
<tr>
	<td><a href='{{.HTMLURL}}'>{{.Login}}</a></td>
</tr>
{{end}}
</table>`

	err := template.Must(template.New("users").
		Parse(templ)).
		Execute(w, users)
	if err != nil {
		fmt.Println("failed to build users view")
	}
}
