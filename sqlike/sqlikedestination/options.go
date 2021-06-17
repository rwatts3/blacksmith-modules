package sqlikedestination

import (
	"database/sql"
	"fmt"

	"github.com/nunchistudio/blacksmith/destination"
	"github.com/nunchistudio/blacksmith/helper/errors"
)

/*
Options is the options the destination can take as an input to be configured.
*/
type Options struct {

	// Realtime indicates if the pubsub adapter of the Blacksmith application shall
	// be used to load events to the destination in realtime or not. When false, the
	// Interval will be used.
	Realtime bool

	// Interval represents an interval or a CRON string at which a job shall be
	// loaded to the destination. It is used as the time-lapse between retries in
	// case of a job failure.
	//
	// Defaults to "@every 1h".
	Interval string

	// MaxRetries indicates the maximum number of retries per job the scheduler will
	// attempt to execute before it succeed. When the limit is reached, the job is
	// marked as "discarded".
	//
	// Defaults to 72.
	MaxRetries uint16

	// Name indicates the identifier of the SQL database which will be used as name
	// for the destination. The computed name is "sqlike(<name>)".
	//
	// Examples: "postgres", "warehouse"
	// Required.
	Name string

	// DB is the database connection created using the package database/sql of the
	// standard library.
	//
	// Required.
	DB *sql.DB

	// Migrations is the relative path where the SQL migration files are located.
	// The path will be used using filepath.Join from the package path/filepath.
	//
	// If not set, migrations for the destination can not be managed with Blacksmith.
	//
	// Example: {"migrations", "mydestination"}
	Migrations []string
}

/*
validate ensures the options passed to initialize the destination are valid.
*/
func (env *Options) validate() error {
	var interval string = destination.Defaults.DefaultSchedule.Interval
	var maxRetries uint16 = destination.Defaults.DefaultSchedule.MaxRetries

	fail := &errors.Error{
		Message:     "sqlike: Failed to load",
		Validations: []errors.Validation{},
	}

	if env == nil {
		fail.Validations = append(fail.Validations, errors.Validation{
			Message: "Options must not be nil",
			Path:    []string{"Options", "Destinations", "sqlike"},
		})

		return fail
	}

	if env.Interval == "" {
		env.Interval = interval
	}

	if env.MaxRetries == 0 {
		env.MaxRetries = maxRetries
	}

	if env.Name == "" {
		fail.Validations = append(fail.Validations, errors.Validation{
			Message: "Database name must be set",
			Path:    []string{"Options", "Destinations", "sqlike", "Name"},
		})
	}

	// Create the computed name of the destination.
	name := fmt.Sprintf("sqlike(%s)", env.Name)

	if env.DB == nil {
		fail.Validations = append(fail.Validations, errors.Validation{
			Message: "Database client connection must be set",
			Path:    []string{"Options", "Destinations", name, "DB"},
		})
	}

	if len(fail.Validations) > 0 {
		return fail
	}

	return nil
}
