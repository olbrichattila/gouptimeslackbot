package main

import "fmt"

type SlackPublisherSpy struct {
	called int
	err    error
}

func NewSlackPublisherSpy() *SlackPublisherSpy {
	return &SlackPublisherSpy{}
}

func (m *SlackPublisherSpy) Send(token, channelID, message string) error {
	m.called++
	return m.err
}

func (m *SlackPublisherSpy) withError(errorMessage string) *SlackPublisherSpy {
	m.err = fmt.Errorf(errorMessage)

	return m
}
