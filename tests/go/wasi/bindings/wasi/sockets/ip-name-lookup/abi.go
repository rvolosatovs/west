// Code generated by wit-bindgen-go. DO NOT EDIT.

package ipnamelookup

import (
	"github.com/bytecodealliance/wasm-tools-go/cm"
	"github.com/wasmCloud/wadge/tests/go/wasi/bindings/wasi/sockets/network"
	"unsafe"
)

// OptionIPAddressShape is used for storage in variant or result types.
type OptionIPAddressShape struct {
	shape [unsafe.Sizeof(cm.Option[network.IPAddress]{})]byte
}
