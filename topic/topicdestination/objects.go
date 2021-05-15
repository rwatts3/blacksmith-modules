package topicdestination

/*
Message represents a message published by the source and received by the message
broker.
*/
type Message struct {
	Body     []byte            `json:"body"`
	Metadata map[string]string `json:"meta"`
}
