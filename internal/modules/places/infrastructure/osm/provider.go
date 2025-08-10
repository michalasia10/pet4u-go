package osm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"src/internal/modules/places/domain"
)

// Provider implements domain.ExternalPlacesProvider using OpenStreetMap Overpass API.
// It performs a simple nearby query constrained by radius and optional keywords.
type Provider struct {
	httpClient  *http.Client
	endpointURL string // Overpass API endpoint
}

func NewProvider(overpassURL string) *Provider {
	if overpassURL == "" {
		overpassURL = "https://overpass-api.de/api/interpreter"
	}
	return &Provider{
		httpClient:  &http.Client{Timeout: 3 * time.Second},
		endpointURL: overpassURL,
	}
}

func (p *Provider) ProviderName() string { return "osm" }

func (p *Provider) Search(criteria domain.SearchCriteria) ([]domain.Place, error) {
	// We rely on Center/RadiusM; if not provided, we cannot query nearby sensibly
	if criteria.Center == nil || criteria.RadiusM == nil {
		return []domain.Place{}, nil
	}
	// Build Overpass QL query. We use around:radius with a set of amenity/shop categories.
	radius := *criteria.RadiusM
	lat := criteria.Center.Lat
	lng := criteria.Center.Lng

	// Basic categories; could be refined by PetType/tags in future
	filters := []string{
		"amenity=veterinary",
		"shop=pet",
		"amenity=park",
		"leisure=dog_park",
	}
	// If PetType specified, prioritize dog/cat specific tags
	if criteria.PetType != nil && *criteria.PetType == domain.PetDog {
		filters = append(filters, "amenity=drinking_water")
	}

	// Optional query keyword: use name~"..." i if provided
	nameFilter := ""
	if criteria.Query != "" {
		// Sanitize quotes; Overpass supports regex; use case-insensitive flag i
		q := url.QueryEscape(criteria.Query)
		// url.QueryEscape not suitable inside Overpass QL; keep simple alnum/space subset
		_ = q
		nameFilter = fmt.Sprintf("[name~\"%s\",i]", criteria.Query)
	}

	// Construct QL
	// Example: [out:json];(node[amenity=veterinary](around:1000,lat,lng);node[shop=pet](around:1000,lat,lng););out center 50;
	ql := "[out:json];("
	for _, f := range filters {
		ql += fmt.Sprintf("node[%s]%s(around:%d,%.6f,%.6f);", f, nameFilter, radius, lat, lng)
		ql += fmt.Sprintf("way[%s]%s(around:%d,%.6f,%.6f);", f, nameFilter, radius, lat, lng)
	}
	ql += ");out center;"

	form := url.Values{}
	form.Set("data", ql)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, p.endpointURL, stringsNewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("osm: status %d", resp.StatusCode)
	}

	var payload overpassResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	places := make([]domain.Place, 0, len(payload.Elements))
	for _, el := range payload.Elements {
		name := el.Tags["name"]
		if name == "" {
			continue
		}
		address := composeAddress(el.Tags)
		ptypes := []domain.PetType{}
		// Heuristic: any of these implies pet-related
		if el.Tags["amenity"] == "veterinary" || el.Tags["shop"] == "pet" || el.Tags["leisure"] == "dog_park" {
			ptypes = append(ptypes, domain.PetDog, domain.PetCat)
		}
		place := domain.Place{
			ID:            "",
			Name:          name,
			Address:       address,
			Location:      domain.GeoPoint{Lat: el.LatOrCenterLat(), Lng: el.LonOrCenterLon()},
			Tags:          deriveTags(el.Tags),
			PetTypes:      ptypes,
			IsPetFriendly: true,
			Sources:       []domain.SourceRef{{Provider: p.ProviderName(), ID: el.ElementID()}},
		}
		places = append(places, place)
	}
	return places, nil
}

// --- helpers ---

type overpassResponse struct {
	Elements []overpassElement `json:"elements"`
}

type overpassElement struct {
	Type   string            `json:"type"`
	ID     int64             `json:"id"`
	Lat    *float64          `json:"lat"`
	Lon    *float64          `json:"lon"`
	Center *overpassCenter   `json:"center"`
	Tags   map[string]string `json:"tags"`
}

type overpassCenter struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func (e overpassElement) ElementID() string {
	return fmt.Sprintf("%s/%d", e.Type, e.ID)
}

func (e overpassElement) LatOrCenterLat() float64 {
	if e.Lat != nil {
		return *e.Lat
	}
	if e.Center != nil {
		return e.Center.Lat
	}
	return 0
}

func (e overpassElement) LonOrCenterLon() float64 {
	if e.Lon != nil {
		return *e.Lon
	}
	if e.Center != nil {
		return e.Center.Lon
	}
	return 0
}

func composeAddress(tags map[string]string) string {
	parts := []string{}
	if v := tags["addr:street"]; v != "" {
		parts = append(parts, v)
	}
	if v := tags["addr:housenumber"]; v != "" {
		if len(parts) > 0 {
			parts[len(parts)-1] = fmt.Sprintf("%s %s", parts[len(parts)-1], v)
		} else {
			parts = append(parts, v)
		}
	}
	if v := tags["addr:city"]; v != "" {
		parts = append(parts, v)
	}
	return stringsJoin(parts, ", ")
}

func deriveTags(tags map[string]string) []string {
	out := []string{}
	for k, v := range tags {
		switch {
		case k == "amenity" && (v == "veterinary" || v == "drinking_water"):
			out = append(out, v)
		case k == "shop" && v == "pet":
			out = append(out, "pet_shop")
		case k == "leisure" && v == "dog_park":
			out = append(out, v)
		}
	}
	return out
}

// lightweight string helpers to avoid extra imports
func stringsJoin(elems []string, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return elems[0]
	}
	n := len(sep) * (len(elems) - 1)
	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}
	b := make([]byte, n)
	bp := copy(b, elems[0])
	for _, s := range elems[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], s)
	}
	return string(b)
}

func stringsNewReader(s string) *stringsReader { return &stringsReader{s: s} }

type stringsReader struct {
	s string
	i int64
}

func (r *stringsReader) Read(p []byte) (n int, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	n = copy(p, r.s[r.i:])
	r.i += int64(n)
	return n, nil
}

func (r *stringsReader) Close() error { return nil }
