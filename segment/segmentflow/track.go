package segmentflow

import (
	"github.com/nunchistudio/blacksmith/flow"
	"github.com/nunchistudio/blacksmith/flow/destination"

	"github.com/nunchistudio/blacksmith-modules/amplitude/amplitudedestination"
	"github.com/nunchistudio/blacksmith-modules/segment/segmentdestination"

	"gopkg.in/segmentio/analytics-go.v3"
)

/*
Track implements the Blacksmith flow.Flow interface for the flow
"track". It holds a common data structure used by triggers and then
loaded to destinations by actions.
*/
type Track struct {
	analytics.Track
}

/*
Options returns the flow options. When disabled, a flow will not be
executed. Therefore no jobs will be created.
*/
func (f *Track) Options() *flow.Options {
	return &flow.Options{
		Enabled: true,
	}
}

/*
Transform is the function being run by when executing the flow from
triggers. It is up to the flow to transform the data from sources'
triggers to destinations' actions.
*/
func (f *Track) Transform(tk *flow.Toolkit) destination.Actions {
	integrations := map[string][]destination.Action{}
	var exists bool

	if f.Track.Context == nil {
		f.Track.Context = &analytics.Context{}
	}

	var productId string
	if _, exists := f.Track.Properties["productId"]; exists {
		got, ok := f.Track.Properties["productId"].(string)
		if ok {
			productId = got
		}
	}

	var quantity float64
	if _, exists := f.Track.Properties["quantity"]; exists {
		got, ok := f.Track.Properties["quantity"].(float64)
		if ok {
			quantity = got
		}
	}

	var price float64
	if _, exists := f.Track.Properties["price"]; exists {
		got, ok := f.Track.Properties["price"].(float64)
		if ok {
			price = got
		}
	}

	var revenue float64
	if _, exists := f.Track.Properties["revenue"]; exists {
		got, ok := f.Track.Properties["revenue"].(float64)
		if ok {
			revenue = got
		}
	}

	_, exists = f.Track.Integrations["Amplitude"]
	if !exists || f.Track.Integrations["Amplitude"] == true {
		integrations["amplitude"] = []destination.Action{
			amplitudedestination.Track{
				Events: []amplitudedestination.Event{
					{
						Event:              f.Track.Event,
						UserId:             f.Track.UserId,
						DeviceId:           f.Track.Context.Device.Id,
						Time:               f.Track.Timestamp.Unix(),
						Context:            f.Track.Context,
						AppVersion:         f.Track.Context.App.Version,
						OSName:             f.Track.Context.OS.Name,
						OSVersion:          f.Track.Context.OS.Version,
						DeviceBrand:        f.Track.Context.Device.Name,
						DeviceManufacturer: f.Track.Context.Device.Manufacturer,
						DeviceModel:        f.Track.Context.Device.Model,
						Carrier:            f.Track.Context.Network.Carrier,
						Country:            f.Track.Context.Location.Country,
						Region:             f.Track.Context.Location.Region,
						City:               f.Track.Context.Location.City,
						Latitude:           f.Track.Context.Location.Latitude,
						Longitude:          f.Track.Context.Location.Longitude,
						IP:                 f.Track.Context.IP,
						ProductId:          productId,
						Quantity:           quantity,
						Price:              price,
						Revenue:            revenue,
					},
				},
			},
		}
	}

	_, exists = f.Track.Integrations["Segment"]
	if !exists || f.Track.Integrations["Segment"] == true {
		integrations["segment"] = []destination.Action{
			segmentdestination.Track{
				Track: f.Track,
			},
		}
	}

	return integrations
}
