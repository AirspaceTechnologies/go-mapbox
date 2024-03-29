package mapbox

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	ProfileDriving        = Profile("mapbox/driving")
	ProfileWalking        = Profile("mapbox/walking")
	ProfileCycling        = Profile("mapbox/cycling")
	ProfileDrivingTraffic = Profile("mapbox/driving-traffic")

	EndpointPlaces          = Endpoint("mapbox.places")
	EndpointPlacesPermanent = Endpoint("mapbox.places-permanent")

	AnnotationDuration   = Annotation("duration")
	AnnotationDistance   = Annotation("distance")
	AnnotationSpeed      = Annotation("speed")
	AnnotationCongestion = Annotation("congestion")

	ApproachUnrestricted = Approach("unrestricted")
	ApproachCurb         = Approach("curb")

	ReverseModeDistance = ReverseMode("distance")
	ReverseModeScore    = ReverseMode("score")

	TypeCountry      = Type("country")
	TypeRegion       = Type("region")
	TypePostcode     = Type("postcode")
	TypeDistrict     = Type("district")
	TypePlace        = Type("place")
	TypeLocality     = Type("locality")
	TypeNeighborhood = Type("neighborhood")
	TypeAddress      = Type("address")
	TypePOI          = Type("poi")
)

type Profile string
type Endpoint string

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

//////////////////////////////////////////////////////////////////

type ReverseMode string

func (r ReverseMode) query() string {
	return string(r)
}

//////////////////////////////////////////////////////////////////

type Types []Type
type Type string

func (t Types) strings() []string {
	res := make([]string, 0, len(t))

	for _, val := range t {
		res = append(res, string(val))
	}

	return res
}

func (t Types) query() string {
	return strings.Join(t.strings(), ",")
}

//////////////////////////////////////////////////////////////////

type DepartureTime time.Time

const (
	DepartureTimeFormat = "2006-01-02T15:04:05Z"
)

func (t DepartureTime) IsZero() bool {
	return time.Time(t).IsZero()
}

func (t DepartureTime) query() string {
	if t.IsZero() {
		return ""
	}
	return time.Time(t).UTC().Format(DepartureTimeFormat)
}

//////////////////////////////////////////////////////////////////
