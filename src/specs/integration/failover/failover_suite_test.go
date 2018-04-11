package failover

import (
	. "pxc-acceptance-tests/vendor/github.com/onsi/ginkgo"
	. "pxc-acceptance-tests/vendor/github.com/onsi/gomega"
	"pxc-acceptance-tests/helpers"
	"testing"
)

var mySQLIntegrationConfig helpers.MySQLIntegrationConfig

func TestFailover(t *testing.T) {

	var err error
	mySQLIntegrationConfig, err = helpers.LoadConfig()
	if err != nil {
		panic("Loading config: " + err.Error())
	}

	err = helpers.ValidateConfig(&mySQLIntegrationConfig)
	if err != nil {
		panic("Validating config: " + err.Error())
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, "PXC Acceptance Tests -- Failover")

}
