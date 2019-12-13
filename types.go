package mapbox

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ProfileDriving        = Profile("mapbox/driving")
	ProfileWalking        = Profile("mapbox/walking")
	ProfileCycling        = Profile("mapbox/cycling")
	ProfileDrivingTraffic = Profile("mapbox/driving-traffic")

	AnnotationDuration   = Annotation("duration")
	AnnotationDistance   = Annotation("distance")
	AnnotationSpeed      = Annotation("speed")
	AnnotationCongestion = Annotation("congestion")

	ApproachUnrestricted = Approach("unrestricted")
	ApproachCurb         = Approach("curb")
)

type Profile string

//////////////////////////////////////////////////////////////////

type Annotations []Annotation
type Annotation string

func (a Annotations) strings() []string {
	res := make([]string, 0, len(a))

	for _, val := range a {
		res = append(res, string(val))
	}

	return res
}

func (a Annotations) query() string {
	return strings.Join(a.strings(), ",")
}

//////////////////////////////////////////////////////////////////

type Approaches []Approach
type Approach string

func (a Approaches) strings() []string {
	res := make([]string, 0, len(a))

	for _, val := range a {
		res = append(res, string(val))
	}

	return res
}

func (a Approaches) query() string {
	return strings.Join(a.strings(), ";")
}

//////////////////////////////////////////////////////////////////

type Sources []int

func (s Sources) strings() []string {
	res := make([]string, 0, len(s))

	for _, val := range s {
		res = append(res, strconv.Itoa(val))
	}

	return res
}

func (s Sources) query() string {
	return strings.Join(s.strings(), ";")
}

//////////////////////////////////////////////////////////////////

type Destinations []int

func (d Destinations) strings() []string {
	res := make([]string, 0, len(d))

	for _, val := range d {
		res = append(res, strconv.Itoa(val))
	}

	return res
}

func (d Destinations) query() string {
	return strings.Join(d.strings(), ";")
}

//////////////////////////////////////////////////////////////////

type FallbackSpeed float64

func (f FallbackSpeed) query() string {
	if f == 0 {
		return ""
	}

	return fmt.Sprintf("%f", f)
}
