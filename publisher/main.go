package main

import (
	"adaptive-moms/parameters"
	"adaptive-moms/shared"
	"fmt"
	"github.com/streadway/amqp"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Parameters struct {
	NumberOfClients  int
	Id               string
	RabbitMQHost     string
	RabbitMQPort     int
	QueueName        string
	NumberOfRequests int
	Mean             float64
	StdDev           float64
	MsgSize          int
	Conn             *amqp.Connection
	Ch               *amqp.Channel
	Queue            amqp.Queue
	Msgs             <-chan amqp.Delivery
}

type Publisher struct {
	Params Parameters
}

func NewPublisher() Publisher {
	return Publisher{}
}

func main() { // Windows
	clientId := "c1"
	p := parameters.LoadParameters()
	publisher := NewPublisher()

	publisher.Initialise(clientId, p)

	// Run a single publisher
	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//publisher.Run(&wg)

	// Run Experiments
	RunExperiments(publisher, p)
}

func (p *Publisher) Run(wg *sync.WaitGroup) {

	// signalise the end of client
	defer wg.Done()

	// Close channels and connections (when finish)
	defer func(Conn *amqp.Connection) {
		err := Conn.Close()
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}(p.Params.Conn)
	defer func(Ch *amqp.Channel) {
		err := Ch.Close()
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}(p.Params.Ch)

	// initialize variables
	err := error(nil)

	// create & fill the message
	msg := make([]uint8, p.Params.MsgSize) // fixed 256 bytes
	for i := 0; i < p.Params.MsgSize; i++ {
		msg[i] = uint8(i % 255)
	}

	for i := 0; i < p.Params.NumberOfRequests; i++ {
		corrId := shared.RandomString(32)

		// make resquests randomly distributed -- experimental purpose -- comment
		interTime := p.Params.Mean + rand.NormFloat64()*p.Params.StdDev
		time.Sleep(time.Duration(interTime) * time.Millisecond)
		err = p.Params.Ch.Publish(
			"",                 // exchange
			p.Params.QueueName, // routing key
			false,              // mandatory
			false,              // immediate

			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: corrId,
				ReplyTo:       p.Params.Queue.Name,
				Body:          msg,
			})
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), "Failed to publish a message")
		}
	}
}

func (p *Publisher) Initialise(id string, params parameters.AllParameters) {
	p.Params.NumberOfClients = params.NumberOfClients
	p.Params.Id = id
	p.Params.RabbitMQHost = params.RabbitMQHost
	p.Params.RabbitMQPort = params.RabbitMQPort
	p.Params.QueueName = params.QueueName
	p.Params.NumberOfRequests = params.NumberOfRequests
	p.Params.Mean = params.Mean
	p.Params.StdDev = params.StdDev
	p.Params.MsgSize = params.MsgSize

	p.configureRabbitMQ(p.Params.RabbitMQHost, p.Params.RabbitMQPort, p.Params.QueueName)
}

func (p *Publisher) configureRabbitMQ(host string, port int, queueName string) {

	err := error(nil)

	p.Params.Conn, err = amqp.Dial("amqp://guest:guest@" + host + ":" + strconv.Itoa(port) + "/")
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Failed to connect to RabbitMQ broker")
	}

	p.Params.Ch, err = p.Params.Conn.Channel()
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Failed to open a channel")
	}

	// Queue - it creates a queue if it does not exist
	p.Params.Queue, err = p.Params.Ch.QueueDeclare(
		queueName, // name
		false,     // durable default is false
		false,     // delete when unused
		false,     // exclusive default is true
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Failed to declare a queue")
	}
}

func RunExperiments(p Publisher, params parameters.AllParameters) {
	wg := sync.WaitGroup{}

	for i := 0; i < p.Params.NumberOfClients; i++ {
		id := "cli-" + strings.TrimSpace(strconv.Itoa(i))
		publisher := NewPublisher()
		publisher.Initialise(id, params)
		go publisher.Run(&wg)
		wg.Add(1)
	}
	fmt.Println("All", p.Params.NumberOfClients, "Clients initialised ...")
	wg.Wait()
	fmt.Println("All", p.Params.NumberOfClients, "clients finished...")
}