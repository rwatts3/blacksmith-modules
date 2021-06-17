package docstoredestination

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/nunchistudio/blacksmith/destination"
	"github.com/nunchistudio/blacksmith/helper/errors"
)

/*
Driver is a custom type allowing the user to only pass supported drivers when
creating the destination.
*/
type Driver string

/*
DriverAWSDynamoDB is used to leverage AWS DynamoDB as the destination's driver.

Environment variables:
  - AWS_ACCESS_KEY_ID (required)
  - AWS_SECRET_ACCESS_KEY (required)
  - AWS_REGION
*/
var DriverAWSDynamoDB Driver = "aws/dynamodb"

/*
DriverAzureCosmosDB is used to leverage Azure CosmosDB as the destination's driver.

Environment variables:
  - MONGO_SERVER_URL (required)
*/
var DriverAzureCosmosDB Driver = "azure/cosmosdb"

/*
DriverMongoDB is used to leverage MongoDB as the destination's driver.

Environment variables:
  - MONGO_SERVER_URL (required)
*/
var DriverMongoDB Driver = "mongodb"

/*
DriverGoogleFirestore is used to leverage Google Firestore as the destination's
driver.

Environment variables:
  - GOOGLE_APPLICATION_CREDENTIALS (required)
*/
var DriverGoogleFirestore Driver = "google/firestore"

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

	// Name indicates the identifier of the document store to use in Blacksmith.
	// The computed name is "docstore(<name>)". This does not have any consequences
	// on the store's name used in the cloud provider.
	//
	// Examples: "mynosql"
	// Required.
	Name string

	// Driver is the driver to leverage for using this destination.
	//
	// Required.
	Driver Driver

	// Connection is the driver's specific connection string to use for writing
	// into the document store.
	//
	// Format for AWS DynamoDB: "<table>"
	// Format for Azure CosmosDB and MongoDB: "<db>/<collection>"
	// Format for Google Firestore: "projects/<project>/databases/(default)/documents/<collection>"
	Connection string

	// Params can be used to add specific configuration per driver.
	//
	// Supported fields for AWS DynamoDB:
	//   url.Values{
	//     "region": {"<region>"}, // Required if environment variable 'AWS_REGION' is not set.
	//     "partition_key": {"<key>"}, // Required. The path to the partition key of a table or an index.
	//     "sort_key": {"<key>"}, // Optional. The path to the sort key of a table or an index.
	//   }
	//
	// Supported fields for Azure CosmosDB and MongoDB:
	//   url.Values{
	//     "id_field": {"<field>"}, // Optional. The field name to use for the "_id" field.
	//   }
	//
	// Supported fields for Google Firestore:
	//   url.Values{
	//     "name_field": {"<field>"}, // Required. The designated field for the primary key.
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
		Message:     "docstore: Failed to load",
		Validations: []errors.Validation{},
	}

	if env == nil {
		fail.Validations = append(fail.Validations, errors.Validation{
			Message: "Options must not be nil",
			Path:    []string{"Options", "Destinations", "docstore"},
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
			Message: "Collection name must be set",
			Path:    []string{"Options", "Destinations", "docstore", "Name"},
		})
	}

	// Create the computed name of the destination.
	name := fmt.Sprintf("docstore(%s)", env.Name)

	if env.Driver == "" {
		fail.Validations = append(fail.Validations, errors.Validation{
			Message: "Collection driver must be set",
			Path:    []string{"Options", "Destinations", name, "Driver"},
		})
	}

	if env.Connection == "" {
		fail.Validations = append(fail.Validations, errors.Validation{
			Message: "Collection connection must be set",
			Path:    []string{"Options", "Destinations", name, "Connection"},
		})
	}

	switch env.Driver {
	case DriverAWSDynamoDB:
		fail.Validations = append(fail.Validations, env.validateDriverAWSDynamoDB(name)...)
	case DriverAzureCosmosDB:
		fail.Validations = append(fail.Validations, env.validateDriverAzureCosmosDB(name)...)
	case DriverGoogleFirestore:
		fail.Validations = append(fail.Validations, env.validateDriverGoogleFirestore(name)...)
	case DriverMongoDB:
		fail.Validations = append(fail.Validations, env.validateDriverMongoDB(name)...)
	}

	if len(fail.Validations) > 0 {
		return fail
	}

	return nil
}

/*
validateDriverAWSDynamoDB is part of the options' validation process and validate
those for the AWS DynamoDB driver.
*/
func (env *Options) validateDriverAWSDynamoDB(name string) []errors.Validation {
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

	// Add a validation error if the AWS region is not set.
	if os.Getenv("AWS_REGION") == "" {
		validations = append(validations, errors.Validation{
			Message: "Environment variable 'AWS_REGION' not set",
			Path:    []string{"Options", "Destinations", name},
		})
	}

	// Add a validation error if the DynamoDB partition key is not set.
	if env.Params.Get("partition_key") == "" {
		validations = append(validations, errors.Validation{
			Message: "'partition_key' must be set",
			Path:    []string{"Options", "Destinations", name, "Params"},
		})
	}

	return validations
}

/*
validateDriverAzureCosmosDB is part of the options' validation process and validate
those for the Azure CosmosDB driver.
*/
func (env *Options) validateDriverAzureCosmosDB(name string) []errors.Validation {
	return env.validateDriverMongoDB(name)
}

/*
validateDriverGoogleFirestore is part of the options' validation process and validate
those for the Google Cloud Storage driver.
*/
func (env *Options) validateDriverGoogleFirestore(name string) []errors.Validation {
	validations := []errors.Validation{}

	// Add a validation error if the Google Cloud account is not set.
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		validations = append(validations, errors.Validation{
			Message: "Environment variable 'GOOGLE_APPLICATION_CREDENTIALS' not set",
			Path:    []string{"Options", "Destinations", name},
		})
	}

	// Add a validation error if the Firestore unique field  is not set.
	if env.Params.Get("name_field") == "" {
		validations = append(validations, errors.Validation{
			Message: "'name_field' must be set",
			Path:    []string{"Options", "Destinations", name, "Params"},
		})
	}

	return validations
}

/*
validateDriverMongoDB is part of the options' validation process and validate
those for the MongoDB driver.
*/
func (env *Options) validateDriverMongoDB(name string) []errors.Validation {
	validations := []errors.Validation{}

	// Add a validation error if the MongoDB server URL is not set.
	if os.Getenv("MONGO_SERVER_URL") == "" {
		validations = append(validations, errors.Validation{
			Message: "Environment variable 'MONGO_SERVER_URL' not set",
			Path:    []string{"Options", "Destinations", name},
		})
	} else if !strings.HasPrefix(os.Getenv("MONGO_SERVER_URL"), "mongodb://") {
		validations = append(validations, errors.Validation{
			Message: "Environment variable 'MONGO_SERVER_URL' not valid",
			Path:    []string{"Options", "Destinations", name},
		})
	}

	// Add a validation error if the MongoDB connection string is not valid.
	if len(strings.Split(env.Connection, "/")) != 2 {
		validations = append(validations, errors.Validation{
			Message: "Connection string not valid",
			Path:    []string{"Options", "Destinations", name, "Connection"},
		})
	}

	return validations
}
