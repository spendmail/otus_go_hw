package rabbitmq

import (
	"encoding/json"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/storage"
	"time"
)

type Config interface {
	GetRabbitDSN() string
	GetExchangeName() string
	GetExchangeKind() string
	GetExchangeDurable() bool
	GetExchangeAutoDelete() bool
	GetExchangeInternal() bool
	GetExchangeNoWait() bool
	GetQueueName() string
	GetQueueDurable() bool
	GetQueueAutoDelete() bool
	GetQueueInternal() bool
	GetQueueNoWait() bool
	GetQueueBindNoWait() bool
	GetQueueBindingKey() string
	GetConsumeConsumer() string
	GetConsumeAutoAck() bool
	GetConsumeExclusive() bool
	GetConsumeNoLocal() bool
	GetConsumeNoWait() bool
	GetPublishMandatory() bool
	GetPublishImmediate() bool
	GetPublishRoutingKey() string
}

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type Client struct {
	Config     Config
	Logger     Logger
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

type Notification struct {
	Id      int64     `json:"id"`
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
	OwnerID int64     `json:"owner_id"`
}

var (
	ErrRabbitmqDial            = errors.New("unable to dial rabbitmq server")
	ErrRabbitmqChanOpen        = errors.New("unable to open rabbitmq channel")
	ErrRabbitmqExchangeDeclare = errors.New("unable to declare rabbitmq exchange")
	ErrRabbitmqQueueDeclare    = errors.New("unable to declare rabbitmq queue")
	ErrRabbitmqQueueBind       = errors.New("unable to bind rabbitmq queue")
	ErrRabbitmqConsume         = errors.New("unable to consume rabbitmq queue")
	ErrRabbitmqPublish         = errors.New("unable to publish a message to rabbitmq")
	ErrRabbitmqConnectionClose = errors.New("unable to to close rabbitmq connection")
	ErrRabbitmqChannelClose    = errors.New("unable to to close rabbitmq channel")
)

func NewClient(config Config, logger Logger) (*Client, error) {

	conn, err := amqp.Dial(config.GetRabbitDSN())
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRabbitmqDial, err.Error())
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRabbitmqChanOpen, err.Error())
	}

	client := &Client{
		config,
		logger,
		conn,
		ch,
	}

	return client, nil
}

func (c *Client) DeclareExchange() error {

	err := c.Channel.ExchangeDeclare(
		c.Config.GetExchangeName(),
		c.Config.GetExchangeKind(),
		c.Config.GetExchangeDurable(),
		c.Config.GetExchangeAutoDelete(),
		c.Config.GetExchangeInternal(),
		c.Config.GetExchangeNoWait(),
		nil,
	)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrRabbitmqExchangeDeclare, err.Error())
	}

	return nil
}

func (c *Client) DeclareQueue() (amqp.Queue, error) {

	queue, err := c.Channel.QueueDeclare(
		c.Config.GetQueueName(),
		c.Config.GetQueueDurable(),
		c.Config.GetQueueAutoDelete(),
		c.Config.GetQueueInternal(),
		c.Config.GetQueueNoWait(),
		nil,
	)
	if err != nil {
		return queue, fmt.Errorf("%w: %s", ErrRabbitmqQueueDeclare, err.Error())
	}

	return queue, nil
}

func (c *Client) BindQueue(queue amqp.Queue) error {

	err := c.Channel.QueueBind(
		queue.Name,
		c.Config.GetQueueBindingKey(),
		c.Config.GetExchangeName(),
		c.Config.GetQueueBindNoWait(),
		nil)

	if err != nil {
		return fmt.Errorf("%w: %s", ErrRabbitmqQueueBind, err.Error())
	}

	return nil
}

func (c *Client) Consume(queue amqp.Queue) (<-chan amqp.Delivery, error) {

	messages, err := c.Channel.Consume(
		queue.Name,
		c.Config.GetConsumeConsumer(),
		c.Config.GetConsumeAutoAck(),
		c.Config.GetConsumeExclusive(),
		c.Config.GetConsumeNoLocal(),
		c.Config.GetConsumeNoWait(),
		nil,
	)

	if err != nil {
		return messages, fmt.Errorf("%w: %s", ErrRabbitmqConsume, err.Error())
	}

	return messages, nil
}

func (c *Client) SendEventNotification(event storage.Event) error {

	notification := &Notification{
		Id:      event.ID,
		Title:   event.Title,
		Date:    event.BeginDate,
		OwnerID: event.OwnerID,
	}

	jsonBody, _ := json.Marshal(&notification)

	err := c.Channel.Publish(
		c.Config.GetExchangeName(),
		c.Config.GetPublishRoutingKey(),
		c.Config.GetPublishMandatory(),
		c.Config.GetPublishImmediate(),
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonBody,
		})

	if err != nil {
		return fmt.Errorf("%w: %s", ErrRabbitmqPublish, err.Error())
	}

	return nil
}

func (c *Client) Close() {
	err := c.Channel.Close()
	if err != nil {
		c.Logger.Error(fmt.Errorf("%w: %s", ErrRabbitmqChannelClose, err.Error()).Error())
	}

	err = c.Connection.Close()
	if err != nil {
		c.Logger.Error(fmt.Errorf("%w: %s", ErrRabbitmqConnectionClose, err.Error()).Error())
	}
}
