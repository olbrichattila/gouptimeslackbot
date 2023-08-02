package main

type loggerSpy struct {
	called      int
	lastMessage string
}

func newLoggerSpy() *loggerSpy {
	return &loggerSpy{}
}

func (l *loggerSpy) Log(message string) {
	l.called++
	l.lastMessage = message
}
