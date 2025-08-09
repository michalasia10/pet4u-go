package application_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAppointmentsApplication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Appointments Application Suite")
}
