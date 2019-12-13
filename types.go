package mapbox

import (
	"fmt"
	"strconv"
)

const (
	ProfileDriving = Profile("mapbox/driving")
	ProfileWalking = Profile("mapbox/walking")
	ProfileCycling = Profile("mapbox/cycling")
	ProfileDrivingTraffic = Profile("mapbox/driving-traffic")

	AnnotationDuration = Annotation("duration")
	AnnotationDistance = Annotation("distance")
	AnnotationSpeed = Annotation("speed")
	AnnotationCongestion = Annotation("congestion")

	ApproachUnrestricted = Approach("unrestricted")
	ApproachCurb = Approach("curb")
)

type Profile string

type Annotations []Annotation
type Annotation string
func (a Annotations) strings() []string {
	res := make([]string, 0, len(a))

	for _, val := range a {
		res = append(res, string(val))
	}

	return res
}

type Approaches []Approach
type Approach string
func (a Approaches) strings() []string {
	res := make([]string, 0, len(a))

	for _, val := range a {
		res = append(res, string(val))
	}

	return res
}

type Sources []int
func (s Sources) strings() []string {
	res := make([]string, 0, len(s))

	for _, val := range s {
		res = append(res, strconv.Itoa(val))
	}

	return res
}

type Destinations []int
func (d Destinations) strings() []string {
	res := make([]string, 0, len(d))

	for _, val := range d {
		res = append(res, strconv.Itoa(val))
	}

	return res
}

type FallbackSpeed float64
func (f FallbackSpeed) String() string {
	if f == 0 {
		return ""
	}

	return fmt.Sprintf("%f", f)
}