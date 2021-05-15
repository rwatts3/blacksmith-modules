package topicdestination

import (
	"fmt"
	"net/url"
	"os"

	"github.com/nunchistudio/blacksmith/flow/destination"
	"github.com/nunchistudio/blacksmith/helper/errors"
)

/*
Driver is a custom type allowing the user to only pass supported drivers when
creating the destination.
*/
type Driver string

/*
DriverAWSSNS is used to leverage AWS SNS as the destination's driver.

Environment variables:
  - AWS_ACCESS_KEY_ID (required)
  - AWS_SECRET_ACCESS_KEY (required)
  - AWS_REGION
*/
var DriverAWSSNS Driver = "aws/sns"

/*
DriverAWSSQS is used to leverage AWS SQS as the destination's driver.

Environment variables:
  - AWS_ACCESS_KEY_ID (required)
  - AWS_SECRET_ACCESS_KEY (required)
  - AWS_REGION
*/
var DriverAWSSQS Driver = "aws/sqs"

/*
DriverAzureServiceBus is used to leverage Azure Service Bus as the destination's
driver.

Environment variables:
  - SERVICEBUS_CONNECTION_STRING (required)
*/
var DriverAzureServiceBus Driver = "azure/servicebus"

/*
DriverGooglePubSub is used to leverage Googe Pub / Sub as the destination's driver.

Environment variables:
  - GOOGLE_APPLICATION_CREDENTIALS (required)
*/
var DriverGooglePubSub Driver = "google/pubsub"

/*
DriverKafka is used to leverage Apache Kafka as the destination's driver.

Environment variables:
  - KAFKA_BROKERS (required)
    Example: "127.0.0.1:9092,127.0.0.1:9093,127.0.0.1:9094"
*/
var DriverKafka Driver = "kafka"

/*
DriverNATS is used to leverage NATS as the destination's driver.

Environment variables:
  - NATS_SERVER_URL (required)
    Example: "nats://127.0.0.1:4222"
*/
var DriverNATS Driver = "nats"

/*
DriverRabbitMQ is used to leverage RabbitMQ as the destination's driver.

Environment variables:
  - RABBIT_SERVER_URL (required)
    Example: "amqp://guest:guest@127.0.0.1:5672/"
*/
var DriverRabbitMQ Driver = "rabbitmq"

/*
Options is the options the destination can take as an input to be configured.
*/
type Options struct {

	// Realtime indicates if the pubsub adapter of the Blacksmith application shall
	// be used to load events to the destination in realtime or not. When false, the
	// Interval will be used.
	Realtime bool

	// Interval represents an interval or a CRON string at which a job shall be
	// loaded to the destination. It is used as the time-lapse between retries in
	// case of a job failure.
	//
	// Defaults to "@every 1h".
	Interval string

	// MaxRetries indicates the maximum number of retries per job the scheduler will
	// attempt to execute before it succeed. When the limit is reached, the job is
	// marked as "discarded".
	//
	// Defaults to 72.
	MaxRetries uint16

	// Name indicates the identifier of the topics to use in Blacksmith. The computed
	// name is "topic(<name>)". This does not have any consequences on the store's
	// name used in the cloud provider.
	//
	// Examples: "mytopic"
	// Required.
	Name string

	// Driver is the driver to leverage for using this destination.
	//
	// Required.
	Driver Driver

	// Connection is the driver's specific connection string to use for publishing
	// messages.
	//
	// Format for AWS SNS: "arn:aws:sns:<region>:<id>:<topic>"
	// Format for AWS SQS: "sqs.<region>.amazonaws.com/<id>/<queue>"
	// Format for Azure Service Bus: "<topic>"
	// Format for Google Pub / Sub: "<project>/<topic>"
	// Format for Apache Kafka: "<topic>"
	// Format for NATS: "<subject>"
	// Format for RabbitMQ: "<exchange>"
	Connection string

	// Params can be used to add specific configuration per driver.
	//
	// Supported fields for AWS SNS / SQS:
	//   url.Values{
	//     "region": {"<region>"}, // Required if environment variable 'AWS_REGION' is not set.
	//   }
	Params url.Values
}

