package osm

import "fmt"

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
