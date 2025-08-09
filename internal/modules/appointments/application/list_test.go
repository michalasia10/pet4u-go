package application_test

import (
	app "src/internal/modules/appointments/application"
	"src/internal/modules/appointments/domain"
	mem "src/internal/modules/appointments/infrastructure/memory"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListUseCase", func() {
	var (
		repo *mem.InMemoryAppointmentRepository
		uc   *app.ListUseCase
	)

	BeforeEach(func() {
		repo = mem.NewInMemoryAppointmentRepository()
		uc = app.NewListUseCase(repo)
	})

	It("returns appointments from repository", func() {
		slot, _ := domain.NewTimeSlot(time.Now().Add(time.Hour), time.Now().Add(2*time.Hour))
		repo.Create(domain.Appointment{ID: "a1", PetID: "p1", SpecialistID: "s1", StartTime: slot.Start, EndTime: slot.End, Status: domain.AppointmentStatusBooked})

		resp, err := uc.Execute()
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.Appointments).To(HaveLen(1))
		Expect(resp.Appointments[0].ID).To(Equal("a1"))
	})
})
