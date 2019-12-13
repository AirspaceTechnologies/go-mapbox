# go-mapbox

API Wrapper for [Mapbox API](https://docs.mapbox.com/api/)

## Usage

### Initialization
```go
mapboxClient, err := mapbox.NewClient(&MapboxConfig{
    Timeout: 30 * time.Second,
    APIKey:  "YOUR_API_KEY_HERE",
})
// error checking ...  
```

### Retrieve a Matrix
```go
request := &mapbox.DirectionsMatrixRequest{
    Profile:       mapbox.ProfileDrivingTraffic,
    Coordinates:   mapbox.Coordinates{
        mapbox.Coordinate{Lat: 33.122508, Lng: -117.306786},
        mapbox.Coordinate{Lat: 32.733810, Lng: -117.193443},
        mapbox.Coordinate{Lat: 33.676084, Lng: -117.867598},
    },

    // optional fields below
    Annotations: mapbox.Annotations{mapbox.AnnotationDistance, mapbox.AnnotationDuration},
    Approaches: mapbox.Approaches{mapbox.ApproachUnrestricted},
    Sources: mapbox.Sources{0},
    FallbackSpeed: 60,
}

response, err := mapboxClient.DirectionsMatrix(context.TODO(), request)
// error checking ... 
```