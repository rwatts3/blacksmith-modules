---
title: SQL databases with Blacksmith
enterprise: false
time: 15
level: Beginner
modules:
  - sqlike
links:
  - name: Go API reference
    url: https://pkg.go.dev/github.com/nunchistudio/blacksmith-modules/sqlike
  - name: Source code
    url: https://github.com/nunchistudio/blacksmith-modules/tree/main/sqlike
---

# SQL databases with Blacksmith

In this tutorial, we are going to configure and use the Go module dedicated to
SQL databases. The module `sqlike` is composed of a package for Loading data
to different SQL databases, which is `sqlikedestination`. This package exposes a
Blacksmith [`destination.Destination`](https://pkg.go.dev/github.com/nunchistudio/blacksmith/flow/destination?tab=doc#Destination).

Any Go SQL driver built on top of the standard library with `database/sql` is
supported. This includes PostgreSQL-compatible, MySQL-compatible, ClickHouse,
Snowflake, and more.

## Registering the destination

To use the SQLike destination you first need to register it in the 
[`*blacksmith.Options`](https://pkg.go.dev/github.com/nunchistudio/blacksmith?tab=doc#Options).

Multiple SQLike destinations can be registered:
```go
package main

import (
  "github.com/nunchistudio/blacksmith"
  "github.com/nunchistudio/blacksmith/flow/destination"

  "github.com/nunchistudio/blacksmith-modules/sqlike/sqlikedestination"
)

func Init() *blacksmith.Options {

  var options = &blacksmith.Options{

    // ...

    Destinations: []destination.Destination{
      sqlikedestination.New(&sqlikedestination.Options{
        DB:   <client>,
        Name: "mydb-a",
      }),
      sqlikedestination.New(&sqlikedestination.Options{
        DB:         <client>,
        Name:       "mydb-b",
        Migrations: []string{"mydb-b", "migrations"},
      }),
    },
  }

  return options
}

```

The destinations are now accessible by using the `sqlike(mydb-a)` and `sqlike(mydb-b)`
identifiers when one is required. The main use case will be for Transforming and
Loading data to the destination.

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
      "sqlike(mydb-a)": {
        sqlikedestination.Run{
          Statements: []sqlikedestination.Statement{
            {
              Query: "INSERT INTO public.actions (name, game, user_id) VALUES ($1, $2, $3);",
              Values: [][]interface{}{
                {"move_up", "mygame", "7923749"},
                {"move_right", "mygame", "9318562"},
                {"move_up", "mygame", "7923749"},
              },
            },
          },
        },
      },
    },
  }, nil
}

```

## Managing migrations for the destination

Destinations registered in a Blacksmith application and leveraging the `sqlike`
module can manage migrations thanks to the Blacksmith CLI.

In the first example given in this document, the destination `sqlike(mydb-b)` has
a directory with `*.sql` files for handling `up` and `down` migrations. To run
the migrations for this destination only, one could run:
```bash
$ blacksmith migrations run --scope "destination:sqlike(mydb-b)"

```

If you haven't built or started your application with the destination yet, you'll
need to add the `--build` flag in order to rebuild the application first:
```bash
$ blacksmith migrations run --scope "destination:sqlike(mydb-b)" \
  --build

```

**Related ressources:**
- Advanced practices >
  [Migrations management](/blacksmith/practices/management/migrations)
- CLI reference >
  [`generate migration`](/blacksmith/cli/generate-migration)
- CLI reference >
  [`migrations ack`](/blacksmith/cli/migrations-ack)
- CLI reference >
  [`migrations run`](/blacksmith/cli/migrations-run)
- CLI reference >
  [`migrations rollback`](/blacksmith/cli/migrations-rollback)
