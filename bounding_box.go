package mapbox

import "fmt"

type BoundingBox struct {
	Min Coordinate
	Max Coordinate
}

func (b BoundingBox) query() string {
	return fmt.Sprintf("%v,%v,%v,%v", b.Min.Lng, b.Min.Lat, b.Max.Lng, b.Max.Lat)
}
