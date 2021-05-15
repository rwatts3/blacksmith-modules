package mailchimpdestination

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/nunchistudio/blacksmith/adapter/store"
	"github.com/nunchistudio/blacksmith/flow/destination"
	"github.com/nunchistudio/blacksmith/helper/errors"

	"gopkg.in/segmentio/analytics-go.v3"
)

/*
Identify implements the Blacksmith destination.Action interface for the action
"identify". It holds the complete job's structure to load into the destination.
*/
type Identify struct {
	env    *Options
	client *http.Client

	Signup
}

/*
String returns the string representation of the action Identify.
*/
func (a Identify) String() string {
	return "identify"
}

/*
Schedule allows the action to override the schedule options of its
destination. Do not override.
*/
func (a Identify) Schedule() *destination.Schedule {
	return nil
}

/*
Marshal is the function being run when the action receives data into
the Identify receiver. It allows to transform and enrich the data before
saving it in the store adapter.
*/
func (a Identify) Marshal(tk *destination.Toolkit) (*destination.Job, error) {

	// Ensure the email address is lowercase. As described in the Mailchimp
	// documentation, the subscriber hash is "the lowercase version of the list
	// member's email address".
	a.Email = strings.ToLower(a.Email)

	// Format the date to ISO 8601 as required by the Mailchimp API.
	formatted := a.TimestampSignup.Format(time.RFC3339)
	a.TimestampSignup, _ = time.Parse(time.RFC3339, formatted)

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
		Version: "v3.0",
		Data:    data,
	}

	// Return the job including the marshaled data.
	return j, nil
}

/*
Load is the function being run by the scheduler to load the data into
the destination. It is in charge of the "L" in the ETL process.
*/
func (a Identify) Load(tk *destination.Toolkit, queue *store.Queue, then chan<- destination.Then) {

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
			var d Identify
			json.Unmarshal(job.Data, &d)

			// Set the appropriate subscription status based on the double opt-in
			// option.
			var status string = "subscribed"
			if a.env.EnableDoubleOptIn == true {
				status = "pending"
			}

			// Create the request body for Mailchimp.
			body := d
			body.StatusIfNew = status

			// Marshal and create a reader for making the HTTP request.
			b, _ := json.Marshal(&body)
			r := bytes.NewReader(b)

			// Create the subscriber hash based on the email address.
			h := md5.Sum([]byte(d.Email))
			s := hex.EncodeToString(h[:])

			// Run the HTTP request against the Mailchimp endpoint using the API Key.
			// Inform the scheduler if any error happened.
			req, _ := http.NewRequest("PUT", "https://"+a.env.DatacenterID+".api.mailchimp.com/3.0/lists/"+a.env.AudienceID+"/members/"+s, r)
			req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("anystring:"+a.env.APIKey)))
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
