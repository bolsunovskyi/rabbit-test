package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

var (
	hostname, username, password string
	port                         int
)

type event struct {
	Step      string `json:"step"`
	SubStep   string `json:"substep"`
	Timestamp int64  `json:"timestamp"`
}

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

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln(err)
	}

	if err := ch.ExchangeDeclare("event-tracker", "topic", true, false, false,
		true, nil); err != nil {
		log.Fatalln(err)
	}

	rand.Seed(time.Now().Unix())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	stepI := 0

	for {

		select {
		case <-c:
			log.Println("exit")
			if err := ch.Close(); err != nil {
				log.Println(err)
			}

			if err := conn.Close(); err != nil {
				log.Println(err)
			}

			os.Exit(0)
		default:
			e := event{
				Step:      steps[stepI % len(steps)],
				SubStep:   substeps[rand.Int31n(int32(len(substeps)))],
				Timestamp: time.Now().Unix(),
			}
			stepI++

			bts, err := json.Marshal(e)
			if err != nil {
				log.Fatalln(err)
			}
			log.Printf("%s - %s\n", e.Step, e.SubStep)
			if err := ch.Publish("event-tracker", e.Step+"."+e.SubStep, false, false, amqp.Publishing{
				ContentType: "application/json",
				Body:        bts,
			}); err != nil {
				log.Fatalln(err)
			}
			time.Sleep(time.Second)
		}

	}

	select {}
}
