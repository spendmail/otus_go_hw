package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Logger LoggerConf
	Http   HttpConf
}

type LoggerConf struct {
	Level   string
	File    string
	Size    int
	Backups int
	Age     int
}

type HttpConf struct {
	Host string
	Port string
}

func NewConfig(path string) Config {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Unable to read config file %s \n", path))
	}

	return Config{
		LoggerConf{
			viper.GetString("logger.level"),
			viper.GetString("logger.file"),
			viper.GetInt("logger.size"),
			viper.GetInt("logger.backups"),
			viper.GetInt("logger.age"),
		},
		HttpConf{
			viper.GetString("http.host"),
			viper.GetString("http.port"),
		},
	}
}
