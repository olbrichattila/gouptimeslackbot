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
	publisher SlackPublisherInterface
}

func newScanner(client upClientInterface, publisher SlackPublisherInterface) *scanner {
	s := &scanner{}
	s.client = client
	s.publisher = publisher

	return s
}

func (s *scanner) Scan(config configAccount) {
	go s.doScan(config)
}

func (s *scanner) doScan(config configAccount) {
	formattedDateTime := time.Now().Format("2006-01-02 15:04:05")
	elapsed, err := s.client.TestUrl(config.MonitorUrl, config.MonitorText)
	if err != nil {
		message := fmt.Sprintf(
			"Host: %s:\nUp bot report:\n\tDate: %s\n\tElapsed: %d miliseconds\n\tError:%v",
			config.MonitorUrl,
			formattedDateTime,
			elapsed,
			err,
		)
		err = s.publisher.Send(config.SlackBotToken, config.SlackChannelId, message)
		if err != nil {
			// this shold be logged instead
			fmt.Println(err)
		}
	}

	if config.SlowWarningLimit > 0 && elapsed > config.SlowWarningLimit {
		message := fmt.Sprintf(
			"Host: %s:\nSlow warning limit reached:\n\tDate: %s\n\tLimit: %d miliseconds\n\tElapsed: %d miliseconds",
			config.MonitorUrl,
			formattedDateTime,
			config.SlowWarningLimit,
			elapsed,
		)
		err = s.publisher.Send(config.SlackBotToken, config.SlackChannelId, message)
		if err != nil {
			// this shold be logged instead
			fmt.Println(err)
		}
	}
}
