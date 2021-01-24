package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/dhconnelly/rtreego"
	"github.com/viert/spatial"
)

const (
	planeSizeLat = 0.00008 // 30 feet
	planeSizeLng = 0.00010 // 30 feet
)

type plane struct {
	id   string
	rect *rtreego.Rect
}

func (p *plane) String() string {
	return fmt.Sprintf("Plane \"%s\" %s", p.id, p.rect)
}

func (p *plane) ID() string {
	return p.id
}

func (p *plane) Bounds() *rtreego.Rect {
	return p.rect
}

func (p *plane) Ref() interface{} {
	return nil
}

func (p *plane) Type() spatial.IndexableType {
	return 1
}

func makePlane(id string, lat float64, lng float64) *plane {
	mb := spatial.MapBounds{
		SouthWestLat: lat,
		SouthWestLng: lng,
		NorthEastLat: lat + planeSizeLat,
		NorthEastLng: lng + planeSizeLng,
	}
	rect := mb.Rects()[0]
	p := &plane{
		id,
		rect,
	}
	return p
}

func main() {
	var wg sync.WaitGroup
	planeID := "RF-350"

	srv := spatial.New(25, 50)
	lst := srv.NewListener(100, 300*time.Millisecond)
	lst.SetBounds(spatial.MapBounds{
		SouthWestLng: -10,
		SouthWestLat: -10,
		NorthEastLng: 10,
		NorthEastLat: 10,
	})

	wg.Add(1)
	go func() {
		for update := range lst.Updates() {
			fmt.Println("got update")
			for _, idxbl := range update {
				fmt.Println("from update", idxbl)
			}
		}
		wg.Done()
	}()

	p := makePlane(planeID, 0, 0)
	srv.Add(p)
	fmt.Println("added")
	time.Sleep(time.Second)

	p = makePlane(planeID, 2, 2)
	fmt.Println("moved")
	srv.Add(p)

	time.Sleep(time.Second)

	p = makePlane(planeID, 12, 12)
	fmt.Println("moved outside")
	srv.Add(p)

	time.Sleep(time.Second)

	lst.Stop()
	wg.Wait()
}
