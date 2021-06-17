package mailchimpdestination

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/nunchistudio/blacksmith/destination"
	"github.com/nunchistudio/blacksmith/helper/logger"

	"github.com/sirupsen/logrus"
)

var _ destination.Destination = &Mailchimp{}

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
					Realtime:          false,
					Interval:          "",
					MaxRetries:        10,
					APIKey:            "fakeapikey",
					DatacenterID:      "dc1",
					AudienceID:        "noop",
					EnableDoubleOptIn: false,
				},
			},
			want: &Mailchimp{
				options: &destination.Options{
					DefaultSchedule: &destination.Schedule{
						Realtime:   false,
						Interval:   destination.Defaults.DefaultSchedule.Interval,
						MaxRetries: 10,
					},
					DefaultVersion: "v3.0",
					Versions: map[string]time.Time{
						"v3.0": time.Time{},
					},
				},
				env: &Options{
					Realtime:          false,
					Interval:          destination.Defaults.DefaultSchedule.Interval,
					MaxRetries:        10,
					APIKey:            "fakeapikey",
					DatacenterID:      "dc1",
					AudienceID:        "noop",
					EnableDoubleOptIn: false,
				},
				client: http.DefaultClient,
			},
			wantFatal: false,
		},
		{
			name: "WithNoMaxRetries",
			args: args{
				env: &Options{
					Realtime:          false,
					Interval:          "@every 1h",
					MaxRetries:        0,
					APIKey:            "fakeapikey",
					DatacenterID:      "dc1",
					AudienceID:        "noop",
					EnableDoubleOptIn: false,
				},
			},
			want: &Mailchimp{
				options: &destination.Options{
					DefaultSchedule: &destination.Schedule{
						Realtime:   false,
						Interval:   "@every 1h",
						MaxRetries: destination.Defaults.DefaultSchedule.MaxRetries,
					},
					DefaultVersion: "v3.0",
					Versions: map[string]time.Time{
						"v3.0": time.Time{},
					},
				},
				env: &Options{
					Realtime:          false,
					Interval:          "@every 1h",
					MaxRetries:        destination.Defaults.DefaultSchedule.MaxRetries,
					APIKey:            "fakeapikey",
					DatacenterID:      "dc1",
					AudienceID:        "noop",
					EnableDoubleOptIn: false,
				},
				client: http.DefaultClient,
			},
			wantFatal: false,
		},
		{
			name: "WithNoAPIKey",
			args: args{
				env: &Options{
					Realtime:          false,
					Interval:          "@every 1h",
					MaxRetries:        0,
					APIKey:            "",
					DatacenterID:      "dc1",
					AudienceID:        "noop",
					EnableDoubleOptIn: false,
				},
			},
			want:      nil,
			wantFatal: true,
		},
		{
			name: "WithNoDatacenterID",
			args: args{
				env: &Options{
					Realtime:          false,
					Interval:          "@every 1h",
					MaxRetries:        0,
					APIKey:            "fakeapikey",
					DatacenterID:      "",
					AudienceID:        "noop",
					EnableDoubleOptIn: false,
				},
			},
			want:      nil,
			wantFatal: true,
		},
		{
			name: "WithNoAudienceID",
			args: args{
				env: &Options{
					Realtime:          false,
					Interval:          "@every 1h",
					MaxRetries:        0,
					APIKey:            "fakeapikey",
					DatacenterID:      "dc1",
					AudienceID:        "",
					EnableDoubleOptIn: false,
				},
			},
			want:      nil,
			wantFatal: true,
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
