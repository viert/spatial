package spatial

import (
	"sync"
	"time"

	"github.com/dhconnelly/rtreego"
	"github.com/viert/spatial/rtree"
)

// Server represents spatial index server
type Server struct {
	tree   *rtree.SafeRtree
	idSubs map[string]map[*Listener]*Listener
	idIdx  map[string]Indexable
	lock   sync.RWMutex
}

// New creates and initializes a new spatial Server
func New(minBranch int, maxBranch int) *Server {
	t := rtree.New(2, minBranch, maxBranch)
	return &Server{
		tree:   t,
		idSubs: make(map[string]map[*Listener]*Listener),
		idIdx:  make(map[string]Indexable),
	}
}

func (s *Server) subscribeID(l *Listener, id string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, found := s.idSubs[id]; !found {
		s.idSubs[id] = make(map[*Listener]*Listener)
	}
	s.idSubs[id][l] = l
}

func (s *Server) unsubscribeID(l *Listener, id string) {
	s.lock.Lock()
	defer s.lock.Unlock()
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

func (s *Server) findObjectsByBoundingBoxes(boxes []*boundingBox, filters ...rtreego.Filter) map[string]Indexable {
	s.lock.RLock()
	defer s.lock.RUnlock()

	results := make(map[string]Indexable)
	for _, box := range boxes {
		rect := box.bounds
		spatials := s.tree.SearchIntersect(rect, filters...)
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
	intersections := s.tree.SearchIntersect(idx.Bounds(), filterBoundingBoxes)
	boxes := make([]boundingBox, 0)
	for _, obj := range intersections {
		if idxbl, ok := obj.(Indexable); ok {
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

	s.lock.RLock()
	curr, found := s.idIdx[obj.ID()]
	s.lock.RUnlock()
	if found {
		// collect listeners to remove obj from
		boxes := s.findBoundingBoxesByObject(curr)
		rmListeners = collectListeners(boxes)
		s.tree.Delete(curr)
	}
	s.lock.Lock()
	s.idIdx[obj.ID()] = obj
	s.lock.Unlock()

	s.tree.Insert(obj)

	boxes := s.findBoundingBoxesByObject(obj)
	addListeners = collectListeners(boxes)

	for l := range rmListeners {
		l.setDirty()
	}
	for l := range addListeners {
		l.setDirty()
	}

	s.lock.RLock()
	lmap, found := s.idSubs[obj.ID()]
	s.lock.RUnlock()
	if found {
		for l := range lmap {
			l.setDirty()
		}
	}
}

// Remove removes a given object from the index and notifies listeners
func (s *Server) Remove(obj Indexable) {
	s.lock.RLock()
	curr, found := s.idIdx[obj.ID()]
	s.lock.RUnlock()
	if found {
		// collect listeners to remove obj from
		boxes := s.findBoundingBoxesByObject(curr)
		listeners := collectListeners(boxes)
		s.tree.Delete(curr)

		for l := range listeners {
			l.setDirty()
		}

		s.lock.RLock()
		lmap, found := s.idSubs[obj.ID()]
		s.lock.RUnlock()
		if found {
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
