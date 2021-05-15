---
title: Publishing messages with Blacksmith
enterprise: false
time: 10
level: Beginner
modules:
  - topic
links:
  - name: Go API reference
    url: https://pkg.go.dev/github.com/nunchistudio/blacksmith-modules/topic
  - name: Source code
    url: https://github.com/nunchistudio/blacksmith-modules/tree/main/topic
---

# Publishing messages with Blacksmith

In this tutorial, we are going to configure and use the Go module dedicated to
Topics of message brokers. The module `topic` is composed of a package for Loading
(_a.k.a._ publishing) data to different brokers, which is `topicdestination`.
This package exposes a Blacksmith
[`destination.Destination`](https://pkg.go.dev/github.com/nunchistudio/blacksmith/flow/destination?tab=doc#Destination).

The following brokers are supported:
- AWS SNS (`DriverAWSSNS`)
- AWS SQS (`DriverAWSSQS`)
- Azure Service Bus (`DriverAzureServiceBus`)
- Google Pub / Sub (`DriverGooglePubSub`)
- Apache Kafka (`DriverKafka`)
- NATS (`DriverNATS`)
- RabbitMQ (`DriverRabbitMQ`)

## Registering the destination

To use the Topic destination you first need to register it in the 
[`*blacksmith.Options`](https://pkg.go.dev/github.com/nunchistudio/blacksmith?tab=doc#Options).

Multiple Topic destinations can be registered:
```go
package main

import (
  "github.com/nunchistudio/blacksmith"
  "github.com/nunchistudio/blacksmith/flow/destination"

  "github.com/nunchistudio/blacksmith-modules/topic/topicdestination"
)

func Init() *blacksmith.Options {

  var options = &blacksmith.Options{

    // ...

    Destinations: []destination.Destination{
      topicdestination.New(&topicdestination.Options{
        Driver:     topicdestination.DriverAWSSNS,
        Name:       "topic-a",
        Connection: "arn:aws:sns:<region>:<id>:<topic>",
        Params: url.Values{
          "region": {"<region>"},
        },
      }),
      topicdestination.New(&topicdestination.Options{
        Driver:     topicdestination.DriverAWSSNS,
        Name:       "topic-b",
        Connection: "arn:aws:sns:<region>:<id>:<topic>",
        Params: url.Values{
          "region": {"<region>"},
        },
      }),
    },
  }

  return options
}

```

The destinations are now accessible by using the `topic(topic-a)` and
`topic(topic-b)` identifiers when one is required. The main use case will
be for Transforming and Loading data to the destination.

## Loading data to the destination

Now that the destination is registered, we can execute its action from a trigger
or from a flow.

In the following example, we call the `Publish` action from a HTTP trigger:
```go
func (t Identify) Extract(tk *source.Toolkit, req *http.Request) (*source.Event, error) {

  // ...

  return &source.Event{
    Context: ctx,
    Data:    data,
    Actions: destination.Actions{
      "topic(topic-a)": []destination.Action{
        topicdestination.Publish{
          Message: topicdestination.Message{
            Body: data,
          },
        },
      },
    },
  }, nil
}

```
