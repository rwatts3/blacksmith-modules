package sqlikedestination

import (
	"encoding/json"

	"github.com/nunchistudio/blacksmith/adapter/store"
	"github.com/nunchistudio/blacksmith/flow/destination"
	"github.com/nunchistudio/blacksmith/helper/errors"
	"github.com/nunchistudio/blacksmith/warehouse"
)

/*
RunOperation implements the Blacksmith destination.Action interface for the
action "run-operation". It holds the complete job's structure to load into
the destination.
*/
type RunOperation struct {
	env *Options
	wh  *warehouse.Warehouse

	// Filename is the path and file name for the SQL file to compile and
	// execute.
	//
	// Example: "./operations/demo.sql"
	// Required.
	Filename string `json:"filename"`

	// Data is a free dictionary of data to pass to the template.
	Data map[string]interface{} `json:"data"`
}

/*
String returns the string representation of the action RunOperation.
*/
func (a RunOperation) String() string {
	return "run-operation"
}

/*
Schedule allows the action to override the schedule options of its
destination. Do not override.
*/
func (a RunOperation) Schedule() *destination.Schedule {
	return nil
}

/*
Marshal is the function being run when the action receives data into
the RunOperation receiver. It allows to transform and enrich the data
before saving it in the store adapter.
*/
func (a RunOperation) Marshal(tk *destination.Toolkit) (*destination.Job, error) {

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
func (a RunOperation) Load(tk *destination.Toolkit, queue *store.Queue, then chan<- destination.Then) {

	// We can go through every events received from the queue and their
	// related jobs. The queue can contain one or many events. The jobs
	// present in the events are specific to this action only.
	for _, event := range queue.Events {
		for _, job := range event.Jobs {

			var run RunOperation
			err := json.Unmarshal(job.Data, &run)
			if err != nil {
				then <- destination.Then{
					Jobs:         []string{job.ID},
					Error:        err,
					ForceDiscard: true,
				}

				continue
			}

			operation, err := a.wh.Compile(run.Filename, run.Data)
			if err != nil {
				then <- destination.Then{
					Jobs:  []string{job.ID},
					Error: err,
				}

				continue
			}

			err = a.wh.Exec(operation)
			then <- destination.Then{
				Jobs:  []string{job.ID},
				Error: err,
			}
		}
	}
}
