package tests

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEchoAnalyzer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Echoanalyzer Suite")
}

var _ = BeforeEach(func() {
	// Expect(m.Up()).To(Succeed())
})

var _ = AfterEach(func() {
	// Expect(m.Down()).To(Succeed())
})

var _ = BeforeSuite(func() {
})

var _ = AfterSuite(func() {
})
