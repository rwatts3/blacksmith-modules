package topicdestination

import (
	"context"
	"encoding/json"

	"github.com/nunchistudio/blacksmith/adapter/store"
	"github.com/nunchistudio/blacksmith/destination"
	"github.com/nunchistudio/blacksmith/helper/errors"

	"gocloud.dev/pubsub"
)

/*
Publish implements the Blacksmith destination.Action interface for the action
"publish". It holds the complete job's structure to load into the destination.
*/
type Publish struct {
	env   *Options
	topic *pubsub.Topic
	ctx   context.Context

	Message Message `json:"message"`
}

/*
String returns the string representation of the action Publish.
*/
func (a Publish) String() string {
	return "publish"
}

/*
Schedule allows the action to override the schedule options of its
destination. Do not override.
*/
func (a Publish) Schedule() *destination.Schedule {
	return nil
}

/*
Marshal is the function being run when the action receives data into
the Publish receiver. It allows to transform and enrich the data before
saving it in the store adapter.
*/
func (a Publish) Marshal(tk *destination.Toolkit) (*destination.Job, error) {

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
func (a Publish) Load(tk *destination.Toolkit, queue *store.Queue, then chan<- destination.Then) {

	// We can go through every events received from the queue and their
	// related jobs. The queue can contain one or many events. The jobs
	// present in the events are specific to this action only.
	//
	// This allows to parse everything needed and make a request to the
	// destination for each event / job.
	for _, event := range queue.Events {
		for _, job := range event.Jobs {
			var pub Publish
			err := json.Unmarshal(job.Data, &pub)
			if err != nil {
				then <- destination.Then{
					Jobs:         []string{job.ID},
					Error:        err,
					ForceDiscard: true,
				}

				continue
			}

			// Format the message.
			msg := &pubsub.Message{
				Body:     pub.Message.Body,
				Metadata: pub.Message.Metadata,
			}

			// Publish the message to the message broker.
			err = a.topic.Send(a.ctx, msg)
			then <- destination.Then{
				Jobs:  []string{job.ID},
				Error: err,
			}
		}
	}
}
