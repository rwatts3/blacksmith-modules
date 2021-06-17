package sqlikedestination

import (
	"github.com/nunchistudio/blacksmith/warehouse"
)

/*
AsWarehouse returns a Blacksmith data warehouse, allowing end-users to run
queries and operations with template SQL on top of their database.
*/
func (d *SQLike) AsWarehouse() (*warehouse.Warehouse, error) {
	opts := &warehouse.Options{
		Name: d.String(),
		DB:   d.env.DB,
	}

	return warehouse.New(opts)
}
