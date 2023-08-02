package main

import "fmt"

type upClientSpy struct {
	called      int
	elapsedTime int
	err         error
}

func newUpClientSpy() *upClientSpy {
	return &upClientSpy{}
}

func (u *upClientSpy) TestUrl(url, contains string) (int, error) {
	u.called++
	return u.elapsedTime, u.err
}

func (u *upClientSpy) withError(errorMessage string) *upClientSpy {
	u.err = fmt.Errorf(errorMessage)
	return u
}

func (u *upClientSpy) withElapsedTime(miliseconds int) *upClientSpy {
	u.elapsedTime = miliseconds
	return u
}
