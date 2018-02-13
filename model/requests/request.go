package requests

import "net/http"

// Request defines everything needed for a struct to be a request type
// these are essentially an experiment to see if I like it
type Request interface {
	Path() string
	FromRequest(*http.Request) error
	Verify() error
}

// VerifyRequest exists essentially so that the compiler forces all request types to conform to the Request interface, it's entirely superfluous
func VerifyRequest(req Request) error {
	return req.Verify()
}
