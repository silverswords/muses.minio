package storage

import (
	"log"

	"github.com/spf13/viper"
)

func (m *minioClient) getConf() *minioClient {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config.yaml")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
	viper.Unmarshal(&m) // 将配置信息绑定到结构体上
	return m
}
