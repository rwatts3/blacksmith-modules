package mailchimpdestination

import (
	"net"
	"time"

	"gopkg.in/segmentio/analytics-go.v3"
)

/*
Signup holds the data to identify a new user via the Mailchimp API.
*/
type Signup struct {
	StatusIfNew     string                  `json:"status_if_new,omitempty"`
	Email           string                  `json:"email_address"`
	FirstName       string                  `json:"first_name,omitempty"`
	LastName        string                  `json:"last_name,omitempty"`
	IPSignup        net.IP                  `json:"ip_signup"`
	TimestampSignup time.Time               `json:"timestamp_signup"`
	Location        *analytics.LocationInfo `json:"location,omitempty"`
}
