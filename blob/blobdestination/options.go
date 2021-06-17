package blobdestination

import (
	"fmt"
	"net/url"
	"os"

	"github.com/nunchistudio/blacksmith/destination"
	"github.com/nunchistudio/blacksmith/helper/errors"
)

/*
Driver is a custom type allowing the user to only pass supported drivers when
creating the destination.
*/
type Driver string

/*
DriverAWSS3 is used to leverage AWS S3 as the destination's driver.

Environment variables:
  - AWS_ACCESS_KEY_ID (required)
  - AWS_SECRET_ACCESS_KEY (required)
  - AWS_REGION
*/
var DriverAWSS3 Driver = "aws/s3"

/*
DriverAzureBlob is used to leverage Azure Blob Storage as the destination's
driver.

Environment variables:
  - AZURE_STORAGE_ACCOUNT (required)
  - AZURE_STORAGE_KEY || AZURE_STORAGE_SAS_TOKEN (required)
*/
var DriverAzureBlob Driver = "azure/blob"

/*
DriverGoogleStorage is used to leverage Google Cloud Storage as the destination's
driver.

Environment variables:
  - GOOGLE_APPLICATION_CREDENTIALS (required)
*/
var DriverGoogleStorage Driver = "google/storage"

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

	// Name indicates the identifier of the bucket to use in Blacksmith. The computed
	// name is "blob(<name>)". This does not have any consequences on the bucket's
	// name used in the cloud provider.
	//
	// Example: "mybucket"
	// Required.
	Name string

	// Driver is the driver to leverage for using this destination.
	//
	// Required.
	Driver Driver

	// Connection is the driver's specific connection string to use for writing
	// into the bucket.
	//
	// Format for AWS S3: "<bucket>"
	// Format for Azure Blob Storage: "<container>"
	// Format for Google Cloud Storage: "<bucket>"
	Connection string

	// Params can be used to add specific configuration per driver.
	//
	// Supported fields for AWS S3:
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
		Message:     "blob: Failed to load",
		Validations: []errors.Validation{},
	}

	if env == nil {
		fail.Validations = append(fail.Validations, errors.Validation{
			Message: "Options must not be nil",
			Path:    []string{"Options", "Destinations", "blob"},
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
			Message: "Bucket name must be set",
			Path:    []string{"Options", "Destinations", "blob", "Name"},
		})
	}

	// Create the computed name of the destination.
	name := fmt.Sprintf("blob(%s)", env.Name)

	if env.Driver == "" {
		fail.Validations = append(fail.Validations, errors.Validation{
			Message: "Bucket driver must be set",
			Path:    []string{"Options", "Destinations", name, "Driver"},
		})
	}

	if env.Connection == "" {
		fail.Validations = append(fail.Validations, errors.Validation{
			Message: "Bucket connection must be set",
			Path:    []string{"Options", "Destinations", name, "Connection"},
		})
	}

	switch env.Driver {
	case DriverAWSS3:
		fail.Validations = append(fail.Validations, env.validateDriverAWSS3(name)...)
	case DriverAzureBlob:
		fail.Validations = append(fail.Validations, env.validateDriverAzureBlob(name)...)
	case DriverGoogleStorage:
		fail.Validations = append(fail.Validations, env.validateDriverGoogleStorage(name)...)
	}

	if len(fail.Validations) > 0 {
		return fail
	}

	return nil
}

/*
validateDriverAWSS3 is part of the options' validation process and validate those
for the AWS S3 driver.
*/
func (env *Options) validateDriverAWSS3(name string) []errors.Validation {
	validations := []errors.Validation{}

	// Do not continue if a custom S3-compatible endpoint is set.
	if env.Params.Get("endpoint") != "" {
		return validations
	}

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
	// work with a S3 bucket.
	if region == "" {
		validations = append(validations, errors.Validation{
			Message: "Environment variable 'AWS_REGION' not set",
			Path:    []string{"Options", "Destinations", name},
		})
	}

	return validations
}

/*
validateDriverAzureBlob is part of the options' validation process and validate
those for the Azure Blob driver.
*/
func (env *Options) validateDriverAzureBlob(name string) []errors.Validation {
	validations := []errors.Validation{}

	// Add a validation error if the Azure storage account is not set.
	if os.Getenv("AZURE_STORAGE_ACCOUNT") == "" {
		validations = append(validations, errors.Validation{
			Message: "Environment variable 'AZURE_STORAGE_ACCOUNT' not set",
			Path:    []string{"Options", "Destinations", name},
		})
	}

	// At least ine of the Azure storage key or token must be set.
	if os.Getenv("AZURE_STORAGE_KEY") == "" && os.Getenv("AZURE_STORAGE_SAS_TOKEN") == "" {
		validations = append(validations, errors.Validation{
			Message: "Environment variable 'AZURE_STORAGE_KEY' or 'AZURE_STORAGE_SAS_TOKEN' not set",
			Path:    []string{"Options", "Destinations", name},
		})
	}

	return validations
}

/*
validateDriverGoogleStorage is part of the options' validation process and validate
those for the Google Cloud Storage driver.
*/
func (env *Options) validateDriverGoogleStorage(name string) []errors.Validation {
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
