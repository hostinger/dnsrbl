package hbl

import (
	"errors"
	"io/ioutil"
	"net"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ExpireTTL string   `yaml:"expire_ttl"`
	AllowList []string `yaml:"allow_list"`
}

func NewConfig(expireTTL string, allowList []string) *Config {
	return &Config{
		ExpireTTL: expireTTL,
		AllowList: allowList,
	}
}

func NewConfigFromFile(cfg string) (*Config, error) {
	file, err := ioutil.ReadFile(cfg)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Config) Validate() error {
	for _, address := range c.AllowList {
		if _, _, err := net.ParseCIDR(address); err != nil {
			return errors.New("Failed to validate configuration file, 'allow_list' argument must be a valid IP address in CIDR notation.")
		}
	}
	if _, err := time.ParseDuration(c.ExpireTTL); err != nil {
		return errors.New("Failed to validate configuration file, 'expire_ttl' must be a valid duration.")
	}
	return nil
}
