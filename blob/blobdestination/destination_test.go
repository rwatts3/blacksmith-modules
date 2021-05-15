package blobdestination

import (
	"context"
	"net/url"
	"reflect"
	"testing"

	"github.com/nunchistudio/blacksmith/flow/destination"
	"github.com/nunchistudio/blacksmith/helper/logger"

	"github.com/sirupsen/logrus"
)

var _ destination.Destination = &Blob{}

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
					Driver:     DriverTest,
					Connection: "conn://fakeurl",
					Params: url.Values{
						"hello": {"world"},
					},
				},
			},
			want: &Blob{
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
					Driver:     DriverTest,
					Connection: "conn://fakeurl",
					Params: url.Values{
						"hello": {"world"},
					},
				},
				ctx: context.Background(),
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
					Driver:     DriverTest,
					Connection: "conn://fakeurl",
					Params: url.Values{
						"hello": {"world"},
					},
				},
			},
			want: &Blob{
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
					Driver:     DriverTest,
					Connection: "conn://fakeurl",
					Params: url.Values{
						"hello": {"world"},
					},
				},
				ctx: context.Background(),
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
					Driver:     DriverTest,
					Connection: "conn://fakeurl",
					Params: url.Values{
						"hello": {"world"},
					},
				},
			},
			want:      nil,
			wantFatal: true,
		},
		{
			name: "WithNoDriver",
			args: args{
				env: &Options{
					Realtime:   false,
					Interval:   "@every 1h",
					MaxRetries: 10,
					Name:       "fakename",
					Driver:     "",
					Connection: "conn://fakeurl",
					Params: url.Values{
						"hello": {"world"},
					},
				},
			},
			want:      nil,
			wantFatal: true,
		},
		{
			name: "WithNoConnection",
			args: args{
				env: &Options{
					Realtime:   false,
					Interval:   "@every 1h",
					MaxRetries: 10,
					Name:       "fakename",
					Driver:     DriverTest,
					Connection: "",
					Params: url.Values{
						"hello": {"world"},
					},
				},
			},
			want:      nil,
			wantFatal: true,
		},
		{
			name: "WithNoParams",
			args: args{
				env: &Options{
					Realtime:   false,
					Interval:   "@every 1h",
					MaxRetries: 10,
					Name:       "fakename",
					Driver:     DriverTest,
					Connection: "conn://fakeurl",
					Params:     nil,
				},
			},
			want: &Blob{
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
					Driver:     DriverTest,
					Connection: "conn://fakeurl",
					Params:     nil,
				},
				ctx: context.Background(),
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
