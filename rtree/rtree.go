package rtree

import (
	"sync"

	"github.com/dhconnelly/rtreego"
)

// SafeRtree is a thread-safe wrapper for rtreego Rtree
type SafeRtree struct {
	tree *rtreego.Rtree
	lock sync.RWMutex
}

// New creates a new SafeRtree
func New(dim int, minBranch int, maxBranch int) *SafeRtree {
	return &SafeRtree{
		tree: rtreego.NewTree(dim, minBranch, maxBranch),
	}
}

// SearchIntersect searches for objects intersecting with a given Rect
func (srt *SafeRtree) SearchIntersect(bb *rtreego.Rect, filters ...rtreego.Filter) []rtreego.Spatial {
	srt.lock.RLock()
	defer srt.lock.RUnlock()
	return srt.tree.SearchIntersect(bb, filters...)
}

// Delete removes an object from Rtree
func (srt *SafeRtree) Delete(obj rtreego.Spatial) bool {
	srt.lock.Lock()
	defer srt.lock.Unlock()
	return srt.tree.Delete(obj)
}

// Insert inserts an object to Rtree
func (srt *SafeRtree) Insert(obj rtreego.Spatial) {
	srt.lock.Lock()
	defer srt.lock.Unlock()
	srt.tree.Insert(obj)
}
