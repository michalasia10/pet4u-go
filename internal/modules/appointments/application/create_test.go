package application_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	app "src/internal/modules/appointments/application"
	"src/internal/modules/appointments/domain"
	mem "src/internal/modules/appointments/infrastructure/memory"
	cport "src/internal/modules/shared/domain/clock"
	idport "src/internal/modules/shared/domain/idgen"
	txport "src/internal/modules/shared/domain/tx"
	cimpl "src/internal/modules/shared/infrastructure/clock"
	idimpl "src/internal/modules/shared/infrastructure/idgen"
	tximpl "src/internal/modules/shared/infrastructure/tx"
	"time"
)

var _ = Describe("CreateUseCase", func() {
	var (
		repo         *mem.InMemoryAppointmentRepository
		idGen        idport.Port
		clock        cport.Port
		tx           txport.Manager
		availability domain.AvailabilityService
		policy       domain.BookingPolicy
		uc           *app.CreateUseCase
	)

	BeforeEach(func() {
		repo = mem.NewInMemoryAppointmentRepository()
		idGen = idimpl.NewTimeIDGen()
		clock = cimpl.NewSystemClock()
		tx = tximpl.NewNoopManager()
		availability = mem.NewInMemoryAvailabilityService(repo)
		policy = mem.NewInMemoryBookingPolicy()
		uc = app.NewCreateUseCase(repo, idGen, clock, tx, availability, policy)
	})

	It("creates an appointment when slot is available and policy allows", func() {
		start := time.Now().Add(2 * time.Hour)
		end := start.Add(30 * time.Minute)
		resp, err := uc.Execute(app.CreateRequest{
			PetID:        "pet_1",
			SpecialistID: "spec_1",
			StartTime:    start,
			EndTime:      end,
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.Appointment.ID).ToNot(BeEmpty())
		Expect(resp.Appointment.PetID).To(Equal("pet_1"))

		t := GinkgoT()
		require.WithinDuration(t, start, resp.Appointment.StartTime, time.Second)
		require.WithinDuration(t, end, resp.Appointment.EndTime, time.Second)
	})

	It("rejects overlapping appointments via availability service", func() {
		// first booking
		start := time.Now().Add(2 * time.Hour)
		end := start.Add(30 * time.Minute)
		_, err := uc.Execute(app.CreateRequest{PetID: "pet_1", SpecialistID: "spec_1", StartTime: start, EndTime: end})
		Expect(err).ToNot(HaveOccurred())

		// overlapping booking for the same specialist should fail
		_, err = uc.Execute(app.CreateRequest{PetID: "pet_2", SpecialistID: "spec_1", StartTime: start.Add(10 * time.Minute), EndTime: end.Add(10 * time.Minute)})
		Expect(err).To(MatchError(mem.ErrSlotTaken))
	})
})
