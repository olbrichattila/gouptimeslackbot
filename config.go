package main

type configInterface interface {
	load()
	getConfigAccounts() *[]configAccount
	getScanFrequency() int
}

type configAccount struct {
	SlackBotToken    string `yaml:"SlackBotToken"`
	SlackChannelId   string `yaml:"SlackChannelId"`
	MonitorUrl       string `yaml:"MonitorUrl"`
	MonitorText      string `yaml:"MonitorText"`
	SlowWarningLimit int    `yaml:"SlowWarningLimit"`
	HttpUserAgent    string `yaml:"HttpUserAgent"`
}

type config struct {
	ScanFrequency int             `yaml:"ScanFrequency"`
	Accounts      []configAccount `yaml:"Accounts"`
}
