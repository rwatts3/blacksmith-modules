package amplitudedestination

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/nunchistudio/blacksmith/adapter/store"
	"github.com/nunchistudio/blacksmith/flow/destination"
	"github.com/nunchistudio/blacksmith/helper/errors"

	"github.com/segmentio/ksuid"
	"gopkg.in/segmentio/analytics-go.v3"
)

/*
Page implements the Blacksmith destination.Action interface for the action
"page". It holds the complete job's structure to load into the destination.
*/
type Page struct {
	env    *Options
	client *http.Client

	APIKey string  `json:"api_key,omitempty"`
	Events []Event `json:"events"`
}

/*
String returns the string representation of the action Page.
*/
func (a Page) String() string {
	return "page"
}

/*
Schedule allows the action to override the schedule options of its
destination. Do not override.
*/
func (a Page) Schedule() *destination.Schedule {
	return nil
}

/*
Marshal is the function being run when the action receives data into
the Page receiver. It allows to transform and enrich the data before
saving it in the store adapter.
*/
func (a Page) Marshal(tk *destination.Toolkit) (*destination.Job, error) {

	// Create a unique insert identifier to ensure data is not duplicated in
	// Amplitude in case of a 500, 502, or 504 HTTP error.
	insertId := ksuid.New().String()
	for i := range a.Events {
		a.Events[i].InsertId = insertId
		a.Events[i].Event = "Viewed page '" + a.Events[i].Event + "'"
	}

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
func (a Page) Load(tk *destination.Toolkit, queue *store.Queue, then chan<- destination.Then) {

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
			var d Page
			json.Unmarshal(job.Data, &d)

			// Apply the context as the "Event Properties".
			for i := range d.Events {
				d.Events[i].Context = &c
			}

			// Create the request body for Amplitude with the Amplitude API Key.
			// We add the Amplitude API key just before loading the data so it is
			// not saved in the store adapter.
			body := d
			body.APIKey = a.env.APIKey

			// Marshal and create a reader for making the HTTP request.
			b, _ := json.Marshal(&body)
			r := bytes.NewReader(b)

			// Run the HTTP request against the Amplitude endpoint using the API Key.
			// Inform the scheduler if any error happened.
			req, _ := http.NewRequest("POST", "https://api2.amplitude.com/2/httpapi", r)
			req.Header.Set("Accept", "*/*")
			req.Header.Set("Content-Type", "application/json")
			res, err := a.client.Do(req)
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

			// Finally, inform the scheduler about the success.
			then <- destination.Then{
				Jobs:  []string{job.ID},
				Error: nil,
			}
		}
	}
}
