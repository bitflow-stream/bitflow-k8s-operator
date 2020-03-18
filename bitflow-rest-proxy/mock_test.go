package main

import (
	"testing"

	"github.com/antongulenko/golib"
)

type RestProxyTestSuite struct {
	golib.AbstractTestSuite
}

func TestMockTestSuite(t *testing.T) {
	new(RestProxyTestSuite).Run(t)
}

func (suite *RestProxyTestSuite) testSomething() {
	// TODO write tests
}
