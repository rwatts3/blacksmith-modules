package segmentdestination

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/nunchistudio/blacksmith/adapter/store"
	"github.com/nunchistudio/blacksmith/flow/destination"
	"github.com/nunchistudio/blacksmith/helper/errors"

	"gopkg.in/segmentio/analytics-go.v3"
)

/*
Alias implements the Blacksmith destination.Action interface for the action
"alias". It holds the complete job's structure to load into the destination.
*/
type Alias struct {
	env    *Options
	client *http.Client

	analytics.Alias
}

/*
String returns the string representation of the action Alias.
*/
func (a Alias) String() string {
	return "alias"
}

/*
Schedule allows the action to override the schedule options of its
destination. Do not override.
*/
func (a Alias) Schedule() *destination.Schedule {
	return nil
}

/*
Marshal is the function being run when the action receives data into
the Alias receiver. It allows to transform and enrich the data before
saving it in the store adapter.
*/
func (a Alias) Marshal(tk *destination.Toolkit) (*destination.Job, error) {

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
		Version: "v1.0",
		Data:    data,
		SentAt:  &a.Timestamp,
	}

	// Return the job including the marshaled data.
	return j, nil
}

/*
Load is the function being run by the scheduler to load the data into
the destination. It is in charge of the "L" in the ETL process.
*/
func (a Alias) Load(tk *destination.Toolkit, queue *store.Queue, then chan<- destination.Then) {

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
			var d analytics.Alias
			json.Unmarshal(job.Data, &d)

			// Create the request body for Segment.
			body := d
			body.Context = &c

			// Marshal and create a reader for making the HTTP request.
			b, _ := json.Marshal(&body)
			r := bytes.NewReader(b)

			// Run the HTTP request against the Segment endpoint using the API Write
			// Key. Inform the scheduler if any error happened.
			req, _ := http.NewRequest("POST", "https://api.segment.io/v1/alias", r)
			req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(a.env.WriteKey)))
			req.Header.Set("Content-Type", "application/json")
			res, err := a.client.Do(req)
			if err != nil {
				then <- destination.Then{
					Jobs: []string{job.ID},
					Error: &errors.Error{
						StatusCode: 500,
						Message:    err.Error(),
					},
				}

				continue
			}

			// Since a non-2xx status code doesn't cause an error, catch HTTP status
			// code to ensure nothing bad happened. As described in the Segment
			// documentation, the API return a status code 200 for all requests.
			// The only exception is if the request is too large or if the JSON is
			// invalid it will respond with a 400.
			if res.StatusCode >= 300 {
				buf := new(bytes.Buffer)
				buf.ReadFrom(res.Body)
				then <- destination.Then{
					Jobs:         []string{job.ID},
					ForceDiscard: res.StatusCode == 400,
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
