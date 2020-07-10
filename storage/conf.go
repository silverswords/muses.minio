package storage

import (
	"log"

	"github.com/spf13/viper"
)

func (b *minioClient) getConf() *minioClient {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config.yaml")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
	viper.Unmarshal(&b) // 将配置信息绑定到结构体上
	return b
}
