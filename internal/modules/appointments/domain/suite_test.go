package domain_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAppointmentsDomain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Appointments Domain Suite")
}
