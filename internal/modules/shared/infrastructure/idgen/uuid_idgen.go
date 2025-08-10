package idgen

import (
	"strings"

	"github.com/google/uuid"
)

// UUIDGen generates RFC 4122 UUIDs. It preserves the prefix convention by
// returning values like: <prefix>_<uuid> or just <uuid> when prefix is empty.
type UUIDGen struct{}

func NewUUIDGen() *UUIDGen { return &UUIDGen{} }

func (UUIDGen) NewID(prefix string) string {
	u := uuid.NewString()
	p := strings.Trim(prefix, "_")
	if p != "" {
		return p + "_" + u
	}
	return u
}
