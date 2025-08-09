package memory_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"src/internal/modules/appointments/domain"
	. "src/internal/modules/appointments/infrastructure/memory"
)

var _ = Describe("InMemory services", func() {
	var repo *InMemoryAppointmentRepository

	BeforeEach(func() {
		repo = NewInMemoryAppointmentRepository()
	})

	It("EnsureAvailable returns ErrSlotTaken when overlap exists", func() {
		start := time.Now().Add(time.Hour)
		end := start.Add(30 * time.Minute)
		// seed conflicting appointment for specialist s1
		repo.Create(domain.Appointment{ID: "a1", SpecialistID: "s1", StartTime: start, EndTime: end, Status: domain.AppointmentStatusBooked})

		svc := NewInMemoryAvailabilityService(repo)
		overlap, _ := domain.NewTimeSlot(start.Add(10*time.Minute), end.Add(10*time.Minute))
		err := svc.EnsureAvailable("s1", overlap)
		Expect(err).To(MatchError(ErrSlotTaken))
	})

	It("EnsureAvailable passes when no overlap for specialist", func() {
		start := time.Now().Add(time.Hour)
		end := start.Add(30 * time.Minute)
		repo.Create(domain.Appointment{ID: "a1", SpecialistID: "s1", StartTime: start, EndTime: end, Status: domain.AppointmentStatusBooked})

		svc := NewInMemoryAvailabilityService(repo)
		nonOverlap, _ := domain.NewTimeSlot(end, end.Add(30*time.Minute))
		err := svc.EnsureAvailable("s1", nonOverlap)
		Expect(err).ToNot(HaveOccurred())
	})

	It("BookingPolicy allows by default", func() {
		pol := NewInMemoryBookingPolicy()
		now := time.Now()
		slot, _ := domain.NewTimeSlot(now.Add(time.Hour), now.Add(2*time.Hour))
		Expect(pol.EnsureEligible("pet_1", slot, now)).To(Succeed())
	})
})
