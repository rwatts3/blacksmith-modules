package sqlikedestination

import (
	"database/sql"
	"encoding/json"

	"github.com/nunchistudio/blacksmith/adapter/store"
	"github.com/nunchistudio/blacksmith/destination"
	"github.com/nunchistudio/blacksmith/helper/errors"
)

/*
RunStatements implements the Blacksmith destination.Action interface for the
action "run-statements". It holds the complete job's structure to load into
the destination.
*/
type RunStatements struct {
	env *Options

	// Statements holds all the prepared statements with their respective query
	// and values to load into the destination.
	Statements []Statement `json:"statements"`
}

/*
String returns the string representation of the action RunStatements.
*/
func (a RunStatements) String() string {
	return "run-statements"
}

/*
Schedule allows the action to override the schedule options of its
destination. Do not override.
*/
func (a RunStatements) Schedule() *destination.Schedule {
	return nil
}

/*
Marshal is the function being run when the action receives data into
the RunStatements receiver. It allows to transform and enrich the data
before saving it in the store adapter.
*/
func (a RunStatements) Marshal(tk *destination.Toolkit) (*destination.Job, error) {

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
func (a RunStatements) Load(tk *destination.Toolkit, queue *store.Queue, then chan<- destination.Then) {

	// Whenever we return, inform the scheduler with the load status.
	var err error
	var discard bool
	defer func() {
		then <- destination.Then{
			Error:        err,
			ForceDiscard: discard,
		}
	}()

	// Load the jobs inside a transaction. If the transaction failed to start
	// there is no need to continue.
	tx, err := a.env.DB.Begin()
	if err != nil {
		return
	}

	// Make sure to rollback the transaction if needed.
	defer tx.Rollback()

	// We can go through every events received from the queue and their
	// related jobs. The queue can contain one or many events. The jobs
	// present in the events are specific to this action only.
	var stop bool
	for _, event := range queue.Events {
		if stop {
			break
		}

		for _, job := range event.Jobs {
			if stop {
				break
			}

			var run RunStatements
			err = json.Unmarshal(job.Data, &run)
			if err != nil {
				discard = true
				stop = true
				break
			}

			// Execute a prepared statement for each one present in the slice
			// of statements.
			var stmt *sql.Stmt
			for _, exec := range run.Statements {
				if stop {
					break
				}

				stmt, err = tx.Prepare(exec.Query)
				if err != nil {
					stop = true
					break
				}

				// Make sure to close the current SQL statement when done.
				defer stmt.Close()

				// Execute the prepared statement with the arguments given.
				for _, row := range exec.Values {
					_, err = stmt.Exec(row...)
					if err != nil {
						stop = true
						break
					}
				}
			}
		}
	}

	// Do not try to commit the transaction if an error occured. In addition,
	// it allows to return the error encountered within the transaction and
	// not the failed transaction itself.
	if err != nil {
		return
	}

	// We can now try to commit the transaction.
	err = tx.Commit()
	return
}
