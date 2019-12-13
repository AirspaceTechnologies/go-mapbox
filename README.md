# go-mapbox-api

API Wrapper for [Mapbox API](https://docs.mapbox.com/api/)

## Usage

### Initialization
```go
mapboxClient, err := NewClient(&MapboxConfig{
    Timeout: 30 * time.Second,
    APIKey:  "YOUR_API_KEY_HERE",
})
// error checking 
```

### Retrieve a matrix
```go
request := DirectionsMatrixRequest{
    Profile:       ProfileDrivingTraffic,
    Coordinates:   Coordinates{
        Coordinate{Lat: 33.122508, Lng: -117.306786},
        Coordinate{Lat: 32.733810, Lng: -117.193443},
        Coordinate{Lat: 33.676084, Lng: -117.867598},
    },

    // optional fields below
    Annotations: Annotations{AnnotationDistance, AnnotationDuration},
    Approaches: Approaches{ApproachUnrestricted},
    Sources: Sources{0},
    FallbackSpeed: 60,
}

response, err := mapbox.DirectionsMatrix(context.TODO(), req)
// error checking ... 
```