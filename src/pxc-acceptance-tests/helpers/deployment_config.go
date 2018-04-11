package helpers

import (
	"fmt"
	"os"

	"encoding/json"
	"io/ioutil"
)

type Component struct {
	Ip        string `json:"ip"`
	SshTunnel string `json:"ssh_tunnel"`
}

type Proxy struct {
	Url         string `json:"url"`
	APIUsername string `json:"api_username"`
	APIPassword string `json:"api_password"`
}

type DbCreds struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
}

type MySQLIntegrationConfig struct {
	DeploymentName     string  `json:deployment_name,omitempty`
	MySQLInstanceGroup string  `json:mysql_instance_group,omitempty`
	BOSH               BOSH    `json:"bosh"`
	Proxy              Proxy   `json:"proxy"`
	DbCreds            DbCreds `json:"db_creds,omitempty"`
}

type BOSH struct {
	CACert       string `json:"ca_cert"`
	Client       string `json:"client"`
	ClientSecret string `json:"client_secret"`
	URL          string `json:"url"`
}

func LoadConfig() (MySQLIntegrationConfig, error) {
	mySQLIntegrationConfig := MySQLIntegrationConfig{}

	path := os.Getenv("CONFIG")
	if path == "" {
		return mySQLIntegrationConfig, fmt.Errorf("Must set $CONFIG to point to an integration config .json file.")
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(buf, &mySQLIntegrationConfig); err != nil {
		panic(err)
	}

	if mySQLIntegrationConfig.DeploymentName == "" {
		mySQLIntegrationConfig.DeploymentName = "pxc"
	}

	if mySQLIntegrationConfig.MySQLInstanceGroup == "" {
		mySQLIntegrationConfig.MySQLInstanceGroup = "mysql"
	}

	return mySQLIntegrationConfig, nil
}

func ValidateConfig(config *MySQLIntegrationConfig) error {

	if config.DbCreds.Host == "" {
		return fmt.Errorf("Field 'db_creds.host' must not be empty")
	}

	if config.DbCreds.Port == 0 {
		return fmt.Errorf("Field 'db_creds.port' must not be empty")
	}

	if config.DbCreds.Username == "" {
		return fmt.Errorf("Field 'db_creds.username' must not be empty")
	}

	if config.DbCreds.Password == "" {
		return fmt.Errorf("Field 'db_creds.password' must not be empty")
	}

	if config.BOSH.CACert == "" {
		return fmt.Errorf("Field 'bosh.ca_cert' must not be empty")
	}

	if config.BOSH.Client == "" {
		return fmt.Errorf("Field 'bosh.client' must not be empty")
	}

	if config.BOSH.ClientSecret == "" {
		return fmt.Errorf("Field 'bosh.client_secret' must not be empty")
	}

	if config.BOSH.URL == "" {
		return fmt.Errorf("Field 'bosh.url' must not be empty")
	}

	if config.Proxy.Url == "" {
		return fmt.Errorf("Field 'proxy.url' must not be empty")
	}

	if config.Proxy.APIUsername == "" {
		return fmt.Errorf("Field 'proxy.api_username' must not be empty")
	}

	if config.Proxy.APIPassword == "" {
		return fmt.Errorf("Field 'proxy.api_password' must not be empty")
	}

	return nil
}
