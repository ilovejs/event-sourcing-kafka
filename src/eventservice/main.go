package main

import (
	"flag"
	"fmt"

	"github.com/Shopify/sarama"

	"github.com/ilovejs/event-sourcing-kafka/src/eventservice/rest"
	"github.com/ilovejs/event-sourcing-kafka/src/lib/configuration"
	"github.com/ilovejs/event-sourcing-kafka/src/lib/msgqueue"
	msgqueue_amqp "github.com/ilovejs/event-sourcing-kafka/src/lib/msgqueue/amqp"
	"github.com/ilovejs/event-sourcing-kafka/src/lib/msgqueue/kafka"
	"github.com/ilovejs/event-sourcing-kafka/src/lib/persistence/dblayer"
	"github.com/streadway/amqp"
)

func main() {
	var eventEmitter msgqueue.EventEmitter
	// filename:string
	confPath := flag.String("conf",
		`.\configuration\config.json`,
		"flag to set the path to the configuration json file")

	flag.Parse()

	// Extract configuration, see function for details:
	config, _ := configuration.ExtractConfiguration(*confPath)

	switch config.MessageBrokerType {
	case "amqp":
		conn, err := amqp.Dial(config.AMQPMessageBroker)
		if err != nil {
			panic(err)
		}

		eventEmitter, err = msgqueue_amqp.NewAMQPEventEmitter(conn, "events")
		if err != nil {
			panic(err)
		}
	case "kafka":
		conf := sarama.NewConfig()
		conf.Producer.Return.Successes = true
		conn, err := sarama.NewClient(config.KafkaMessageBrokers, conf)
		if err != nil {
			panic(err)
		}

		eventEmitter, err = kafka.NewKafkaEventEmitter(conn)
		if err != nil {
			panic(err)
		}
	default:
		panic("Bad message broker type: " + config.MessageBrokerType)
	}

	fmt.Println("Connecting to database")
	dbhandler, _ := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)

	fmt.Println("Serving API")
	//RESTful API start
	err := rest.ServeAPI(config.RestfulEndpoint, dbhandler, eventEmitter)
	if err != nil {
		panic(err)
	}
}
