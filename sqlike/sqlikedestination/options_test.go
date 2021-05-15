package sqlikedestination

import (
	"database/sql"
	"testing"
)

func TestOptions_validate(t *testing.T) {
	tests := []struct {
		name    string
		fields  *Options
		wantErr bool
	}{
		{
			name:    "WithEmptyOptions",
			fields:  &Options{},
			wantErr: true,
		},
		{
			name: "WithNoInterval",
			fields: &Options{
				Realtime:   false,
				Interval:   "",
				MaxRetries: 10,
				Name:       "fakename",
				DB:         &sql.DB{},
				Migrations: []string{"relative", "path"},
			},
			wantErr: false,
		},
		{
			name: "WithNoMaxRetries",
			fields: &Options{
				Realtime:   false,
				Interval:   "@every 1h",
				MaxRetries: 0,
				Name:       "fakename",
				DB:         &sql.DB{},
				Migrations: []string{"relative", "path"},
			},
			wantErr: false,
		},
		{
			name: "WithNoName",
			fields: &Options{
				Realtime:   false,
				Interval:   "@every 1h",
				MaxRetries: 10,
				Name:       "",
				DB:         &sql.DB{},
				Migrations: []string{"relative", "path"},
			},
			wantErr: true,
		},
		{
			name: "WithNoDB",
			fields: &Options{
				Realtime:   false,
				Interval:   "@every 1h",
				MaxRetries: 10,
				Name:       "fakename",
				DB:         nil,
				Migrations: []string{"relative", "path"},
			},
			wantErr: true,
		},
		{
			name: "WithNoMigrations",
			fields: &Options{
				Realtime:   false,
				Interval:   "@every 1h",
				MaxRetries: 10,
				Name:       "fakename",
				DB:         &sql.DB{},
				Migrations: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := tt.fields
			if err := env.validate(); (err != nil) != tt.wantErr {
				t.Errorf("Options.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
