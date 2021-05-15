package segmentflow

import (
	"github.com/nunchistudio/blacksmith/flow"
	"github.com/nunchistudio/blacksmith/flow/destination"

	"github.com/nunchistudio/blacksmith-modules/amplitude/amplitudedestination"
	"github.com/nunchistudio/blacksmith-modules/segment/segmentdestination"

	"gopkg.in/segmentio/analytics-go.v3"
)

/*
Page implements the Blacksmith flow.Flow interface for the flow
"page". It holds a common data structure used by triggers and then
loaded to destinations by actions.
*/
type Page struct {
	analytics.Page
}

/*
Options returns the flow options. When disabled, a flow will not be
executed. Therefore no jobs will be created.
*/
func (f *Page) Options() *flow.Options {
	return &flow.Options{
		Enabled: true,
	}
}

/*
Transform is the function being run by when executing the flow from
triggers. It is up to the flow to transform the data from sources'
triggers to destinations' actions.
*/
func (f *Page) Transform(tk *flow.Toolkit) destination.Actions {
	integrations := map[string][]destination.Action{}
	var exists bool

	if f.Page.Context == nil {
		f.Page.Context = &analytics.Context{}
	}

	_, exists = f.Page.Integrations["Amplitude"]
	if !exists || f.Page.Integrations["Amplitude"] == true {
		integrations["amplitude"] = []destination.Action{
			amplitudedestination.Page{
				Events: []amplitudedestination.Event{
					{
						Event:              f.Page.Name,
						UserId:             f.Page.UserId,
						DeviceId:           f.Page.Context.Device.Id,
						Time:               f.Page.Timestamp.Unix(),
						Context:            f.Page.Context,
						AppVersion:         f.Page.Context.App.Version,
						OSName:             f.Page.Context.OS.Name,
						OSVersion:          f.Page.Context.OS.Version,
						DeviceBrand:        f.Page.Context.Device.Name,
						DeviceManufacturer: f.Page.Context.Device.Manufacturer,
						DeviceModel:        f.Page.Context.Device.Model,
						Carrier:            f.Page.Context.Network.Carrier,
						Country:            f.Page.Context.Location.Country,
						Region:             f.Page.Context.Location.Region,
						City:               f.Page.Context.Location.City,
						Latitude:           f.Page.Context.Location.Latitude,
						Longitude:          f.Page.Context.Location.Longitude,
						IP:                 f.Page.Context.IP,
					},
				},
			},
		}
	}

	_, exists = f.Page.Integrations["Segment"]
	if !exists || f.Page.Integrations["Segment"] == true {
		integrations["segment"] = []destination.Action{
			segmentdestination.Page{
				Page: f.Page,
			},
		}
	}

	return integrations
}
