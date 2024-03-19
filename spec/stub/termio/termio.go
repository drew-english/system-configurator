package termio

import (
	"bytes"

	"github.com/drew-english/system-configurator/pkg/termio"
)

func CaptureTermOut(action func()) (string, string) {
	stdout, stderr := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	tmpTerm := termio.New()
	tmpTerm.Out = stdout
	tmpTerm.ErrOut = stderr

	originalTerm := termio.DefaultIO
	termio.DefaultIO = tmpTerm
	action()
	termio.DefaultIO = originalTerm

	return stdout.String(), stderr.String()
}
