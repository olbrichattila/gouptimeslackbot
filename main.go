// Package main provides the entry point for the application.
// This package is a slack uptime boot, see README.md
package main

import (
	"fmt"
	"time"
)

const defaultScanFrequency = 60
const defaultHTTPUserAgent = "GolangUptimeBot/1.0"

type app struct {
	client         upClientInterface
	slackPublisher slackPublisherInterface
	config         configInterface
	scanner        scannerInterface
	// logger         loggerInterface //@TODO Add logger
}

func newApp() *app {
	client := newUpClient(&request{})
	publisher := newSlackPublisher()
	logger := newLogger()

	return &app{
		client:         client,
		slackPublisher: publisher,
		config:         resolveConfig(),
		scanner:        newScanner(client, publisher, logger),
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
