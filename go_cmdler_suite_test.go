package cmdler_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGoCmdler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoCmdler Suite")
}
