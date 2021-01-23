package spatial

import (
	"sync"
	"time"

	"github.com/dhconnelly/rtreego"
)

// Server represents spatial index server
type Server struct {
	tree   *rtreego.Rtree
	idSubs map[string]map[*Listener]*Listener
	idIdx  map[string]Indexable
	lock   sync.RWMutex
}

// New creates and initializes a new spatial Server
func New(minBranch int, maxBranch int, updateChanSize int) *Server {
	t := rtreego.NewTree(2, minBranch, maxBranch)
	return &Server{
		tree:   t,
		idSubs: make(map[string]map[*Listener]*Listener),
		idIdx:  make(map[string]Indexable),
	}
}

func (s *Server) subscribeID(l *Listener, id string) {
	if _, found := s.idSubs[id]; !found {
		s.idSubs[id] = make(map[*Listener]*Listener)
	}
	s.idSubs[id][l] = l
}

func (s *Server) unsubscribeID(l *Listener, id string) {
	if idmap, found := s.idSubs[id]; found {
		if _, found = idmap[l]; found {
			delete(idmap, l)
		}
	}
}

func (s *Server) findObjectsByIDs(ids map[string]bool) map[string]Indexable {
	s.lock.RLock()
	defer s.lock.RUnlock()

	results := make(map[string]Indexable)
	for id := range ids {
		if obj, found := s.idIdx[id]; found {
			results[obj.ID()] = obj
		}
	}
	return results
}

func (s *Server) findObjectsByBoundingBoxes(boxes []*boundingBox) map[string]Indexable {
	s.lock.RLock()
	defer s.lock.RUnlock()

	results := make(map[string]Indexable)
	for _, box := range boxes {
		rect := box.bounds
		spatials := s.tree.SearchIntersect(rect)
		for _, sp := range spatials {
			if idxbl, ok := sp.(Indexable); ok {
				if idxbl.Type() > 0 {
					results[idxbl.ID()] = idxbl
				}
			}
		}
	}

	return results
}

func (s *Server) findBoundingBoxesByObject(idx Indexable) []boundingBox {
	intersections := s.tree.SearchIntersect(idx.Bounds())
	boxes := make([]boundingBox, 0)
	for _, obj := range intersections {
		idxbl, ok := obj.(Indexable)
		if ok && idxbl.Type() == itBoundingBox {
			if box, ok := idxbl.(*boundingBox); ok {
				boxes = append(boxes, *box)
			}
		}
	}
	return boxes
}

// Add adds a new object if it doesn't exist (checking by it's ID())
// or modifies existing one, and notifies listeners
func (s *Server) Add(obj Indexable) {
	var rmListeners map[*Listener]*Listener
	var addListeners map[*Listener]*Listener

	s.lock.Lock()
	defer s.lock.Unlock()

	if curr, found := s.idIdx[obj.ID()]; found {
		// collect listeners to remove obj from
		boxes := s.findBoundingBoxesByObject(curr)
		rmListeners = collectListeners(boxes)
		s.tree.Delete(curr)
	}
	s.idIdx[obj.ID()] = obj

	s.tree.Insert(obj)

	boxes := s.findBoundingBoxesByObject(obj)
	addListeners = collectListeners(boxes)

	for l := range rmListeners {
		l.setDirty()
	}
	for l := range addListeners {
		l.setDirty()
	}

	if lmap, found := s.idSubs[obj.ID()]; found {
		for l := range lmap {
			l.setDirty()
		}
	}
}

// Remove removes a given object from the index and notifies listeners
func (s *Server) Remove(obj Indexable) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if curr, found := s.idIdx[obj.ID()]; found {
		// collect listeners to remove obj from
		boxes := s.findBoundingBoxesByObject(curr)
		listeners := collectListeners(boxes)
		s.tree.Delete(curr)

		for l := range listeners {
			l.setDirty()
		}

		if lmap, found := s.idSubs[obj.ID()]; found {
			for l := range lmap {
				l.setDirty()
			}
		}
	}
}

// NewListener creates and returns a new listener
func (s *Server) NewListener(chSize int, interval time.Duration) *Listener {
	return newListener(s, chSize, interval)
}
