package handler

import (
	"testing"

	"github.com/I1820/gm/router"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

// Suite is a test suite for APIs.
type Suite struct {
	suite.Suite
	engine *echo.Echo
}

// SetupSuite initiates tm test suite
func (suite *Suite) SetupSuite() {
	suite.engine = router.App()

	lh := LoRa{}
	lh.Register(suite.engine.Group(""))
}

// Let's test APIs!
func TestService(t *testing.T) {
	suite.Run(t, new(Suite))
}
