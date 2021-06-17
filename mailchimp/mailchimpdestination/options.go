package mailchimpdestination

import (
	"github.com/nunchistudio/blacksmith/destination"
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

	// APIKey is the Mailchimp API Key to use for loading data into Mailchimp. You
	// can create and copy-paste your Mailchimp API Key from 'Account Settings' >
	// 'Extras' > 'API Keys'.
	//
	// Required.
	APIKey string

	// DatacenterID is the datacenter identifier of your Mailchimp account. You
	// can find it in the Mailchimp URL in your browser when you are logged in.
	// It is the 'us1' in 'https://us1.admin.mailchimp.com/lists/'.
	//
	// Required.
	DatacenterID string

	// AudienceID is the audience identifier to connect to. You can find your
	// Audience ID in your Mailchimp Settings pane under the Audiences tab. Go
	// to 'Manage Audiences' > 'Settings' and click on 'Audience Name & Defaults'.
	// The Audience ID will be on the right side.
	//
	// Required.
	AudienceID string

	// EnableDoubleOptIn is an optional flag to control whether a double opt-in
	// confirmation message is sent when subscribing new users. When enabled,
	// the status of a new subscriber will be set to 'pending' until the
	// subscription has been confirmed by the subscriber (via email). When
	// disabled, the subscription will automatically be set to 'subscribed'
	// and no email confirmation is sent to the subscriber.
	//
	// Defaults to false.
	EnableDoubleOptIn bool
}

/*
validate ensures the options passed to initialize the destination are valid.
*/
func (env *Options) validate() error {
	var interval string = destination.Defaults.DefaultSchedule.Interval
	var maxRetries uint16 = destination.Defaults.DefaultSchedule.MaxRetries

	fail := &errors.Error{
		Message:     "destination/mailchimp: Failed to load",
		Validations: []errors.Validation{},
	}

	if env == nil {
		fail.Validations = append(fail.Validations, errors.Validation{
			Message: "Options must not be nil",
			Path:    []string{"Options", "Destinations", "mailchimp"},
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
			Message: "Mailchimp API key must be set",
			Path:    []string{"Options", "Destinations", "mailchimp", "APIKey"},
		})
	}

	if env.DatacenterID == "" {
		fail.Validations = append(fail.Validations, errors.Validation{
			Message: "Mailchimp datacenter must be set",
			Path:    []string{"Options", "Destinations", "mailchimp", "DatacenterID"},
		})
	}

	if env.AudienceID == "" {
		fail.Validations = append(fail.Validations, errors.Validation{
			Message: "Mailchimp audience must be set",
			Path:    []string{"Options", "Destinations", "mailchimp", "AudienceID"},
		})
	}

	if len(fail.Validations) > 0 {
		return fail
	}

	return nil
}
