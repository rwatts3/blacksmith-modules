package blobdestination

import (
	"context"
	"encoding/json"

	"github.com/nunchistudio/blacksmith/adapter/store"
	"github.com/nunchistudio/blacksmith/destination"
	"github.com/nunchistudio/blacksmith/helper/errors"

	"gocloud.dev/blob"
)

/*
Write implements the Blacksmith destination.Action interface for the action
"write". It holds the complete job's structure to load into the destination.
*/
type Write struct {
	env    *Options
	ctx    context.Context
	bucket *blob.Bucket

	Filename string `json:"filename"`
	Content  []byte `json:"content"`
}

/*
String returns the string representation of the action Write.
*/
func (a Write) String() string {
	return "write"
}

/*
Schedule allows the action to override the schedule options of its
destination. Do not override.
*/
func (a Write) Schedule() *destination.Schedule {
	return nil
}

/*
Marshal is the function being run when the action receives data into
the Write receiver. It allows to transform and enrich the data before
saving it in the store adapter.
*/
func (a Write) Marshal(tk *destination.Toolkit) (*destination.Job, error) {

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
		Data: data,
	}

	// Return the job including the marshaled data.
	return j, nil
}

/*
Load is the function being run by the scheduler to load the data into
the destination. It is in charge of the "L" in the ETL process.
*/
func (a Write) Load(tk *destination.Toolkit, queue *store.Queue, then chan<- destination.Then) {

	// We can go through every events received from the queue and their
	// related jobs. The queue can contain one or many events. The jobs
	// present in the events are specific to this action only.
	//
	// This allows to parse everything needed and make a request to the
	// destination for each event / job.
	for _, event := range queue.Events {
		for _, job := range event.Jobs {

			// Unmarshal the `data` key of the job.
			var payload Write
			err := json.Unmarshal(job.Data, &payload)
			if err != nil {
				then <- destination.Then{
					Jobs:         []string{job.ID},
					Error:        err,
					ForceDiscard: true,
				}

				continue
			}

			// Try to open a new writer with the bucket.
			writer, err := a.bucket.NewWriter(a.ctx, payload.Filename, nil)
			if err != nil {
				then <- destination.Then{
					Jobs:  []string{job.ID},
					Error: err,
				}

				continue
			}

			// Try to write the content into the bucket with given filename.
			// If the writer didn't return an error when closing it is safe to
			// assume the content is has successfully been written.
			writer.Write(payload.Content)
			err = writer.Close()
			then <- destination.Then{
				Jobs:  []string{job.ID},
				Error: err,
			}
		}
	}
}
