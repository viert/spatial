package object

import (
	"github.com/dhconnelly/rtreego"
	"github.com/viert/spatial"
)

// Object is an Indexable implementation helper type
type Object struct {
	id      string
	ref     interface{}
	bounds  *rtreego.Rect
	objType spatial.IndexableType
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
func (o *Object) Type() spatial.IndexableType {
	return o.objType
}

// New creates a new instance of Object
func New(id string, objType spatial.IndexableType, bounds *rtreego.Rect, ref interface{}) *Object {
	return &Object{
		id:      id,
		objType: objType,
		bounds:  bounds,
		ref:     ref,
	}
}
