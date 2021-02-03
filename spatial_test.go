package spatial

import (
	"math"
	"testing"
	"time"

	"github.com/dhconnelly/rtreego"
)

const (
	itUserObject  IndexableType = 1
	itUserObject2 IndexableType = 2

	testObjectID = "obj1"
)

var (
	testBounds = MapBounds{
		SouthWestLng: -10,
		SouthWestLat: -10,
		NorthEastLng: 10,
		NorthEastLat: 10,
	}
)

type (
	object struct {
		id      string
		rect    *rtreego.Rect
		objType IndexableType
	}
)

func (o *object) ID() string {
	return o.id
}

func (o *object) Type() IndexableType {
	return o.objType
}

func (o *object) Bounds() *rtreego.Rect {
	return o.rect
}

func (o *object) Ref() interface{} {
	return nil
}

func newObject(oType IndexableType, id string, lat float64, lng float64) *object {
	p := rtreego.Point{lng, lat}
	rect, _ := rtreego.NewRect(p, []float64{0.1, 0.1})

	return &object{id, rect, oType}
}

func newRectObject(oType IndexableType, id string, rect *rtreego.Rect) *object {
	return &object{id, rect, oType}
}

func timeout(n time.Duration) <-chan time.Time {
	t := time.NewTimer(n * time.Millisecond)
	return t.C
}

func getUpdates(ch <-chan []Indexable) []Indexable {
	for {
		select {
		case <-timeout(50):
			return nil
		case u := <-ch:
			return u
		}
	}
}

func eq(a, b float64) bool {
	tolerance := 0.0001
	diff := math.Abs(a - b)
	return diff < tolerance
}

func checkCoords(rect *rtreego.Rect, x float64, y float64) bool {
	return eq(rect.PointCoord(0), x) && eq(rect.PointCoord(1), y)
}

func TestBoundsUpdate(t *testing.T) {
	var updates []Indexable
	var ok bool

	srv := New(25, 50)
	lst := srv.NewListener(100, 10*time.Millisecond)
	defer lst.Stop()

	lst.SetBounds(testBounds)
	ch := lst.Updates()

	// make sure there are no updates yet
	updates = getUpdates(ch)
	if updates != nil {
		t.Errorf("unexpected update: %s", updates)
		return
	}

	obj := newObject(itUserObject, testObjectID, 0, 0)
	srv.Add(obj)

	updates = getUpdates(ch)
	if updates == nil {
		t.Errorf("one update expected, got nil")
		return
	}

	if len(updates) != 1 {
		t.Errorf("one update expected, got %d", len(updates))
		return
	}

	if obj, ok = updates[0].(*object); !ok {
		t.Errorf("error asserting update type as object")
		return
	}

	if !checkCoords(obj.rect, 0, 0) {
		t.Errorf("expected coords %.3f/%.3f, got %.3f/%.3f",
			0.0, 0.0,
			obj.rect.PointCoord(0), obj.rect.PointCoord(1),
		)
		return
	}

	obj = newObject(itUserObject, testObjectID, 3, 3)
	srv.Add(obj)

	updates = getUpdates(ch)
	if updates == nil {
		t.Errorf("one update expected, got nil")
		return
	}

	if len(updates) != 1 {
		t.Errorf("one update expected, got %d", len(updates))
		return
	}

	if obj, ok = updates[0].(*object); !ok {
		t.Errorf("error asserting update type as object")
		return
	}

	if !checkCoords(obj.rect, 3, 3) {
		t.Errorf("expected coords %.3f/%.3f, got %.3f/%.3f",
			3.0, 3.0,
			obj.rect.PointCoord(0), obj.rect.PointCoord(1),
		)
		return
	}

	// move outside
	obj = newObject(itUserObject, testObjectID, 12, 12)
	srv.Add(obj)

	updates = getUpdates(ch)
	if updates == nil {
		t.Errorf("zero updates expected, got nil")
		return
	}

	if len(updates) != 0 {
		t.Errorf("zero updates expected, got %d", len(updates))
		return
	}

}

