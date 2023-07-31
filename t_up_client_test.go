package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

const (
	Reset  = "\x1b[0m"
	Green  = "\x1b[32m"
	Yellow = "\x1b[33m"
)

type upClientTestSuite struct {
	suite.Suite
}

func TestUpClientRunner(t *testing.T) {
	suite.Run(t, new(upClientTestSuite))
}

func (t *upClientTestSuite) SetupTest() {
	fmt.Println(Yellow+"Running up client test: "+Green, t.T().Name()+Reset)
}

func (t *upClientTestSuite) TestElapsedMilisecondsMeasured() {
	mockRequestSpy := newMockRequestSpy()
	client := newUpClient(mockRequestSpy)

	startTime := time.Now()
	time.Sleep(time.Duration(100) * time.Millisecond)

	elapsed := client.elapsedTimeInMiliseconds(startTime)

	t.Greater(elapsed, 0)
}

func (t *upClientTestSuite) TestUpClientReturnsTrueIfNoContent() {
	mockRequestSpy := newMockRequestSpy().withDelay(10)
	client := newUpClient(mockRequestSpy)

	elapsed, err := client.TestUrl("https://google.com", "")

	t.Nil(err)

	t.Greater(elapsed, 0)
}

func (t *upClientTestSuite) TestUpClientReturnsErrorIfPageNotExist() {
	mockRequestSpy := newMockRequestSpy().withError().withDelay(10)
	client := newUpClient(mockRequestSpy)
	elapsed, err := client.TestUrl("https://dssddasda.com", "")

	t.NotNil(err)

	t.Greater(elapsed, 0)
}

func (t *upClientTestSuite) TestUpClientReturnsErrorIfContentDoesNotMatch() {
	mockRequestSpy := newMockRequestSpy().withResponse("<html").withDelay(10)
	client := newUpClient(mockRequestSpy)
	elapsed, err := client.TestUrl("https://google.com", "itisnotmatch")

	t.NotNil(err)

	t.Greater(elapsed, 0)
}

func (t *upClientTestSuite) TestUpClientReturnsNoErrorIfContentMatches() {
	mockRequestSpy := newMockRequestSpy().withResponse("<html").withDelay(10)
	client := newUpClient(mockRequestSpy)
	elapsed, err := client.TestUrl("https://google.com", "<html")

	t.Nil(err)

	t.Greater(elapsed, 0)
}
