package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"golang.org/x/exp/slices"
)

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string `json:"title"`
	State     string `json:"state"`
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string    // in Markdown format
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

const githubDomain string = "https://api.github.com"

var (
	cmds      []string = []string{"open", "close", "update", "read"}
	githubPat string   = os.Getenv("GITHUB_PAT")
	editor    string   = os.Getenv("TEXT_EDITOR")
)

func main() {
	cmd := os.Args[1]
	if !slices.Contains(cmds, cmd) {
		log.Fatalln("invalid cmd")
	}

	owner := os.Args[2]
	repo := os.Args[3]

	var id int64
	var err error
	if slices.Contains([]string{"close", "update", "read"}, cmd) {
		id, err = strconv.ParseInt(os.Args[4], 10, 64)
		if err != nil {
			log.Fatalln("missing id when calling close, update or read")
		}
	}

	switch cmd {
	case "read":
		resp, err := readIssue(owner, repo, id)
		if err != nil {
			log.Fatalf("issue %v not readable\n", id)
		}
		fmt.Printf("%v\n", resp)
	case "open":
		issue := &Issue{}
		issue, err = edit(issue)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := openIssue(owner, repo, issue)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", resp)
	case "close":
		issue, err := readIssue(owner, repo, id)
		if err != nil {
			log.Fatal(err)
		}

		issue.State = "closed"
		resp, err := updateIssue(owner, repo, issue)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%v\n", resp)
	case "update":
		issue, err := readIssue(owner, repo, id)
		if err != nil {
			log.Fatal(err)
		}

		issue, err = edit(issue)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := updateIssue(owner, repo, issue)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", resp)
	}
}

func request(url string, method string, body *Issue) (*Issue, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal %v to []byte", data)
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, errors.New("failed to create request")
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+githubPat)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("request failed")
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read response")
	}
	defer resp.Body.Close()

	json.Unmarshal(respBody, &body)

	return body, nil
}

func openIssue(owner, repo string, issue *Issue) (*Issue, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues", githubDomain, owner, repo)
	resp, err := request(url, "POST", issue)
	if err != nil {
		return nil, fmt.Errorf("failed to open issue: %s", err)
	}

	return resp, nil
}

func updateIssue(owner, repo string, issue *Issue) (*Issue, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues/%d", githubDomain, owner, repo, issue.Number)
	resp, err := request(url, "PATCH", issue)
	if err != nil {
		return nil, fmt.Errorf("failed to open issue %v", issue)
	}

	return resp, nil
}

func readIssue(owner, repo string, issueId int64) (*Issue, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues/%d", githubDomain, owner, repo, issueId)
	resp, err := request(url, "GET", &Issue{})
	if err != nil {
		return nil, fmt.Errorf("failed to open issue %v", issueId)
	}

	return resp, nil
}

func edit(issue *Issue) (*Issue, error) {
	// open in text editor and create a new issue from text file
	issueJson, err := json.MarshalIndent(issue, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed to edit issue %v", issue)
	}

	filename := fmt.Sprintf("issue-%d.txt", issue.Number)
	os.WriteFile(filename, issueJson, 0644)
	defer os.Remove(filename)

	cmd := exec.Command(editor, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to open %s with file %s", editor, filename)
	}

	updatedJson, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.New("failed to open updated issue")
	}

	json.Unmarshal(updatedJson, issue)
	fmt.Printf("%v\n", issue)

	return issue, nil
}