func TestWatchID(t *testing.T) {
	var updates []Indexable
	var ok bool

	srv := New(25, 50)
	lst := srv.NewListener(100, 10*time.Millisecond)
	defer lst.Stop()

	lst.SetBounds(testBounds)
	lst.SubscribeID(testObjectID)

	ch := lst.Updates()

	// make sure there are no updates yet
	updates = getUpdates(ch)
	if updates != nil {
		t.Errorf("unexpected update: %s", updates)
		return
	}

	obj := newObject(itUserObject, testObjectID, 0, 0)
	srv.Add(obj)

	updates = getUpdates(ch)
	if updates == nil {
		t.Errorf("one update expected, got nil")
		return
	}

	if len(updates) != 1 {
		t.Errorf("one update expected, got %d", len(updates))
		return
	}

	if obj, ok = updates[0].(*object); !ok {
		t.Errorf("error asserting update type as object")
		return
	}

	if !checkCoords(obj.rect, 0, 0) {
		t.Errorf("expected coords %.3f/%.3f, got %.3f/%.3f",
			0.0, 0.0,
			obj.rect.PointCoord(0), obj.rect.PointCoord(1),
		)
		return
	}

	obj = newObject(itUserObject, testObjectID, 3, 3)
	srv.Add(obj)

	updates = getUpdates(ch)
	if updates == nil {
		t.Errorf("one update expected, got nil")
		return
	}

	if len(updates) != 1 {
		t.Errorf("one update expected, got %d", len(updates))
		return
	}

	if obj, ok = updates[0].(*object); !ok {
		t.Errorf("error asserting update type as object")
		return
	}

	if !checkCoords(obj.rect, 3, 3) {
		t.Errorf("expected coords %.3f/%.3f, got %.3f/%.3f",
			3.0, 3.0,
			obj.rect.PointCoord(0), obj.rect.PointCoord(1),
		)
		return
	}

	// move outside
	obj = newObject(itUserObject, testObjectID, 12, 12)
	srv.Add(obj)

	updates = getUpdates(ch)
	if updates == nil {
		t.Errorf("one update expected, got nil")
		return
	}

	if len(updates) != 1 {
		t.Errorf("one update expected, got %d", len(updates))
		return
	}

	if obj, ok = updates[0].(*object); !ok {
		t.Errorf("error asserting update type as object")
		return
	}

	if !checkCoords(obj.rect, 12, 12) {
		t.Errorf("expected coords %.3f/%.3f, got %.3f/%.3f",
			12.0, 12.0,
			obj.rect.PointCoord(0), obj.rect.PointCoord(1),
		)
		return
	}

}

func TestFilters(t *testing.T) {
	srv := New(25, 50)
	lst := srv.NewListener(100, time.Second)
	defer lst.Stop()

	obj := newObject(itUserObject, testObjectID, 0, 0)
	srv.Add(obj)

	obj2 := newObject(itUserObject2, "obj2", 0.1, 0.1)
	srv.Add(obj2)

	p := rtreego.Point{-1.0, -1.0}
	rect, _ := rtreego.NewRect(p, []float64{2.0, 2.0})

	boxObj := newRectObject(itBoundingBox, "box1", rect)
	srv.Add(boxObj)

	p = rtreego.Point{-30.0, -30.0}
	rect, _ = rtreego.NewRect(p, []float64{60.0, 60.0})
	bounds := newBoundingBox(rect, lst)
	srv.Add(bounds)

	results := srv.findObjectsByBoundingBoxes([]*boundingBox{bounds})

	if len(results) != 2 {
		t.Errorf("expected exactly 2 objects, but %d were found", len(results))
		return
	}

	for id := range results {
		if id != obj.id && id != obj2.id {
			t.Errorf("unexpected object with id \"%s\" found", id)
		}
	}

	results = srv.findObjectsByBoundingBoxes(
		[]*boundingBox{bounds},
		FilterByTypes([]IndexableType{itUserObject}),
	)

	if len(results) != 1 {
		t.Errorf("expected exactly 1 object, but %d were found", len(results))
		return
	}

	if _, found := results[obj.id]; !found {
		t.Errorf("object %s is expected to be in results", obj.id)
	}

	results = srv.findObjectsByBoundingBoxes(
		[]*boundingBox{bounds},
		FilterByTypes([]IndexableType{itUserObject2}),
	)

	if len(results) != 1 {
		t.Errorf("expected exactly 1 object, but %d were found", len(results))
		return
	}

	if _, found := results[obj2.id]; !found {
		t.Errorf("object %s is expected to be in results", obj.id)
	}

}
