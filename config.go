package main

type configInterface interface {
	load()
	getConfigAccounts() *[]configAccount
	getScanFrequency() int
}

type configAccount struct {
	SlackBotToken    string `yaml:"SlackBotToken"`
	SlackChannelID   string `yaml:"SlackChannelID"`
	MonitorURL       string `yaml:"MonitorURL"`
	MonitorText      string `yaml:"MonitorText"`
	SlowWarningLimit int    `yaml:"SlowWarningLimit"`
	HTTPUserAgent    string `yaml:"HTTPUserAgent"`
}

type config struct {
	ScanFrequency int             `yaml:"ScanFrequency"`
	Accounts      []configAccount `yaml:"Accounts"`
}
