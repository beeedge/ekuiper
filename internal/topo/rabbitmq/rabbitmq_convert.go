package rabbitmq

import (
	"fmt"

	"github.com/lf-edge/ekuiper/internal/conf"
	"github.com/lf-edge/ekuiper/pkg/api"
	"github.com/streadway/amqp"
)

type rabbitMQConvertSink struct {
	Username     string
	Password     string
	URL          string
	Exchange     string
	ExchangeType string
	RoutingKey   string
	conn         *amqp.Connection
	channel      *amqp.Channel
}

func (rcs *rabbitMQConvertSink) Configure(props map[string]interface{}) error {
	if i, ok := props["username"]; ok {
		if u, ok := i.(string); ok {
			rcs.Username = u
		} else {
			return fmt.Errorf("Not valid username %v.", i)
		}
	}

	if i, ok := props["password"]; ok {
		if p, ok := i.(string); ok {
			rcs.Password = p
		} else {
			return fmt.Errorf("Not valid password %v.", i)
		}
	}

	if i, ok := props["url"]; ok {
		if u, ok := i.(string); ok {
			rcs.URL = u
		} else {
			return fmt.Errorf("Not valid url %v.", i)
		}
	}

	if i, ok := props["exchange"]; ok {
		if e, ok := i.(string); ok {
			rcs.Exchange = e
		} else {
			return fmt.Errorf("Not valid exchange %v.", i)
		}
	}

	if i, ok := props["exchangeType"]; ok {
		if e, ok := i.(string); ok {
			rcs.ExchangeType = e
		} else {
			return fmt.Errorf("Not valid exchangeType %v.", i)
		}
	}

	if i, ok := props["routingKey"]; ok {
		if r, ok := i.(string); ok {
			rcs.RoutingKey = r
		} else {
			return fmt.Errorf("Not valid routingKey %v.", i)
		}
	}

	conf.Log.Debugf("Initialized with configurations %#v.", rcs)
	return nil
}

func (rcs *rabbitMQConvertSink) Open(ctx api.StreamContext) error {
	logger := ctx.GetLogger()
	logger.Infof("Sink connet to rabbitmq.")
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", rcs.Username, rcs.Password, rcs.URL))
	if err != nil {
		return err
	}
	rcs.conn = conn
	logger.Infof("Sink declare a channel.")
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	rcs.channel = channel
	return nil
}

func (rcs *rabbitMQConvertSink) Collect(ctx api.StreamContext, data interface{}) error {
	logger := ctx.GetLogger()
	if err := rcs.channel.ExchangeDeclare(
		rcs.Exchange,
		rcs.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}
	bytes, _, err := ctx.TransformOutput(data)
	if err != nil {
		return err
	}
	// Set routingKey
	switch i := data.(type) {
	case []map[string]interface{}:
		for _, value := range i {
			rcs.RoutingKey = value["topic"].(string)
		}
	}
	logger.Infof("routingKey", rcs.RoutingKey)

	if err := rcs.channel.Publish(
		rcs.Exchange,
		string(rcs.RoutingKey),
		false,
		false,
		amqp.Publishing{
			Body: bytes,
		},
	); err != nil {
		return err
	}

	return nil
}

func (rcs *rabbitMQConvertSink) Close(ctx api.StreamContext) error {
	rcs.channel.Close()
	rcs.conn.Close()
	return nil
}

func GetConvertSink() *rabbitMQConvertSink {
	return &rabbitMQConvertSink{}
}
