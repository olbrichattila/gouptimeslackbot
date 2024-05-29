package main

type configInterface interface {
	load()
	getConfigAccounts() *[]configAccount
}

type configAccount struct {
	ScanFrequency           int    `yaml:"ScanFrequency"`
	SlackBotToken           string `yaml:"SlackBotToken"`
	SlackChannelID          string `yaml:"SlackChannelID"`
	MonitorURL              string `yaml:"MonitorURL"`
	MonitorText             string `yaml:"MonitorText"`
	SlowWarningLimit        int    `yaml:"SlowWarningLimit"`
	HTTPUserAgent           string `yaml:"HTTPUserAgent"`
	RepeatNotificationDelay int    `yaml:"RepeatNotificationDelay"`
}

type config struct {
	Accounts []configAccount `yaml:"Accounts"`
}
