package domain_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "src/internal/modules/appointments/domain"
)

var _ = Describe("TimeSlot", func() {
	It("validates Start < End", func() {
		t0 := time.Now()
		_, err := NewTimeSlot(t0, t0)
		Expect(err).To(MatchError(ErrInvalidTimeSlot))
	})

	It("computes overlap correctly", func() {
		t0 := time.Now()
		a, _ := NewTimeSlot(t0, t0.Add(30*time.Minute))
		b, _ := NewTimeSlot(t0.Add(10*time.Minute), t0.Add(40*time.Minute))
		c, _ := NewTimeSlot(t0.Add(30*time.Minute), t0.Add(60*time.Minute))
		Expect(a.Overlaps(b)).To(BeTrue())
		Expect(a.Overlaps(c)).To(BeFalse())
	})
})
