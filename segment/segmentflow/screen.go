package segmentflow

import (
	"github.com/nunchistudio/blacksmith/flow"
	"github.com/nunchistudio/blacksmith/flow/destination"

	"github.com/nunchistudio/blacksmith-modules/amplitude/amplitudedestination"
	"github.com/nunchistudio/blacksmith-modules/segment/segmentdestination"

	"gopkg.in/segmentio/analytics-go.v3"
)

/*
Screen implements the Blacksmith flow.Flow interface for the flow
"screen". It holds a common data structure used by triggers and then
loaded to destinations by actions.
*/
type Screen struct {
	analytics.Screen
}

/*
Options returns the flow options. When disabled, a flow will not be
executed. Therefore no jobs will be created.
*/
func (f *Screen) Options() *flow.Options {
	return &flow.Options{
		Enabled: true,
	}
}

/*
Transform is the function being run by when executing the flow from
triggers. It is up to the flow to transform the data from sources'
triggers to destinations' actions.
*/
func (f *Screen) Transform(tk *flow.Toolkit) destination.Actions {
	integrations := map[string][]destination.Action{}
	var exists bool

	if f.Screen.Context == nil {
		f.Screen.Context = &analytics.Context{}
	}

	_, exists = f.Screen.Integrations["Amplitude"]
	if !exists || f.Screen.Integrations["Amplitude"] == true {
		integrations["amplitude"] = []destination.Action{
			amplitudedestination.Screen{
				Events: []amplitudedestination.Event{
					{
						Event:              f.Screen.Name,
						UserId:             f.Screen.UserId,
						DeviceId:           f.Screen.Context.Device.Id,
						Time:               f.Screen.Timestamp.Unix(),
						Context:            f.Screen.Context,
						AppVersion:         f.Screen.Context.App.Version,
						OSName:             f.Screen.Context.OS.Name,
						OSVersion:          f.Screen.Context.OS.Version,
						DeviceBrand:        f.Screen.Context.Device.Name,
						DeviceManufacturer: f.Screen.Context.Device.Manufacturer,
						DeviceModel:        f.Screen.Context.Device.Model,
						Carrier:            f.Screen.Context.Network.Carrier,
						Country:            f.Screen.Context.Location.Country,
						Region:             f.Screen.Context.Location.Region,
						City:               f.Screen.Context.Location.City,
						Latitude:           f.Screen.Context.Location.Latitude,
						Longitude:          f.Screen.Context.Location.Longitude,
						IP:                 f.Screen.Context.IP,
					},
				},
			},
		}
	}

	_, exists = f.Screen.Integrations["Segment"]
	if !exists || f.Screen.Integrations["Segment"] == true {
		integrations["segment"] = []destination.Action{
			segmentdestination.Screen{
				Screen: f.Screen,
			},
		}
	}

	return integrations
}
