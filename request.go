package main

import (
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
