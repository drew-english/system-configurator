package alternate_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAlternate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Alternate Suite")
}
