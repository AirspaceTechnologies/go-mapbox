package mapbox

type DirectionsResponse struct {
	Code   string  `json:"code"`
	UUID   string  `json:"uuid,omitempty"`
	Routes []Route `json:"routes"`
}

type Route struct {
	Duration        float64    `json:"duration"`
	Distance        float64    `json:"distance"`
	WeightName      string     `json:"weight_name"`
	Weight          float64    `json:"weight"`
	DurationTypical float64    `json:"duration_typical,omitempty"`
	WeightTypical   float64    `json:"weight_typical,omitempty"`
	Geometry        string     `json:"geometry,omitempty"`
	Legs            []RouteLeg `json:"legs"`
	VoiceLocale     string     `json:"voiceLocale,omitempty"`
	Waypoints       []Waypoint `json:"waypoints,omitempty"`
}

// RouteLeg represents a leg of the route between two waypoints.
type RouteLeg struct {
	Distance           float64              `json:"distance"`           // The distance traveled by the leg, in meters.
	Duration           float64              `json:"duration"`           // The estimated travel time, in seconds.
	Summary            string               `json:"summary"`            // A summary of the leg, containing the names of the significant roads.
	Weight             float64              `json:"weight"`             // The weight of the leg. The weight value is similar to the duration but includes additional factors like traffic.
	Steps              []Step               `json:"steps"`              // An array of RouteStep objects, each representing a step in the leg.
	Annotation         DirectionsAnnotation `json:"annotation"`         // Additional details about the leg.
	Admins             []Admin              `json:"admins"`             // Array of administrative region objects traversed by the leg.
	VoiceInstructions  []VoiceInstruction   `json:"voiceInstructions"`  // An array of VoiceInstruction objects.
	BannerInstructions []BannerInstruction  `json:"bannerInstructions"` // An array of BannerInstruction objects.
	ViaWaypoints       []ViaWaypoint        `json:"via_waypoints"`
}

// Step represents a single step in a leg of a route, containing maneuver instructions and distance/duration.
type Step struct {
	Distance      float64        `json:"distance"`      // The distance for this step in meters.
	Duration      float64        `json:"duration"`      // The estimated travel time for this step in seconds.
	Geometry      string         `json:"geometry"`      // An encoded polyline string or GeoJSON LineString representing the step geometry.
	Name          string         `json:"name"`          // The name of the road or path used in the step.
	Maneuver      Maneuver       `json:"maneuver"`      // The maneuver required to move from this step to the next.
	Mode          string         `json:"mode"`          // The travel mode of the step.
	Weight        float64        `json:"weight"`        // Similar to duration but includes additional factors like traffic.
	Intersections []Intersection `json:"intersections"` // An array of Intersection objects.
}

// Maneuver contains information about the required maneuver for a step, including type and bearing.
type Maneuver struct {
	BearingAfter  float64   `json:"bearing_after"`  // The clockwise angle from true north to the direction of travel after the maneuver.
	BearingBefore float64   `json:"bearing_before"` // The clockwise angle from true north to the direction of travel before the maneuver.
	Location      []float64 `json:"location"`       // A [longitude, latitude] pair describing the location of the maneuver.
	Type          string    `json:"type"`           // A string signifying the type of maneuver. Example: "turn".
	Modifier      string    `json:"modifier"`       // An additional modifier to provide more detail. Example: "left".
	Instruction   string    `json:"instruction"`    // Verbal instruction for the maneuver.
}

// Annotation contains additional details about each point along the route leg.
type DirectionsAnnotation struct {
	Distance   []float64  `json:"distance"`   // Array of distances between each pair of coordinates.
	Duration   []float64  `json:"duration"`   // Array of expected travel times from each coordinate to the next.
	Speed      []float64  `json:"speed"`      // Array of travel speeds.
	Congestion []string   `json:"congestion"` // Array of congestion levels.
	Maxspeed   []Maxspeed `json:"maxspeed"`
}

// Admin represents administrative region information.
type Admin struct {
	ISO_3166_1_alpha3 string `json:"iso_3166_1_alpha3"` // The ISO 3166-1 alpha-3 country code.
	ISO_3166_1        string `json:"iso_3166_1"`        // The ISO 3166-1 alpha-2 country code.
}

// VoiceInstruction represents a single voice instruction for navigation.
type VoiceInstruction struct {
	DistanceAlongGeometry float64 `json:"distanceAlongGeometry"` // The distance from the current step at which to provide the instruction.
	Announcement          string  `json:"announcement"`          // The verbal instruction.
	SSMLAnnouncement      string  `json:"ssmlAnnouncement"`      // The instruction in SSML format for text-to-speech engines.
}

