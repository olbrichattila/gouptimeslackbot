package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type requestInterface interface {
	get(string) (string, error)
}

type request struct {
}

// @todo add http user agent from config
func (r *request) get(url string) (string, error) {
	resp, err := http.Get(url)
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
