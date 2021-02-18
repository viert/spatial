package spatial

import "github.com/dhconnelly/rtreego"

// IndexableType is a type for enums describing types of objects in index tree.
// Negative numbers are reserved for internal use
type IndexableType int

const (
	itBoundingBox IndexableType = -1
)

// Indexable interface
type Indexable interface {
	rtreego.Spatial
	ID() string
	Ref() interface{}
	Type() IndexableType
}

// Object is an Indexable implementation helper type
type Object struct {
	id      string
	ref     interface{}
	bounds  *rtreego.Rect
	objType IndexableType
}

// ID implements Indexable
func (o *Object) ID() string {
	return o.id
}

// Bounds implements Indexable
func (o *Object) Bounds() *rtreego.Rect {
	return o.bounds
}

// Ref implements Indexable
func (o *Object) Ref() interface{} {
	return o.ref
}

// Type implements Indexable
func (o *Object) Type() IndexableType {
	return o.objType
}

// NewObject creates a new instance of Object
func NewObject(id string, objType IndexableType, bounds *rtreego.Rect, ref interface{}) *Object {
	return &Object{
		id:      id,
		objType: objType,
		bounds:  bounds,
		ref:     ref,
	}
}
