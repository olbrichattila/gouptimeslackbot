package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type scanTestSuite struct {
	suite.Suite
}

func TestScanRunner(t *testing.T) {
	suite.Run(t, new(scanTestSuite))
}

func (t *scanTestSuite) SetupTest() {
	fmt.Println(Yellow+"Running test scan env: "+Green, t.T().Name()+Reset)
}

func (t *scanTestSuite) TestItWorks() {
	// @TODO add tests
	t.True(true)
}
