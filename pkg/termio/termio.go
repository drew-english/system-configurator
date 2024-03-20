// Handles terminal input/output.
// Based on GH CLI implementation at https://github.com/cli/cli/blob/trunk/pkg/iostreams/iostreams.go
package termio

import (
	"os"

	"golang.org/x/term"
)

var (
	DefaultIO = New()
)

func IsTerminal(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}

func IsInteractive() bool {
	return DefaultIO.IsInteractive()
}

func StdinIsTerminal() bool {
	return DefaultIO.StdinIsTerminal()
}

func StdoutIsTerminal() bool {
	return DefaultIO.StdoutIsTerminal()
}

func StderrIsTerminal() bool {
	return DefaultIO.StderrIsTerminal()
}

func Style() *style {
	return DefaultIO.Style()
}

func WithNeverPrompt(v bool) *IO {
	io := DefaultIO.clone()
	io.SetNeverPrompt(v)
	return io
}

func Print(s string) {
	DefaultIO.Print(s)
}

func PrintErr(s string) {
	DefaultIO.PrintErr(s)
}

func Warn(s string) {
	DefaultIO.Warn(s)
}

func Error(s string) {
	DefaultIO.Error(s)
}

func Warnf(s string, args ...any) {
	DefaultIO.Warnf(s, args...)
}

func Errorf(s string, args ...any) {
	DefaultIO.Errorf(s, args...)
}