// BannerInstruction represents a visual instruction for navigation.
type BannerInstruction struct {
	DistanceAlongGeometry float64      `json:"distanceAlongGeometry"` // The distance from the current step at which to show the instruction.
	Primary               Instruction  `json:"primary"`               // The primary instruction for this step.
	Secondary             *Instruction `json:"secondary,omitempty"`   // An optional secondary instruction.
}

// Instruction contains the details of a navigation instruction.
type Instruction struct {
	Text       string      `json:"text"`       // The instruction text.
	Type       string      `json:"type"`       // The type of maneuver.
	Modifier   string      `json:"modifier"`   // An additional modifier to provide more detail.
	Components []Component `json:"components"` // Components of the instruction.
}

// Component represents a part of the instruction, useful for highlighting parts of the text.
type Component struct {
	Text string `json:"text"` // The component text.
	Type string `json:"type"` // The type of component, e.g., "text" or "icon".
}

// Intersection represents an intersection along a step.
type Intersection struct {
	Location      []float64 `json:"location"`       // The location of the intersection [longitude, latitude].
	Bearings      []int     `json:"bearings"`       // The bearings at the intersection, in degrees.
	Entry         []bool    `json:"entry"`          // A boolean flag indicating the availability of the corresponding bearing.
	In            int       `json:"in,omitempty"`   // The index into the bearings/entry array that denotes the incoming bearing to the intersection.
	Out           int       `json:"out,omitempty"`  // The index into the bearings/entry array that denotes the outgoing bearing from the intersection.
	AdminIndex    int       `json:"admin_index"`    // The index into the leg's admins array of the administrative region containing this intersection.
	GeometryIndex int       `json:"geometry_index"` // The index into the leg geometry of this intersection's coordinate.
}

// AdminCrossing is an administrative boundary change along a leg, located at
// the first intersection inside the new region. Mapbox emits one admins entry
// per boundary traversed, so consecutive entries can share a country (e.g.
// state lines) — compare From/To ISO codes to detect country changes. Duration
// (seconds) and Distance (meters) are cumulative from the start of the leg;
// they are exact when the request included distance/duration annotations and
// fall back to completed-step sums otherwise.
type AdminCrossing struct {
	Location []float64 // The location of the crossing [longitude, latitude].
	From     Admin
	To       Admin
	Duration float64
	Distance float64
}

// AdminCrossings locates administrative region changes along the leg's step
// intersections. Legs traversing fewer than two regions produce none.
func (l RouteLeg) AdminCrossings() []AdminCrossing {
	if len(l.Admins) < 2 {
		return nil
	}

	var crossings []AdminCrossing
	var stepDuration, stepDistance float64
	prevAdmin := -1

	for _, step := range l.Steps {
		for _, intersection := range step.Intersections {
			idx := intersection.AdminIndex
			if idx != prevAdmin && prevAdmin >= 0 && idx >= 0 && idx < len(l.Admins) && len(intersection.Location) >= 2 {
				duration, distance, ok := l.annotationTotals(intersection.GeometryIndex)
				if !ok {
					duration, distance = stepDuration, stepDistance
				}

				crossings = append(crossings, AdminCrossing{
					Location: intersection.Location,
					From:     l.Admins[prevAdmin],
					To:       l.Admins[idx],
					Duration: duration,
					Distance: distance,
				})
			}
			prevAdmin = idx
		}

		stepDuration += step.Duration
		stepDistance += step.Distance
	}

	return crossings
}

// annotationTotals sums the per-segment duration/distance annotations for the
// leg geometry up to the given coordinate index. ok is false when the leg has
// no annotations (not requested) or the index is out of range.
func (l RouteLeg) annotationTotals(geometryIndex int) (duration, distance float64, ok bool) {
	durations := l.Annotation.Duration
	distances := l.Annotation.Distance
	if geometryIndex <= 0 || geometryIndex > len(durations) || geometryIndex > len(distances) {
		return 0, 0, false
	}

	for i := 0; i < geometryIndex; i++ {
		duration += durations[i]
		distance += distances[i]
	}
	return duration, distance, true
}

type ViaWaypoint struct {
	WaypointIndex     int     `json:"waypoint_index"`
	DistanceFromStart float64 `json:"distance_from_start"`
	GeometryIndex     int     `json:"geometry_index"`
}

type Maxspeed struct {
	Speed int    `json:"speed"`
	Unit  string `json:"unit"`
}
