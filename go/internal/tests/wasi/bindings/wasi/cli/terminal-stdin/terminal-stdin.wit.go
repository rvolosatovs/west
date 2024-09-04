// Code generated by wit-bindgen-go. DO NOT EDIT.

// Package terminalstdin represents the imported interface "wasi:cli/terminal-stdin@0.2.0".
//
// An interface providing an optional `terminal-input` for stdin as a
// link-time authority.
package terminalstdin

import (
	terminalinput "github.com/rvolosatovs/west/go/internal/tests/wasi/bindings/wasi/cli/terminal-input"
	"github.com/ydnar/wasm-tools-go/cm"
)

// GetTerminalStdin represents the imported function "get-terminal-stdin".
//
// If stdin is connected to a terminal, return a `terminal-input` handle
// allowing further interaction with it.
//
//	get-terminal-stdin: func() -> option<terminal-input>
//
//go:nosplit
func GetTerminalStdin() (result cm.Option[terminalinput.TerminalInput]) {
	wasmimport_GetTerminalStdin(&result)
	return
}

//go:wasmimport wasi:cli/terminal-stdin@0.2.0 get-terminal-stdin
//go:noescape
func wasmimport_GetTerminalStdin(result *cm.Option[terminalinput.TerminalInput])
