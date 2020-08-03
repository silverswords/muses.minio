package storage

import (
	"log"

	"github.com/spf13/viper"
)

type config struct {
	Clients []map[string]string
}

type configInfo struct {
	configName string
	configPath string
}

func (b *Bucket) getConfig() *config {
	var config config
	viper.SetConfigName(b.configName)
	viper.AddConfigPath(b.configPath)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalln(err)
	}

	return &config
}
