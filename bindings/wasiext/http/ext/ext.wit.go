// Code generated by wit-bindgen-go. DO NOT EDIT.

// Package ext represents the imported interface "wasiext:http/ext@0.1.0".
package ext

import (
	"github.com/bytecodealliance/wasm-tools-go/cm"
	"github.com/wasmCloud/west/bindings/wasi/http/types"
)

// NewResponseOutparam represents the imported function "new-response-outparam".
//
//	new-response-outparam: func() -> tuple<response-outparam, future-incoming-response>
//
//go:nosplit
func NewResponseOutparam() (result cm.Tuple[types.ResponseOutparam, types.FutureIncomingResponse]) {
	wasmimport_NewResponseOutparam(&result)
	return
}

//go:wasmimport wasiext:http/ext@0.1.0 new-response-outparam
//go:noescape
func wasmimport_NewResponseOutparam(result *cm.Tuple[types.ResponseOutparam, types.FutureIncomingResponse])

// NewIncomingRequest represents the imported function "new-incoming-request".
//
//	new-incoming-request: func(req: outgoing-request) -> incoming-request
//
//go:nosplit
func NewIncomingRequest(req types.OutgoingRequest) (result types.IncomingRequest) {
	req0 := cm.Reinterpret[uint32](req)
	result0 := wasmimport_NewIncomingRequest((uint32)(req0))
	result = cm.Reinterpret[types.IncomingRequest]((uint32)(result0))
	return
}

//go:wasmimport wasiext:http/ext@0.1.0 new-incoming-request
//go:noescape
func wasmimport_NewIncomingRequest(req0 uint32) (result0 uint32)
