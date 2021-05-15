package topicdestination

import (
	"context"
	"fmt"
	"strings"

	"github.com/nunchistudio/blacksmith/flow/destination"
	"github.com/nunchistudio/blacksmith/helper/errors"
	"github.com/nunchistudio/blacksmith/helper/logger"

	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/awssnssqs"
	_ "gocloud.dev/pubsub/azuresb"
	_ "gocloud.dev/pubsub/gcppubsub"
	_ "gocloud.dev/pubsub/kafkapubsub"
	_ "gocloud.dev/pubsub/natspubsub"
	_ "gocloud.dev/pubsub/rabbitpubsub"
)

/*
Topic implements the Blacksmith destination.Destination interface for working
with message brokers.
*/
type Topic struct {
	options *destination.Options
	env     *Options
	ctx     context.Context
	topic   *pubsub.Topic
}

/*
New returns a valid Blacksmith destination.Destination for a message broker.
*/
func New(env *Options) destination.Destination {
	if err := env.validate(); err != nil {
		logger.Default.Fatal(err)
		return nil
	}

	return &Topic{
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
when creating the Topic destination.
*/
func (d *Topic) String() string {
	return fmt.Sprintf("topic(%s)", d.env.Name)
}

/*
Init is part of the destination.WithHooks interface. It allows to properly open
the connection with the message broker. It is called when starting the scheduler
service.
*/
func (d *Topic) Init(tk *destination.Toolkit) error {
	var err error
	var topic *pubsub.Topic

	// Open the topic given the driver and the URL passed by the user.
	url := d.env.Connection + "?" + d.env.Params.Encode()
	switch d.env.Driver {
	case DriverAWSSNS:
		topic, err = pubsub.OpenTopic(d.ctx, "awssns:///"+url)
	case DriverAWSSQS:
		splitted := strings.Split(url, ":")
		transformed := splitted[2] + "." + splitted[3] + ".amazonaws.com/" + splitted[4] + "/" + splitted[5]
		topic, err = pubsub.OpenTopic(d.ctx, "awssqs://"+transformed)
	case DriverAzureServiceBus:
		topic, err = pubsub.OpenTopic(d.ctx, "azuresb://"+url)
	case DriverGooglePubSub:
		topic, err = pubsub.OpenTopic(d.ctx, "gcppubsub://"+url)
	case DriverKafka:
		topic, err = pubsub.OpenTopic(d.ctx, "kafka://"+url)
	case DriverNATS:
		topic, err = pubsub.OpenTopic(d.ctx, "nats://"+url)
	case DriverRabbitMQ:
		topic, err = pubsub.OpenTopic(d.ctx, "rabbit://"+url)
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

	d.topic = topic
	return nil
}

/*
Shutdown is part of the destination.WithHooks interface. It allows to properly
close the connection with the message broker. It is called when shutting down the
scheduler service.
*/
func (d *Topic) Shutdown(tk *destination.Toolkit) error {
	if d.topic != nil {
		err := d.topic.Shutdown(d.ctx)
		if err != nil {
			return &errors.Error{
				Message: fmt.Sprintf("%s: Failed to properly close connection with topic", d.String()),
			}
		}
	}

	return nil
}

/*
Options returns common destination options for a blob storage. They will be
shared across every actions of this destination, except when overridden.
*/
func (d *Topic) Options() *destination.Options {
	return d.options
}

/*
Actions return a list of actions the destination Blob is able to handle.
*/
func (d *Topic) Actions() map[string]destination.Action {
	return map[string]destination.Action{
		"publish": Publish{
			env:   d.env,
			topic: d.topic,
			ctx:   d.ctx,
		},
	}
}
