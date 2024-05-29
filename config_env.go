package main

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type envConfig struct {
	loadEnv bool
	config  *config
}

func newConfig(laodEnv bool) *envConfig {
	c := &envConfig{
		loadEnv: laodEnv,
		config:  &config{},
	}
	c.load()
	return c
}

func (c *envConfig) load() {
	if c.loadEnv {
		err := c.loadDotEnv()
		if err != nil {
			panic("Error loading configuration: " + err.Error())
		}
	}

	account := configAccount{
		SlackBotToken:           os.Getenv("SLACK_BOT_TOKEN"),
		SlackChannelID:          os.Getenv("SLACK_CHANNEL_ID"),
		MonitorURL:              os.Getenv("MONITOR_URL"),
		MonitorText:             os.Getenv("MONITOR_TEXT"),
		HTTPUserAgent:           os.Getenv("HTTP_USER_AGENT"),
		SlowWarningLimit:        c.asInt("SLOW_WARNING_LIMIT", 3000),
		RepeatNotificationDelay: c.asInt("REPEAT_NOTIFICATION_DELAY", 3600),
		ScanFrequency:           c.asInt("SCAN_FREQUENCY", 60),
	}

	c.config.Accounts = []configAccount{account}
}

func (c *envConfig) getConfigAccounts() *[]configAccount {
	return &c.config.Accounts
}

func (c *envConfig) asInt(env string, def int) int {
	value := os.Getenv(env)
	if value == "" {
		return def
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return def
	}

	return result
}

func (c *envConfig) loadDotEnv() error {
	if fileExists("./.env") {
		if err := godotenv.Load(); err != nil {
			return err
		}
	}

	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
