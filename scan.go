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
	lastReportTime time.Time
	errorCount     int
}

type uptimeInfos struct {
	errorInfo       uptimeInfo
	slowWarningInfo uptimeInfo
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
	go func() {
		now := time.Now()
		uptime := uptimeInfos{
			errorInfo:       uptimeInfo{lastReportTime: now, errorCount: 0},
			slowWarningInfo: uptimeInfo{lastReportTime: now, errorCount: 0},
		}

		for {
			skipErrorMessageSending := s.skipSending(uptime.errorInfo)
			skipSlowWarningMessageSending := s.skipSending(uptime.slowWarningInfo)
			errorMessageSent, slowWarningMessageSent := s.doScan(config, skipErrorMessageSending, skipSlowWarningMessageSending)

			s.updateMessageSentStatus(&uptime.errorInfo, errorMessageSent, skipErrorMessageSending)
			s.updateMessageSentStatus(&uptime.slowWarningInfo, slowWarningMessageSent, skipSlowWarningMessageSending)

			s.sendAggregateMessage(config, &uptime.errorInfo, float64(config.RepeatNotificationDelay), "Page load error aggretate message")
			s.sendAggregateMessage(config, &uptime.slowWarningInfo, float64(config.RepeatNotificationDelay), "Slow page load warning aggregate message")

			time.Sleep(time.Duration(config.ScanFrequency) * time.Second)
		}
	}()
}

func (s *scanner) updateMessageSentStatus(uptimeinfo *uptimeInfo, messageSent, skipSending bool) {
	if messageSent {
		uptimeinfo.errorCount++
		if !skipSending {
			uptimeinfo.lastReportTime = time.Now()
		}
	}
}

func (s *scanner) sendAggregateMessage(config configAccount, uptimeinfo *uptimeInfo, delay float64, message string) {
	if uptimeinfo.errorCount > 1 && time.Since(uptimeinfo.lastReportTime).Seconds() > delay {
		message := fmt.Sprintf(
			"%s\nHost: %s:\nUp bot report:\n\tStarded Date: %s\n\tDate: %s\n\tOccured  %d times\nThe message is sent after %.0f seconds",
			message,
			config.MonitorURL,
			uptimeinfo.lastReportTime.Format("2006-01-02 15:04:05"),
			time.Now().Format("2006-01-02 15:04:05"),
			uptimeinfo.errorCount,
			delay,
		)

		uptimeinfo.errorCount = 0
		uptimeinfo.lastReportTime = time.Now()

		_ = s.publisher.Send(config.SlackBotToken, config.SlackChannelID, message)
	}
}

func (s *scanner) skipSending(uptimeInfo uptimeInfo) bool {
	return uptimeInfo.errorCount > 0
}

func (s *scanner) doScan(config configAccount, skipSendingErrorReport, skipSendingSlowWarningReport bool) (bool, bool) {
	errorMessageSent := false
	slowWarningMessageSent := false
	formattedDateTime := time.Now().Format("2006-01-02 15:04:05")
	elapsed, err := s.client.TestURL(config.HTTPUserAgent, config.MonitorURL, config.MonitorText)
	if err != nil {
		errorMessageSent = true
		if !skipSendingErrorReport {
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
		slowWarningMessageSent = true
		if !skipSendingSlowWarningReport {
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

	return errorMessageSent, slowWarningMessageSent
}
