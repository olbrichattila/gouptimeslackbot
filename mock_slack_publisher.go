package main

import "fmt"

type slackPublisherSpy struct {
	called int
	err    error
}

func newSlackPublisherSpy() *slackPublisherSpy {
	return &slackPublisherSpy{}
}

func (m *slackPublisherSpy) Send(_, _, _ string) error {
	m.called++
	return m.err
}

func (m *slackPublisherSpy) withError(errorMessage string) *slackPublisherSpy {
	m.err = fmt.Errorf(errorMessage)

	return m
}
