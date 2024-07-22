package mapbox

import (
	"fmt"
	"strings"
)

type Coordinate struct {
	Lat float64 `json:"latitude"`
	Lng float64 `json:"longitude"`
}

func (c Coordinate) IsZero() bool {
	return c.Lat == 0 && c.Lng == 0
}

type Coordinates []Coordinate

// https://docs.mapbox.com/api/#coordinate-format
func (c Coordinate) WGS84Format() string {
	var b strings.Builder
	b.Grow(21) // 10(lat) + 10(lng) + 1(comma)

	fmt.Fprintf(&b, "%v,%v", c.Lng, c.Lat)

	return b.String()
}

// https://docs.mapbox.com/api/#coordinate-format
func (c Coordinates) WGS84Format() string {
	var b strings.Builder
	b.Grow(len(c) * 20) // 10(lat) + 10(lng)

	for _, coordinate := range c {
		fmt.Fprintf(&b, "%v;", coordinate.WGS84Format())
	}

	result := b.String()
	return result[:len(result)-1] // remove trailing ';'
}

////////////////////////////////////////////////////////////////////////////////

type ExtendedCoordinate struct {
	Latitude       float64         `json:"latitude"`
	Longitude      float64         `json:"longitude"`
	Accuracy       string          `json:"accuracy"`
	RoutablePoints []RoutablePoint `json:"routable_points,omitempty"`
}

type RoutablePoint struct {
	Name      string
	Latitude  float64
	Longitude float64
}

////////////////////////////////////////////////////////////////////////////////

const (
	GeometryLngIdx = 0
	GeometryLatIdx = 1
)

type Geometry struct {
	Coordinates []float64 `json:"coordinates"` // [lng, lat]
	Type        string    `json:"type"`
}

func (g Geometry) Latitude() float64 {
	return g.Coordinates[GeometryLatIdx]
}

func (g Geometry) Longitude() float64 {
	return g.Coordinates[GeometryLngIdx]
}
