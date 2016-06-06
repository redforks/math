package decimal_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDecimal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Decimal Suite")
}
