// Code generated by wit-bindgen-go. DO NOT EDIT.

// Package timezone represents the imported interface "wasi:clocks/timezone@0.2.1".
package timezone

import (
	wallclock "github.com/wasmCloud/wadge/tests/go/wasi/bindings/wasi/clocks/wall-clock"
)

// TimezoneDisplay represents the record "wasi:clocks/timezone@0.2.1#timezone-display".
//
// Information useful for displaying the timezone of a specific `datetime`.
//
// This information may vary within a single `timezone` to reflect daylight
// saving time adjustments.
//
//	record timezone-display {
//		utc-offset: s32,
//		name: string,
//		in-daylight-saving-time: bool,
//	}
type TimezoneDisplay struct {
	// The number of seconds difference between UTC time and the local
	// time of the timezone.
	//
	// The returned value will always be less than 86400 which is the
	// number of seconds in a day (24*60*60).
	//
	// In implementations that do not expose an actual time zone, this
	// should return 0.
	UtcOffset int32

	// The abbreviated name of the timezone to display to a user. The name
	// `UTC` indicates Coordinated Universal Time. Otherwise, this should
	// reference local standards for the name of the time zone.
	//
	// In implementations that do not expose an actual time zone, this
	// should be the string `UTC`.
	//
	// In time zones that do not have an applicable name, a formatted
	// representation of the UTC offset may be returned, such as `-04:00`.
	Name string

	// Whether daylight saving time is active.
	//
	// In implementations that do not expose an actual time zone, this
	// should return false.
	InDaylightSavingTime bool
}

// Display represents the imported function "display".
//
// Return information needed to display the given `datetime`. This includes
// the UTC offset, the time zone name, and a flag indicating whether
// daylight saving time is active.
//
// If the timezone cannot be determined for the given `datetime`, return a
// `timezone-display` for `UTC` with a `utc-offset` of 0 and no daylight
// saving time.
//
//	display: func(when: datetime) -> timezone-display
//
//go:nosplit
func Display(when wallclock.DateTime) (result TimezoneDisplay) {
	when0, when1 := lower_DateTime(when)
	wasmimport_Display((uint64)(when0), (uint32)(when1), &result)
	return
}

//go:wasmimport wasi:clocks/timezone@0.2.1 display
//go:noescape
func wasmimport_Display(when0 uint64, when1 uint32, result *TimezoneDisplay)

// UtcOffset represents the imported function "utc-offset".
//
// The same as `display`, but only return the UTC offset.
//
//	utc-offset: func(when: datetime) -> s32
//
//go:nosplit
func UtcOffset(when wallclock.DateTime) (result int32) {
	when0, when1 := lower_DateTime(when)
	result0 := wasmimport_UtcOffset((uint64)(when0), (uint32)(when1))
	result = (int32)((uint32)(result0))
	return
}

//go:wasmimport wasi:clocks/timezone@0.2.1 utc-offset
//go:noescape
func wasmimport_UtcOffset(when0 uint64, when1 uint32) (result0 uint32)
