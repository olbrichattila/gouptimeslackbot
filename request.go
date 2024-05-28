package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type requestInterface interface {
	get(string, string) (string, error)
}

type request struct {
}

func (r *request) get(userAgent, url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if userAgent == "" {
		userAgent = defaultHTTPUserAgent
	}

	req.Header.Set("User-Agent", userAgent)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return "", fmt.Errorf("Response status code is not withing accepted range %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
