package spatial

import (
	"fmt"
	"sync/atomic"

	"github.com/dhconnelly/rtreego"
)

type boundingBox struct {
	id       string
	bounds   *rtreego.Rect
	listener *Listener
}

var (
	bboxAutoID uint64 = 0
)

func newBoundingBox(bounds *rtreego.Rect, lstr *Listener) *boundingBox {
	id := atomic.AddUint64(&bboxAutoID, 1)
	return &boundingBox{
		id:       fmt.Sprintf("bbx:%d", id),
		bounds:   bounds,
		listener: lstr,
	}
}

func (b boundingBox) ID() string {
	return b.id
}

func (b boundingBox) Bounds() *rtreego.Rect {
	return b.bounds
}

func (b boundingBox) Type() IndexableType {
	return itBoundingBox
}

func (b boundingBox) Ref() interface{} {
	return nil
}

func collectListeners(boxes []boundingBox) map[*Listener]*Listener {
	lmap := make(map[*Listener]*Listener)
	for _, box := range boxes {
		lmap[box.listener] = box.listener
	}
	return lmap
}
