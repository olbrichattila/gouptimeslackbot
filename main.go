// Package main provides the entry point for the application.
// This package is a slack uptime boot, see README.md
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

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
	// publisher := newSlackPublisherSpy() // Swap this for testing
	logger := newLogger()

	return &app{
		client:         client,
		slackPublisher: publisher,
		config:         resolveConfig(),
		scanner:        newScanner(client, publisher, logger),
	}
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	app := newApp()

	accounts := app.config.getConfigAccounts()

	for _, config := range *accounts {
		app.scanner.Scan(config)
	}

	<-sigs
}

func resolveConfig() configInterface {
	if fileExists(yamlFileName) {
		fmt.Println("Using yaml config file...", yamlFileName)
		return newYamlConfig()
	}

	fmt.Println("Using .env as config...")
	return newConfig(true)
}
