---
title: NoSQL document storage with Blacksmith
enterprise: false
time: 10
level: Beginner
modules:
  - docstore
links:
  - name: Go API reference
    url: https://pkg.go.dev/github.com/nunchistudio/blacksmith-modules/docstore
  - name: Source code
    url: https://github.com/nunchistudio/blacksmith-modules/tree/main/docstore
---

# NoSQL document storage with Blacksmith

In this tutorial, we are going to configure and use the Go module dedicated to
Document stores. The module `docstore` is composed of a package for Loading data
to different stores, which is `docstoredestination`. This package exposes a Blacksmith
[`destination.Destination`](https://pkg.go.dev/github.com/nunchistudio/blacksmith/flow/destination?tab=doc#Destination).

The following document stores are supported:
- AWS DynamoDB (`DriverAWSDynamoDB`)
- Azure CosmosDB (`DriverAzureCosmosDB`)
- MongoDB (`DriverMongoDB`)
- Google Firestore (`DriverGoogleFirestore`)

## Registering the destination

To use the Docstore destination you first need to register it in the 
[`*blacksmith.Options`](https://pkg.go.dev/github.com/nunchistudio/blacksmith?tab=doc#Options).

Multiple Docstore destinations can be registered:
```go
package main

import (
  "github.com/nunchistudio/blacksmith"
  "github.com/nunchistudio/blacksmith/flow/destination"

  "github.com/nunchistudio/blacksmith-modules/docstore/docstoredestination"
)

func Init() *blacksmith.Options {

  var options = &blacksmith.Options{

    // ...

    Destinations: []destination.Destination{
      docstoredestination.New(&docstoredestination.Options{
        Driver:     docstoredestination.DriverGoogleFirestore,
        Name:       "docstore-a",
        Connection: "projects/myproject/databases/(default)/documents/mycollection",
        Params: url.Values{
          "name_field": {"<field>"},
        },
      }),
      docstoredestination.New(&docstoredestination.Options{
        Driver:     docstoredestination.DriverAzureCosmosDB,
        Name:       "docstore-b",
        Connection: "mydb/mycollection",
      }),
    },
  }

  return options
}

```

The destinations are now accessible by using the `docstore(docstore-a)` and
`docstore(docstore-b)` identifiers when one is required. The main use case will
be for Transforming and Loading data to the destination.

## Loading data to the destination

Now that the destination is registered, we can execute its action from a trigger
or from a flow.

In the following example, we call the `Put` action from a HTTP trigger:
```go
func (t Identify) Extract(tk *source.Toolkit, req *http.Request) (*source.Event, error) {

  // ...

  return &source.Event{
    Context: ctx,
    Data:    data,
    Actions: destination.Actions{
      "docstore(docstore-a)": []destination.Action{
        docstoredestination.Put{
          Document: map[string]interface{}{
            "key": "value",
          },
        },
      },
    },
  }, nil
}

```
