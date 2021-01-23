package spatial

import (
	"sync"
	"time"
)

// Listener is an object watching for objects in bounding boxes and with specific ids
// and sends updates through a channel
type Listener struct {
	lock           sync.RWMutex
	srv            *Server
	ch             chan []Indexable
	boxes          []*boundingBox
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
		watchIds:       make(map[string]bool),
		updateInterval: interval,
		stopped:        false,
		dirty:          false,
	}
	go lstr.loop()
	return lstr
}

// SetBounds sets bounds to listen to
func (l *Listener) SetBounds(mb MapBounds) {
	l.lock.Lock()
	l.srv.lock.Lock()
	defer l.lock.Unlock()
	defer l.srv.lock.Unlock()

	for _, box := range l.boxes {
		l.srv.tree.Delete(box)
	}

	rects := mb.Rects()
	boxes := make([]*boundingBox, len(rects))
	for i, rect := range rects {
		box := newBoundingBox(rect, l)
		l.srv.tree.Insert(box)
		boxes[i] = box
	}
	l.boxes = boxes
}

// Stop stops the listener, closes all the channels so it's free to cleanup by GC
func (l *Listener) Stop() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.stopped = true
}

// SubscribeID adds a specific id to watch
func (l *Listener) SubscribeID(id string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.watchIds[id] = true
	l.srv.subscribeID(l, id)
}

// UnsubscribeID unsubscribes from a specific id
func (l *Listener) UnsubscribeID(id string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	delete(l.watchIds, id)
	l.srv.unsubscribeID(l, id)
}

func (l *Listener) setDirty() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.dirty = true
}

// Updates returns the update channel
func (l *Listener) Updates() <-chan []Indexable {
	return l.ch
}

func (l *Listener) loop() {
	t := time.NewTicker(l.updateInterval)
	defer t.Stop()

	for range t.C {
		l.lock.RLock()
		if l.stopped {
			l.lock.RUnlock()
			break
		}

		if l.dirty {
			objmap := make(map[string]Indexable)

			for key, obj := range l.srv.findObjectsByIDs(l.watchIds) {
				objmap[key] = obj
			}

			for key, obj := range l.srv.findObjectsByBoundingBoxes(l.boxes) {
				objmap[key] = obj
			}

			objects := make([]Indexable, 0)
			for _, obj := range objmap {
				objects = append(objects, obj)
			}

			l.ch <- objects
			l.dirty = false
		}
		l.lock.RUnlock()
	}

	close(l.ch)
}
