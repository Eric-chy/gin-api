package mq

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

//MessageBody AMQP消息的body。类型将在Request header上设置
type MessageBody struct {
	Data []byte
	Type string
}

//Message 消息
type Message struct {
	Queue         string
	ReplyTo       string
	ContentType   string
	CorrelationID string
	Priority      uint8
	Body          MessageBody
}

//Connection 创建连接
type Connection struct {
	name     string
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	queues   []string
	err      chan error
}

var (
	connectionPool = make(map[string]*Connection)
)

//NewConnection 从连接池中取出或者新建后放进消息池
func NewConnection(name, exchange string, queues []string) *Connection {
	if c, ok := connectionPool[name]; ok {
		return c
	}
	c := &Connection{
		exchange: exchange,
		queues:   queues,
		err:      make(chan error),
	}
	connectionPool[name] = c
	return c
}

//GetConnection 从连接池中取出一个连接
func GetConnection(name string) *Connection {
	return connectionPool[name]
}

func (c *Connection) Connect() error {
	var err error
	amqpURI := "amqp://guest:guest@localhost:5672/"
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("Error in creating rabbitmq connection with %s : %s", amqpURI, err.Error())
	}
	go func() {
		<-c.conn.NotifyClose(make(chan *amqp.Error)) //Listen to NotifyClose
		c.err <- errors.New("Connection Closed")
	}()
	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}
	if err := c.channel.ExchangeDeclare(
		c.exchange, // name
		"direct",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return fmt.Errorf("Error in Exchange Declare: %s", err)
	}
	return nil
}

func (c *Connection) BindQueue() error {
	for _, q := range c.queues {
		if _, err := c.channel.QueueDeclare(q, true, false, false, false, nil); err != nil {
			return fmt.Errorf("error in declaring the queue %s", err)
		}
		if err := c.channel.QueueBind(q, "my_routing_key", c.exchange, false, nil); err != nil {
			return fmt.Errorf("Queue  Bind error: %s", err)
		}
	}
	return nil
}

//Reconnect 重连
func (c *Connection) Reconnect() error {
	if err := c.Connect(); err != nil {
		return err
	}
	if err := c.BindQueue(); err != nil {
		return err
	}
	return nil
}

//Publish 发布到队列
func (c *Connection) Publish(m Message) error {
	select { //non blocking channel - if there is no error will go to default where we do nothing
	case err := <-c.err:
		if err != nil {
			c.Reconnect()
		}
	default:
	}
	p := amqp.Publishing{
		Headers:       amqp.Table{"type": m.Body.Type},
		ContentType:   m.ContentType,
		CorrelationId: m.CorrelationID,
		Body:          m.Body.Data,
		ReplyTo:       m.ReplyTo,
	}
	if err := c.channel.Publish(c.exchange, m.Queue, false, false, p); err != nil {
		return fmt.Errorf("Error in Publishing: %s", err)
	}
	return nil
}

//Consume consumes the messages from the queues and passes it as map of chan of amqp.Delivery
func (c *Connection) Consume() (map[string]<-chan amqp.Delivery, error) {
	m := make(map[string]<-chan amqp.Delivery)
	for _, q := range c.queues {
		deliveries, err := c.channel.Consume(q, "", false, false, false, false, nil)
		if err != nil {
			return nil, err
		}
		m[q] = deliveries
	}
	return m, nil
}

//HandleConsumedDeliveries 处理队列中消息。单个消费者
func (c *Connection) HandleConsumedDeliveries(q string, delivery <-chan amqp.Delivery, fn func(Connection, string, <-chan amqp.Delivery)) {
	for {
		go fn(*c, q, delivery)
		if err := <-c.err; err != nil {
			c.Reconnect()
			deliveries, err := c.Consume()
			if err != nil {
				panic(err) //raising panic if consume fails even after reconnecting
			}
			delivery = deliveries[q]
		}
	}
}

//入队例子
func send() {
	conn := NewConnection("my-producer", "my-exchange", []string{"queue-1", "queue-2"})
	if err := conn.Connect(); err != nil {
		panic(err)
	}
	if err := conn.BindQueue(); err != nil {
		panic(err)
	}
	for _, q := range conn.queues {
		m := Message{
			Queue: q,
			//set the necessary fields
		}
		if err := conn.Publish(m); err != nil {
			panic(err)
		}
	}
}

//出队例子
func receive() {
	forever := make(chan bool)
	conn := NewConnection("my-consumer-1", "my-exchange", []string{"queue-1", "queue-2"})
	if err := conn.Connect(); err != nil {
		panic(err)
	}
	if err := conn.BindQueue(); err != nil {
		panic(err)
	}
	deliveries, err := conn.Consume()
	if err != nil {
		panic(err)
	}
	for q, d := range deliveries {
		go conn.HandleConsumedDeliveries(q, d, messageHandler)
	}
	<-forever
}

func messageHandler(c Connection, q string, deliveries <-chan amqp.Delivery) {
	for d := range deliveries {
		m := Message{
			Queue:         q,
			Body:          MessageBody{Data: d.Body, Type: d.Headers["type"].(string)},
			ContentType:   d.ContentType,
			Priority:      d.Priority,
			CorrelationID: d.CorrelationId,
		}
		//handle the custom message
		log.Println("Got message from queue ", m.Queue)
		d.Ack(false)
	}
}
