package main

import (
	"testing"

	"github.com/antongulenko/golib"
	"github.com/stretchr/testify/suite"
)

type RestProxyTestSuite struct {
	golib.AbstractTestSuite
}

func TestMockTestSuite(t *testing.T) {
	suite.Run(t, new(RestProxyTestSuite))
}

func (suite *RestProxyTestSuite) TestSomething() {
	// TODO write tests
}
