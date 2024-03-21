package termio

import (
	"fmt"
	"io"
	"os"
)

type IO struct {
	cfg *Config

	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer

	neverPrompt bool
}

func New() *IO {
	return NewWithConfig(ConfigFromEnv())
}

func NewWithConfig(cfg *Config) *IO {
	return &IO{
		cfg:    cfg,
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
}

func (io *IO) Print(s string) {
	fmt.Fprint(io.Out, s)
}

func (io *IO) Printf(s string, args ...any) {
	fmt.Fprintf(io.Out, s, args...)
}

func (io *IO) PrintErr(s string) {
	fmt.Fprint(io.ErrOut, s)
}

func (io *IO) Warn(s string) {
	io.PrintErr(io.Style().Yellow("WARNING: ") + s)
}

func (io *IO) Warnf(s string, args ...any) {
	io.Warn(fmt.Sprintf(s, args...))
}

func (io *IO) Error(s string) {
	io.PrintErr(io.Style().Red("ERROR: ") + s)
}

func (io *IO) Errorf(s string, args ...any) {
	io.Error(fmt.Sprintf(s, args...))
}

func (io *IO) IsInteractive() bool {
	if io.neverPrompt {
		return false
	}

	return io.StdinIsTerminal() && io.StdoutIsTerminal()
}

func (io *IO) StdinIsTerminal() bool {
	if f, ok := io.In.(*os.File); ok {
		return IsTerminal(f)
	}

	return false
}

func (io *IO) StdoutIsTerminal() bool {
	if f, ok := io.Out.(*os.File); ok {
		return IsTerminal(f)
	}

	return false
}

func (io *IO) StderrIsTerminal() bool {
	if f, ok := io.ErrOut.(*os.File); ok {
		return IsTerminal(f)
	}

	return false
}

func (io *IO) SetNeverPrompt(v bool) {
	io.neverPrompt = v
}

func (io *IO) Style() *style {
	return NewStyle(!io.cfg.ColorDisabled, io.cfg.Color256Enabled, io.cfg.TrueColorEnabled)
}

func (io *IO) clone() *IO {
	return &IO{
		cfg:         io.cfg,
		In:          io.In,
		Out:         io.Out,
		ErrOut:      io.ErrOut,
		neverPrompt: io.neverPrompt,
	}
}
