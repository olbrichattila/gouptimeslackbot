package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type scanTestSuite struct {
	suite.Suite
	account *configAccount
}

func TestScanRunner(t *testing.T) {
	suite.Run(t, new(scanTestSuite))
}

func (t *scanTestSuite) SetupTest() {
	fmt.Println(Yellow+"Running test scan env: "+Green, t.T().Name()+Reset)
	t.account = &configAccount{
		SlackBotToken:           "token",
		SlackChannelID:          "channelId",
		MonitorURL:              "http://test.com",
		MonitorText:             "<html",
		SlowWarningLimit:        30,
		RepeatNotificationDelay: 20,
		ScanFrequency:           60,
	}
}

func (t *scanTestSuite) TestItMessageSentIfLoadingTimeExceededWarningLimit() {
	upClientSpy := newUpClientSpy().withElapsedTime(100)
	publisherSpy := newSlackPublisherSpy()
	loggerSpy := newLoggerSpy()

	scanner := newScanner(upClientSpy, publisherSpy, loggerSpy)

	scanner.doScan(*t.account, false)

	t.Equal(1, upClientSpy.called)

	t.Equal(1, publisherSpy.called)

	t.Equal(0, loggerSpy.called)
}

func (t *scanTestSuite) TestItMessageNotSentIfLoadingTimeUnderWarningLimit() {
	upClientSpy := newUpClientSpy().withElapsedTime(20)
	publisherSpy := newSlackPublisherSpy()
	loggerSpy := newLoggerSpy()

	scanner := newScanner(upClientSpy, publisherSpy, loggerSpy)

	scanner.doScan(*t.account, false)

	t.Equal(1, upClientSpy.called)

	t.Equal(0, publisherSpy.called)

	t.Equal(0, loggerSpy.called)
}

func (t *scanTestSuite) TestItMessageSentIfPageDidNotLoad() {
	upClientSpy := newUpClientSpy().withElapsedTime(10).withError("cannot connect")
	publisherSpy := newSlackPublisherSpy()
	loggerSpy := newLoggerSpy()

	scanner := newScanner(upClientSpy, publisherSpy, loggerSpy)

	scanner.doScan(*t.account, false)

	t.Equal(1, upClientSpy.called)

	t.Equal(1, publisherSpy.called)

	t.Equal(0, loggerSpy.called)
}

func (t *scanTestSuite) TestIfTwoMessageSentIfPageDidNotLoadAndTimeAlsoExceeded() {
	upClientSpy := newUpClientSpy().withElapsedTime(100).withError("cannot connect")
	publisherSpy := newSlackPublisherSpy()
	loggerSpy := newLoggerSpy()

	scanner := newScanner(upClientSpy, publisherSpy, loggerSpy)

	scanner.doScan(*t.account, false)

	t.Equal(1, upClientSpy.called)

	t.Equal(2, publisherSpy.called)

	t.Equal(0, loggerSpy.called)
}

func (t *scanTestSuite) TestSendMessageErrorsAreLogged() {
	upClientSpy := newUpClientSpy().withElapsedTime(100).withError("cannot connect")
	publisherSpy := newSlackPublisherSpy().withError("cannot send slack message")
	loggerSpy := newLoggerSpy()

	scanner := newScanner(upClientSpy, publisherSpy, loggerSpy)

	scanner.doScan(*t.account, false)

	t.Equal(1, upClientSpy.called)

	t.Equal(2, publisherSpy.called)

	t.Equal(2, loggerSpy.called)

	t.Equal("cannot send slack message", loggerSpy.lastMessage)
}

func (t *scanTestSuite) TestIfSkipSendingFlagIsConsidered() {
	upClientSpy := newUpClientSpy().withElapsedTime(100)
	publisherSpy := newSlackPublisherSpy()
	loggerSpy := newLoggerSpy()

	scanner := newScanner(upClientSpy, publisherSpy, loggerSpy)

	scanner.doScan(*t.account, true)

	t.Equal(1, upClientSpy.called)

	t.Equal(0, publisherSpy.called)

	t.Equal(0, loggerSpy.called)
}
