package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

const RABBITMQ_SERVER = "amqp://root:123456@47.115.211.122:5672/"

type RabbitMQ struct {
	channel  *amqp.Channel
	conn     *amqp.Connection
	Name     string
	exchange string
}

// 创建一个新的rabbitmq.RabbitMq结构体
func New(s string) *RabbitMQ {
	connection, err := amqp.Dial(s)
	if err != nil {
		panic(err)
	}
	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	queue, err := channel.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	mq := new(RabbitMQ)
	mq.channel = channel
	mq.conn = connection
	mq.Name = queue.Name
	return mq
}

// 将自己的消息队列和一个exchange绑定
func (q *RabbitMQ) Bind(exchange string) {
	err := q.channel.QueueBind(
		q.Name,
		"",
		exchange,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	q.exchange = exchange
}

// 向消息队列发送消息
func (q *RabbitMQ) Send(queue string, body interface{}) {
	str, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	err = q.channel.Publish("", queue, false, false, amqp.Publishing{
		ReplyTo: q.Name,
		Body:    []byte(str),
	})
	if err != nil {
		panic(err)
	}
}

// 向exchange发送消息
func (q *RabbitMQ) Publish(exchange string, body interface{}) {
	str, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	err = q.channel.Publish(exchange, "", false, false, amqp.Publishing{
		ReplyTo: q.Name,
		Body:    []byte(str),
	})
	if err != nil {
		log.Fatal(err)
	}
}

// 生成一个接收消息的go channel
func (q *RabbitMQ) Consume() <-chan amqp.Delivery {
	consume, err := q.channel.Consume(q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	return consume
}

func (q *RabbitMQ) Close() {
	q.channel.Close()
	q.conn.Close()
}
