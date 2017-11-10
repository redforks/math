package rate_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rate Suite")
}
