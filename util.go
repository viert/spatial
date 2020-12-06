package spatial

import "github.com/dhconnelly/rtreego"

func MakeBounds(
	southWestLat float64,
	southWestLng float64,
	northEastLat float64,
	northEastLng float64,
) (*rtreego.Rect, error) {
	p := rtreego.Point{southWestLng, southWestLat}
	return rtreego.NewRect(p, []float64{northEastLng - southWestLng, northEastLat - southWestLat})
}
