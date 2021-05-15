---
title: Loading data to Mailchimp
enterprise: false
time: 10
level: Beginner
modules:
  - mailchimp
  - segment
links:
  - name: Segment Clone
    url: /blacksmith/tutorials/fragment
---

# Mailchimp with Blacksmith

In this tutorial, we are going to configure and use the Go module dedicated to the
Mailchimp API. The module `mailchimp` is composed of a package for Loading data
to Mailchimp, which is `mailchimpdestination`. This package exposes a Blacksmith
[`destination.Destination`](https://pkg.go.dev/github.com/nunchistudio/blacksmith/flow/destination?tab=doc#Destination)
following the [Segment Specification](https://segment.com/docs/connections/spec/)
for Customer Data.

## Registering the destination

To use the Mailchimp destination you first need to register it in the 
[`*blacksmith.Options`](https://pkg.go.dev/github.com/nunchistudio/blacksmith?tab=doc#Options).

```go
package main

import (
  "github.com/nunchistudio/blacksmith"
  "github.com/nunchistudio/blacksmith/flow/destination"

  "github.com/nunchistudio/blacksmith-modules/mailchimp/mailchimpdestination"
)

func Init() *blacksmith.Options {

  var options = &blacksmith.Options{

    // ...

    Destinations: []destination.Destination{
      mailchimpdestination.New(&mailchimpdestination.Options{
        APIKey:            os.Getenv("MAILCHIMP_API_KEY"),
        DatacenterID:      os.Getenv("MAILCHIMP_DATACENTER"),
        AudienceID:        os.Getenv("MAILCHIMP_AUDIENCE"),
        EnableDoubleOptIn: true,
      }),
    },
  }

  return options
}

```

The destination is now accessible by using the `mailchimp` identifier when one is
required. The main use case will be for Transforming and Loading data to the
destination.

## Loading data to the destination

Now that the destination is registered, we can execute its actions from a trigger
or from a flow.

In the following example, we call the `Identify` action from a HTTP trigger:
```go
package rest

import (
  "encoding/json"
  "net/http"
  "strings"
  "time"

  "github.com/nunchistudio/blacksmith/flow/destination"
  "github.com/nunchistudio/blacksmith/flow/source"
  "github.com/nunchistudio/blacksmith/helper/errors"

  "github.com/nunchistudio/blacksmith-modules/mailchimp/mailchimpdestination"

  "gopkg.in/segmentio/analytics-go.v3"
)

/*
Identify implements the Blacksmith source.Trigger interface for the trigger
"identify". It holds the complete payload structure sent by an event and that
will be received by the gateway.
*/
type Identify struct {
  env *Options

  analytics.Identify
}

/*
String returns the string representation of the trigger Identify.
*/
func (t Identify) String() string {
  return "identify"
}

/*
Mode allows to register the trigger as a HTTP route. This means, every
time a "POST" request is executed against the route "/v1/identify", the
Extract function will run.
*/
func (t Identify) Mode() *source.Mode {
  return &source.Mode{
    Mode: source.ModeHTTP,
    UsingHTTP: &source.Route{
      Methods:  []string{"POST"},
      Path:     t.env.Prefix + "/v1/identify",
      ShowMeta: t.env.ShowMeta,
      ShowData: t.env.ShowData,
    },
  }
}

/*
Extract is the function being run when the HTTP route is triggered. It is
in charge of the "E" in the ETL process: Extract the data from the source.

The function allows to return data to flows or directly to actions. It is
the "T" in the ETL process: it transforms the payload from the source's
trigger to given destinations' actions.
*/
func (t Identify) Extract(tk *source.Toolkit, req *http.Request) (*source.Event, error) {

  // Create an empty payload, catch unwanted fields, and unmarshal it.
  // Return an error if any occured.
  var payload Identify
  decoder := json.NewDecoder(req.Body)
  decoder.DisallowUnknownFields()
  err := decoder.Decode(&payload)
  if err != nil {
    return nil, &errors.Error{
      StatusCode: 400,
      Message:    "Bad Request",
      Validations: []errors.Validation{
        {
          Message: err.Error(),
          Path:    []string{"analytics", "Identify"},
        },
      },
    }
  }

  // Add the current timestamp if none was provided.
  if payload.Timestamp.IsZero() {
    payload.Timestamp = time.Now().UTC()
  }

  // Validate the payload using the Segment official library.
  err = payload.Validate()
  if err != nil {
    fail := err.(analytics.FieldError)
    return nil, &errors.Error{
      StatusCode: 400,
      Message:    "Bad Request",
      Validations: []errors.Validation{
        {
          Message: fail.Name + " must be set",
          Path:    append(strings.Split(fail.Type, "."), fail.Name),
        },
      },
    }
  }

  // Try to marshal the context from the request payload.
  var ctx []byte
  if payload.Context != nil {
    ctx, err = payload.Context.MarshalJSON()
    if err != nil {
      return nil, &errors.Error{
        StatusCode: 400,
        Message:    "Bad Request",
      }
    }
  }

  // Try to marshal the data from the request payload.
  var data []byte
  if payload.Traits != nil {
    data, err = json.Marshal(&payload.Traits)
    if err != nil {
      return nil, &errors.Error{
        StatusCode: 400,
        Message:    "Bad Request",
      }
    }
  }

  // Return the context, data, and a collection of actions to run, including
  // the Mailchimp action.
  return &source.Event{
    Version: "v1.0",
    Context: ctx,
    Data:    data,
    Actions: destination.Actions{
      "mailchimp": []destination.Action{
        mailchimpdestination.Identify{
          Email: payload.Traits["email"],
          // ...
        },
      },
    },
    SentAt:  &payload.Timestamp,
  }, nil
}

```

### Using the Segment flow

The `mailchimp` identifier is used in the Segment flow. This flow Loads Customer
Data to web services following the Segment Specification, such as Mailchimp. If
the destination is indeed registered in the Blacksmith application like we did in
the first step, the scheduler will create the appropriate jobs.

Therefore, the `Extract` function only need to execute the Segment flow and shall
look like this:
```go
func (t Identify) Extract(tk *source.Toolkit, req *http.Request) (*source.Event, error) {

  // ...

  return &source.Event{
    Version: "v1.0",
    Context: ctx,
    Data:    data,
    Flows: []flow.Flow{
      &segmentflow.Identify{
        Identify: payload.Identify,
      },
    },
    SentAt: &payload.Timestamp,
  }, nil
}

```
