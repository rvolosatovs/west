// Code generated by wit-bindgen-go. DO NOT EDIT.

// Package outgoinghandler represents the imported interface "wasi:http/outgoing-handler@0.2.1".
//
// This interface defines a handler of outgoing HTTP Requests. It should be
// imported by components which wish to make HTTP Requests.
package outgoinghandler

import (
	"github.com/bytecodealliance/wasm-tools-go/cm"
	"github.com/wasmCloud/wadge/bindings/wasi/http/types"
)

// Handle represents the imported function "handle".
//
// This function is invoked with an outgoing HTTP Request, and it returns
// a resource `future-incoming-response` which represents an HTTP Response
// which may arrive in the future.
//
// The `options` argument accepts optional parameters for the HTTP
// protocol's transport layer.
//
// This function may return an error if the `outgoing-request` is invalid
// or not allowed to be made. Otherwise, protocol errors are reported
// through the `future-incoming-response`.
//
//	handle: func(request: outgoing-request, options: option<request-options>) -> result<future-incoming-response,
//	error-code>
//
//go:nosplit
func Handle(request types.OutgoingRequest, options cm.Option[types.RequestOptions]) (result cm.Result[ErrorCodeShape, types.FutureIncomingResponse, types.ErrorCode]) {
	request0 := cm.Reinterpret[uint32](request)
	options0, options1 := lower_OptionRequestOptions(options)
	wasmimport_Handle((uint32)(request0), (uint32)(options0), (uint32)(options1), &result)
	return
}

//go:wasmimport wasi:http/outgoing-handler@0.2.1 handle
//go:noescape
func wasmimport_Handle(request0 uint32, options0 uint32, options1 uint32, result *cm.Result[ErrorCodeShape, types.FutureIncomingResponse, types.ErrorCode])
