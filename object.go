package spatial

import "github.com/dhconnelly/rtreego"

// Object represents a 2D object on map
type Object struct {
	rect *rtreego.Rect
	ref  interface{}
	id   string
}

// ID to implement Indexable
func (o *Object) ID() string {
	return o.id
}

// Bounds to implement Indexable
func (o *Object) Bounds() *rtreego.Rect {
	return o.rect
}

// Ref returns a ref to an object
func (o *Object) Ref() interface{} {
	return o.ref
}

// NewObject creates a new indexable object
func newObject(
	lat float64,
	lng float64,
	width float64,
	height float64,
	id string,
	ref interface{},
) (*Object, error) {
	pt := rtreego.Point{lng, lat}
	rect, err := rtreego.NewRect(pt, []float64{width, height})

	if err != nil {
		return nil, err
	}

	return &Object{
		rect: rect,
		ref:  ref,
		id:   id,
	}, nil
}
