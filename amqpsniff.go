package main

import (
	"encoding/json"
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/streadway/amqp"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type flags struct {
	dsn      string
	queue    string
	bindings []string
	output   string
	limit    int
}

func main() {
	kingpin.CommandLine.Help = "AMQP sniffer creates queue and dumps incoming messages."
	f := flags{}

	kingpin.Flag("dsn", "Connection credentials, env var AMQP_DSN, example: amqp://guest:guest@rabbit:5672").
		Envar("AMQP_DSN").Required().StringVar(&f.dsn)
	kingpin.Flag("queue", "Queue name, example: my-queue").Required().StringVar(&f.queue)
	kingpin.Flag("bindings", "Bindings, colon-separated exchange and routing key, example: users:*.user.created").
		Required().StringsVar(&f.bindings)
	kingpin.Flag("output", "Path to output JSONL file").Default("messages.jsonl").StringVar(&f.output)
	kingpin.Flag("limit", "Number messages to collect before exiting").Default("100").IntVar(&f.limit)

	kingpin.Parse()

	conn, err := amqp.Dial(f.dsn)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	dump, err := os.OpenFile(f.output, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-exit

		_ = dump.Close()
		_ = ch.Close()
		_ = conn.Close()
	}()

	_, err = ch.QueueDeclare(f.queue, false, true, true, false, nil)

	for _, binding := range f.bindings {
		b := strings.SplitN(binding, ":", 2)

		if err := ch.QueueBind(f.queue, b[1], b[0], false, nil); err != nil {
			log.Fatalf("failed to bind queue: %v", err)
		}
	}

	msgs, err := ch.Consume(f.queue, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	enc := json.NewEncoder(dump)

	i := 0

	for msg := range msgs {
		i++

		if i%10 == 0 {
			fmt.Println(i)
		} else {
			fmt.Print(".")
		}

		var v interface{}
		err = json.Unmarshal(msg.Body, &v)
		if err != nil { // Not a valid JSON Payload
			v = string(msg.Body)
		}

		err = enc.Encode(v)
		if err != nil {
			log.Fatal(err)
		}

		if i > f.limit {
			break
		}
	}

	fmt.Println(i)
}
