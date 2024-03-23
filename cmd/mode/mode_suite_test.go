package mode_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMode(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mode Suite")
}
