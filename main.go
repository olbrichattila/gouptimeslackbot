package main

import (
	"fmt"
	"time"
)

const defaultScanFrequency = 60

type app struct {
	client         upClientInterface
	slackPublisher SlackPublisherInterface
	config         configInterface
}

func newApp() *app {
	return &app{
		client:         newUpClient(&request{}),
		slackPublisher: NewSlackPublisher(),
		config:         resolveConfig(),
	}
}

func main() {
	app := newApp()

	frequency := app.config.getScanFrequency()
	accounts := app.config.getConfigAccounts()
	for {
		for _, config := range *accounts {
			go doScan(app, config)
		}
		time.Sleep(time.Duration(frequency) * time.Second)
	}
}

/*
 * @todo move this out from here into a tetable format
 * add http response code check and report other then 20*
 * limit number of messages if there are too many and aggregate them
 */
func doScan(app *app, config configAccount) {
	elapsed, err := app.client.TestUrl(config.MonitorUrl, config.MonitorText)
	if err != nil {
		// @todo add event date / time to the message
		message := fmt.Sprintf("%s:\nUp boat report:\n\tElapsed: %d\n\tError:%v", config.MonitorUrl, elapsed, err)
		err = app.slackPublisher.Send(config.SlackBotToken, config.SlackChannelId, message)
		if err != nil {
			fmt.Println(err)
		}
	}

	if config.SlowWarningLimit > 0 && elapsed > config.SlowWarningLimit {
		// @todo add event date / time to the message
		message := fmt.Sprintf("%s:\nSlow warning limit reached:\n\tLimit: %d\n\tElapsed: %d", config.MonitorUrl, config.SlowWarningLimit, elapsed)
		err = app.slackPublisher.Send(config.SlackBotToken, config.SlackChannelId, message)
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