/*
validate ensures the options passed to initialize the destination are valid.
*/
func (env *Options) validate() error {
	var interval string = destination.Defaults.DefaultSchedule.Interval
	var maxRetries uint16 = destination.Defaults.DefaultSchedule.MaxRetries

	fail := &errors.Error{
		Message:     "topic: Failed to load",
		Validations: []errors.Validation{},
	}

	if env == nil {
		fail.Validations = append(fail.Validations, errors.Validation{
			Message: "Options must not be nil",
			Path:    []string{"Options", "Destinations", "topic"},
		})

		return fail
	}

	if env.Interval == "" {
		env.Interval = interval
	}

	if env.MaxRetries == 0 {
		env.MaxRetries = maxRetries
	}

	if env.Name == "" {
		fail.Validations = append(fail.Validations, errors.Validation{
			Message: "Topic name must be set",
			Path:    []string{"Options", "Destinations", "topic", "Name"},
		})
	}

	// Create the computed name of the destination.
	name := fmt.Sprintf("topic(%s)", env.Name)

	if env.Driver == "" {
		fail.Validations = append(fail.Validations, errors.Validation{
			Message: "Topic driver must be set",
			Path:    []string{"Options", "Destinations", name, "Driver"},
		})
	}

	if env.Connection == "" {
		fail.Validations = append(fail.Validations, errors.Validation{
			Message: "Topic connection must be set",
			Path:    []string{"Options", "Destinations", name, "Connection"},
		})
	}

	switch env.Driver {
	case DriverAWSSNS, DriverAWSSQS:
		fail.Validations = append(fail.Validations, env.validateDriverAWSSNSSQS(name)...)
	case DriverAzureServiceBus:
		fail.Validations = append(fail.Validations, env.validateDriverAzureServiceBus(name)...)
	case DriverGooglePubSub:
		fail.Validations = append(fail.Validations, env.validateDriverGooglePubSub(name)...)
	case DriverKafka:
		fail.Validations = append(fail.Validations, env.validateDriverKafka(name)...)
	case DriverNATS:
		fail.Validations = append(fail.Validations, env.validateDriverNATS(name)...)
	case DriverRabbitMQ:
		fail.Validations = append(fail.Validations, env.validateDriverRabbitMQ(name)...)
	}

	if len(fail.Validations) > 0 {
		return fail
	}

	return nil
}

/*
validateDriverAWSSNSSQS is part of the options' validation process and validate
those for the AWS SNS / SQS drivers.
*/
func (env *Options) validateDriverAWSSNSSQS(name string) []errors.Validation {
	validations := []errors.Validation{}

	// Add a validation error if the AWS access key is not set.
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		validations = append(validations, errors.Validation{
			Message: "Environment variable 'AWS_ACCESS_KEY_ID' not set",
			Path:    []string{"Options", "Destinations", name},
		})
	}

	// Add a validation error if the AWS secret access key is not set.
	if os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		validations = append(validations, errors.Validation{
			Message: "Environment variable 'AWS_SECRET_ACCESS_KEY' not set",
			Path:    []string{"Options", "Destinations", name},
		})
	}

	// Try to find the desired region to use from the environment variable and
	// then from the params.
	var region = os.Getenv("AWS_REGION")
	if region == "" {
		region = env.Params.Get("region")
	}

	// Add a validation error if no region has been found since it's required to
	// work with DynamoDB.
	if region == "" {
		validations = append(validations, errors.Validation{
			Message: "Environment variable 'AWS_REGION' not set",
			Path:    []string{"Options", "Destinations", name},
		})
	}

	return validations
}

/*
validateDriverAzureServiceBus is part of the options' validation process and validate
those for the Azure Service Bus driver.
*/
func (env *Options) validateDriverAzureServiceBus(name string) []errors.Validation {
	validations := []errors.Validation{}

	// Add a validation error if the Azure connection string is not set.
	if os.Getenv("SERVICEBUS_CONNECTION_STRING") == "" {
		validations = append(validations, errors.Validation{
			Message: "Environment variable 'SERVICEBUS_CONNECTION_STRING' not set",
			Path:    []string{"Options", "Destinations", name},
		})
	}

	return validations
}

/*
validateDriverGooglePubSub is part of the options' validation process and validate
those for the Google Cloud Storage driver.
*/
func (env *Options) validateDriverGooglePubSub(name string) []errors.Validation {
	validations := []errors.Validation{}

	// Add a validation error if the Google Cloud account is not set.
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		validations = append(validations, errors.Validation{
			Message: "Environment variable 'GOOGLE_APPLICATION_CREDENTIALS' not set",
			Path:    []string{"Options", "Destinations", name},
		})
	}

	return validations
}

/*
validateDriverKafka is part of the options' validation process and validate those
for the Apache Kafka driver.
*/
func (env *Options) validateDriverKafka(name string) []errors.Validation {
	validations := []errors.Validation{}

	// Add a validation error if the Kafka Brokers is not set.
	if os.Getenv("KAFKA_BROKERS") == "" {
		validations = append(validations, errors.Validation{
			Message: "Environment variable 'KAFKA_BROKERS' not set",
			Path:    []string{"Options", "Destinations", name},
		})
	}

	return validations
}

/*
validateDriverNATS is part of the options' validation process and validate those
for the NATS driver.
*/
func (env *Options) validateDriverNATS(name string) []errors.Validation {
	validations := []errors.Validation{}

	// Add a validation error if the Kafka Brokers is not set.
	if os.Getenv("NATS_SERVER_URL") == "" {
		validations = append(validations, errors.Validation{
			Message: "Environment variable 'NATS_SERVER_URL' not set",
			Path:    []string{"Options", "Destinations", name},
		})
	}

	return validations
}

/*
validateDriverRabbitMQ is part of the options' validation process and validate
those for the RabbitMQ driver.
*/
func (env *Options) validateDriverRabbitMQ(name string) []errors.Validation {
	validations := []errors.Validation{}

	// Add a validation error if the Kafka Brokers is not set.
	if os.Getenv("RABBIT_SERVER_URL") == "" {
		validations = append(validations, errors.Validation{
			Message: "Environment variable 'RABBIT_SERVER_URL' not set",
			Path:    []string{"Options", "Destinations", name},
		})
	}

	return validations
}
