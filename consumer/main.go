package main

import (
	"flag"
	"github.com/streadway/amqp"
	"fmt"
	"log"
)

var (
	hostname, username, password string
	port                         int
)

func main() {
	flag.StringVar(&hostname, "h", "localhost", "rabbit hostname")
	flag.StringVar(&username, "u", "user", "rabbit user")
	flag.StringVar(&password, "up", "password", "rabbit password")
	flag.IntVar(&port, "p", 5672, "rabbit port")
	flag.Parse()

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, hostname, port))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln(err)
	}
	defer ch.Close()


	if err := ch.ExchangeDeclare("event-tracker", "topic", true, false, false,
		true, nil); err != nil {
		log.Fatalln(err)
	}

	q, err := ch.QueueDeclare("mail-service", true, true, false, false, nil)
	if err != nil {
		log.Fatalln(err)
	}

	if err := ch.QueueBind(q.Name, "StopAD.*", "event-tracker", false, nil); err != nil {
		log.Fatalln(err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalln(err)
	}

	for d := range msgs {
		log.Printf(" [x] %s", d.Body)
	}
}
