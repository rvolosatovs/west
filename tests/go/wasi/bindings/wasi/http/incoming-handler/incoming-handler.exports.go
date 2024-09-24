// Code generated by wit-bindgen-go. DO NOT EDIT.

package incominghandler

import (
	"github.com/wasmCloud/wadge/tests/go/wasi/bindings/wasi/http/types"
)

// Exports represents the caller-defined exports from "wasi:http/incoming-handler@0.2.1".
var Exports struct {
	// Handle represents the caller-defined, exported function "handle".
	//
	// This function is invoked with an incoming HTTP Request, and a resource
	// `response-outparam` which provides the capability to reply with an HTTP
	// Response. The response is sent by calling the `response-outparam.set`
	// method, which allows execution to continue after the response has been
	// sent. This enables both streaming to the response body, and performing other
	// work.
	//
	// The implementor of this function must write a response to the
	// `response-outparam` before returning, or else the caller will respond
	// with an error on its behalf.
	//
	//	handle: func(request: incoming-request, response-out: response-outparam)
	Handle func(request types.IncomingRequest, responseOut types.ResponseOutparam)
}
