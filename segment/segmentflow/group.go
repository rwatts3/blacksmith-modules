package segmentflow

import (
	"github.com/nunchistudio/blacksmith/destination"
	"github.com/nunchistudio/blacksmith/flow"

	"github.com/nunchistudio/blacksmith-modules/amplitude/amplitudedestination"
	"github.com/nunchistudio/blacksmith-modules/segment/segmentdestination"

	"gopkg.in/segmentio/analytics-go.v3"
)

/*
Group implements the Blacksmith flow.Flow interface for the flow
"group". It holds a common data structure used by triggers and then
loaded to destinations by actions.
*/
type Group struct {
	analytics.Group
}

/*
Options returns the flow options. When disabled, a flow will not be
executed. Therefore no jobs will be created.
*/
func (f *Group) Options() *flow.Options {
	return &flow.Options{
		Enabled: true,
	}
}

/*
Transform is the function being run by when executing the flow from
triggers. It is up to the flow to transform the data from sources'
triggers to destinations' actions.
*/
func (f *Group) Transform(tk *flow.Toolkit) destination.Actions {
	integrations := map[string][]destination.Action{}
	var exists bool

	if f.Group.Context == nil {
		f.Group.Context = &analytics.Context{}
	}

	var groupType string
	if _, exists := f.Group.Traits["industry"]; exists {
		got, ok := f.Group.Traits["industry"].(string)
		if ok {
			groupType = got
		}
	}

	_, exists = f.Group.Integrations["Amplitude"]
	if !exists || f.Group.Integrations["Amplitude"] == true {
		integrations["amplitude"] = []destination.Action{
			amplitudedestination.Group{
				Identification: []amplitudedestination.Identification{
					{
						UserId:          f.Group.UserId,
						GroupType:       groupType,
						GroupValue:      f.Group.GroupId,
						GroupProperties: f.Group.Traits,
					},
				},
			},
		}
	}

	_, exists = f.Group.Integrations["Segment"]
	if !exists || f.Group.Integrations["Segment"] == true {
		integrations["segment"] = []destination.Action{
			segmentdestination.Group{
				Group: f.Group,
			},
		}
	}

	return integrations
}
