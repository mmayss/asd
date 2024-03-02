package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

func GetAllCommits(client http.Client, token string, repos []Repo) []Entry {
	var waitGroup sync.WaitGroup
	ch := make(chan Entry, len(repos))

	for _, repo := range repos {
		waitGroup.Add(1)
		go getCommits(client, token, repo, &waitGroup, ch)
	}

	waitGroup.Wait()
	close(ch)

	var result []Entry
	for res := range ch {
		result = append(result, res)
	}

	return result
}

func getCommits(client http.Client, token string, repo Repo, wg *sync.WaitGroup, ch chan<- Entry) {
	defer wg.Done()
	url := fmt.Sprintf("https://api.github.com/repos/samisul/%s/commits", repo.Name)

	resp, err := sendRequest(client, url, token)

	if err != nil {
		panic(err)
	}

	logEntry := resp.Header.Get("link")
	var links []string
	if logEntry != "" {
		for _, link := range strings.Split(logEntry, ",") {
			url := strings.Replace(strings.Split(link, ";")[0], "<", "", -1)
			url = strings.Replace(url, ">", "", -1)
			url = strings.TrimSpace(url)
			links = append(links, url)
		}
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var result []response

	json.Unmarshal(body, &result)

	for _, link := range links {
		resp, err := sendRequest(client, link, token)
		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		var tempResult []response

		if err := json.Unmarshal(body, &tempResult); err != nil {
			fmt.Println("Can not unmarshal JSON")
		}

		result = append(result, tempResult...)
	}

	ch <- Entry{CommitCount: len(result), Weight: repo.Weight}
}

func sendRequest(client http.Client, url string, token string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("since", time.Now().AddDate(0, 0, -7).Format(time.RFC3339))
	q.Add("per_page", "100")
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return client.Do(req)
}

type response struct {
	Commit struct {
		Committer struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"committer"`
		Message string `json:"message"`
	} `json:"commit"`
}

type Entry struct {
	CommitCount int
	Weight      int
}
