// Code generated by wit-bindgen-go. DO NOT EDIT.

package timezone

import (
	wallclock "github.com/rvolosatovs/west/bindings/wasi/clocks/wall-clock"
)

func lower_DateTime(v wallclock.DateTime) (f0 uint64, f1 uint32) {
	f0 = (uint64)(v.Seconds)
	f1 = (uint32)(v.Nanoseconds)
	return
}
