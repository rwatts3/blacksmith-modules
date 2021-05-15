package amplitudedestination

import (
	"github.com/nunchistudio/blacksmith/flow/destination"
	"github.com/nunchistudio/blacksmith/helper/errors"
)

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

	// APIKey is the Amplitude API Key to use for loading data into Amplitude. You
	// can copy-paste your Amplitude API Key from the Amplitude Settings page.
	//
	// Required.
	APIKey string
}

/*
validate ensures the options passed to initialize the destination are valid.
*/
func (env *Options) validate() error {
	var interval string = destination.Defaults.DefaultSchedule.Interval
	var maxRetries uint16 = destination.Defaults.DefaultSchedule.MaxRetries

	fail := &errors.Error{
		Message:     "destination/amplitude: Failed to load",
		Validations: []errors.Validation{},
	}

	if env == nil {
		fail.Validations = append(fail.Validations, errors.Validation{
			Message: "Options must not be nil",
			Path:    []string{"Options", "Destinations", "amplitude"},
		})

		return fail
	}

	if env.Interval == "" {
		env.Interval = interval
	}

	if env.MaxRetries == 0 {
		env.MaxRetries = maxRetries
	}

	if env.APIKey == "" {
		fail.Validations = append(fail.Validations, errors.Validation{
			Message: "Amplitude API key must be set",
			Path:    []string{"Options", "Destinations", "amplitude", "APIKey"},
		})
	}

	if len(fail.Validations) > 0 {
		return fail
	}

	return nil
}
