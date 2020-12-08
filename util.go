package spatial

import "github.com/dhconnelly/rtreego"

const (
	eastmostLongintude = 179.9999999
	northmostLatitude  = 89.9999999
)

// MapBounds represents a world coordinates bounding box
type MapBounds struct {
	SouthWestLng float64
	SouthWestLat float64
	NorthEastLng float64
	NorthEastLat float64
}

// Rect converts MapBounds to rtregoo Rect
func (mb MapBounds) Rect() *rtreego.Rect {
	point := rtreego.Point{mb.SouthWestLng, mb.SouthWestLat}
	width := mb.NorthEastLng - mb.SouthWestLng
	height := mb.NorthEastLat - mb.SouthWestLat
	rect, _ := rtreego.NewRect(point, []float64{width, height})
	return rect
}

func split(c MapBounds) []MapBounds {
	boxes := make([]MapBounds, 1)
	boxes[0] = c

	if c.SouthWestLng > c.NorthEastLng {
		temp := make([]MapBounds, 0)
		for _, box := range boxes {
			// western box
			temp = append(temp, MapBounds{box.SouthWestLng, box.SouthWestLat, eastmostLongintude, box.NorthEastLat})
			// eastern box
			temp = append(temp, MapBounds{-eastmostLongintude, box.SouthWestLat, box.NorthEastLng, box.NorthEastLat})
		}
		boxes = temp
	}

	if c.SouthWestLat > c.NorthEastLat {
		temp := make([]MapBounds, 0)
		for _, box := range boxes {
			// northern box
			temp = append(temp, MapBounds{box.SouthWestLng, box.SouthWestLat, box.NorthEastLng, northmostLatitude})
			// southern box
			temp = append(temp, MapBounds{box.SouthWestLng, -northmostLatitude, box.NorthEastLng, box.NorthEastLat})
		}
		boxes = temp
	}

	return boxes
}

// GetBoundingBoxes returns a list of Rects supporint latitude/longitude wrapping
func getBoundingBoxes(mb MapBounds) []*rtreego.Rect {
	boxes := split(mb)
	rects := make([]*rtreego.Rect, len(boxes))
	for i, box := range boxes {
		rects[i] = box.Rect()
	}
	return rects
}
