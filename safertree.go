package spatial

import (
	"github.com/dhconnelly/rtreego"
	"sync"
)

type SafeRtree struct {
	tree *rtreego.Rtree
	lock sync.RWMutex
}

func NewSafeRtree(dim int, minBranch int, maxBranch int) *SafeRtree {
	return &SafeRtree{
		tree: rtreego.NewTree(dim, minBranch, maxBranch),
	}
}

func (srt *SafeRtree) SearchIntersect(bb *rtreego.Rect, filters ...rtreego.Filter) []rtreego.Spatial {
	srt.lock.RLock()
	defer srt.lock.RUnlock()
	return srt.tree.SearchIntersect(bb, filters...)
}

func (srt *SafeRtree) Delete(obj rtreego.Spatial) bool {
	srt.lock.Lock()
	defer srt.lock.Unlock()
	return srt.tree.Delete(obj)
}

func (srt *SafeRtree) Insert(obj rtreego.Spatial) {
	srt.lock.Lock()
	defer srt.lock.Unlock()
	srt.tree.Insert(obj)
}