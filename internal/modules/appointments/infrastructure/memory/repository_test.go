package memory_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"src/internal/modules/appointments/domain"
	. "src/internal/modules/appointments/infrastructure/memory"
)

var _ = Describe("InMemoryAppointmentRepository", func() {
	It("creates and lists appointments", func() {
		repo := NewInMemoryAppointmentRepository()
		slot, _ := domain.NewTimeSlot(time.Now().Add(time.Hour), time.Now().Add(90*time.Minute))
		a := domain.Appointment{ID: "a1", PetID: "p1", SpecialistID: "s1", StartTime: slot.Start, EndTime: slot.End, Status: domain.AppointmentStatusBooked}
		_, err := repo.Create(a)
		Expect(err).ToNot(HaveOccurred())

		items, err := repo.List()
		Expect(err).ToNot(HaveOccurred())
		Expect(items).To(HaveLen(1))
		Expect(items[0].ID).To(Equal("a1"))
	})
})
