package main

import (
	"fmt"
	"time"
)

type mockRequestSpy struct {
	err      error
	response string
	delay    int
}

func newMockRequestSpy() *mockRequestSpy {
	return &mockRequestSpy{}
}

func (r *mockRequestSpy) get(_, _ string) (string, error) {
	if r.delay > 0 {
		time.Sleep(time.Duration(r.delay) * time.Millisecond)
	}

	return r.response, r.err
}

func (r *mockRequestSpy) withError() *mockRequestSpy {
	r.err = fmt.Errorf("custom test error")
	return r
}

func (r *mockRequestSpy) withResponse(t string) *mockRequestSpy {
	r.response = t
	return r
}

func (r *mockRequestSpy) withDelay(delay int) *mockRequestSpy {
	r.delay = delay
	return r
}
