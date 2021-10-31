package config

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var ErrConfigRead = errors.New("unable to read config file")

type Config struct {
	Logger    LoggerConf
	HTTP      HTTPConf
	Storage   StorageConf
	GRPC      GRPCConf
	Rabbitmq  RabbitmqConf
	Exchange  ExchangeConf
	Queue     QueueConf
	Consume   ConsumeConf
	Publish   PublishConf
	Scheduler SchedulerConf
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

type GRPCConf struct {
	Host string
	Port string
}

type StorageConf struct {
	Implementation string
	DSN            string
}

type RabbitmqConf struct {
	DSN string
}

type ExchangeConf struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
}

type QueueConf struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	BindNoWait bool
	BindingKey string
}

type ConsumeConf struct {
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
}

type PublishConf struct {
	Mandatory  bool
	Immediate  bool
	RoutingKey string
}

type SchedulerConf struct {
	RemindIn int
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
		GRPCConf{
			viper.GetString("grpc.host"),
			viper.GetString("grpc.port"),
		},
		RabbitmqConf{
			viper.GetString("rabbitmq.dsn"),
		},
		ExchangeConf{
			viper.GetString("exchange.name"),
			viper.GetString("exchange.kind"),
			viper.GetBool("exchange.durable"),
			viper.GetBool("exchange.autoDelete"),
			viper.GetBool("exchange.internal"),
			viper.GetBool("exchange.noWait"),
		},
		QueueConf{
			viper.GetString("queue.name"),
			viper.GetBool("queue.durable"),
			viper.GetBool("queue.autoDelete"),
			viper.GetBool("queue.exclusive"),
			viper.GetBool("queue.noWait"),
			viper.GetBool("queue.bindNoWait"),
			viper.GetString("queue.bindingKey"),
		},
		ConsumeConf{
			viper.GetString("consume.consumer"),
			viper.GetBool("consume.autoAck"),
			viper.GetBool("consume.exclusive"),
			viper.GetBool("consume.noLocal"),
			viper.GetBool("consume.noWait"),
		},
		PublishConf{
			viper.GetBool("publish.mandatory"),
			viper.GetBool("publish.immediate"),
			viper.GetString("publish.routingKey"),
		},
		SchedulerConf{
			viper.GetInt("scheduler.remindIn"),
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

func (c *Config) GetHTTPHost() string {
	return c.HTTP.Host
}

func (c *Config) GetHTTPPort() string {
	return c.HTTP.Port
}

func (c *Config) GetGrpcHost() string {
	return c.GRPC.Host
}

func (c *Config) GetGrpcPort() string {
	return c.GRPC.Port
}

func (c *Config) GetStorageImplementation() string {
	return c.Storage.Implementation
}

func (c *Config) GetStorageDSN() string {
	return c.Storage.DSN
}

func (c *Config) GetRabbitDSN() string {
	return c.Rabbitmq.DSN
}

func (c *Config) GetExchangeName() string {
	return c.Exchange.Name
}

func (c *Config) GetExchangeKind() string {
	return c.Exchange.Kind
}

func (c *Config) GetExchangeDurable() bool {
	return c.Exchange.Durable
}

func (c *Config) GetExchangeAutoDelete() bool {
	return c.Exchange.AutoDelete
}

func (c *Config) GetExchangeInternal() bool {
	return c.Exchange.Internal
}

func (c *Config) GetExchangeNoWait() bool {
	return c.Exchange.NoWait
}

func (c *Config) GetQueueName() string {
	return c.Queue.Name
}

func (c *Config) GetQueueDurable() bool {
	return c.Queue.Durable
}

func (c *Config) GetQueueAutoDelete() bool {
	return c.Queue.AutoDelete
}

func (c *Config) GetQueueInternal() bool {
	return c.Queue.Exclusive
}

func (c *Config) GetQueueNoWait() bool {
	return c.Queue.NoWait
}

func (c *Config) GetQueueBindNoWait() bool {
	return c.Queue.BindNoWait
}

func (c *Config) GetQueueBindingKey() string {
	return c.Queue.BindingKey
}

func (c *Config) GetConsumeConsumer() string {
	return c.Consume.Consumer
}

func (c *Config) GetConsumeAutoAck() bool {
	return c.Consume.AutoAck
}

func (c *Config) GetConsumeExclusive() bool {
	return c.Consume.Exclusive
}

func (c *Config) GetConsumeNoLocal() bool {
	return c.Consume.NoLocal
}

func (c *Config) GetConsumeNoWait() bool {
	return c.Consume.NoWait
}

func (c *Config) GetPublishMandatory() bool {
	return c.Publish.Mandatory
}

func (c *Config) GetPublishImmediate() bool {
	return c.Publish.Immediate
}

func (c *Config) GetPublishRoutingKey() string {
	return c.Publish.RoutingKey
}

func (c *Config) GetSchedulerRemindIn() int {
	return c.Scheduler.RemindIn
}
