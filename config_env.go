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
		c.loadDotEnv()
	}

	account := configAccount{
		SlackBotToken:    os.Getenv("SLACK_BOT_TOKEN"),
		SlackChannelId:   os.Getenv("SLACK_CHANNEL_ID"),
		MonitorUrl:       os.Getenv("MONITOR_URL"),
		MonitorText:      os.Getenv("MONITOR_TEXT"),
		HttpUserAgent:    os.Getenv("HTTP_USER_AGENT"),
		SlowWarningLimit: c.asInt("SLOW_WARNING_LIMIT"),
	}

	c.config.ScanFrequency = c.asInt("SCAN_FREQUENCY")
	c.config.Accounts = []configAccount{account}
}

func (c *envConfig) getConfigAccounts() *[]configAccount {
	return &c.config.Accounts
}

func (c *envConfig) getScanFrequency() int {
	frequency := c.config.ScanFrequency
	if frequency == 0 {
		return defaultScanFrequency
	}

	return frequency
}

func (c *envConfig) asInt(env string) int {
	value := os.Getenv(env)
	if value == "" {
		return 0
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return 0
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
