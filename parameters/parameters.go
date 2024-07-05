package parameters

import (
	"adaptive-moms/shared"
	"github.com/spf13/viper"
	"time"
)

type AllParameters struct {
	Alfa                  float64
	Beta                  float64
	Direction             float64
	DeffuzificationMethod string
	HysteresisBand        float64
	Kp                    float64
	Ki                    float64
	Kd                    float64
	Max                   float64
	Min                   float64
	MembershipFunction    string
	PC                    float64
	SetPoint              float64
	DeltaTime             time.Duration
	RabbitMQHost          string
	RabbitMQPort          int
	QueueName             string
	MonitorTime           time.Duration
	NumberOfRequests      int
	Mean                  float64
	StdDev                float64
	MsgSize               int
	NumberOfClients       int
	OutputFile            string
	DockerDir             string
	Deadzone              float64
	ControllerType        string
}

func LoadParameters() AllParameters {
	r := AllParameters{}

	// Set the file name of the configuration file
	fileName := "config"
	viper.SetConfigName(fileName)

	// Set the path to look for the configuration file
	filePath := "C:\\Users\\user\\go\\adaptive-moms\\data"
	viper.AddConfigPath(filePath)

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	// Read the configuration file
	err := viper.ReadInConfig()
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// Get values from the configuration file or environment variables
	r.Direction = viper.GetFloat64("Direction")
	r.HysteresisBand = viper.GetFloat64("HysteresisBand")
	r.Kp = viper.GetFloat64("Kp")
	r.Ki = viper.GetFloat64("Ki")
	r.Kd = viper.GetFloat64("Kd")
	r.Max = viper.GetFloat64("Max")
	r.Min = viper.GetFloat64("Min")
	r.PC = viper.GetFloat64("PC")
	r.SetPoint = viper.GetFloat64("SetPoint")
	r.DeltaTime = viper.GetDuration("DeltaTime") * time.Second
	r.RabbitMQHost = viper.GetString("RabbitMQHost")
	r.RabbitMQPort = viper.GetInt("RabbitMQPort")
	r.QueueName = viper.GetString("QueueName")
	r.NumberOfRequests = viper.GetInt("NumberOfRequests")
	r.Mean = viper.GetFloat64("Mean")
	r.StdDev = viper.GetFloat64("StdDev")
	r.NumberOfClients = viper.GetInt("NumberOfClients")
	r.MonitorTime = viper.GetDuration("MonitorTime") * time.Second
	r.OutputFile = viper.GetString("OutputFile")
	r.DockerDir = viper.GetString("DockerDir")
	r.Deadzone = viper.GetFloat64("Deadzone")
	r.Alfa = viper.GetFloat64("Alfa")
	r.Beta = viper.GetFloat64("Beta")
	r.DeffuzificationMethod = viper.GetString("DeffuzificationMethod")
	r.MembershipFunction = viper.GetString("MembershipFunction")
	r.ControllerType = viper.GetString("ControllerType")

	return r
}
