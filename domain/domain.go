package domain

import "time"

//// Struct for event logs
type Event struct {
	/// Host,it is the IP address of the client
	Host string `json:"host"`
	/// User identifier for the client
	User string `json:"user"`
	/// Time of the rquest
	Time time.Time `json:"time"`
	/// Method: GET or POST
	Method string `json:"method"`
	/// URL: Resource requested by the client
	URL string `json:"url"`
	/// Protocol: HTTP/1, HTTP/2
	Protocol string `json:"protocol"`
	/// Status of request: 200(Successful), 503(Service Unavailable error)
	Status int `json:"status"`
	/// Total bytes transferred to client
	Bytes int `json:"bytes"`
}
