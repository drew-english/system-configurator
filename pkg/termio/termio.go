// Handles terminal input/output.
// Based on GH CLI implementation at https://github.com/cli/cli/blob/trunk/pkg/iostreams/iostreams.go
package termio

import (
	"os"

	"golang.org/x/term"
)

var (
	defaultIO = New()
)

func IsTerminal(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}

func IsInteractive() bool {
	return defaultIO.IsInteractive()
}

func StdinIsTerminal() bool {
	return defaultIO.StdinIsTerminal()
}

func StdoutIsTerminal() bool {
	return defaultIO.StdoutIsTerminal()
}

func StderrIsTerminal() bool {
	return defaultIO.StderrIsTerminal()
}

func Style() *style {
	return defaultIO.Style()
}

func WithNeverPrompt(v bool) *IO {
	io := defaultIO.clone()
	io.SetNeverPrompt(v)
	return io
}

func Print(s string) {
	defaultIO.Print(s)
}

func PrintErr(s string) {
	defaultIO.PrintErr(s)
}
