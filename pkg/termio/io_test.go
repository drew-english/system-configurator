package termio_test

import (
	"bytes"
	"os"
	"syscall"

	"github.com/drew-english/system-configurator/pkg/termio"
	"github.com/google/go-cmp/cmp/cmpopts"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var termFile = func(fd uint32) *os.File {
	f, err := os.OpenFile("/dev/tty", syscall.O_RDWR, os.FileMode(fd))
	Expect(err).ToNot(HaveOccurred())
	return os.NewFile(uintptr(f.Fd()), "/dev/tty")
}

var _ = Describe("IO", func() {
	var io *termio.IO

	JustBeforeEach(func() {
		io = termio.New()
	})

	Describe("New", func() {
		It("returns an IO instance", func() {
			var io *termio.IO = termio.New()
			Expect(io).ToNot(BeNil())
			Expect(io).To(BeComparableTo(&termio.IO{
				In:     os.Stdin,
				Out:    os.Stdout,
				ErrOut: os.Stderr,
			}, cmpopts.IgnoreUnexported(termio.IO{}, os.File{})))
		})
	})

	Describe("NewWithConfig", func() {
		It("returns an IO instance", func() {
			var cfg *termio.Config = &termio.Config{
				ColorDisabled:    false,
				Color256Enabled:  false,
				TrueColorEnabled: false,
			}
			var io *termio.IO = termio.NewWithConfig(cfg)
			Expect(io).ToNot(BeNil())
			Expect(io).To(BeComparableTo(&termio.IO{
				In:     os.Stdin,
				Out:    os.Stdout,
				ErrOut: os.Stderr,
			}, cmpopts.IgnoreUnexported(termio.IO{}, os.File{})))
		})
	})

	Describe("Print", func() {
		buf := bytes.NewBuffer([]byte{})

		JustBeforeEach(func() {
			io.Out = buf
		})

		It("prints to the output", func() {
			io.Print("testing1234")
			Expect(buf.String()).To(Equal("testing1234"))
		})
	})

	Describe("Printf", func() {
		buf := bytes.NewBuffer([]byte{})

		JustBeforeEach(func() {
			io.Out = buf
		})

		It("prints to the output", func() {
			io.Printf("testing1234: %s", "foo")
			Expect(buf.String()).To(Equal("testing1234: foo"))
		})
	})

	Describe("PrintErr", func() {
		buf := bytes.NewBuffer([]byte{})

		JustBeforeEach(func() {
			io.ErrOut = buf
		})

		It("prints to the error output", func() {
			io.PrintErr("testing1234")
			Expect(buf.String()).To(Equal("testing1234"))
		})
	})

	Describe("Warn", func() {
		buf := bytes.NewBuffer([]byte{})

		JustBeforeEach(func() {
			io.ErrOut = buf
		})

		It("prints to the error output", func() {
			io.Warn("testing1234")
			Expect(buf.String()).To(Equal(io.Style().Yellow("WARNING: ") + "testing1234"))
		})
	})

	Describe("Warnf", func() {
		buf := bytes.NewBuffer([]byte{})

		JustBeforeEach(func() {
			io.ErrOut = buf
		})

		It("prints to the error output", func() {
			io.Warnf("testing1234: %s", "foo")
			Expect(buf.String()).To(Equal(io.Style().Yellow("WARNING: ") + "testing1234: foo"))
		})
	})

	Describe("Error", func() {
		buf := bytes.NewBuffer([]byte{})

		JustBeforeEach(func() {
			io.ErrOut = buf
		})

		It("prints to the error output", func() {
			io.Error("testing1234")
			Expect(buf.String()).To(Equal(io.Style().Red("ERROR: ") + "testing1234"))
		})
	})

	Describe("Errorf", func() {
		buf := bytes.NewBuffer([]byte{})

		JustBeforeEach(func() {
			io.ErrOut = buf
		})

		It("prints to the error output", func() {
			io.Errorf("testing1234: %s", "foo")
			Expect(buf.String()).To(Equal(io.Style().Red("ERROR: ") + "testing1234: foo"))
		})
	})

	Describe("IsInteractive", func() {
		originalStdout := os.Stdout
		originalStdin := os.Stdin

		subject := func() bool {
			return io.IsInteractive()
		}

		BeforeEach(func() {
			os.Stdout = termFile(0)
			os.Stdin = termFile(1)
		})

		JustBeforeEach(func() {
			io.SetNeverPrompt(false)
		})

		AfterEach(func() {
			os.Stdout.Close()
			os.Stdin.Close()
			os.Stdout = originalStdout
			os.Stdin = originalStdin
		})

		It("returns true when stdin and stdout are terminals", func() {
			Expect(subject()).To(BeTrue())
		})

		Context("when neverPrompt is true", func() {
			JustBeforeEach(func() {
				io.SetNeverPrompt(true)
			})

			It("returns false", func() {
				Expect(subject()).To(BeFalse())
			})
		})

		Context("when stdin is not terminal", func() {
			BeforeEach(func() {
				os.Stdin = nil
			})

			It("returns false", func() {
				Expect(subject()).To(BeFalse())
			})
		})

		Context("when stdout is not terminal", func() {
			BeforeEach(func() {
				os.Stdout = nil
			})

			It("returns false", func() {
				Expect(subject()).To(BeFalse())
			})
		})
	})

	Describe("StdinIsTerminal", func() {
		originalStdin := os.Stdin

		subject := func() bool {
			return io.StdinIsTerminal()
		}

		BeforeEach(func() {
			os.Stdin = termFile(1)
		})

		AfterEach(func() {
			os.Stdin.Close()
			os.Stdin = originalStdin
		})

		It("returns true", func() {
			Expect(subject()).To(BeTrue())
		})

		Context("when stdin is not a terminal", func() {
			BeforeEach(func() {
				os.Stdin = nil
			})

			It("returns false", func() {
				Expect(subject()).To(BeFalse())
			})
		})
	})

	Describe("StdoutIsTerminal", func() {
		originalStdout := os.Stdout

		subject := func() bool {
			return io.StdoutIsTerminal()
		}

		BeforeEach(func() {
			os.Stdout = termFile(1)
		})

		AfterEach(func() {
			os.Stdout.Close()
			os.Stdout = originalStdout
		})

		It("returns true", func() {
			Expect(subject()).To(BeTrue())
		})

		Context("when stdout is not a terminal", func() {
			BeforeEach(func() {
				os.Stdout = nil
			})

			It("returns false", func() {
				Expect(subject()).To(BeFalse())
			})
		})
	})

	Describe("StderrIsTerminal", func() {
		originalStderr := os.Stderr

		subject := func() bool {
			return io.StderrIsTerminal()
		}

		BeforeEach(func() {
			os.Stderr = termFile(1)
		})

		AfterEach(func() {
			os.Stderr.Close()
			os.Stderr = originalStderr
		})

		It("returns true", func() {
			Expect(subject()).To(BeTrue())
		})

		Context("when stderr is not a terminal", func() {
			BeforeEach(func() {
				os.Stderr = nil
			})

			It("returns false", func() {
				Expect(subject()).To(BeFalse())
			})
		})
	})
})
