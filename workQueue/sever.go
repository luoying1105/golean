package workQueue

import (
	"fullstack/tools"
	"github.com/streadway/amqp"
	"log"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func mqConnect() (*amqp.Connection, *amqp.Channel, error) {

	conn, err := amqp.Dial(mqUrl)
	tools.FailOnError(err, "failed to connect tp rabbitmq")
	channel, err := conn.Channel()
	tools.FailOnError(err, "failed to open a channel")
	return conn, channel, err
}

func PushData(msgContent *string) {

	conn, channel, err := mqConnect()
	tools.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	defer channel.Close()

	err = channel.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	tools.FailOnError(err, "Failed to declare a queue")

	err = channel.Publish(
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(*msgContent),
		})
	log.Printf(" [x] Sent %s", *msgContent)
	tools.FailOnError(err, "Failed to publish a message")

}

func Receive() {
	conn, channel, err := mqConnect()
	tools.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	defer channel.Close()

	err = channel.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	tools.FailOnError(err, "Failed to declare an exchange")

	q, err := channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	tools.FailOnError(err, "Failed to declare a queue")

	err = channel.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil,
	)
	tools.FailOnError(err, "Failed to bind a queue")

	msgs, err := channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	tools.FailOnError(err, "")
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
