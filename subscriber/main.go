package main

import (
	"adaptive-moms/controllers"
	"adaptive-moms/parameters"
	"adaptive-moms/shared"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"math"
	_ "net/http/pprof"
	"os"
	"strconv"
	"time"
)

type Parameters struct {
	Conn           *amqp091.Connection
	Ch             *amqp091.Channel
	Queue          amqp091.Queue
	Msgs           <-chan amqp091.Delivery
	RabbitMQHost   string
	RabbitMQPort   int
	QueueName      string
	ControllerType string
	Controller     controllers.Controller
	PC             float64
	MonitorTime    time.Duration
	SetPoints      []int
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

	// Create & initialise the controller
	controller := controllers.NewController(p.ControllerType)
	controller.Initialise(p)

	// Create & initialise the subscriber
	subscriber := NewSubscriber()
	subscriber.Initialise(p, controller)

	subscriber.Run(p)
}

func (s *Subscriber) Run(p parameters.AllParameters) {

	switch p.ExecutionType {
	case shared.OpenLoop:
		s.RunOpenLoop(p)
	case shared.MonitoredOpenLoop:
		s.RunMonitoredOpenLoop(p)
	case shared.ClosedLoop:
		s.RunClosedLoop(p)
	case shared.ExperimentClosedLoop:
		s.RunExperimentClosedLoop(p)
	default:
		shared.ErrorHandler(shared.GetFunction(), "Unknown ´Execution Type'")
	}
}

func (s *Subscriber) Initialise(p parameters.AllParameters, c controllers.Controller) {
	s.Params.RabbitMQHost = p.RabbitMQHostSub
	s.Params.RabbitMQPort = p.RabbitMQPort
	s.Params.QueueName = p.QueueName
	s.Params.PC = p.PC
	s.Params.MonitorTime = time.Duration(p.MonitorTime)
	s.Params.SetPoints = p.SetPoints
	s.Params.Controller = c

	s.configureRabbitMQ(p)
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

func (s Subscriber) RunOpenLoop(p parameters.AllParameters) {

	// Close channels and connections (when finish)
	defer func(Conn *amqp091.Connection) {
		err := Conn.Close()
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}(s.Params.Conn)
	defer func(Ch *amqp091.Channel) {
		err := Ch.Close()
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}(s.Params.Ch)

	fmt.Printf("Subscriber running [%v] ...\n", p.ExecutionType)
	for d := range s.Params.Msgs {
		err := d.Ack(false) // send ack to broker
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		fmt.Println(d.Body)
	}
}

func (s *Subscriber) RunMonitoredOpenLoop(p parameters.AllParameters) {
	n := 0 // number of receive messages

	// configure monitor timer
	tt := time.Tick(s.Params.MonitorTime)

	// receive messages
	fmt.Printf("Subscriber running [%v] ...\n", p.ExecutionType)
	err := error(nil)
	for {
		n++ // increment the number of received messages
		select {
		case d := <-s.Params.Msgs:
			err = d.Ack(false) // send ack to broker
			if err != nil {
				shared.ErrorHandler(shared.GetFunction(), err.Error())
			}
			n++ // increment the number of received messages
		case <-tt:
			// inspect queue
			s.Params.Queue, err = s.Params.Ch.QueueInspect(s.Params.QueueName)
			if err != nil {
				shared.ErrorHandler(shared.GetFunction(), "Impossible to inspect the queue")
				os.Exit(0)
			}

			// calculate rate
			rate := float64(n) / float64(s.Params.MonitorTime.Seconds())
			fmt.Printf("%.0f;%.2f;%v\n", s.Params.PC, rate, s.Params.SetPoints[0])
			n = 0
		}
	}
}

func (s *Subscriber) RunClosedLoop(p parameters.AllParameters) {

	// initialise the counter of received messages
	n := 0

	// configure timer
	tt := time.Tick(s.Params.MonitorTime)

	// receive message
	fmt.Printf("Subscriber running [%v] ...\n", p.ExecutionType)
	err := error(nil)
	for {
		select {
		case d := <-s.Params.Msgs:
			err = d.Ack(false) // send ack to broker
			if err != nil {
				shared.ErrorHandler(shared.GetFunction(), err.Error())
			}
			n++ // increment the number of received messages
		case <-tt:
			// inspect queue
			s.Params.Queue, err = s.Params.Ch.QueueInspect(s.Params.QueueName)
			if err != nil {
				shared.ErrorHandler(shared.GetFunction(), "Impossible to inspect the queue")
				os.Exit(0)
			}

			// calculate rate
			rate := float64(n) / float64(s.Params.MonitorTime.Seconds())

			// show current data (* not save information *)
			fmt.Printf("%.0f;%.2f;%v\n", s.Params.PC, rate, s.Params.SetPoints[0])

			// compute new pc
			newPC := s.Params.Controller.Update(float64(s.Params.SetPoints[0]), rate, s.Params.PC)
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

func (s *Subscriber) RunExperimentClosedLoop(p parameters.AllParameters) {
	err := error(nil)

	// define and open csv file to record experiment results
	dataFileName := p.OutputFile
	df, err := os.Create(p.DockerDir + "\\" + dataFileName)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	defer df.Close()

	// initialise the counter of received messages
	n := 0

	// configure timer
	tt := time.Tick(s.Params.MonitorTime)

	// configure current setpoint
	currentSetpoint := 0

	// configure current sample
	currentSample := 0

	// receive messages
	fmt.Printf("Subscriber running [%v] ...\n", p.ExecutionType)
	for {
		select {
		case d := <-s.Params.Msgs:
			err = d.Ack(false) // send ack to broker
			if err != nil {
				shared.ErrorHandler(shared.GetFunction(), err.Error())
			}
			n++ // increment the number of received messages
		case <-tt:
			// update current sample
			currentSample++

			// inspect queue
			s.Params.Queue, err = s.Params.Ch.QueueInspect(s.Params.QueueName)
			if err != nil {
				shared.ErrorHandler(shared.GetFunction(), "Impossible to inspect the queue")
				os.Exit(0)
			}

			// calculate rate
			rate := float64(n) / float64(s.Params.MonitorTime.Seconds())

			// register experiment data
			fmt.Fprintf(df, "%.0f;%.2f;%v\n", s.Params.PC, rate, s.Params.SetPoints[currentSetpoint])

			// show experiment data
			fmt.Printf("%.0f;%.2f;%v\n", s.Params.PC, rate, s.Params.SetPoints[currentSetpoint])

			// compute new pc
			newPC := s.Params.Controller.Update(float64(s.Params.SetPoints[currentSetpoint]), rate, s.Params.PC)
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

			// check current sample size & current set point
			if currentSample >= p.SampleSizePerLevel {
				currentSample = 0
				currentSetpoint++
				if currentSetpoint >= len(p.SetPoints) { // stop condition
					return
				}
			}
		}
	}
}

func (c *Subscriber) configureRabbitMQ(params parameters.AllParameters) {
	err := error(nil)

	// create connection
	c.Params.Conn, err = amqp091.Dial("amqp://guest:guest@" + params.RabbitMQHostSub + ":" + strconv.Itoa(params.RabbitMQPort) + "/")
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
		params.QueueName, // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
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
		int(params.PC), // prefetch count
		0,              // prefetch size
		true,           // global
	)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Failed to set QoS")
	}
	return
}
