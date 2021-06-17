package segmentflow

import (
	"github.com/nunchistudio/blacksmith/destination"
	"github.com/nunchistudio/blacksmith/flow"

	"github.com/nunchistudio/blacksmith-modules/amplitude/amplitudedestination"
	"github.com/nunchistudio/blacksmith-modules/segment/segmentdestination"

	"gopkg.in/segmentio/analytics-go.v3"
)

/*
Alias implements the Blacksmith flow.Flow interface for the flow
"alias". It holds a common data structure used by triggers and then
loaded to destinations by actions.
*/
type Alias struct {
	analytics.Alias
}

/*
Options returns the flow options. When disabled, a flow will not be
executed. Therefore no jobs will be created.
*/
func (f *Alias) Options() *flow.Options {
	return &flow.Options{
		Enabled: true,
	}
}

/*
Transform is the function being run by when executing the flow from
triggers. It is up to the flow to transform the data from sources'
triggers to destinations' actions.
*/
func (f *Alias) Transform(tk *flow.Toolkit) destination.Actions {
	integrations := map[string][]destination.Action{}
	var exists bool

	if f.Alias.Context == nil {
		f.Alias.Context = &analytics.Context{}
	}

	_, exists = f.Alias.Integrations["Amplitude"]
	if !exists || f.Alias.Integrations["Amplitude"] == true {
		integrations["amplitude"] = []destination.Action{
			amplitudedestination.Alias{
				Mapping: []amplitudedestination.UserMap{
					{
						UserId:       f.Alias.PreviousId,
						GlobalUserId: f.Alias.UserId,
					},
				},
			},
		}
	}

	_, exists = f.Alias.Integrations["Segment"]
	if !exists || f.Alias.Integrations["Segment"] == true {
		integrations["segment"] = []destination.Action{
			segmentdestination.Alias{
				Alias: f.Alias,
			},
		}
	}

	return integrations
}
