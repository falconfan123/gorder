package consumer

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"

	"github.com/falconfan123/gorder/common/broker"
	"github.com/falconfan123/gorder/common/genproto/orderpb"
	"github.com/falconfan123/gorder/payment/app"
	"github.com/falconfan123/gorder/payment/app/command"
	"github.com/goccy/go-json"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	app app.Application
}

func NewConsumer(app app.Application) *Consumer {
	return &Consumer{
		app: app,
	}
}

func (c *Consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Warnf("fail to consume: queue=%s, err=%v", q.Name, err)
	}

	var forever chan struct{}

	//Continuously reading messages from msgs,  in an infinite loop
	go func() {
		for msg := range msgs {
			c.handleMessage(msg, q)
		}
	}()
	<-forever

}

func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue) {
	logrus.Infof("Payment receive a message from%s, msg = %v", q.Name, string(msg.Body))
	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
	tr := otel.Tracer("rabbitmq")
	_, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
	defer span.End()

	o := &orderpb.Order{}
	if err := json.Unmarshal(msg.Body, o); err != nil {
		logrus.Infof("fail to unmarshal msg to order, err=%v", err)
		_ = msg.Nack(false, false)
		return
	}
	if _, err := c.app.Commands.CreatePayment.Handle(ctx, command.CreatePayment{Order: o}); err != nil {
		//TODO: retry
		logrus.Infof("fail to create order, err=%v", err)
		_ = msg.Nack(false, false)
		return
	}

	span.AddEvent("payment.created")
	//create payment link
	_ = msg.Ack(false)
	logrus.Info("consume success")
}
