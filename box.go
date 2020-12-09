package spatial

import "github.com/dhconnelly/rtreego"

type watchBox struct {
	rect     *rtreego.Rect
	id       string
	srv      *Server
	listener *Listener
}

var (
	wbAutoinc = 0
)

// ID to implement Indexable
func (w *watchBox) ID() string {
	return w.id
}

// Bounds to implement Indexable
func (w *watchBox) Bounds() *rtreego.Rect {
	return w.rect
}

func (w *watchBox) Ref() interface{} {
	return nil
}
