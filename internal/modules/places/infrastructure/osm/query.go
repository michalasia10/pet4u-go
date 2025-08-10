package osm

import (
	"fmt"
	"strings"
)

// QueryBuilder builds a simple Overpass QL around query.
type QueryBuilder struct {
	radius     int
	lat        float64
	lng        float64
	nameFilter string
	filters    []string
}

func NewQueryBuilder() *QueryBuilder { return &QueryBuilder{filters: []string{}} }

func (b *QueryBuilder) Around(radius int, lat, lng float64) *QueryBuilder {
	b.radius = radius
	b.lat = lat
	b.lng = lng
	return b
}

func (b *QueryBuilder) WithFilters(filters ...string) *QueryBuilder {
	b.filters = append(b.filters, filters...)
	return b
}

// NameRegex sets a case-insensitive regex filter on tag name.
func (b *QueryBuilder) NameRegex(regex string) *QueryBuilder {
	if strings.TrimSpace(regex) != "" {
		b.nameFilter = fmt.Sprintf("[name~\"%s\",i]", regex)
	}
	return b
}

func (b *QueryBuilder) Build() string {
	var sb strings.Builder
	sb.WriteString("[out:json];(")
	for _, f := range b.filters {
		fmt.Fprintf(&sb, "node[%s]%s(around:%d,%.6f,%.6f);", f, b.nameFilter, b.radius, b.lat, b.lng)
		fmt.Fprintf(&sb, "way[%s]%s(around:%d,%.6f,%.6f);", f, b.nameFilter, b.radius, b.lat, b.lng)
	}
	sb.WriteString(");out center;")
	return sb.String()
}
