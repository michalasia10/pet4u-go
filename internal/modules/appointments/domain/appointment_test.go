package domain_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "src/internal/modules/appointments/domain"
)

var _ = Describe("NewAppointment", func() {
	It("rejects past bookings", func() {
		now := time.Now()
		slot, _ := NewTimeSlot(now.Add(-10*time.Minute), now.Add(10*time.Minute))
		_, err := NewAppointment(now, 0, "pet_1", "spec_1", slot)
		Expect(err).To(MatchError(ErrBookingInPast))
	})

	It("enforces minimal lead time", func() {
		now := time.Now()
		slot, _ := NewTimeSlot(now.Add(5*time.Minute), now.Add(35*time.Minute))
		_, err := NewAppointment(now, 10*time.Minute, "pet_1", "spec_1", slot)
		Expect(err).To(MatchError(ErrTooShortLeadTime))
	})

	It("requires non-empty pet and specialist IDs", func() {
		now := time.Now()
		slot, _ := NewTimeSlot(now.Add(20*time.Minute), now.Add(50*time.Minute))
		_, err := NewAppointment(now, 5*time.Minute, "", "spec_1", slot)
		Expect(err).To(MatchError(ErrInvalidPetID))
		_, err = NewAppointment(now, 5*time.Minute, "pet_1", "", slot)
		Expect(err).To(MatchError(ErrInvalidSpecialist))
	})
})
