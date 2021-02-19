package helper

import (
	"math"

	"github.com/dhconnelly/rtreego"
)

func nmToLatLon(latSizeNM float64, lngSizeNM float64, atLatitude float64) []float64 {
	// for latitude 60nm is 1˚
	latSize := (1.0 / 60) * latSizeNM

	// for longitude it depends on current latitude
	// first let's convert latitude degree to radians
	latitudeRad := (atLatitude * 2 * math.Pi) / 360

	// calculate size
	lngSize := (1.0 / 60) * lngSizeNM
	// and make a latitude correction
	lngSize = lngSize / math.Abs(math.Cos(latitudeRad))
	return []float64{latSize, lngSize}
}

// Square makes square bounds of a given size
// Lat and Lng represent top left angle of the square
func Square(lat float64, lng float64, sizeNM float64) *rtreego.Rect {
	p := rtreego.Point{lat, lng}
	sizes := nmToLatLon(sizeNM, sizeNM, lat)
	rect, _ := rtreego.NewRect(p, sizes)
	return rect
}

// SquareCentered makes square bounds of a given size
// Lat and Lng represent center of the square
func SquareCentered(lat float64, lng float64, sizeNM float64) *rtreego.Rect {
	sizes := nmToLatLon(sizeNM/2, sizeNM/2, lat)
	lat = lat - sizes[0]
	lng = lng - sizes[1]
	return Square(lat, lng, sizeNM)
}
