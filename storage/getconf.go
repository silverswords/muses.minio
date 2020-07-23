package storage

import (
	"log"

	"github.com/spf13/viper"
)

type config struct {
	minioClients []*minioClient
}

func getConfig() *config {
	var config config
	viper.SetConfigName("config")
	viper.AddConfigPath("./config.yaml")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}

	err = viper.Unmarshal(&config) // 将配置信息绑定到结构体上
	if err != nil {
		log.Fatalln(err)
	}

	return &config
}
