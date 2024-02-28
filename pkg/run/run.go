package run

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type (
	RunCmd interface {
		Run() error
		Output() ([]byte, error)
	}

	// CmdError provides more visibility into why an exec.Cmd had failed
	CmdError struct {
		Args   []string
		Err    error
		Stderr *bytes.Buffer
	}

	cmdWrap struct {
		*exec.Cmd
	}
)

// Find an executable within path.
// Provides a hook for testing
var Find = exec.LookPath

// Generate a run command.
// Provides a hook for testing
var Command = func(name string, arg ...string) RunCmd {
	return &cmdWrap{exec.Command(name, arg...)}
}

func (c *cmdWrap) Run() error {
	var stderr bytes.Buffer
	c.Stderr = &stderr

	if err := c.Cmd.Run(); err != nil {
		return CmdError{c.Args, err, &stderr}
	}

	return nil
}

func (c *cmdWrap) Output() ([]byte, error) {
	var stderr bytes.Buffer
	c.Stderr = &stderr

	out, err := c.Cmd.Output()
	if err != nil {
		return out, CmdError{c.Args, err, &stderr}
	}

	return out, nil
}

func (e CmdError) Error() string {
	msg := e.Stderr.String()
	if msg != "" && !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	return fmt.Sprintf("%s%s: %s", msg, e.Args[0], e.Err)
}

func (e CmdError) Unwrap() error {
	return e.Err
}
