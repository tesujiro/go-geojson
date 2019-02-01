package geojson

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type memberKind int

const (
	GeometryObject = iota
	FeatureObject
	FeatureCollectionObject
)

type Member struct {
	memberKind     memberKind
	Type           string          `json:"type"`
	CoordinatesRaw json.RawMessage `json:"coordinates,omitempty"`
	CoordinatesObj interface{}     `json:"-"`
	GeometryRaw    json.RawMessage `json:"geometry,omitempty"`
	GeometryObj    *Member         `json:"-"`
	//GeometriesRaw  josn.RawMessage   `json:"geometries,omitempty"`
	//GeometriesObj  []Member          `json:"-"`
	Properties map[string]string `json:"properties,omitempty"`
	//BBox       [4]float64        `json:"bbox,omitempty"`
}

type Point [2]float64
type MultiPoint []Point

func (p Point) String() string {
	return fmt.Sprintf("[%v %v]", p[0], p[1])
}

type LineString []Point
type MultiLineString []LineString

type Polygon []LineString
type MultiPolygon []Polygon

/*
func (s LineString) String() string {
	var helper func(LineString) string
	helper = func(a LineString) string {
		if len(a) == 0 {
			return ""
		}
		return fmt.Sprintf(" %s%s", a[0], helper(a[1:]))
	}
	if len(s) == 0 {
		return "[]"
	}
	return fmt.Sprintf("[%v%v]", s[0], helper(s[1:]))
}
*/

func NewMember(b []byte) (*Member, error) {
	var member Member
	err := json.Unmarshal(b, &member)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal error: %v", err)
	}
	err = member.setProperties()
	if err != nil {
		return nil, fmt.Errorf("%v:%v", err, member)
	}
	return &member, nil
}

func NewMembers(b []byte) ([]*Member, error) {
	var members []*Member
	err := json.Unmarshal(b, &members)
	if err != nil {
		return nil, err
	}
	for _, member := range members {
		err := member.setProperties()
		if err != nil {
			return nil, fmt.Errorf("%v:%v", err, member)
		}
	}
	return members, nil
}

func (member *Member) setProperties() error {
	err := member.setObjectType()
	if err != nil {
		return fmt.Errorf("%v:%v", err, member)
	}
	return nil
}

func (member *Member) setObjectType() error {
	switch member.Type {
	case "Point", "LineString", "Polygon":
		member.memberKind = GeometryObject
		err := member.setCoordinatesObject()
		if err != nil {
			return err
		}
	case "MultiPoint", "MultiLineString", "MultiPolygon":
		member.memberKind = GeometryObject
		//TODO
	case "GeometryCollection":
		member.memberKind = GeometryObject
		//TODO
	case "Feature":
		member.memberKind = FeatureObject
		err := member.setGeometryObject()
		if err != nil {
			return err
		}
	case "FeatureCollection":
		member.memberKind = FeatureCollectionObject
		//TODO
	default:
		return fmt.Errorf("Unknown type: %v", member.Type)
	}

	return nil
}

func (member *Member) setCoordinatesObject() error {
	var object interface{}
	switch member.Type {
	case "Point":
		object = new(Point)
	case "LineString":
		object = new(LineString)
	case "Polygon":
		object = new(Polygon)
	default:
		return fmt.Errorf("Unknown type: %v", member.Type)
	}
	err := json.Unmarshal(member.CoordinatesRaw, &object)
	if err != nil {
		return fmt.Errorf("Unmarshal error:%v coordinates:%s", err, member.CoordinatesRaw)
	}
	//fmt.Printf("object:%v\n", object)
	member.CoordinatesObj = object
	return nil
}

func (member *Member) setGeometryObject() error {
	geometry, err := NewMember(member.GeometryRaw)
	if err != nil {
		return fmt.Errorf("error:%v geometry:%s", err, member.GeometryRaw)
	}
	member.GeometryObj = geometry
	return nil
}

func (member *Member) String() string {
	switch member.memberKind {
	case GeometryObject:
		return fmt.Sprintf("type:%v coordinates:%v", member.Type, reflect.ValueOf(member.CoordinatesObj).Elem())
	case FeatureObject:
		return fmt.Sprintf("type:%v geometry:%v properties:%v", member.Type, member.GeometryObj, member.Properties)
	case FeatureCollectionObject:
		//TODO
		return ""
	default:
		return fmt.Sprintf("Unknown Object Type:%v", member.memberKind)
	}
}
