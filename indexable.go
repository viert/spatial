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
