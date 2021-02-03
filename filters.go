package spatial

import (
	"github.com/dhconnelly/rtreego"
)

// FilterByTypes creates an rtreego.Filter to filter search results by indexable object type
func FilterByTypes(types []IndexableType) rtreego.Filter {
	typeMap := make(map[IndexableType]bool)
	for _, t := range types {
		typeMap[t] = true
	}

	return func(results []rtreego.Spatial, obj rtreego.Spatial) (bool, bool) {
		idxbl, ok := obj.(Indexable)
		if !ok {
			return true, false
		}

		itype := idxbl.Type()
		_, found := typeMap[itype]
		return !found, false
	}
}

var (
	filterBoundingBoxes = FilterByTypes([]IndexableType{itBoundingBox})
)
