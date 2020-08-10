package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Accounts Accounts `yaml:"accounts"`
	Redis    Redis    `yaml:"redis"`
}

type Redis struct {
	Address string `yaml:"address"`
	DB      int    `yaml:"db"`
}

type Account struct {
	ID                         string   `yaml:"id"`
	Token                      Token    `yaml:"token"`
	HomeTimelineFetchInterval  int      `yaml:"home_timeline_fetch_interval"`
	DirectmessageFetchInterval int      `yaml:"directmessage_fetch_interval"`
	Webhooks                   []string `yaml:"webhooks"`
}

type Accounts []Account

func (a Accounts) Find(userID string) *Account {
	for _, account := range a {
		if account.ID == userID {
			return &account
		}
	}

	return nil
}

type Token struct {
	ConsumerKey       string `yaml:"consumer_key"`
	ConsumerSecret    string `yaml:"consumer_secret"`
	AccessToken       string `yaml:"access_token"`
	AccessTokenSecret string `yaml:"access_token_secret"`
}

func NewConfig(configPath string) (*Config, error) {
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	config := Config{}

	if err = yaml.Unmarshal(buf, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if len(config.Accounts) == 0 {
		return nil, fmt.Errorf("accounts are required")
	}

	m := make(map[string]Account)

	for _, account := range config.Accounts {
		if _, ok := m[account.ID]; ok {
			return nil, fmt.Errorf("deplicate account IDs")
		}

		m[account.ID] = account
	}

	return &config, nil
}
