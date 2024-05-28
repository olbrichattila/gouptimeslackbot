package main

import (
	"fmt"
	"time"
)

type scannerInterface interface {
	Scan(config configAccount)
}

type scanner struct {
	client    upClientInterface
	publisher slackPublisherInterface
	logger    loggerInterface
}

func newScanner(
	client upClientInterface,
	publisher slackPublisherInterface,
	logger loggerInterface,
) *scanner {
	s := &scanner{}
	s.client = client
	s.publisher = publisher
	s.logger = logger

	return s
}

func (s *scanner) Scan(config configAccount) {
	go s.doScan(config)
}

func (s *scanner) doScan(config configAccount) {
	formattedDateTime := time.Now().Format("2006-01-02 15:04:05")
	elapsed, err := s.client.TestURL(config.HTTPUserAgent, config.MonitorURL, config.MonitorText)
	if err != nil {
		message := fmt.Sprintf(
			"Host: %s:\nUp bot report:\n\tDate: %s\n\tElapsed: %d miliseconds\n\tError:%v",
			config.MonitorURL,
			formattedDateTime,
			elapsed,
			err,
		)
		err = s.publisher.Send(config.SlackBotToken, config.SlackChannelID, message)
		if err != nil {
			s.logger.Log(err.Error())
		}
	}

	if config.SlowWarningLimit > 0 && elapsed > config.SlowWarningLimit {
		message := fmt.Sprintf(
			"Host: %s:\nSlow warning limit reached:\n\tDate: %s\n\tLimit: %d miliseconds\n\tElapsed: %d miliseconds",
			config.MonitorURL,
			formattedDateTime,
			config.SlowWarningLimit,
			elapsed,
		)
		err = s.publisher.Send(config.SlackBotToken, config.SlackChannelID, message)
		if err != nil {
			s.logger.Log(err.Error())
		}
	}
}
