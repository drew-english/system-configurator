package run

import (
	"bytes"
	"errors"
	"regexp"

	"github.com/drew-english/system-configurator/pkg/run"
)

type (
	commandStub interface {
		run.RunCmd
		matches(string) bool
		matched() bool
		setMatched([]string)
		pattern() string
		execEffects([]string)
	}

	baseCommandStub struct {
		regex      *regexp.Regexp
		matchedCmd []string
		effects    []commandEffect
	}

	successCommandStub struct {
		*baseCommandStub
		stdout string
	}

	errorCommandStub struct {
		*baseCommandStub
		exitStatus int
		stderr     string
	}
)

func (s *baseCommandStub) matches(line string) bool {
	if !s.matched() && s.regex.MatchString(line) {
		return true
	}

	return false
}

func (s *baseCommandStub) matched() bool {
	return len(s.matchedCmd) > 0
}

func (s *baseCommandStub) pattern() string {
	return s.regex.String()
}

func (s *baseCommandStub) execEffects(args []string) {
	for _, effect := range s.effects {
		effect(args)
	}
}

func (s *baseCommandStub) setMatched(cmd []string) {
	s.matchedCmd = cmd
}

func (s *baseCommandStub) Run() error {
	panic("not implemented, use derived command stubs")
}

func (s *baseCommandStub) Output() ([]byte, error) {
	panic("not implemented, use derived command stubs")
}

func (s *successCommandStub) Run() error {
	return nil
}

func (s *successCommandStub) Output() ([]byte, error) {
	return []byte(s.stdout), nil
}

func (s *errorCommandStub) Run() error {
	return run.CmdError{
		Args:   s.matchedCmd,
		Stderr: bytes.NewBuffer([]byte(s.stderr)),
		Err:    errors.New("generic error"),
	}
}
