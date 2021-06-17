package mailchimpdestination

import (
	"net/http"
	"time"

	"github.com/nunchistudio/blacksmith/destination"
	"github.com/nunchistudio/blacksmith/helper/logger"
)

/*
statusForceDiscards holds informations if a job must be discarded based on
the status code returned by the Mailchimp API.

Reference: https://mailchimp.com/developer/marketing/docs/errors/
*/
var statusForceDiscards = map[int]bool{
	400: true,
	401: false,
	403: false,
	404: false,
	405: true,
	414: true,
	422: true,
	429: false,
	500: false,
	503: false,
}

/*
Mailchimp implements the Blacksmith destination.Destination interface for the
destination "mailchimp".
*/
type Mailchimp struct {
	options *destination.Options
	env     *Options
	client  *http.Client
}

/*
New returns a valid Blacksmith destination.Destination for Mailchimp.
*/
func New(env *Options) destination.Destination {
	if err := env.validate(); err != nil {
		logger.Default.Fatal(err)
		return nil
	}

	return &Mailchimp{
		options: &destination.Options{
			DefaultSchedule: &destination.Schedule{
				Realtime:   env.Realtime,
				Interval:   env.Interval,
				MaxRetries: env.MaxRetries,
			},
			DefaultVersion: "v3.0",
			Versions: map[string]time.Time{
				"v3.0": time.Time{},
			},
		},
		env:    env,
		client: http.DefaultClient,
	}
}

/*
String returns the string representation of the destination Mailchimp.
*/
func (d *Mailchimp) String() string {
	return "mailchimp"
}

/*
Options returns common destination options for Mailchimp. They will be
shared across every actions of this destination, except when overridden.
*/
func (d *Mailchimp) Options() *destination.Options {
	return d.options
}

/*
Actions return a list of actions the destination Mailchimp is able to handle.
*/
func (d *Mailchimp) Actions() map[string]destination.Action {
	return map[string]destination.Action{
		"identify": Identify{
			env:    d.env,
			client: d.client,
		},
	}
}
