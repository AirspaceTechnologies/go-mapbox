package mapbox

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestSearchboxBrand(t *testing.T) {
	tests := map[string][]string{
		"null":                            nil,
		`"null"`:                          {`"null"`},
		`[]`:                              {},
		`["Best Breakfast in Town Café"]`: {"Best Breakfast in Town Café"},
		`["Pancakes Eggs", "Waffles"]`:    {"Pancakes Eggs", "Waffles"},
		`{"place": "ok"}`:                 {`{"place": "ok"}`},
	}

	for name, expected := range tests {
		t.Run(fmt.Sprintf(`Testing '%s'`, name), func(t *testing.T) {
			var dest SearchboxBrand
			if err := json.Unmarshal([]byte(name), &dest); err != nil {
				t.Fatal(err)
			}

			converted := []string(dest)
			if !reflect.DeepEqual(converted, expected) {
				t.Errorf("expected %v to equal %v", converted, expected)
			}
		})
	}
}
