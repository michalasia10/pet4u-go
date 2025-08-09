package idgen

import (
    "strings"
    "time"
)

type TimeIDGen struct{}

func NewTimeIDGen() *TimeIDGen { return &TimeIDGen{} }

func (TimeIDGen) NewID(prefix string) string {
    p := strings.Trim(prefix, "_")
    if p != "" { p = p + "_" }
    return p + time.Now().UTC().Format("20060102150405.000000000")
}


