package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/tsconn23/linotwit/internal/bootstrap/flags"
)

type CredentialInfo struct {
	AccessToken string
	AccessSecret string
	ConsumerKey string
	ConsumerSecret string
}

type SubscriptionInfo struct {
	Handles []string
	Hashtags []string
}

type ConfigInfo struct {
	Credentials CredentialInfo
	Subscriptions SubscriptionInfo
}

type Loader struct {
	directory string
	fileName string
	flags    flags.Common
}

func NewLoader(flags flags.Common) Loader {
	loader := Loader{
		directory:  flags.ConfigDirectory(),
		fileName:   flags.ConfigFileName(),
		flags:      flags,
	}
	return loader
}

func (l Loader) Process(cfg *ConfigInfo) error {
	filePath := l.directory + "/" + l.fileName
	if _, err := toml.DecodeFile(filePath, &cfg); err != nil {
		return fmt.Errorf("could not load configuration file (%s): %s", filePath, err.Error())
	}
	return nil
}