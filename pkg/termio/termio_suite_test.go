package termio_test

import (
	"os"
	"syscall"
	"testing"

	"github.com/drew-english/system-configurator/pkg/termio"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTermio(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Termio Suite")
}

var _ = Describe("Termio", func() {
	It("exposes top level functions delegating to a default IO instance", func() {
		defaultIO := termio.New()
		Expect(termio.IsInteractive()).To(Equal(defaultIO.IsInteractive()))
		Expect(termio.StdinIsTerminal()).To(Equal(defaultIO.StdinIsTerminal()))
		Expect(termio.StdoutIsTerminal()).To(Equal(defaultIO.StdoutIsTerminal()))
		Expect(termio.StderrIsTerminal()).To(Equal(defaultIO.StderrIsTerminal()))
		Expect(termio.Style()).ToNot(BeNil())
		Expect(termio.Warn).ToNot(BeNil())
		Expect(termio.Error).ToNot(BeNil())
		Expect(termio.Print).ToNot(BeNil())
		Expect(termio.PrintErr).ToNot(BeNil())
	})

	Describe("IsTerminal", func() {
		var file *os.File

		BeforeEach(func() {
			f, err := os.OpenFile("/dev/tty", syscall.O_RDWR, 0)
			Expect(err).ToNot(HaveOccurred())
			file = os.NewFile(uintptr(f.Fd()), "/dev/tty")
		})

		AfterEach(func() {
			file.Close()
		})

		It("returns true", func() {
			Expect(termio.IsTerminal(file)).To(BeTrue())
		})

		Context("when the file is not a terminal", func() {
			BeforeEach(func() {
				file, _ = os.Create("test.txt")
			})

			AfterEach(func() {
				file.Close()
				os.Remove("test.txt")
			})

			It("returns false", func() {
				Expect(termio.IsTerminal(file)).To(BeFalse())
			})
		})
	})

	Describe("WithNeverPrompt", func() {
		It("returns a new IO instance with neverPrompt set to true", func() {
			io := termio.WithNeverPrompt(true)
			Expect(io.IsInteractive()).To(BeFalse())
		})
	})
})
