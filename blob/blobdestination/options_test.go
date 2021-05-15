package blobdestination

import (
	"net/url"
	"testing"
)

var DriverTest Driver = "fake"

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
				Driver:     DriverTest,
				Connection: "conn://fakeurl",
				Params: url.Values{
					"hello": {"world"},
				},
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
				Driver:     DriverTest,
				Connection: "conn://fakeurl",
				Params: url.Values{
					"hello": {"world"},
				},
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
				Driver:     DriverTest,
				Connection: "conn://fakeurl",
				Params: url.Values{
					"hello": {"world"},
				},
			},
			wantErr: true,
		},
		{
			name: "WithNoDriver",
			fields: &Options{
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
			wantErr: true,
		},
		{
			name: "WithNoConnection",
			fields: &Options{
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
			wantErr: true,
		},
		{
			name: "WithNoParams",
			fields: &Options{
				Realtime:   false,
				Interval:   "@every 1h",
				MaxRetries: 10,
				Name:       "fakename",
				Driver:     DriverTest,
				Connection: "conn://fakeurl",
				Params:     nil,
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
