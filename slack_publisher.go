package main

import (
	"github.com/slack-go/slack"
)

type SlackPublisherInterface interface {
	Send(string, string, string) error
}

type SlackPublisher struct {
}

func NewSlackPublisher() *SlackPublisher {
	return &SlackPublisher{}
}

func (m *SlackPublisher) Send(token, channelID, message string) error {
	api := slack.New(token)

	msgOptions := slack.MsgOptionText(message, false)

	_, _, err := api.PostMessage(channelID, msgOptions)

	return err
}
