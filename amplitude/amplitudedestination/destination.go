package amplitudedestination

import (
	"net/http"
	"time"

	"github.com/nunchistudio/blacksmith/flow/destination"
	"github.com/nunchistudio/blacksmith/helper/logger"
)

/*
statusForceDiscards holds informations if a job must be discarded based on
the status code returned by the Amplitude API.

Reference: https://developers.amplitude.com/docs/http-api-v2#response-format
*/
var statusForceDiscards = map[int]bool{
	400: true,
	413: true,
	422: true,
	429: false,
	500: false,
	502: false,
	503: false,
	504: false,
}

/*
Amplitude implements the Blacksmith destination.Destination interface for the
destination "amplitude".
*/
type Amplitude struct {
	options *destination.Options
	env     *Options
	client  *http.Client
}

/*
New returns a valid Blacksmith destination.Destination for Amplitude.
*/
func New(env *Options) destination.Destination {
	if err := env.validate(); err != nil {
		logger.Default.Fatal(err)
		return nil
	}

	return &Amplitude{
		options: &destination.Options{
			DefaultSchedule: &destination.Schedule{
				Realtime:   env.Realtime,
				Interval:   env.Interval,
				MaxRetries: env.MaxRetries,
			},
			DefaultVersion: "v2.0",
			Versions: map[string]time.Time{
				"v2.0": time.Time{},
			},
		},
		env:    env,
		client: http.DefaultClient,
	}
}

/*
String returns the string representation of the destination Amplitude.
*/
func (d *Amplitude) String() string {
	return "amplitude"
}

/*
Options returns common destination options for Amplitude. They will be
shared across every actions of this destination, except when overridden.
*/
func (d *Amplitude) Options() *destination.Options {
	return d.options
}

/*
Actions return a list of actions the destination Amplitude is able to handle.
*/
func (d *Amplitude) Actions() map[string]destination.Action {
	return map[string]destination.Action{
		"identify": Identify{
			env:    d.env,
			client: d.client,
		},
		"track": Track{
			env:    d.env,
			client: d.client,
		},
		"group": Group{
			env:    d.env,
			client: d.client,
		},
		"alias": Alias{
			env:    d.env,
			client: d.client,
		},
		"page": Page{
			env:    d.env,
			client: d.client,
		},
		"screen": Screen{
			env:    d.env,
			client: d.client,
		},
	}
}
