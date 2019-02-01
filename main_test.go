package geojson

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewMember(t *testing.T) {
	var tests = []struct {
		json     []byte
		ok       bool
		expected *Member
	}{
		{
			json: []byte(`{"type": "Point", "coordinates": [1.23, 4.56]}`),
			ok:   true,
			expected: &Member{
				memberKind:     GeometryObject,
				Type:           "Point",
				CoordinatesRaw: []byte(`[1.23, 4.56]`),
				CoordinatesObj: &Point{1.23, 4.56},
			},
		},
		{
			json: []byte(`{"type": "XXX", "coordinates": [1.23, 4.56]}`),
			ok:   false,
		},
		{
			json: []byte(`{"type": "Point", "coordinates": [1.23, 4.56],xxxxxx`),
			ok:   false,
		},
		{
			json: []byte(`{"type": "LineString", "coordinates": [[1.23, 4.56],[7.89,10.12]]}`),
			ok:   true,
			expected: &Member{
				memberKind:     GeometryObject,
				Type:           "LineString",
				CoordinatesRaw: []byte(`[[1.23, 4.56],[7.89,10.12]]`),
				CoordinatesObj: &LineString{Point{1.23, 4.56}, Point{7.89, 10.12}},
			},
		},
		{
			json: []byte(`{"type": "Polygon", "coordinates": [[[1.23, 4.56],[7.89,10.12],[3.45,6.78],[1.23,4.56]]]}`),
			ok:   true,
			expected: &Member{
				memberKind:     GeometryObject,
				Type:           "Polygon",
				CoordinatesRaw: []byte(`[[[1.23, 4.56],[7.89,10.12],[3.45,6.78],[1.23,4.56]]]`),
				CoordinatesObj: &Polygon{LineString{Point{1.23, 4.56}, Point{7.89, 10.12}, Point{3.45, 6.78}, Point{1.23, 4.56}}},
			},
		},
		{
			json: []byte(`{"type": "Feature", "geometry": {"type": "Point", "coordinates": [1.23, 4.56]}, "properties": {"name": "point:A"}}`),
			ok:   true,
			expected: &Member{
				memberKind:  FeatureObject,
				Type:        "Feature",
				GeometryRaw: []byte(`{"type": "Point", "coordinates": [1.23, 4.56]}`),
				GeometryObj: &Member{
					memberKind:     GeometryObject,
					Type:           "Point",
					CoordinatesRaw: []byte(`[1.23, 4.56]`),
					CoordinatesObj: &Point{1.23, 4.56},
				},
				Properties: map[string]string{"name": "point:A"},
			},
		},
		{
			json: []byte(`{"type": "Feature", "geometry": {"type": "point", "coordinates": [1.23, 4.56]}, "properties": {"name": "point:A"}}`),
			ok:   false,
		},
	}

	for _, test := range tests {
		m, err := NewMember(test.json)
		if test.ok && err != nil {
			t.Fatalf("err:%v json:%s\n", err, test.json)
		} else if !test.ok && err == nil {
			t.Fatalf("expected: err actual: ok json:%s\n", test.json)
		}
		if test.expected != nil && !reflect.DeepEqual(test.expected, m) {
			fmt.Printf("member:%v\n", m)
			t.Errorf("json: %s\nexpected: %v\nactual: %v", test.json, test.expected, m)
			t.Errorf("reflect.DeepEqual(memberKind):%v\n", reflect.DeepEqual(test.expected.memberKind, m.memberKind))
			t.Errorf("reflect.DeepEqual(Type):%v\n", reflect.DeepEqual(test.expected.Type, m.Type))
			t.Errorf("reflect.DeepEqual(CoordinatesRaw):%v\n", reflect.DeepEqual(test.expected.CoordinatesRaw, m.CoordinatesRaw))
			t.Errorf("reflect.DeepEqual(CoordinatesObj):%v\n", reflect.DeepEqual(test.expected.CoordinatesObj, m.CoordinatesObj))
			t.Errorf("reflect.DeepEqual(GeometryRaw):%v\n", reflect.DeepEqual(test.expected.GeometryRaw, m.GeometryRaw))
			t.Errorf("reflect.DeepEqual(GeometryObj):%v\n", reflect.DeepEqual(test.expected.GeometryObj, m.GeometryObj))
		}
	}
}

func TestNewMembers(t *testing.T) {
	var tests = []struct {
		json     []byte
		ok       bool
		expected []*Member
	}{
		{
			json: []byte(`[
		{"type": "Point", "coordinates": [1.23, 4.56]},
		{"type": "LineString", "coordinates": [[1.23, 4.56],[7.89,10.12]]},
		{"type": "Polygon", "coordinates": [[[1.23, 4.56],[7.89,10.12],[3.45,6.78],[1.23,4.56]]]},
		{"type": "Feature", "geometry": {"type": "Point", "coordinates": [1.23, 4.56]}, "properties": {"name": "point:A"}}
		]`),
			ok: true,
			expected: []*Member{
				&Member{
					memberKind:     GeometryObject,
					Type:           "Point",
					CoordinatesRaw: []byte(`[1.23, 4.56]`),
					CoordinatesObj: &Point{1.23, 4.56},
				},
				&Member{
					memberKind:     GeometryObject,
					Type:           "LineString",
					CoordinatesRaw: []byte(`[[1.23, 4.56],[7.89,10.12]]`),
					CoordinatesObj: &LineString{Point{1.23, 4.56}, Point{7.89, 10.12}},
				},
				&Member{
					memberKind:     GeometryObject,
					Type:           "Polygon",
					CoordinatesRaw: []byte(`[[[1.23, 4.56],[7.89,10.12],[3.45,6.78],[1.23,4.56]]]`),
					CoordinatesObj: &Polygon{LineString{Point{1.23, 4.56}, Point{7.89, 10.12}, Point{3.45, 6.78}, Point{1.23, 4.56}}},
				},
				&Member{
					memberKind:  FeatureObject,
					Type:        "Feature",
					GeometryRaw: []byte(`{"type": "Point", "coordinates": [1.23, 4.56]}`),
					GeometryObj: &Member{
						memberKind:     GeometryObject,
						Type:           "Point",
						CoordinatesRaw: []byte(`[1.23, 4.56]`),
						CoordinatesObj: &Point{1.23, 4.56},
					},
					Properties: map[string]string{"name": "point:A"},
				},
			},
		},
		{
			json: []byte(`[
		{"type": "point", "coordinates": [1.23, 4.56]}
		]`),
			ok: false,
		},
		{
			json: []byte(`[
		{"type": "Point", "coordinates": [1.23, 4.56],xxxxx
		]`),
			ok: false,
		},
	}
	for _, test := range tests {
		m, err := NewMembers(test.json)
		if test.ok && err != nil {
			t.Fatalf("err:%v json:%s\n", err, test.json)
		} else if !test.ok && err == nil {
			t.Fatalf("expected: err actual: ok json:%s\n", test.json)
		}
		if test.expected != nil && !reflect.DeepEqual(test.expected, m) {
			fmt.Printf("member:%v\n", m)
			t.Errorf("json: %s\nexpected: %v\nactual: %v", test.json, test.expected, m)
		}
	}
}
