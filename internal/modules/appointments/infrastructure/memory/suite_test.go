package memory_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAppointmentsInfraMemorySuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Appointments Infra Memory Suite")
}
