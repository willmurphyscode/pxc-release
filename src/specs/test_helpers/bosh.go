package test_helpers

import (
	boshdir "pxc-acceptance-tests/vendor/github.com/cloudfoundry/bosh-cli/director"
	boshuaa "pxc-acceptance-tests/vendor/github.com/cloudfoundry/bosh-cli/uaa"
	boshlog "pxc-acceptance-tests/vendor/github.com/cloudfoundry/bosh-utils/logger"

	"fmt"
	"os"
)

func BuildDirector() (boshdir.Director, error) {

	logger := boshlog.NewLogger(boshlog.LevelError)
	factory := boshdir.NewFactory(logger)

	// Build a Director config from address-like string.
	// HTTPS is required and certificates are always verified.
	config, err := boshdir.NewConfigFromURL(boshEnvironment())
	if err != nil {
		return nil, fmt.Errorf("building director config: %s", err)
	}

	// Configure custom trusted CA certificates.
	// If nothing is provided default system certificates are used.
	config.CACert = boshCaCert()

	// Allow Director to fetch UAA tokens when necessary.
	uaa, err := buildUAA()
	if err != nil {
		return nil, fmt.Errorf("building uaa: %s", err)
	}

	config.TokenFunc = boshuaa.NewClientTokenSession(uaa).TokenFunc

	return factory.New(config, boshdir.NewNoopTaskReporter(), boshdir.NewNoopFileReporter())
}

func boshEnvironment() string {
	return os.Getenv("BOSH_ENVIRONMENT")
}

func boshClient() string {
	return os.Getenv("BOSH_CLIENT")
}

func boshClientSecret() string {
	return os.Getenv("BOSH_CLIENT_SECRET")
}

func boshCaCert() string {
	return os.Getenv("BOSH_CA_CERT")
}

func buildUAA() (boshuaa.UAA, error) {
	logger := boshlog.NewLogger(boshlog.LevelError)
	factory := boshuaa.NewFactory(logger)

	// Build a UAA config from a URL.
	// HTTPS is required and certificates are always verified.

	config, err := boshuaa.NewConfigFromURL(fmt.Sprintf("https://%s:8443", boshEnvironment()))
	if err != nil {
		return nil, fmt.Errorf("ERROR build uaa config: %s", err)
	}

	// Set client credentials for authentication.
	// Machine level access should typically use a client instead of a particular user.
	config.Client = boshClient()
	config.ClientSecret = boshClientSecret()

	// Configure trusted CA certificates.
	// If nothing is provided default system certificates are used.
	config.CACert = boshCaCert()

	fmt.Println("about to build factory")
	return factory.New(config)
}
