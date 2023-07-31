package main

import (
	"fmt"
	"strings"
	"time"
)

type upClientInterface interface {
	TestUrl(string, string) (int, error)
}

type upClient struct {
	request requestInterface
}

func newUpClient(request requestInterface) *upClient {
	return &upClient{request: request}
}

func (u *upClient) TestUrl(url, contains string) (int, error) {
	startTime := time.Now()
	resp, err := u.request.get(url)
	if err != nil {
		return u.elapsedTimeInMiliseconds(startTime), err
	}

	elapsedMiliseconds := u.elapsedTimeInMiliseconds(startTime)
	if contains == "" {
		return elapsedMiliseconds, nil
	}

	if !strings.Contains(resp, contains) {
		return elapsedMiliseconds, fmt.Errorf("The page does not contain %s", contains)
	}

	return elapsedMiliseconds, nil
}

func (u *upClient) elapsedTimeInMiliseconds(startTime time.Time) int {
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)

	return int(elapsedTime.Milliseconds())
}
