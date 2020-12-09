package spatial

import "time"

// UpdateAction is a enum describing update actions
type UpdateAction int

// Listener listens to updates in index
type Listener struct {
	ch      chan []Indexable
	srv     *Server
	boxes   []*watchBox
	dirty   bool
	stopped bool
}

// Updates returns the update channel of the listener
func (l *Listener) Updates() <-chan []Indexable {
	return l.ch
}

// Unsubscribe closes the channel and removes listener from the index
func (l *Listener) Unsubscribe() {
	l.stopped = true
	l.srv.removeBoxes(l.boxes)
	close(l.ch)
}

func (l *Listener) loop() {
	t := time.NewTicker(l.srv.interval)
	defer t.Stop()

	for range t.C {
		if l.stopped {
			break
		}
		if l.dirty {
			objmap := l.srv.findObjects(l.boxes)
			objects := make([]Indexable, 0)
			for _, obj := range objmap {
				objects = append(objects, obj)
			}
			l.ch <- objects
			l.dirty = false
		}
	}
}
