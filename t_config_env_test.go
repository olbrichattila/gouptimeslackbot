package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type configTestSuite struct {
	suite.Suite
}

func TestConfigRunner(t *testing.T) {
	suite.Run(t, new(configTestSuite))
}

func (t *configTestSuite) SetupTest() {
	fmt.Println(Yellow+"Running test config env: "+Green, t.T().Name()+Reset)
}

func (t *configTestSuite) TestConfigReturnsCorrectDefaultValuesFromEnv() {
	config := newConfig(false)
	t.Equal(60, config.getScanFrequency())

	accounts := config.getConfigAccounts()

	for _, a := range *accounts {
		t.Equal("", a.MonitorText)
		t.Equal("", a.MonitorURL)
		t.Equal("", a.SlackBotToken)
		t.Equal("", a.SlackChannelID)
		t.Equal(0, a.SlowWarningLimit)
		t.Equal("", a.HTTPUserAgent)
	}
}

func (t *configTestSuite) TestConfigReturnsCorrectValuesFromEnv() {
	t.setCustomEnvValues()
	config := newConfig(false)
	t.Equal(35, config.getScanFrequency())

	accounts := config.getConfigAccounts()

	for _, a := range *accounts {
		t.Equal("text", a.MonitorText)
		t.Equal("url", a.MonitorURL)
		t.Equal("token", a.SlackBotToken)
		t.Equal("channel id", a.SlackChannelID)
		t.Equal(1500, a.SlowWarningLimit)
		t.Equal("TestUserAgent/1.0", a.HTTPUserAgent)
	}
}

func (t *configTestSuite) setCustomEnvValues() {
	os.Setenv("SLACK_BOT_TOKEN", "token")
	os.Setenv("SLACK_CHANNEL_ID", "channel id")
	os.Setenv("MONITOR_URL", "url")
	os.Setenv("MONITOR_TEXT", "text")
	os.Setenv("SLOW_WARNING_LIMIT", "1500")
	os.Setenv("SCAN_FREQUENCY", "35")
	os.Setenv("HTTP_USER_AGENT", "TestUserAgent/1.0")
}
