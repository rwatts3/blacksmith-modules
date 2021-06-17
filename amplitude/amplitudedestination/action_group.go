package amplitudedestination

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/nunchistudio/blacksmith/adapter/store"
	"github.com/nunchistudio/blacksmith/destination"
	"github.com/nunchistudio/blacksmith/helper/errors"

	"gopkg.in/segmentio/analytics-go.v3"
)

/*
Group implements the Blacksmith destination.Action interface for the action
"group". It holds the complete job's structure to load into the destination.
*/
type Group struct {
	env    *Options
	client *http.Client

	APIKey         string           `json:"api_key,omitempty"`
	Identification []Identification `json:"identification"`
}

/*
String returns the string representation of the action Group.
*/
func (a Group) String() string {
	return "group"
}

/*
Schedule allows the action to override the schedule options of its
destination. Do not override.
*/
func (a Group) Schedule() *destination.Schedule {
	return nil
}

/*
Marshal is the function being run when the action receives data into
the Group receiver. It allows to transform and enrich the data before
saving it in the store adapter.
*/
func (a Group) Marshal(tk *destination.Toolkit) (*destination.Job, error) {

	// Try to marshal the data passed directly to the receiver.
	data, err := json.Marshal(&a)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: 400,
			Message:    "Bad Request",
		}
	}

	// Create a job with the data. Since the 'Context' key is not
	// set, the one from the event will automatically be applied.
	j := &destination.Job{
		Version: "v2.0",
		Data:    data,
	}

	// Return the job including the marshaled data.
	return j, nil
}

/*
Load is the function being run by the scheduler to load the data into
the destination. It is in charge of the "L" in the ETL process.
*/
func (a Group) Load(tk *destination.Toolkit, queue *store.Queue, then chan<- destination.Then) {

	// We can go through every events received from the queue and their
	// related jobs. The queue can contain one or many events. The jobs
	// present in the events are specific to this action only.
	//
	// This allows to parse everything needed and make a request to the
	// destination for each event / job.
	for _, event := range queue.Events {
		for _, job := range event.Jobs {

			// Unmarshal the `context` key of the job.
			var c analytics.Context
			json.Unmarshal(job.Context, &c)

			// Unmarshal the `data` key of the job.
			var d Group
			json.Unmarshal(job.Data, &d)

			// Create the request query params for this Amplitude endpoint dedicated
			// to user grouping. We add the Amplitude API key just before loading
			// the data so it is not saved in the store adapter.
			b, _ := json.Marshal(d.Identification)
			data := url.Values{
				"api_key":        {a.env.APIKey},
				"identification": {string(b[:])},
			}

			// Run the HTTP request against the Amplitude endpoint using the API Key.
			// Inform the scheduler if any error happened.
			res, err := a.client.PostForm("https://api.amplitude.com/groupidentify", data)
			if err != nil {
				then <- destination.Then{
					Jobs:         []string{job.ID},
					ForceDiscard: true,
					Error: &errors.Error{
						StatusCode: 500,
						Message:    err.Error(),
					},
				}

				continue
			}

			// Since a non-2xx status code doesn't cause an error, catch HTTP status
			// code to ensure nothing bad happened.
			//
			// TODO: Handle status code not registered in the map.
			if res.StatusCode >= 300 {
				buf := new(bytes.Buffer)
				buf.ReadFrom(res.Body)
				then <- destination.Then{
					Jobs:         []string{job.ID},
					ForceDiscard: statusForceDiscards[res.StatusCode],
					Error: &errors.Error{
						StatusCode: res.StatusCode,
						Message:    buf.String(),
					},
				}

				continue
			}

			// Finally, inform the scheduler about the success. We also register
			// a "Track" event for the "Group" event so it appears on the Amplitude
			// history.
			then <- destination.Then{
				Jobs:  []string{job.ID},
				Error: nil,
				OnSucceeded: []destination.Action{
					Track{
						Events: []Event{
							{
								Event:              "Group",
								UserId:             d.Identification[0].UserId,
								DeviceId:           c.Device.Id,
								Time:               event.ReceivedAt.Unix(),
								Context:            &c,
								AppVersion:         c.App.Version,
								OSName:             c.OS.Name,
								OSVersion:          c.OS.Version,
								DeviceBrand:        c.Device.Name,
								DeviceManufacturer: c.Device.Manufacturer,
								DeviceModel:        c.Device.Model,
								Carrier:            c.Network.Carrier,
								Country:            c.Location.Country,
								Region:             c.Location.Region,
								City:               c.Location.City,
								Latitude:           c.Location.Latitude,
								Longitude:          c.Location.Longitude,
								IP:                 c.IP,
							},
						},
					},
				},
			}
		}
	}
}
