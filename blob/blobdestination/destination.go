package blobdestination

import (
	"context"
	"fmt"

	"github.com/nunchistudio/blacksmith/destination"
	"github.com/nunchistudio/blacksmith/helper/errors"
	"github.com/nunchistudio/blacksmith/helper/logger"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/azureblob"
	_ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/s3blob"
)

/*
Blob implements the Blacksmith destination.Destination interface for working
with blob storages.
*/
type Blob struct {
	options *destination.Options
	env     *Options
	ctx     context.Context
	bucket  *blob.Bucket
}

/*
New returns a valid Blacksmith destination.Destination for a blob storage.
*/
func New(env *Options) destination.Destination {
	if err := env.validate(); err != nil {
		logger.Default.Fatal(err)
		return nil
	}

	return &Blob{
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
when creating the Blob destination.
*/
func (d *Blob) String() string {
	return fmt.Sprintf("blob(%s)", d.env.Name)
}

/*
Init is part of the destination.WithHooks interface. It allows to properly open
the connection with the bucket. It is called when starting the scheduler service.
*/
func (d *Blob) Init(tk *destination.Toolkit) error {
	var err error
	var bucket *blob.Bucket

	// Open the bucket given the driver and the URL passed by the user.
	url := d.env.Connection + "?" + d.env.Params.Encode()
	switch d.env.Driver {
	case DriverAWSS3:
		bucket, err = blob.OpenBucket(d.ctx, "s3://"+url)
	case DriverAzureBlob:
		bucket, err = blob.OpenBucket(d.ctx, "azblob://"+url)
	case DriverGoogleStorage:
		bucket, err = blob.OpenBucket(d.ctx, "gs://"+url)
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

	d.bucket = bucket
	return nil
}

/*
Shutdown is part of the destination.WithHooks interface. It allows to properly
close the connection with the bucket. It is called when shutting down the
scheduler service.
*/
func (d *Blob) Shutdown(tk *destination.Toolkit) error {
	if d.bucket != nil {
		err := d.bucket.Close()
		if err != nil {
			return &errors.Error{
				Message: fmt.Sprintf("%s: Failed to properly close connection with bucket", d.String()),
			}
		}
	}

	return nil
}

/*
Options returns common destination options for a blob storage. They will be
shared across every actions of this destination, except when overridden.
*/
func (d *Blob) Options() *destination.Options {
	return d.options
}

/*
Actions return a list of actions the destination Blob is able to handle.
*/
func (d *Blob) Actions() map[string]destination.Action {
	return map[string]destination.Action{
		"write": Write{
			env:    d.env,
			ctx:    d.ctx,
			bucket: d.bucket,
		},
	}
}
