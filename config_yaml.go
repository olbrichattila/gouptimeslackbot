package main

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

const yamlFileName = "./config.yaml"

type yamlConfig struct {
	config *config
}

func newYamlConfig() *yamlConfig {
	var c yamlConfig
	c.load()
	return &c
}

func (c *yamlConfig) load() {
	c.config = &config{}
	c.loadConfig()
}

func (c *yamlConfig) getConfigAccounts() *[]configAccount {
	return &c.config.Accounts
}

func (c *yamlConfig) getScanFrequency() int {
	frequency := c.config.ScanFrequency
	if frequency == 0 {
		return defaultScanFrequency
	}

	return frequency
}

func (c *yamlConfig) createExampleConfig() {
	conf := newConfig(true)

	a := []configAccount{conf.config.Accounts[0], conf.config.Accounts[0]}
	y := &config{Accounts: a}
	file, err := os.Create(yamlFileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	if err := encoder.Encode(y); err != nil {
		fmt.Println("Error encoding to YAML:", err)
		return
	}
}

func (c *yamlConfig) loadConfig() {
	file, err := ioutil.ReadFile(yamlFileName)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	if err := yaml.Unmarshal(file, c.config); err != nil {
		fmt.Println("Error unmarshaling YAML:", err)
		return
	}
}
