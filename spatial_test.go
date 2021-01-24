package spatial

import (
	"math"
	"testing"
	"time"

	"github.com/dhconnelly/rtreego"
)

const (
	itUserObject IndexableType = 1
	testObjectID               = "obj1"
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
		id   string
		rect *rtreego.Rect
	}
)

func (o *object) ID() string {
	return o.id
}

func (o *object) Type() IndexableType {
	return itUserObject
}

func (o *object) Bounds() *rtreego.Rect {
	return o.rect
}

func (o *object) Ref() interface{} {
	return nil
}

func newObject(id string, lat float64, lng float64) *object {
	p := rtreego.Point{lng, lat}
	rect, _ := rtreego.NewRect(p, []float64{0.1, 0.1})

	return &object{id, rect}
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

	obj := newObject(testObjectID, 0, 0)
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

	obj = newObject(testObjectID, 3, 3)
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
	obj = newObject(testObjectID, 12, 12)
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

	obj := newObject(testObjectID, 0, 0)
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

	obj = newObject(testObjectID, 3, 3)
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
	obj = newObject(testObjectID, 12, 12)
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
