package sqlikedestination

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/nunchistudio/blacksmith/flow/destination"
	"github.com/nunchistudio/blacksmith/helper/logger"

	"github.com/sirupsen/logrus"
)

var _ destination.Destination = &SQLike{}

func TestNew(t *testing.T) {
	var fatal bool
	logger.Default.Level = logrus.PanicLevel
	logger.Default.ExitFunc = func(int) {
		fatal = true
	}
	defer func() {
		logger.Default.ExitFunc = nil
	}()

	type args struct {
		env *Options
	}
	tests := []struct {
		name      string
		args      args
		want      destination.Destination
		wantFatal bool
	}{
		{
			name: "WithNilOptions",
			args: args{
				env: nil,
			},
			want:      nil,
			wantFatal: true,
		},
		{
			name: "WithEmptyOptions",
			args: args{
				env: &Options{},
			},
			want:      nil,
			wantFatal: true,
		},
		{
			name: "WithNoInterval",
			args: args{
				env: &Options{
					Realtime:   false,
					Interval:   "",
					MaxRetries: 10,
					Name:       "fakename",
					DB:         &sql.DB{},
					Migrations: []string{"relative", "path"},
				},
			},
			want: &SQLike{
				options: &destination.Options{
					DefaultSchedule: &destination.Schedule{
						Realtime:   false,
						Interval:   destination.Defaults.DefaultSchedule.Interval,
						MaxRetries: 10,
					},
				},
				env: &Options{
					Realtime:   false,
					Interval:   destination.Defaults.DefaultSchedule.Interval,
					MaxRetries: 10,
					Name:       "fakename",
					DB:         &sql.DB{},
					Migrations: []string{"relative", "path"},
				},
			},
			wantFatal: false,
		},
		{
			name: "WithNoMaxRetries",
			args: args{
				env: &Options{
					Realtime:   false,
					Interval:   "@every 1h",
					MaxRetries: 0,
					Name:       "fakename",
					DB:         &sql.DB{},
					Migrations: []string{"relative", "path"},
				},
			},
			want: &SQLike{
				options: &destination.Options{
					DefaultSchedule: &destination.Schedule{
						Realtime:   false,
						Interval:   "@every 1h",
						MaxRetries: destination.Defaults.DefaultSchedule.MaxRetries,
					},
				},
				env: &Options{
					Realtime:   false,
					Interval:   "@every 1h",
					MaxRetries: destination.Defaults.DefaultSchedule.MaxRetries,
					Name:       "fakename",
					DB:         &sql.DB{},
					Migrations: []string{"relative", "path"},
				},
			},
			wantFatal: false,
		},
		{
			name: "WithNoName",
			args: args{
				env: &Options{
					Realtime:   false,
					Interval:   "@every 1h",
					MaxRetries: 10,
					Name:       "",
					DB:         &sql.DB{},
					Migrations: []string{"relative", "path"},
				},
			},
			want:      nil,
			wantFatal: true,
		},
		{
			name: "WithNoDB",
			args: args{
				env: &Options{
					Realtime:   false,
					Interval:   "@every 1h",
					MaxRetries: 10,
					Name:       "fakename",
					DB:         nil,
					Migrations: []string{"relative", "path"},
				},
			},
			want:      nil,
			wantFatal: true,
		},
		{
			name: "WithNoMigrations",
			args: args{
				env: &Options{
					Realtime:   false,
					Interval:   "@every 1h",
					MaxRetries: 10,
					Name:       "fakename",
					DB:         &sql.DB{},
					Migrations: nil,
				},
			},
			want: &SQLike{
				options: &destination.Options{
					DefaultSchedule: &destination.Schedule{
						Realtime:   false,
						Interval:   "@every 1h",
						MaxRetries: 10,
					},
				},
				env: &Options{
					Realtime:   false,
					Interval:   "@every 1h",
					MaxRetries: 10,
					Name:       "fakename",
					DB:         &sql.DB{},
					Migrations: nil,
				},
			},
			wantFatal: false,
		},
	}
	for _, tt := range tests {
		fatal = false

		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.env)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(fatal, tt.wantFatal) {
				t.Errorf("New() should panic")
			}
		})
	}
}
