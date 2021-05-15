package segmentdestination

import (
	"net/http"
	"time"

	"github.com/nunchistudio/blacksmith/flow/destination"
	"github.com/nunchistudio/blacksmith/helper/logger"
)

/*
Segment implements the Blacksmith destination.Destination interface for the
destination "segment".
*/
type Segment struct {
	options *destination.Options
	env     *Options
	client  *http.Client
}

/*
New returns a valid Blacksmith destination.Destination for Segment.
*/
func New(env *Options) destination.Destination {
	if err := env.validate(); err != nil {
		logger.Default.Fatal(err)
		return nil
	}

	return &Segment{
		options: &destination.Options{
			DefaultSchedule: &destination.Schedule{
				Realtime:   env.Realtime,
				Interval:   env.Interval,
				MaxRetries: env.MaxRetries,
			},
			DefaultVersion: "v1.0",
			Versions: map[string]time.Time{
				"v1.0": time.Time{},
			},
		},
		env:    env,
		client: http.DefaultClient,
	}
}

/*
String returns the string representation of the destination Segment.
*/
func (d *Segment) String() string {
	return "segment"
}

/*
Options returns common destination options for Segment. They will be
shared across every actions of this destination, except when overridden.
*/
func (d *Segment) Options() *destination.Options {
	return d.options
}

/*
Actions return a list of actions the destination Segment is able to handle.
*/
func (d *Segment) Actions() map[string]destination.Action {
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
