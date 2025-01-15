package tests

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEchoVision(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EchoVision Suite")
}
