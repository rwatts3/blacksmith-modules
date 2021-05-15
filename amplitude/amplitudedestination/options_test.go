package amplitudedestination

import (
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
				APIKey:     "fakeapikey",
			},
			wantErr: false,
		},
		{
			name: "WithNoMaxRetries",
			fields: &Options{
				Realtime:   false,
				Interval:   "@every 1h",
				MaxRetries: 0,
				APIKey:     "fakeapikey",
			},
			wantErr: false,
		},
		{
			name: "WithNoAPIKey",
			fields: &Options{
				Realtime:   false,
				Interval:   "@every 1h",
				MaxRetries: 0,
				APIKey:     "",
			},
			wantErr: true,
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
