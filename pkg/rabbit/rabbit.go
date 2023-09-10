package rabbit

import (
	"encoding/json"
	"rabbitmain/internal"
	"rabbitmain/pkg/entity"

	"github.com/streadway/amqp"
)

type rabbit struct {
	Ch        *amqp.Channel
	Conn      *amqp.Connection
	RepoMongo internal.RepoMongoInterface
}

func NewRabbit(uri string, repo internal.RepoMongoInterface) (*rabbit, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &rabbit{
		Ch:        ch,
		Conn:      conn,
		RepoMongo: repo,
	}, nil
}

func (e *rabbit) Consume(queueName string) (<-chan amqp.Delivery, error) {
	queue, err := e.Ch.Consume(
		entity.QueueRequestCollects, "", false, false, false, false, nil,
	)
	if err != nil {
		return nil, err
	}
	return queue, nil
}

func (e *rabbit) Publish(message interface{}) error {
	vv, err := json.Marshal(message)
	if err != nil {
		return nil
	}
	e.RepoMongo.Register(message)
	return e.Ch.Publish(entity.ExchangeRequestCollects, entity.QueueRequestCollects, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        vv,
	})
}

func (e *rabbit) RePublish(message []byte) error {
	e.RepoMongo.Register(message)
	return e.Ch.Publish(entity.ExchangeRequestCollects, entity.QueueRequestCollects, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        message,
	})
}
