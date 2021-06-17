package sqlikedestination

import (
	"fmt"
	"path/filepath"

	"github.com/nunchistudio/blacksmith/adapter/wanderer"
	"github.com/nunchistudio/blacksmith/destination"
	"github.com/nunchistudio/blacksmith/helper/errors"
	"github.com/nunchistudio/blacksmith/helper/logger"
	"github.com/nunchistudio/blacksmith/warehouse"

	"github.com/nunchistudio/blacksmith-modules/sqlike"
)

/*
SQLike implements the Blacksmith destination.Destination interface for any and
every SQL destinations.
*/
type SQLike struct {
	options *destination.Options
	env     *Options
	wh      *warehouse.Warehouse
}

/*
New returns a valid Blacksmith destination.Destination for a SQL-like database.
*/
func New(env *Options) destination.Destination {
	if err := env.validate(); err != nil {
		logger.Default.Fatal(err)
		return nil
	}

	return &SQLike{
		options: &destination.Options{
			DefaultSchedule: &destination.Schedule{
				Realtime:   env.Realtime,
				Interval:   env.Interval,
				MaxRetries: env.MaxRetries,
			},
		},
		env: env,
	}
}

/*
String returns the string representation of the destination passed by the user
when creating the SQLike destination.
*/
func (d *SQLike) String() string {
	return fmt.Sprintf("sqlike(%s)", d.env.Name)
}

/*
Init is part of the destination.WithHooks interface. The SQL client is already
initialized and passed in the destination's options. But we still need to
save the warehouse for future use.
*/
func (d *SQLike) Init(tk *destination.Toolkit) error {
	wh, err := d.AsWarehouse()
	if err != nil {
		return err
	}

	d.wh = wh
	return nil
}

/*
Shutdown is part of the destination.WithHooks interface. It allows to properly
close the connection pool with the database. It is called when shutting down the
scheduler service or after running migrations.
*/
func (d *SQLike) Shutdown(tk *destination.Toolkit) error {
	if d.env.DB != nil {
		err := d.env.DB.Close()
		if err != nil {
			return &errors.Error{
				Message: fmt.Sprintf("%s: Failed to properly close connection with database", d.String()),
			}
		}
	}

	return nil
}

/*
Options returns common destination options for a SQL-like database. They will be
shared across every actions of this destination, except when overridden.
*/
func (d *SQLike) Options() *destination.Options {
	return d.options
}

/*
Actions return a list of actions the destination SQLike is able to handle.
*/
func (d *SQLike) Actions() map[string]destination.Action {
	return map[string]destination.Action{
		"run-statements": RunStatements{
			env: d.env,
		},
		"run-operation": RunOperation{
			env: d.env,
			wh:  d.wh,
		},
	}
}

/*
Migrate is the implementation of the wanderer.WithMigrate interface for the
destination SQLike. It allows the destination, and all of its actions, to have
a migration logic. This is the function called whenever a migration needs to
run or to rollback.

It leverages the sqlike package for running the migration within a SQL
transaction, using the standard database/sql package.
*/
func (d *SQLike) Migrate(tk *wanderer.Toolkit, migration *wanderer.Migration) error {
	if d.env.Migrations == nil || len(d.env.Migrations) == 0 {
		return nil
	}

	return sqlike.RunMigration(d.env.DB, filepath.Join(d.env.Migrations...), migration)
}

/*
Migrations is the implementation of the wanderer.WithMigrations interface for
the destination SQLike. It allows the destination to have migrations.

It leverages the sqlike package for finding compatible SQL files within a
directory.
*/
func (d *SQLike) Migrations(tk *wanderer.Toolkit) ([]*wanderer.Migration, error) {
	if d.env.Migrations == nil || len(d.env.Migrations) == 0 {
		return []*wanderer.Migration{}, nil
	}

	return sqlike.LoadMigrations(filepath.Join(d.env.Migrations...))
}
