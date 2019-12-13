package mapbox

import (
	"fmt"
	"strings"
)

type Coordinate struct {
	Lat float64
	Lng float64
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
