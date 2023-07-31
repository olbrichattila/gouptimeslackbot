package main

import (
	"fmt"
	"time"
)

const defaultScanFrequency = 60

type app struct {
	client        upClientInterface
	slackMessager SlackPublisherInterface
	config        configInterface
}

func newApp() *app {
	return &app{
		client:        newUpClient(&request{}),
		slackMessager: NewSlackPublisher(),
		config:        resolveConfig(),
	}
}

func main() {
	app := newApp()

	frequency := app.config.getScanFrequency()
	for {
		accounts := app.config.getConfigAccounts()
		for _, config := range *accounts {
			go doScan(app, &config)
		}
		time.Sleep(time.Duration(frequency) * time.Second)
	}
}

func doScan(app *app, config *configAccount) {
	elapsed, err := app.client.TestUrl(config.MonitorUrl, config.MonitorText)
	if err != nil {
		message := fmt.Sprintf("Up boat report:\n\tElapsed: %d\n\tError:%v", elapsed, err)
		err = app.slackMessager.Send(config.SlackBotToken, config.SlackChannelId, message)
		if err != nil {
			fmt.Println(err)
		}
	}

	if config.SlowWarningLimit > 0 && elapsed > config.SlowWarningLimit {
		message := fmt.Sprintf("Slow warning limit reached:\n\tLimit: %d\n\tElapsed: %d", config.SlowWarningLimit, elapsed)
		err = app.slackMessager.Send(config.SlackBotToken, config.SlackChannelId, message)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func resolveConfig() configInterface {
	if fileExists(yamlFileName) {
		fmt.Println("Using yaml config file...", yamlFileName)
		return newYamlConfig()
	}

	fmt.Println("Using .env as config...")
	return newConfig(true)
}
