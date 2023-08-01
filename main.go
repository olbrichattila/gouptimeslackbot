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
	scanner        scannerInterface
}

func newApp() *app {
	client := newUpClient(&request{})
	publisher := NewSlackPublisher()

	return &app{
		client:         client,
		slackPublisher: publisher,
		config:         resolveConfig(),
		scanner:        newScanner(client, publisher),
	}
}

func main() {
	app := newApp()

	frequency := app.config.getScanFrequency()
	accounts := app.config.getConfigAccounts()
	for {
		for _, config := range *accounts {
			app.scanner.Scan(config)
		}
		time.Sleep(time.Duration(frequency) * time.Second)
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
