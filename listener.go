package spatial

import (
	"fmt"

	"github.com/dhconnelly/rtreego"
)

var (
	listenerInc = 0
)

// UpdateAction is a enum describing update actions
type UpdateAction int

// UpdateAction enum definition
const (
	UAAdd UpdateAction = iota
	UARemove
	UAUpdate
)

// Update is an object transmitting via listener's chan
type Update struct {
	Object Indexable
	Action UpdateAction
}

func (u Update) String() string {
	action := ""
	switch u.Action {
	case UAAdd:
		action = "Add"
	case UARemove:
		action = "Remove"
	case UAUpdate:
		action = "Update"
	}
	return fmt.Sprintf(action+" object %v, %v", u.Object.ID(), u.Object.Bounds())
}

// Listener listens to updates in index
type Listener struct {
	bounds *rtreego.Rect
	u      chan Update
	id     string
	srv    *Server
}

// ID implements Indexable
func (l *Listener) ID() string {
	return l.id
}

// Bounds implements Indexable
func (l *Listener) Bounds() *rtreego.Rect {
	return l.bounds
}

// Updates returns the data channel of the listener
func (l *Listener) Updates() <-chan Update {
	return l.u
}

func (l *Listener) remove(obj Indexable) {
	l.u <- Update{obj, UARemove}
}

func (l *Listener) update(obj Indexable) {
	l.u <- Update{obj, UAUpdate}
}

func (l *Listener) add(obj Indexable) {
	l.u <- Update{obj, UAAdd}
}

// Unsubscribe closes the channel and removes listener from the index
func (l *Listener) Unsubscribe() {
	l.srv.removeListener(l.id)
	close(l.u)
}
