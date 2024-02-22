package sys_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSys(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sys Suite")
}
