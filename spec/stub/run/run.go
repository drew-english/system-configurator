// Package to stub run commands.
// Based on GH CLI implementation at https://github.com/cli/cli/blob/trunk/internal/run/stub.go
package run

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/drew-english/system-configurator/lib/run"
)

type (
	commandStubManager struct {
		stubs []commandStub
	}

	commandEffect func([]string)
)

// StubCommand installs a catch-all for exec.Command. It returns a tear down function
// to be called at the end of tests to clean up the stubs and ensure they were called.
func StubCommand() (*commandStubManager, func(testing.TB)) {
	cs := &commandStubManager{}

	restoreOriginalCommand := registerCommandHook(cs)
	return cs, func(t testing.TB) {
		defer restoreOriginalCommand()

		var unmatched []string
		for _, stub := range cs.stubs {
			if stub.matched() {
				continue
			}

			unmatched = append(unmatched, stub.pattern())
		}

		if len(unmatched) == 0 {
			cs.stubs = []commandStub{}
			return
		}

		t.Helper()
		t.Errorf("unmatched stubs (%d): %s", len(unmatched), strings.Join(unmatched, "\n"))
	}
}

func StubFind(pattern string, errResult error) func() {
	originalFind := run.Find

	run.Find = func(name string) (string, error) {
		re := regexp.MustCompile(pattern)
		if re.MatchString(name) {
			return "", errResult
		}

		panic(fmt.Sprintf("Find not stubbed for `%s`", name))
	}

	return func() {
		run.Find = originalFind
	}
}

func registerCommandHook(cs *commandStubManager) func() {
	originalCommand := run.Command

	run.Command = func(name string, arg ...string) run.RunCmd {
		args := append([]string{name}, arg...)
		stub := cs.find(args)
		if stub == nil {
			panic(fmt.Sprintf("no exec stub for `%s`", strings.Join(args, " ")))
		}

		stub.execEffects(args)
		stub.setMatched(append([]string{name}, arg...))
		return stub
	}

	return func() {
		run.Command = originalCommand
	}
}

// Register a command stub that is successful. Use effects to perform actions when the command is called.
func (cs *commandStubManager) Register(pattern string, output string, effects ...commandEffect) {
	if len(pattern) < 1 {
		panic("cannot use empty regexp pattern")
	}

	cs.stubs = append(cs.stubs, &successCommandStub{
		baseCommandStub: &baseCommandStub{
			regex:   regexp.MustCompile(pattern),
			effects: effects,
		},
		stdout: output,
	})
}

// RegisterError registers a command stub that returns an error. Use effects to perform actions when the command is called.
func (cs *commandStubManager) RegisterError(pattern string, exitStatus int, stderr string, effects ...commandEffect) {
	if len(pattern) < 1 {
		panic("cannot use empty regexp pattern")
	}

	cs.stubs = append(cs.stubs, &errorCommandStub{
		baseCommandStub: &baseCommandStub{
			regex:   regexp.MustCompile(pattern),
			effects: effects,
		},
		exitStatus: exitStatus,
		stderr:     stderr,
	})
}

func (cs *commandStubManager) find(args []string) commandStub {
	line := strings.Join(args, " ")
	for _, stub := range cs.stubs {
		if !stub.matched() && stub.matches(line) {
			return stub
		}
	}

	return nil
}
