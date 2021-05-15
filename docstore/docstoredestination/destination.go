package docstoredestination

import (
	"context"
	"fmt"

	"github.com/nunchistudio/blacksmith/flow/destination"
	"github.com/nunchistudio/blacksmith/helper/errors"
	"github.com/nunchistudio/blacksmith/helper/logger"

	"gocloud.dev/docstore"
	_ "gocloud.dev/docstore/awsdynamodb"
	_ "gocloud.dev/docstore/gcpfirestore"
	_ "gocloud.dev/docstore/mongodocstore"
)

/*
Docstore implements the Blacksmith destination.Destination interface for working
with NoSQL document stores.
*/
type Docstore struct {
	options    *destination.Options
	env        *Options
	ctx        context.Context
	collection *docstore.Collection
}

/*
New returns a valid Blacksmith destination.Destination for a document store.
*/
func New(env *Options) destination.Destination {
	if err := env.validate(); err != nil {
		logger.Default.Fatal(err)
		return nil
	}

	return &Docstore{
		options: &destination.Options{
			DefaultSchedule: &destination.Schedule{
				Realtime:   env.Realtime,
				Interval:   env.Interval,
				MaxRetries: env.MaxRetries,
			},
		},
		env: env,
		ctx: context.Background(),
	}
}

/*
String returns the string representation of the destination passed by the user
when creating the Docstore destination.
*/
func (d *Docstore) String() string {
	return fmt.Sprintf("docstore(%s)", d.env.Name)
}

/*
Init is part of the destination.WithHooks interface. It allows to properly open
the connection with the document store. It is called when starting the scheduler
service.
*/
func (d *Docstore) Init(tk *destination.Toolkit) error {
	var err error
	var collection *docstore.Collection

	// Open the document store given the driver and the URL passed by the user.
	url := d.env.Connection + "?" + d.env.Params.Encode()
	switch d.env.Driver {
	case DriverAWSDynamoDB:
		collection, err = docstore.OpenCollection(d.ctx, "dynamodb://"+url)
	case DriverAzureCosmosDB:
		collection, err = docstore.OpenCollection(d.ctx, "mongo://"+url)
	case DriverGoogleFirestore:
		collection, err = docstore.OpenCollection(d.ctx, "firestore://"+url)
	case DriverMongoDB:
		collection, err = docstore.OpenCollection(d.ctx, "mongo://"+url)
	default:
		return &errors.Error{
			Message: fmt.Sprintf("%s: Driver not supported", d.String()),
		}
	}

	if err != nil {
		return &errors.Error{
			Message: fmt.Sprintf("%s: %s", d.String(), err.Error()),
		}
	}

	d.collection = collection
	return nil
}

/*
Shutdown is part of the destination.WithHooks interface. It allows to properly
close the connection with the collection. It is called when shutting down the
scheduler service.
*/
func (d *Docstore) Shutdown(tk *destination.Toolkit) error {
	if d.collection != nil {
		err := d.collection.Close()
		if err != nil {
			return &errors.Error{
				Message: fmt.Sprintf("%s: Failed to properly close connection with collection", d.String()),
			}
		}
	}

	return nil
}

/*
Options returns common destination options for a blob storage. They will be
shared across every actions of this destination, except when overridden.
*/
func (d *Docstore) Options() *destination.Options {
	return d.options
}

/*
Actions return a list of actions the destination Docstore is able to handle.
*/
func (d *Docstore) Actions() map[string]destination.Action {
	return map[string]destination.Action{
		"put": Put{
			env:        d.env,
			ctx:        d.ctx,
			collection: d.collection,
		},
	}
}
