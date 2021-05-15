---
title: Blob storage with Blacksmith
enterprise: false
time: 10
level: Beginner
modules:
  - blob
links:
  - name: Go API reference
    url: https://pkg.go.dev/github.com/nunchistudio/blacksmith-modules/blob
  - name: Source code
    url: https://github.com/nunchistudio/blacksmith-modules/tree/main/blob
---

# Blob storage with Blacksmith

In this tutorial, we are going to configure and use the Go module dedicated to
Blob stores. The module `blob` is composed of a package for Loading data to
different stores, which is `blobdestination`. This package exposes a Blacksmith
[`destination.Destination`](https://pkg.go.dev/github.com/nunchistudio/blacksmith/flow/destination?tab=doc#Destination).

The following blob stores are supported:
- AWS S3-compatible (`DriverAWSS3`)
- Azure Blob Storage (`DriverAzureBlob`)
- Google Cloud Storage (`DriverGoogleStorage`)

## Registering the destination

To use the Blob destination you first need to register it in the 
[`*blacksmith.Options`](https://pkg.go.dev/github.com/nunchistudio/blacksmith?tab=doc#Options).

Multiple Blob destinations can be registered:
```go
package main

import (
  "github.com/nunchistudio/blacksmith"
  "github.com/nunchistudio/blacksmith/flow/destination"

  "github.com/nunchistudio/blacksmith-modules/blob/blobdestination"
)

func Init() *blacksmith.Options {

  var options = &blacksmith.Options{

    // ...

    Destinations: []destination.Destination{
      blobdestination.New(&blobdestination.Options{
        Driver:     blobdestination.DriverAzureBlob,
        Name:       "bucket-a",
        Connection: "mybucket-a",
      }),
      blobdestination.New(&blobdestination.Options{
        Driver:     blobdestination.DriverAzureBlob,
        Name:       "bucket-b",
        Connection: "mybucket-b",
      }),
    },
  }

  return options
}

```

The destinations are now accessible by using the `blob(bucket-a)` and `blob(bucket-b)`
identifiers when one is required. The main use case will be for Transforming and
Loading data to the destination.

## Loading data to the destination

Now that the destination is registered, we can execute its action from a trigger
or from a flow.

In the following example, we call the `Write` action from a HTTP trigger:
```go
func (t Identify) Extract(tk *source.Toolkit, req *http.Request) (*source.Event, error) {

  // ...

  return &source.Event{
    Context: ctx,
    Data:    data,
    Actions: destination.Actions{
      "blob(bucket-a)": []destination.Action{
        blobdestination.Write{
          Filename: "myevent.json",
          Content:  data,
        },
      },
    },
  }, nil
}

```
