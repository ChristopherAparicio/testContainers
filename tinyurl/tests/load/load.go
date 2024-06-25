package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// Do load testing against the tinyurl service
func main() {
	httpClient := http.DefaultClient

	var sync sync.WaitGroup
	sync.Add(100)
	for range 100 {
		go postAndReadURL(&sync, httpClient)
	}

	sync.Wait()
	fmt.Println("Done")
}

type CreateResponseBody struct {
	OriginalURL  string `json:"originalUrl"`
	ShortenedURL string `json:"shortenedUrl"`
}

func postAndReadURL(sync *sync.WaitGroup, client *http.Client) error {
	defer sync.Done()

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/create", strings.NewReader(fmt.Sprintf(`{
		"originalUrl":"http://google.fr/%s"
	}`, uuid.New().String())),
	)

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	resp.Body.Close()

	var createResponseBody CreateResponseBody
	if err := json.Unmarshal(body, &createResponseBody); err != nil {
		return err
	}

	shortenedURL := createResponseBody.ShortenedURL

	req, err = http.NewRequest(http.MethodGet, "http://localhost:8080/"+shortenedURL, nil)
	if err != nil {
		return err
	}

	resp, err = client.Do(req)
	if err != nil {
		return err
	}

	fmt.Println(resp.StatusCode)
	return nil
}
