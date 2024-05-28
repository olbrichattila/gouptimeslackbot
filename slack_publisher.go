package main

import (
	"github.com/slack-go/slack"
)

type slackPublisherInterface interface {
	Send(string, string, string) error
}

type slackPublisher struct {
}

func newSlackPublisher() *slackPublisher {
	return &slackPublisher{}
}

func (m *slackPublisher) Send(token, channelID, message string) error {
	api := slack.New(token)

	msgOptions := slack.MsgOptionText(message, false)

	_, _, err := api.PostMessage(channelID, msgOptions)

	return err
}
