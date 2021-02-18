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
	meta    map[string]string
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

// Meta is Object meta getter
func (o *Object) Meta(key string) string {
	return o.meta[key]
}

// HasMetaKey is Object meta key checker
func (o *Object) HasMetaKey(key string) bool {
	_, found := o.meta[key]
	return found
}

// NewObject creates a new instance of Object
func NewObject(id string, objType IndexableType, bounds *rtreego.Rect, ref interface{}, meta map[string]string) *Object {
	if meta == nil {
		meta = make(map[string]string)
	}

	return &Object{
		id:      id,
		objType: objType,
		bounds:  bounds,
		ref:     ref,
		meta:    meta,
	}
}
