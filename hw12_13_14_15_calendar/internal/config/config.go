package config

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var ErrConfigRead = errors.New("unable to read config file")

type Config struct {
	Logger  LoggerConf
	HTTP    HTTPConf
	Storage StorageConf
}

type LoggerConf struct {
	Level   string
	File    string
	Size    int
	Backups int
	Age     int
}

type HTTPConf struct {
	Host string
	Port string
}

type StorageConf struct {
	Implementation string
	DSN            string
}

func NewConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConfigRead, path)
	}

	return &Config{
		LoggerConf{
			viper.GetString("logger.level"),
			viper.GetString("logger.file"),
			viper.GetInt("logger.size"),
			viper.GetInt("logger.backups"),
			viper.GetInt("logger.age"),
		},
		HTTPConf{
			viper.GetString("http.host"),
			viper.GetString("http.port"),
		},
		StorageConf{
			viper.GetString("storage.implementation"),
			viper.GetString("storage.dsn"),
		},
	}, nil
}

func (c *Config) GetLoggerLevel() string {
	return c.Logger.Level
}

func (c *Config) GetLoggerFile() string {
	return c.Logger.File
}

func (c *Config) GetLoggerSize() int {
	return c.Logger.Size
}

func (c *Config) GetLoggerBackups() int {
	return c.Logger.Backups
}

func (c *Config) GetLoggerAge() int {
	return c.Logger.Age
}

func (c *Config) GetServerHost() string {
	return c.HTTP.Host
}

func (c *Config) GetServerPort() string {
	return c.HTTP.Port
}

func (c *Config) GetStorageImplementation() string {
	return c.Storage.Implementation
}

func (c *Config) GetStorageDSN() string {
	return c.Storage.DSN
}
