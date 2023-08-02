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
		SlackBotToken:    "token",
		SlackChannelId:   "channelId",
		MonitorUrl:       "http://test.com",
		MonitorText:      "<html",
		SlowWarningLimit: 30,
	}
}

func (t *scanTestSuite) TestItMessageSentIfLoadingTimeExceededWarningLimit() {
	upClientSpy := newUpClientSpy().withElapsedTime(100)
	publisherSpy := NewSlackPublisherSpy()

	scanner := newScanner(upClientSpy, publisherSpy)

	scanner.doScan(*t.account)

	t.Equal(1, upClientSpy.called)

	t.Equal(1, publisherSpy.called)
}

func (t *scanTestSuite) TestItMessageNotSentIfLoadingTimeUnderWarningLimit() {
	upClientSpy := newUpClientSpy().withElapsedTime(20)
	publisherSpy := NewSlackPublisherSpy()

	scanner := newScanner(upClientSpy, publisherSpy)

	scanner.doScan(*t.account)

	t.Equal(1, upClientSpy.called)

	t.Equal(0, publisherSpy.called)
}

func (t *scanTestSuite) TestItMessageSentIfPageDidNotLoad() {
	upClientSpy := newUpClientSpy().withElapsedTime(10).withError("cannot connect")
	publisherSpy := NewSlackPublisherSpy()

	scanner := newScanner(upClientSpy, publisherSpy)

	scanner.doScan(*t.account)

	t.Equal(1, upClientSpy.called)

	t.Equal(1, publisherSpy.called)
}

func (t *scanTestSuite) TestIfTwoMessageSentIfPageDidNotLoadAndTimeAlsoExceeded() {
	upClientSpy := newUpClientSpy().withElapsedTime(100).withError("cannot connect")
	publisherSpy := NewSlackPublisherSpy()

	scanner := newScanner(upClientSpy, publisherSpy)

	scanner.doScan(*t.account)

	t.Equal(1, upClientSpy.called)

	t.Equal(2, publisherSpy.called)
}
