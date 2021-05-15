package segmentflow

import (
	"github.com/nunchistudio/blacksmith/flow"
	"github.com/nunchistudio/blacksmith/flow/destination"

	"github.com/nunchistudio/blacksmith-modules/amplitude/amplitudedestination"
	"github.com/nunchistudio/blacksmith-modules/mailchimp/mailchimpdestination"
	"github.com/nunchistudio/blacksmith-modules/segment/segmentdestination"

	"gopkg.in/segmentio/analytics-go.v3"
)

/*
Identify implements the Blacksmith flow.Flow interface for the flow
"identify". It holds a common data structure used by triggers and then
loaded to destinations by actions.
*/
type Identify struct {
	analytics.Identify
}

/*
Options returns the flow options. When disabled, a flow will not be
executed. Therefore no jobs will be created.
*/
func (f *Identify) Options() *flow.Options {
	return &flow.Options{
		Enabled: true,
	}
}

/*
Transform is the function being run by when executing the flow from
triggers. It is up to the flow to transform the data from sources'
triggers to destinations' actions.
*/
func (f *Identify) Transform(tk *flow.Toolkit) destination.Actions {
	integrations := map[string][]destination.Action{}
	var exists bool

	if f.Identify.Context == nil {
		f.Identify.Context = &analytics.Context{}
	}

	var email string
	if _, exists := f.Identify.Traits["email"]; exists {
		got, ok := f.Identify.Traits["email"].(string)
		if ok {
			email = got
		}
	}

	var firstName string
	if _, exists := f.Identify.Traits["firstName"]; exists {
		got, ok := f.Identify.Traits["firstName"].(string)
		if ok {
			firstName = got
		}
	}

	var lastName string
	if _, exists := f.Identify.Traits["lastName"]; exists {
		got, ok := f.Identify.Traits["lastName"].(string)
		if ok {
			lastName = got
		}
	}

	_, exists = f.Identify.Integrations["Amplitude"]
	if !exists || f.Identify.Integrations["Amplitude"] == true {
		integrations["amplitude"] = []destination.Action{
			amplitudedestination.Identify{
				Events: []amplitudedestination.Event{
					{
						UserId:             f.Identify.UserId,
						DeviceId:           f.Identify.Context.Device.Id,
						Time:               f.Identify.Timestamp.Unix(),
						Traits:             f.Identify.Traits,
						Context:            f.Identify.Context,
						AppVersion:         f.Identify.Context.App.Version,
						OSName:             f.Identify.Context.OS.Name,
						OSVersion:          f.Identify.Context.OS.Version,
						DeviceBrand:        f.Identify.Context.Device.Name,
						DeviceManufacturer: f.Identify.Context.Device.Manufacturer,
						DeviceModel:        f.Identify.Context.Device.Model,
						Carrier:            f.Identify.Context.Network.Carrier,
						Country:            f.Identify.Context.Location.Country,
						Region:             f.Identify.Context.Location.Region,
						City:               f.Identify.Context.Location.City,
						Latitude:           f.Identify.Context.Location.Latitude,
						Longitude:          f.Identify.Context.Location.Longitude,
						IP:                 f.Identify.Context.IP,
					},
				},
			},
		}
	}

	_, exists = f.Identify.Integrations["Mailchimp"]
	if !exists || f.Identify.Integrations["Mailchimp"] == true {
		integrations["mailchimp"] = []destination.Action{
			mailchimpdestination.Identify{
				Signup: mailchimpdestination.Signup{
					Email:           email,
					FirstName:       firstName,
					LastName:        lastName,
					IPSignup:        f.Identify.Context.IP,
					TimestampSignup: f.Identify.Timestamp,
					Location: &analytics.LocationInfo{
						Latitude:  f.Identify.Context.Location.Latitude,
						Longitude: f.Identify.Context.Location.Longitude,
					},
				},
			},
		}
	}

	_, exists = f.Identify.Integrations["Segment"]
	if !exists || f.Identify.Integrations["Segment"] == true {
		integrations["segment"] = []destination.Action{
			segmentdestination.Identify{
				Identify: f.Identify,
			},
		}
	}

	return integrations
}
