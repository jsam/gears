package datagears_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDatagears(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Datagears Suite")
}
