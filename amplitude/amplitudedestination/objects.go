package amplitudedestination

import (
	"net"

	"gopkg.in/segmentio/analytics-go.v3"
)

/*
Event holds the data of a event for the Amplitude API. It is used by the
Identify, Track, Page, and Screen actions.
*/
type Event struct {
	Event              string             `json:"event_type"`
	InsertId           string             `json:"insert_id,omitempty"`
	UserId             string             `json:"user_id"`
	DeviceId           string             `json:"device_id"`
	Time               int64              `json:"time,omitempty"`
	Traits             analytics.Traits   `json:"user_properties"`
	Context            *analytics.Context `json:"event_properties"`
	AppVersion         string             `json:"app_version,omitempty"`
	Platform           string             `json:"platform,omitempty"`
	OSName             string             `json:"os_name,omitempty"`
	OSVersion          string             `json:"os_version,omitempty"`
	DeviceBrand        string             `json:"device_brand,omitempty"`
	DeviceManufacturer string             `json:"device_manufacturer,omitempty"`
	DeviceModel        string             `json:"device_model,omitempty"`
	Carrier            string             `json:"carrier,omitempty"`
	Country            string             `json:"country,omitempty"`
	Region             string             `json:"region,omitempty"`
	City               string             `json:"city,omitempty"`
	Latitude           float64            `json:"location_lat,omitempty"`
	Longitude          float64            `json:"location_lng,omitempty"`
	DMA                string             `json:"dma,omitempty"`
	Language           string             `json:"language,omitempty"`
	Paying             string             `json:"paying,omitempty"`
	StartVersion       string             `json:"start_version,omitempty"`
	IP                 net.IP             `json:"ip,omitempty"`
	ProductId          string             `json:"productId,omitempty"`
	Quantity           float64            `json:"quantity,omitempty"`
	Price              float64            `json:"price,omitempty"`
	Revenue            float64            `json:"revenue,omitempty"`
}

/*
Identification struct holds the data needed to identify a user within a group
in Amplitude. It is used by the Group action.
*/
type Identification struct {
	UserId          string           `json:"user_id"`
	GroupType       string           `json:"group_type,omitempty"`
	GroupValue      string           `json:"group_value"`
	GroupProperties analytics.Traits `json:"group_properties"`
}

/*
UserMap struct holds the data needed to map users in Amplitude. It is used
by the Alias action.
*/
type UserMap struct {
	UserId       string `json:"user_id"`
	GlobalUserId string `json:"global_user_id"`
}
