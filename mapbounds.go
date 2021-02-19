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
func (mb *MapBounds) rect() *rtreego.Rect {
	point := rtreego.Point{mb.SouthWestLat, mb.SouthWestLng}
	lngSize := mb.NorthEastLng - mb.SouthWestLng
	latSize := mb.NorthEastLat - mb.SouthWestLat
	rect, _ := rtreego.NewRect(point, []float64{latSize, lngSize})
	return rect
}

func (mb *MapBounds) split() []*MapBounds {
	boxes := make([]*MapBounds, 1)
	boxes[0] = mb

	if mb.SouthWestLng > mb.NorthEastLng {
		temp := make([]*MapBounds, 0)
		for _, box := range boxes {
			// western box
			temp = append(temp, &MapBounds{box.SouthWestLng, box.SouthWestLat, eastmostLongintude, box.NorthEastLat})
			// eastern box
			temp = append(temp, &MapBounds{-eastmostLongintude, box.SouthWestLat, box.NorthEastLng, box.NorthEastLat})
		}
		boxes = temp
	}

	if mb.SouthWestLat > mb.NorthEastLat {
		temp := make([]*MapBounds, 0)
		for _, box := range boxes {
			// northern box
			temp = append(temp, &MapBounds{box.SouthWestLng, box.SouthWestLat, box.NorthEastLng, northmostLatitude})
			// southern box
			temp = append(temp, &MapBounds{box.SouthWestLng, -northmostLatitude, box.NorthEastLng, box.NorthEastLat})
		}
		boxes = temp
	}

	return boxes
}

// Rects returns a list of Rects supporint latitude/longitude wrapping
func (mb *MapBounds) Rects() []*rtreego.Rect {
	boxes := mb.split()
	rects := make([]*rtreego.Rect, len(boxes))
	for i, box := range boxes {
		rects[i] = box.rect()
	}
	return rects
}
