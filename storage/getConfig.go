package storage

import (
	"github.com/spf13/viper"
)

type config struct {
	Clients []map[string]string
}

type configInfo struct {
	configName string
	configPath string
}

func (b *Bucket) getConfig() (*config, error) {
	var config config
	viper.SetConfigName(b.configName)
	viper.AddConfigPath(b.configPath)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
