package main

import (
	"fmt"
	"sync"

	"github.com/viert/spatial"
)

const (
	planeSizeLat = 0.00008 // 30 feet
	planeSizeLng = 0.00010 // 30 feet
)

type plane struct {
	id string
}

func main() {
	var wg sync.WaitGroup

	srv := spatial.New(25, 50, 100)

	listener := srv.Subscribe(spatial.MapBounds{
		SouthWestLng: -10.0,
		SouthWestLat: -10.0,
		NorthEastLng: 10.0,
		NorthEastLat: 10.0,
	})

	wg.Add(1)
	go func() {
		for msg := range listener.Updates() {
			fmt.Println(msg)
		}
		fmt.Println("done reading updates")
		wg.Done()
	}()

	p := &plane{"U-AQ7E"}

	_, err := srv.Add(0, 0, planeSizeLng, planeSizeLat, p.id, p)
	if err != nil {
		panic(err)
	}

	_, err = srv.Add(3, 3, planeSizeLng, planeSizeLat, p.id, p)
	if err != nil {
		panic(err)
	}

	_, err = srv.Add(12, 12, planeSizeLng, planeSizeLat, p.id, p)
	if err != nil {
		panic(err)
	}

	listener.Unsubscribe()

	wg.Wait()
}
