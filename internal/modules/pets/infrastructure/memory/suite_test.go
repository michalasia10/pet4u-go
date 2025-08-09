package memory_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPetsInfraMemorySuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pets Infra Memory Suite")
}
