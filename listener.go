package spatial

import (
	"sync"
	"time"

	"github.com/dhconnelly/rtreego"
)

// Listener is an object watching for objects in bounding boxes and with specific ids
// and sends updates through a channel
type Listener struct {
	lock           sync.RWMutex
	srv            *Server
	ch             chan []Indexable
	boxes          []*boundingBox
	filter         rtreego.Filter
	watchIds       map[string]bool
	updateInterval time.Duration
	dirty          bool
	stopped        bool
}

func newListener(srv *Server, chSize int, interval time.Duration) *Listener {
	lstr := &Listener{
		srv:            srv,
		ch:             make(chan []Indexable, chSize),
		boxes:          make([]*boundingBox, 0),
		filter:         nil,
		watchIds:       make(map[string]bool),
		updateInterval: interval,
		stopped:        false,
		dirty:          false,
	}
	go lstr.loop()
	return lstr
}

func (l *Listener) disposeBoxes() {
	// No locking here as l.boxes is always overwritten and never changed
	for _, box := range l.boxes {
		l.srv.tree.Delete(box)
	}
}

func (l *Listener) unsubscribeAll() {
	// Hoping that l.srv.unsubscribeID never calls back to the listener
	// or at least not to a method that tries to acquire the lock
	l.lock.RLock()
	defer l.lock.RUnlock()
	for id := range l.watchIds {
		l.srv.unsubscribeID(l, id)
	}
}

// SetBounds sets bounds to listen to
func (l *Listener) SetBounds(mb MapBounds) {
	l.disposeBoxes()

	rects := mb.Rects()
	boxes := make([]*boundingBox, len(rects))
	for i, rect := range rects {
		box := newBoundingBox(rect, l)
		l.srv.tree.Insert(box)
		boxes[i] = box
	}
	l.boxes = boxes
}

// SetTypes sets a filter to listen for objects of specified types only
// Does not apply for ID subscriptions
func (l *Listener) SetTypes(types []IndexableType) {
	if len(types) == 0 {
		l.filter = nil
	} else {
		l.filter = FilterByTypes(types)
	}
}

// Stop stops the listener, closes all the channels so it's free to cleanup by GC
func (l *Listener) Stop() {
	l.disposeBoxes()
	l.unsubscribeAll()
	l.stopped = true
}

// SubscribeID adds a specific id to watch
func (l *Listener) SubscribeID(id string) {
	l.lock.Lock()
	l.watchIds[id] = true
	l.lock.Unlock()
	l.srv.subscribeID(l, id)
}

// UnsubscribeID unsubscribes from a specific id
func (l *Listener) UnsubscribeID(id string) {
	l.lock.Lock()
	delete(l.watchIds, id)
	l.lock.Unlock()
	l.srv.unsubscribeID(l, id)
}

// ForceUpdate forces the dirty flag on
func (l *Listener) ForceUpdate() {
    l.setDirty()
}

func (l *Listener) setDirty() {
	l.dirty = true
}

// Updates returns the update channel
func (l *Listener) Updates() <-chan []Indexable {
	return l.ch
}

func (l *Listener) loop() {
	var rmap map[string]Indexable
	t := time.NewTicker(l.updateInterval)
	defer t.Stop()

	for range t.C {
		if l.stopped {
			break
		}

		if l.dirty {
			objmap := make(map[string]Indexable)

			l.lock.RLock()
			for key, obj := range l.srv.findObjectsByIDs(l.watchIds) {
				objmap[key] = obj
			}
			l.lock.RUnlock()

			if l.filter == nil {
				rmap = l.srv.findObjectsByBoundingBoxes(l.boxes)
			} else {
				rmap = l.srv.findObjectsByBoundingBoxes(l.boxes, l.filter)
			}

			for key, obj := range rmap {
				objmap[key] = obj
			}

			objects := make([]Indexable, 0)
			for _, obj := range objmap {
				objects = append(objects, obj)
			}

			l.ch <- objects
			l.dirty = false
		}
	}

	close(l.ch)
}
