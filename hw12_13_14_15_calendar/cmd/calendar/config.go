package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Logger LoggerConf
}

type LoggerConf struct {
	Level   string
	File    string
	Size    int
	Backups int
	Age     int
}

func NewConfig() Config {
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Unable to read config file %s \n", configFile))
	}

	return Config{
		LoggerConf{
			viper.GetString("logger.level"),
			viper.GetString("logger.file"),
			viper.GetInt("logger.size"),
			viper.GetInt("logger.backups"),
			viper.GetInt("logger.age"),
		},
	}
}
