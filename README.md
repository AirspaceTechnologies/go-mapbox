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
    DepartureTime: mapbox.DepartureTime(time.Now()),
}

response, err := mapboxClient.DirectionsMatrix(context.TODO(), request)
// error checking ... 
```

### Forward Geocode

```go
request := &mapbox.ForwardGeocodeRequest{
    SearchText:   "6005 Hidden Valley Rd, Suite 280, Carlsbad, CA 92011"

    // optional fields below
    Autocomplete: false,
    BBox: mapbox.BoundingBox{
        Min: mapbox.Coordinate{
            Lat: 33.121217,
            Lng: -117.310429,
        }, Max: mapbox.Coordinate{
            Lat: 33.124973,
            Lng: -117.305054,
        }},
    Country:    "us",
    Language:   "en",
    Limit:      1,
    Proximity:    Coordinate{Lat: 33.121217, Lng: -117.310429,},
    Types:        mapbox.Types{mapbox.TypeCountry},
}

response, err := mapboxClient.ForwardGeocode(context.TODO(), request)
// error checking ... 
```

### Reverse Geocode

```go
request := &mapbox.ReverseGeocodeRequest{
    Coordinates:   mapbox.Coordinates{
        mapbox.Coordinate{Lat: 33.122508, Lng: -117.306786}
    },

    // optional fields below
    Country:     "us",
    Language:    "en",
    Limit:       1,
    Types:       mapbox.Types{mapbox.TypeCountry, mapbox.TypeAddress},
}

response, err := mapboxClient.ReverseGeocode(context.TODO(), request)
// error checking ... 
```

### Reverse Geocode Batch

```go
requests := ReverseGeocodeBatchRequest{
    mapbox.ReverseGeocodeRequest{
        Coordinates:   mapbox.Coordinates{
            mapbox.Coordinate{Lat: 33.122508, Lng: -117.306786}
        },
        Types:       mapbox.Types{mapbox.TypeCountry, mapbox.TypeAddress},
        Language:   "en",
    },
     mapbox.ReverseGeocodeRequest{
        Coordinates:   mapbox.Coordinates{
            mapbox.Coordinate{Lat: 32.733810, Lng: -117.193443}
        },
        Types:       mapbox.Types{mapbox.TypeCountry, mapbox.TypeAddress},
        Language:   "en",
    },
}

responses, err := mapboxClient.ReverseGeocodeBatch(context.TODO(), request)
if err != nil {
    // handle error... 
}

for i, response := range responses.Batch {
    // responses are in request order, i.e. request[i] ==> responses.Batch[i]
}
```

### Reverse Searchbox

```go
request := &mapbox.ReverseGeocodeRequest{
    Coordinates:   mapbox.Coordinates{
        mapbox.Coordinate{Lat: 33.122508, Lng: -117.306786}
    },

    // optional fields below
    Country:     "us",
    Language:    "en",
    Limit:       1,
    Types:       mapbox.Types{mapbox.TypePOI},
}

response, err := mapboxClient.SearchboxReverse(context.TODO(), request)
// error checking ... 
```

### Retrieve Directions

```go
request := &mapbox.DirectionsRequest{
    Profile:       mapbox.ProfileDrivingTraffic,
    Coordinates:   mapbox.Coordinates{
        mapbox.Coordinate{Lat: 33.122508, Lng: -117.306786},
        mapbox.Coordinate{Lat: 32.733810, Lng: -117.193443},
    },

    // optional fields below
    Annotations: mapbox.Annotations{mapbox.AnnotationDistance, mapbox.AnnotationDuration},
}

response, err := mapboxClient.Directions(context.TODO(), request)
// error checking ...
```
