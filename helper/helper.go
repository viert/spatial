package helper

import (
	"math"

	"github.com/dhconnelly/rtreego"
)

func nmToLatLon(x float64, y float64, atLatitude float64) []float64 {
	// for latitude 60nm is 1˚
	latSize := (1 / 60) * y

	// for longitude it depends on current latitude

	// first let's convert latitude degree to radians
	latitudeRad := (atLatitude * 2 * math.Pi) / 360

	// calculate size
	lngSize := (1 / 60) * x
	// and make a latitude correction
	lngSize = lngSize / math.Cos(latitudeRad)
	return []float64{latSize, lngSize}
}

// Square makes square bounds of a given size
// Lat and Lng represent top left angle of the square
func Square(lat float64, lng float64, sizeNM float64) *rtreego.Rect {
	p := rtreego.Point{lat, lng}
	rect, _ := rtreego.NewRect(p, nmToLatLon(sizeNM, sizeNM, lat))
	return rect
}

// SquareCentered makes square bounds of a given size
// Lat and Lng represent center of the square
func SquareCentered(lat float64, lng float64, sizeNM float64) *rtreego.Rect {
	lat = lat - sizeNM/2
	lng = lng - sizeNM/2
	return Square(lat, lng, sizeNM)
}
