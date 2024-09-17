// Code generated by wit-bindgen-go. DO NOT EDIT.

package types

import (
	"github.com/bytecodealliance/wasm-tools-go/cm"
	wallclock "github.com/rvolosatovs/west/examples/go/http/bindings/wasi/clocks/wall-clock"
	"unsafe"
)

// DateTimeShape is used for storage in variant or result types.
type DateTimeShape struct {
	shape [unsafe.Sizeof(wallclock.DateTime{})]byte
}

// MetadataHashValueShape is used for storage in variant or result types.
type MetadataHashValueShape struct {
	shape [unsafe.Sizeof(MetadataHashValue{})]byte
}

// TupleListU8BoolShape is used for storage in variant or result types.
type TupleListU8BoolShape struct {
	shape [unsafe.Sizeof(cm.Tuple[cm.List[uint8], bool]{})]byte
}

func lower_DateTime(v wallclock.DateTime) (f0 uint64, f1 uint32) {
	f0 = (uint64)(v.Seconds)
	f1 = (uint32)(v.Nanoseconds)
	return
}

func lower_NewTimestamp(v NewTimestamp) (f0 uint32, f1 uint64, f2 uint32) {
	f0 = (uint32)(v.Tag())
	switch f0 {
	case 2: // timestamp
		v1, v2 := lower_DateTime(*v.Timestamp())
		f1 = (uint64)(v1)
		f2 = (uint32)(v2)
	}
	return
}

// DescriptorStatShape is used for storage in variant or result types.
type DescriptorStatShape struct {
	shape [unsafe.Sizeof(DescriptorStat{})]byte
}

// OptionDirectoryEntryShape is used for storage in variant or result types.
type OptionDirectoryEntryShape struct {
	shape [unsafe.Sizeof(cm.Option[DirectoryEntry]{})]byte
}
