package main

import (
	"flag"
	"github.com/streadway/amqp"
	"gorm.io/driver/postgres"
	"log"
	"os"
	"service/database"
	"service/router"
	"service/utils"
)

var (
	queueName = flag.String("queue_name", "authorQueue", "Name of the queue that this service will connect to.")
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func connectToChannel(connection *amqp.Connection) *amqp.Channel {
	ch, err := connection.Channel()
	failOnError(err, "Failed to open a channel")
	return ch
}

func openConnection() *amqp.Connection {
	address, ok := os.LookupEnv("RABBITMQ_ADDRESS")
	if !ok {
		address = "amqp://guest:guest@localhost:55001/"
	}
	log.Printf("Connecting to RabbitMQ at: %s", address)
	conn, err := amqp.Dial(address)
	failOnError(err, "Failed to connect to RabbitMQ")
	return conn
}

func createQueue(channel *amqp.Channel) amqp.Queue {
	q, err := channel.QueueDeclare(
		*queueName, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue")
	return q
}

func main() {
	conn := openConnection()
	channel := connectToChannel(conn)
	queue := createQueue(channel)

	messages, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	failOnError(err, "Failed to register a consumer")

	dbConnectorString, ok := os.LookupEnv("DB_CONNECTOR")
	if !ok {
		dbConnectorString = "postgres://postgres:postgrespw@localhost:55002"
	}
	log.Printf("Connecting to database at: %s\n", dbConnectorString)
	postgresDialector := postgres.Open(dbConnectorString)
	connector := database.NewConnection(postgresDialector)
	routeManager := router.NewRouteManager(connector)

	forever := make(chan bool)

	go func() {
		for message := range messages {
			body := message.Body
			event := utils.DecodeEvent(body)
			log.Printf("Received a message: %s", event.String())
			response := utils.BuildResponse(routeManager.RouteEvent(event))

			err := channel.Publish(
				"", message.ReplyTo,
				false, // mandatory
				false, // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: message.CorrelationId,
					Body:          utils.EncodeResponseToByte(response),
				})
			failOnError(err, "Failed to publish a message")
			message.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	defer channel.Close()
	defer conn.Close()
}
