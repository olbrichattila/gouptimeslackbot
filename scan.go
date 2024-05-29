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

type uptimeInfo struct {
	lastInitiated time.Time
	errorCount    int
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
	// TODO separate slow warning message from uptime error
	go func() {
		uptime := uptimeInfo{lastInitiated: time.Now(), errorCount: 0}

		for {
			skipSending := uptime.errorCount > 0
			isSent := s.doScan(config, skipSending)
			if isSent {
				uptime.errorCount++
				if !skipSending {
					uptime.lastInitiated = time.Now()
				}
			}

			if uptime.errorCount > 1 && time.Since(uptime.lastInitiated).Seconds() > float64(config.RepeatNotificationDelay) {
				message := fmt.Sprintf(
					"Host: %s:\nUp bot report:\n\tStarded Date: %s\n\tDate: %s\n\tOccured  %d times",
					config.MonitorURL,
					uptime.lastInitiated.Format("2006-01-02 15:04:05"),
					time.Now().Format("2006-01-02 15:04:05"),
					uptime.errorCount,
				)

				uptime.errorCount = 0
				uptime.lastInitiated = time.Now()

				_ = s.publisher.Send(config.SlackBotToken, config.SlackChannelID, message)
			}

			time.Sleep(time.Duration(config.ScanFrequency) * time.Second)
		}
	}()
}

func (s *scanner) doScan(config configAccount, skipSending bool) bool {
	isMessageSent := false
	formattedDateTime := time.Now().Format("2006-01-02 15:04:05")
	elapsed, err := s.client.TestURL(config.HTTPUserAgent, config.MonitorURL, config.MonitorText)
	if err != nil {
		isMessageSent = true
		if !skipSending {
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
	}

	if config.SlowWarningLimit > 0 && elapsed > config.SlowWarningLimit {
		isMessageSent = true
		if !skipSending {
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

	return isMessageSent
}
