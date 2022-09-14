package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"portalnesia.com/api/config"
)

const (
	CHANNEL_SEND_EMAIL    = "send_email"
	CHANNEL_SEND_MILIS    = "send_milis"
	CHANNEL_SEND_WHATSAPP = "send_whatsapp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func RabbitmqQueue(channel string, msg interface{}) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", os.Getenv("RABBIT_USER"), os.Getenv("RABBIT_PASS"), os.Getenv("DB_HOST"), os.Getenv("RABBIT_PORT")))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		channel, //name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body, err := json.Marshal(msg)
	failOnError(err, "Failed to json encoded messages")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if config.NODE_ENV != "test" {
		err = ch.PublishWithContext(
			ctx,
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		failOnError(err, "Failed to publish messages")
	}
}
