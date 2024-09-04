//go:generate cargo build -p west-sys --release
//go:generate go run github.com/ydnar/wasm-tools-go/cmd/wit-bindgen-go@v0.1.5 generate -w imports -o bindings ../wit

package west

// #cgo darwin,amd64  LDFLAGS: -L${SRCDIR}/../lib/x86_64-darwin -lwest
// #cgo darwin,arm64  LDFLAGS: -L${SRCDIR}/../lib/aarch64-darwin -lwest
// #cgo linux,amd64   LDFLAGS: -L${SRCDIR}/../lib/x86_64-linux -lwest
// #cgo linux,arm64   LDFLAGS: -L${SRCDIR}/../lib/aarch64-linux -lwest
// #cgo windows,amd64 LDFLAGS: -L${SRCDIR}/../lib/x86_64-windows -lwest
// #include "../include/west.h"
// #include <stdlib.h>
import "C"

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"
)

var (
	errorHandlerMu sync.RWMutex
	errorHandler   atomic.Value

	instanceMu sync.RWMutex
	instance   *Instance
)

func init() {
	errorHandler.Store(func(err error) {
		log.Fatalf("failed to call instance: %s", err)
	})
}

func setErrorHandler(f func(error)) func(error) {
	return errorHandler.Swap(f).(func(error))
}

func SetErrorHandler(f func(error)) func(error) {
	errorHandlerMu.Lock()
	defer errorHandlerMu.Unlock()
	return setErrorHandler(f)
}

func WithErrorHandler(handler func(error), f func()) {
	errorHandlerMu.Lock()
	defer errorHandlerMu.Unlock()

	handler = setErrorHandler(handler)
	defer setErrorHandler(handler)

	f()
}

func CurrentErrorHandler() func(error) {
	return errorHandler.Load().(func(error))
}

func WithCurrentErrorHandler(f func(func(error))) {
	errorHandlerMu.RLock()
	defer errorHandlerMu.RUnlock()
	f(CurrentErrorHandler())
}

func setInstance(i *Instance) *Instance {
	i, instance = instance, i
	return i
}

func SetInstance(i *Instance) *Instance {
	instanceMu.Lock()
	defer instanceMu.Unlock()

	return setInstance(i)
}

func WithInstance(i *Instance, f func()) {
	instanceMu.Lock()
	defer instanceMu.Unlock()

	i = setInstance(i)
	defer setInstance(i)

	f()
}

func withCurrentInstance[T any](f func(*Instance) T, handleErr func(error)) T {
	instanceMu.RLock()
	defer instanceMu.RUnlock()
	if instance == nil {
		instanceMu.RUnlock()
		func() {
			instanceMu.Lock()
			defer instanceMu.Unlock()

			if instance == nil {
				i, err := NewInstance(nil)
				if err != nil {
					handleErr(err)
				}
				instance = i
			}
		}()
		instanceMu.RLock()
	}
	return f(instance)
}

func WithCurrentInstance[T any](f func(*Instance) T) T {
	return withCurrentInstance(f, func(err error) {
		log.Fatal(err)
	})
}

func RunTest(t *testing.T, f func()) {
	WithErrorHandler(
		func(err error) {
			t.Fatalf("failed to call instance: %s", err)
		},
		func() {
			withCurrentInstance(
				func(_ *Instance) struct{} {
					f()
					return struct{}{}
				},
				func(err error) {
					t.Fatal(err)
				},
			)
		},
	)
}

type Instance struct {
	ptr unsafe.Pointer
}

type Config struct {
	Wasm []byte
}

func NewInstance(conf *Config) (*Instance, error) {
	var pinner runtime.Pinner
	defer pinner.Unpin()

	c := C.default_config()
	if conf != nil {
		if len(conf.Wasm) > 0 {
			ptr := unsafe.SliceData(conf.Wasm)
			pinner.Pin(ptr)

			c.wasm = C.List_u8{
				ptr: (*C.uchar)(ptr),
				len: C.ulong(len(conf.Wasm)),
			}
		}
	}
	ptr := C.instance_new(c)
	if ptr == nil {
		n := C.error_len()
		buf := make([]C.char, n)
		if n = C.error_take(unsafe.SliceData(buf), n); n > 0 {
			err := errors.New(C.GoStringN(unsafe.SliceData(buf), C.int(n)))
			return nil, fmt.Errorf("failed to create an instance: %w", err)
		} else {
			return nil, errors.New("failed to create an instance")
		}
	}
	instance := &Instance{ptr: ptr}
	runtime.SetFinalizer(instance, func(instance *Instance) {
		C.instance_free(instance.ptr)
	})
	return instance, nil
}

func (i Instance) Call(instance string, name string, args ...unsafe.Pointer) error {
	instanceC := C.CString(instance)
	defer C.free(unsafe.Pointer(instanceC))
	nameC := C.CString(name)
	defer C.free(unsafe.Pointer(nameC))

	if !C.instance_call(i.ptr, instanceC, nameC, unsafe.SliceData(args)) {
		n := C.error_len()
		buf := make([]C.char, n)
		if n = C.error_take(unsafe.SliceData(buf), n); n > 0 {
			err := errors.New(C.GoStringN(unsafe.SliceData(buf), C.int(n)))
			return fmt.Errorf("failed to call function on an instance: %w", err)
		} else {
			return errors.New("failed to call function on an instance")
		}
	}
	return nil
}
