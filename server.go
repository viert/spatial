package spatial

import (
	"fmt"
	"sync"
	"time"

	"github.com/dhconnelly/rtreego"
)

const (
	userObjectPrefix = "u:"
)

// Indexable is an interface for objects stored in spacial index
type Indexable interface {
	ID() string
	Bounds() *rtreego.Rect
	Ref() interface{}
}

// Server is the main 2D index server object
type Server struct {
	tree     *rtreego.Rtree
	idIdx    map[string]Indexable
	chSize   int
	lock     sync.RWMutex
	interval time.Duration
}

// New creates a new server. minBranch and maxBranch are the RTree branching properties
// Refer to https://github.com/dhconnelly/rtreego
func New(minBranch int, maxBranch int, updateChanSize int, notifyInterval time.Duration) *Server {
	t := rtreego.NewTree(2, minBranch, maxBranch)
	return &Server{
		tree:     t,
		idIdx:    make(map[string]Indexable),
		chSize:   updateChanSize,
		interval: notifyInterval,
	}
}

func (s *Server) update(obj Indexable) Indexable {
	var result Indexable

	id := obj.ID()
	if existing, found := s.idIdx[id]; found {
		s.tree.Delete(existing)
		result = existing
	}
	s.tree.Insert(obj)
	s.idIdx[id] = obj

	return result
}

// Add adds an object of a given size and given coordinates to index or modifies
// an existing one if the object with the same ID is present
func (s *Server) Add(
	lat float64,
	lng float64,
	width float64,
	height float64,
	id string,
	ref interface{},
) (*Object, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	var rmLstrs map[*Listener]*Listener
	var addLstrs map[*Listener]*Listener

	obj, err := newObject(lat, lng, width, height, id, ref)
	if err != nil {
		return nil, err
	}

	prev := s.update(obj)
	if prev != nil {
		rmLstrs = s.findListeners(prev)
		for _, lst := range rmLstrs {
			lst.dirty = true
		}
	}

	addLstrs = s.findListeners(obj)
	for _, lst := range addLstrs {
		lst.dirty = true
	}

	return obj, nil
}

// Remove removes the object by id and returns true if it was actually deleted
func (s *Server) Remove(id string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	obj, found := s.idIdx[id]

	if found {
		lstrs := s.findListeners(obj)
		for _, l := range lstrs {
			l.dirty = true
		}
		s.tree.Delete(obj)
		delete(s.idIdx, id)
	}
	return found
}

// Subscribe returns a listener with a channel transmitting index updates
func (s *Server) Subscribe(bounds MapBounds) *Listener {
	s.lock.Lock()
	defer s.lock.Unlock()

	rects := getBoundingBoxes(bounds)
	boxes := make([]*watchBox, len(rects))

	lstr := &Listener{
		ch:      make(chan []Indexable, s.chSize),
		srv:     s,
		dirty:   true,
		stopped: false,
	}

	for i, rect := range rects {
		wbAutoinc++
		id := fmt.Sprintf("wb:%d", wbAutoinc)
		boxes[i] = &watchBox{
			rect:     rect,
			id:       id,
			srv:      s,
			listener: lstr,
		}
		s.tree.Insert(boxes[i])
		s.idIdx[id] = boxes[i]
	}

	lstr.boxes = boxes
	go lstr.loop()

	return lstr
}

func (s *Server) findListeners(obj Indexable) map[*Listener]*Listener {
	objs := s.tree.SearchIntersect(obj.Bounds())

	lmap := make(map[*Listener]*Listener)

	for _, obj := range objs {
		wb, ok := obj.(*watchBox)
		if ok {
			lmap[wb.listener] = wb.listener
		}
	}

	return lmap
}

func (s *Server) findObjects(boxes []*watchBox) map[string]Indexable {
	objects := make(map[string]Indexable)
	for _, box := range boxes {
		indexables := s.tree.SearchIntersect(box.Bounds())
		for _, idxbl := range indexables {
			obj, ok := idxbl.(*Object)
			if ok {
				objects[obj.id] = obj
			}
		}
	}
	return objects
}

func (s *Server) removeBoxes(boxes []*watchBox) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, wb := range boxes {
		if _, found := s.idIdx[wb.id]; found {
			s.tree.Delete(wb)
			delete(s.idIdx, wb.id)
		}
	}
}
