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

	AnnotationDuration   = Annotation("duration")
	AnnotationDistance   = Annotation("distance")
	AnnotationSpeed      = Annotation("speed")
	AnnotationCongestion = Annotation("congestion")

	ApproachUnrestricted = Approach("unrestricted")
	ApproachCurb         = Approach("curb")

	TypeCountry          = Type("country")
	TypeRegion           = Type("region")
	TypePostcode         = Type("postcode")
	TypeDistrict         = Type("district")
	TypePlace            = Type("place")
	TypeLocality         = Type("locality")
	TypeNeighborhood     = Type("neighborhood")
	TypeStreet           = Type("street")
	TypeBlock            = Type("block")
	TypeAddress          = Type("address")
	TypeSecondaryAddress = Type("secondary_address")
	TypePOI              = Type("poi")
	TypePOILandmark      = Type("poi.landmark")

	ExcludeMotorway      = Exclude("motorway")
	ExcludeToll          = Exclude("toll")
	ExcludeFerry         = Exclude("ferry")
	ExcludeUnpaved       = Exclude("unpaved")
	ExcludeCashOnlyTolls = Exclude("cash_only_tolls")

	GeometriesGeoJSON   = Geometries("geojson")
	GeometriesPolyline  = Geometries("polyline")
	GeometriesPolyline6 = Geometries("polyline6")

	IncludeHov2 = Include("hov2")
	IncludeHov3 = Include("hov3")
	IncludeHot  = Include("hot")

	OverviewFull       = Overview("full")
	OverviewSimplified = Overview("simplified")
	OverviewFalse      = Overview("false")

	VoiceUnitsImpreial = VoiceUnits("imperial")
	VoiceUnitsMetric   = VoiceUnits("metric")
)

type Profile string
type Endpoint string
type Geometries string
type Overview string
type VoiceUnits string

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

type Excludes []Exclude
type Exclude string

func (e Excludes) strings() []string {
	res := make([]string, 0, len(e))

	for _, val := range e {
		res = append(res, string(val))
	}

	return res
}

func (e Excludes) query() string {
	return strings.Join(e.strings(), ",")
}

//////////////////////////////////////////////////////////////////

type Includes []Include
type Include string

func (i Includes) strings() []string {
	res := make([]string, 0, len(i))

	for _, val := range i {
		res = append(res, string(val))
	}

	return res
}

func (i Includes) query() string {
	return strings.Join(i.strings(), ",")
}

//////////////////////////////////////////////////////////////////

type DirectionWaypoints []DirectionWaypoint
type DirectionWaypoint string

func (d DirectionWaypoints) strings() []string {
	res := make([]string, 0, len(d))

	for _, val := range d {
		res = append(res, string(val))
	}

	return res
}

func (d DirectionWaypoints) query() string {
	return strings.Join(d.strings(), ";")
}

//////////////////////////////////////////////////////////////////

type WaypointNames []WaypointName
type WaypointName string

func (w WaypointNames) strings() []string {
	res := make([]string, 0, len(w))

	for _, val := range w {
		res = append(res, string(val))
	}

	return res
}

func (w WaypointNames) query() string {
	return strings.Join(w.strings(), ";")
}

//////////////////////////////////////////////////////////////////

type WaypointTargets []WaypointTarget
type WaypointTarget string

func (w WaypointTargets) strings() []string {
	res := make([]string, 0, len(w))

	for _, val := range w {
		res = append(res, string(val))
	}

	return res
}

func (w WaypointTargets) query() string {
	return strings.Join(w.strings(), ";")
}

//////////////////////////////////////////////////////////////////

type DepartAt time.Time

const (
	DepartAtFormat = "2006-01-02T15:04:05Z"
)

func (t DepartAt) IsZero() bool {
	return time.Time(t).IsZero()
}

func (t DepartAt) query() string {
	if t.IsZero() {
		return ""
	}
	return time.Time(t).UTC().Format(DepartAtFormat)
}

//////////////////////////////////////////////////////////////////

type ArriveBy time.Time

const (
	ArriveByFormat = "2006-01-02T15:04:05Z"
)

func (t ArriveBy) IsZero() bool {
	return time.Time(t).IsZero()
}

func (t ArriveBy) query() string {
	if t.IsZero() {
		return ""
	}
	return time.Time(t).UTC().Format(ArriveByFormat)
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
