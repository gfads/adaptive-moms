package main

import (
	"adaptive-moms/controllers"
	"adaptive-moms/parameters"
	"adaptive-moms/shared"
	"fmt"
	"github.com/streadway/amqp"
	"math"
	_ "net/http/pprof"
	"os"
	"strconv"
	"sync"
	"time"
)

type Parameters struct {
	Conn           *amqp.Connection
	Ch             *amqp.Channel
	Queue          amqp.Queue
	Msgs           <-chan amqp.Delivery
	RabbitMQHost   string
	RabbitMQPort   int
	QueueName      string
	ControllerType string
	Controller     controllers.Controller
	PC             float64
	MonitorTime    time.Duration
	SetPoint       float64
}

type Subscriber struct {
	Params Parameters
}

func NewSubscriber() Subscriber {
	return Subscriber{}
}

func main() {
	// load configuration parameters
	p := parameters.LoadParameters()

	// define and open csv file to record experiment results
	dataFileName := p.OutputFile
	df, err := os.Create(p.DockerDir + "\\" + dataFileName)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	defer df.Close()

	// Create & initialise the controller
	controller := controllers.NewController(p.ControllerType)
	controller.Initialise(p)

	// Create & initialise the subscriber
	subscriber := NewSubscriber()
	subscriber.Initialise(p, controller)

	// Run closed loop, i.e., run with a controller
	subscriber.RunClosedLoop(df)
}

func (s *Subscriber) Initialise(p parameters.AllParameters, c controllers.Controller) {
	s.Params.RabbitMQHost = p.RabbitMQHost
	s.Params.RabbitMQPort = p.RabbitMQPort
	s.Params.QueueName = p.QueueName
	s.Params.PC = p.PC
	s.Params.MonitorTime = time.Duration(p.MonitorTime)
	s.Params.SetPoint = p.SetPoint
	s.Params.Controller = c

	s.configureRabbitMQ(s.Params.RabbitMQHost, s.Params.RabbitMQPort, s.Params.QueueName, s.Params.PC)
}

func (s Subscriber) RunOpenLoop(wg *sync.WaitGroup) {
	defer wg.Done()

	// Close channels and connections (when finish)
	defer func(Conn *amqp.Connection) {
		err := Conn.Close()
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}(s.Params.Conn)
	defer func(Ch *amqp.Channel) {
		err := Ch.Close()
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}(s.Params.Ch)

	for d := range s.Params.Msgs {
		err := d.Ack(false) // send ack to broker
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		fmt.Println(d.Body)
	}
}

func (s *Subscriber) Warmup() {

	fmt.Println("Begin of Warming up...")

	// configure pc to zero
	err := s.Params.Ch.Qos(
		0,    // prefetch count
		0,    // prefetch size
		true, // global TODO default is false
	)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Failed to set QoS")
	}
	for i := 0; i < shared.WarmupMessages; i++ {
		//for i := 0; i < 10; i++ {
		d := <-s.Params.Msgs
		err := d.Ack(false) // send ack to broker
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}

	// return pc back to initial pc
	err = s.Params.Ch.Qos(
		int(s.Params.PC), // prefetch count
		0,                // prefetch size
		true,             // global TODO default is false
	)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Failed to set QoS")
	}
	fmt.Println("End of Warming up...")
}

func (s *Subscriber) RunMonitoredOpenLoop(wg *sync.WaitGroup) {

	defer wg.Done()
	n := 0 // number of receive messages

	// receive messages
	tt := time.Tick(s.Params.MonitorTime * time.Second)
	for d := range s.Params.Msgs {
		err := d.Ack(false) // send ack to broker
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		n++ // increment the number of received messages
		select {
		case <-tt:
			// inspect queue
			s.Params.Queue, err = s.Params.Ch.QueueInspect(s.Params.QueueName)
			if err != nil {
				shared.ErrorHandler(shared.GetFunction(), "Impossible to inspect the queue")
				os.Exit(0)
			}

			// calculate rate
			rate := float64(n) / float64(s.Params.MonitorTime)
			fmt.Printf("%d;%.2f;%v;%d\n", s.Params.PC, rate, s.Params.MonitorTime, n)
			n = 0
		default:
		}
	}
}

func (s *Subscriber) RunClosedLoop(df *os.File) {

	// initialise the counter of received messages
	n := 0

	// configure timer
	tt := time.Tick(s.Params.MonitorTime)

	// receive messages
	for d := range s.Params.Msgs {
		err := d.Ack(false) // send ack to broker
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		n++ // increment the number of received messages
		select {
		case <-tt:
			// inspect queue
			s.Params.Queue, err = s.Params.Ch.QueueInspect(s.Params.QueueName)
			if err != nil {
				shared.ErrorHandler(shared.GetFunction(), "Impossible to inspect the queue")
				os.Exit(0)
			}

			// calculate rate
			rate := float64(n) / float64(s.Params.MonitorTime.Seconds())

			// register experiment data
			fmt.Fprintf(df, "%.0f;%.2f;%.2f\n", s.Params.PC, rate, s.Params.SetPoint)

			// compute new pc
			newPC := s.Params.Controller.Update(s.Params.SetPoint, rate, s.Params.PC)
			err := s.Params.Ch.Qos(
				int(newPC), // prefetch count
				0,          // prefetch size
				true,       // global - default is false
			)
			if err != nil {
				shared.ErrorHandler(shared.GetFunction(), "Failed to set QoS")
			}

			// update pc
			s.Params.PC = math.Round(newPC)

			// reset no. of received messages
			n = 0
		default: // receive next message
		}
	}
}

func (c *Subscriber) configureRabbitMQ(host string, port int, queueName string, pc float64) {
	err := error(nil)

	// create connection
	c.Params.Conn, err = amqp.Dial("amqp://guest:guest@" + c.Params.RabbitMQHost + ":" + strconv.Itoa(c.Params.RabbitMQPort) + "/")
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Failed to connect to RabbitMQ broker")
	}

	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Failed to connect to RabbitMQ")
	}

	// create channel
	c.Params.Ch, err = c.Params.Conn.Channel()
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Failed to open a channel")
	}

	// declare queues
	c.Params.Queue, err = c.Params.Ch.QueueDeclare(
		c.Params.QueueName, // name
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Failed to declare a queue")
	}

	// create a consumer
	c.Params.Msgs, err = c.Params.Ch.Consume(
		c.Params.Queue.Name, // queue
		"",                  // consumer
		false,               // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Failed to register a consumer")
	}

	// configure initial QoS of Req channel
	err = c.Params.Ch.Qos(
		int(pc), // prefetch count
		0,       // prefetch size
		true,    // global TODO default is false
	)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Failed to set QoS")
	}
	return
}
