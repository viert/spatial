package spatial

import (
	"fmt"
	"sync"

	"github.com/dhconnelly/rtreego"
)

const (
	userObjectPrefix = "u:"
)

// Indexable is an interface for objects stored in spacial index
type Indexable interface {
	ID() string
	Bounds() *rtreego.Rect
}

// Server is the main 2D index server object
type Server struct {
	tree   *rtreego.Rtree
	idIdx  map[string]Indexable
	chSize int
	lock   sync.RWMutex
}

// New creates a new server. minBranch and maxBranch are the RTree branching properties
// Refer to https://github.com/dhconnelly/rtreego
func New(minBranch int, maxBranch int, updateChanSize int) *Server {
	t := rtreego.NewTree(2, minBranch, maxBranch)
	return &Server{
		tree:   t,
		idIdx:  make(map[string]Indexable),
		chSize: updateChanSize,
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
	}
	addLstrs = s.findListeners(obj)

	if rmLstrs == nil {
		for _, l := range addLstrs {
			l.add(obj)
		}
	} else {
		addOnly := make([]*Listener, 0)
		removeOnly := make([]*Listener, 0)
		update := make([]*Listener, 0)

		for _, al := range addLstrs {
			if _, found := rmLstrs[al]; found {
				update = append(addOnly, al)
			} else {
				addOnly = append(addOnly, al)
			}
		}

		for _, rl := range rmLstrs {
			if _, found := addLstrs[rl]; !found {
				removeOnly = append(removeOnly, rl)
			}
		}

		for _, l := range addOnly {
			l.add(obj)
		}
		for _, l := range removeOnly {
			l.remove(obj)
		}
		for _, l := range update {
			l.update(obj)
		}
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
			l.remove(obj)
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
	lstr := new(Listener)
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

	lstr.u = make(chan Update, s.chSize)
	lstr.srv = s
	lstr.boxes = boxes

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

func (s *Server) removeListenerWatchBoxes(listener *Listener) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, wb := range listener.boxes {
		if _, found := s.idIdx[wb.id]; found {
			s.tree.Delete(wb)
			delete(s.idIdx, wb.id)
		}
	}
}
