// Code generated by wit-bindgen-go. DO NOT EDIT.

package tcp

import (
	"github.com/bytecodealliance/wasm-tools-go/cm"
	"github.com/wasmCloud/wadge/bindings/wasi/io/streams"
	"github.com/wasmCloud/wadge/bindings/wasi/sockets/network"
	"unsafe"
)

// TupleTCPSocketInputStreamOutputStreamShape is used for storage in variant or result types.
type TupleTCPSocketInputStreamOutputStreamShape struct {
	shape [unsafe.Sizeof(cm.Tuple3[TCPSocket, streams.InputStream, streams.OutputStream]{})]byte
}

// TupleInputStreamOutputStreamShape is used for storage in variant or result types.
type TupleInputStreamOutputStreamShape struct {
	shape [unsafe.Sizeof(cm.Tuple[streams.InputStream, streams.OutputStream]{})]byte
}

// IPSocketAddressShape is used for storage in variant or result types.
type IPSocketAddressShape struct {
	shape [unsafe.Sizeof(network.IPSocketAddress{})]byte
}

func lower_IPv4Address(v network.IPv4Address) (f0 uint32, f1 uint32, f2 uint32, f3 uint32) {
	f0 = (uint32)(v[0])
	f1 = (uint32)(v[1])
	f2 = (uint32)(v[2])
	f3 = (uint32)(v[3])
	return
}

func lower_IPv4SocketAddress(v network.IPv4SocketAddress) (f0 uint32, f1 uint32, f2 uint32, f3 uint32, f4 uint32) {
	f0 = (uint32)(v.Port)
	f1, f2, f3, f4 = lower_IPv4Address(v.Address)
	return
}

func lower_IPv6Address(v network.IPv6Address) (f0 uint32, f1 uint32, f2 uint32, f3 uint32, f4 uint32, f5 uint32, f6 uint32, f7 uint32) {
	f0 = (uint32)(v[0])
	f1 = (uint32)(v[1])
	f2 = (uint32)(v[2])
	f3 = (uint32)(v[3])
	f4 = (uint32)(v[4])
	f5 = (uint32)(v[5])
	f6 = (uint32)(v[6])
	f7 = (uint32)(v[7])
	return
}

func lower_IPv6SocketAddress(v network.IPv6SocketAddress) (f0 uint32, f1 uint32, f2 uint32, f3 uint32, f4 uint32, f5 uint32, f6 uint32, f7 uint32, f8 uint32, f9 uint32, f10 uint32) {
	f0 = (uint32)(v.Port)
	f1 = (uint32)(v.FlowInfo)
	f2, f3, f4, f5, f6, f7, f8, f9 = lower_IPv6Address(v.Address)
	f10 = (uint32)(v.ScopeID)
	return
}

func lower_IPSocketAddress(v network.IPSocketAddress) (f0 uint32, f1 uint32, f2 uint32, f3 uint32, f4 uint32, f5 uint32, f6 uint32, f7 uint32, f8 uint32, f9 uint32, f10 uint32, f11 uint32) {
	f0 = (uint32)(v.Tag())
	switch f0 {
	case 0: // ipv4
		v1, v2, v3, v4, v5 := lower_IPv4SocketAddress(*v.IPv4())
		f1 = (uint32)(v1)
		f2 = (uint32)(v2)
		f3 = (uint32)(v3)
		f4 = (uint32)(v4)
		f5 = (uint32)(v5)
	case 1: // ipv6
		v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11 := lower_IPv6SocketAddress(*v.IPv6())
		f1 = (uint32)(v1)
		f2 = (uint32)(v2)
		f3 = (uint32)(v3)
		f4 = (uint32)(v4)
		f5 = (uint32)(v5)
		f6 = (uint32)(v6)
		f7 = (uint32)(v7)
		f8 = (uint32)(v8)
		f9 = (uint32)(v9)
		f10 = (uint32)(v10)
		f11 = (uint32)(v11)
	}
	return
}
