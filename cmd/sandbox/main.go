package main

import (
	"fmt"
	"sync"
	"time"

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

	srv := spatial.New(25, 50, 100, time.Second)

	listener := srv.Subscribe(spatial.MapBounds{
		SouthWestLng: -10.0,
		SouthWestLat: -10.0,
		NorthEastLng: 10.0,
		NorthEastLat: 10.0,
	})

	wg.Add(1)
	go func() {
		for msg := range listener.Updates() {
			fmt.Printf("update count: %d\n", len(msg))
			for _, idxbl := range msg {
				ref := idxbl.Ref()
				p, ok := ref.(*plane)
				if ok {
					fmt.Println(p)
				}
			}
		}
		fmt.Println("done reading updates")
		wg.Done()
	}()

	p := &plane{"U-AQ7E"}

	fmt.Println("added")
	_, err := srv.Add(0, 0, planeSizeLng, planeSizeLat, p.id, p)
	if err != nil {
		panic(err)
	}
	time.Sleep(2 * time.Second)

	fmt.Println("moved")
	_, err = srv.Add(3, 3, planeSizeLng, planeSizeLat, p.id, p)
	if err != nil {
		panic(err)
	}
	time.Sleep(2 * time.Second)

	fmt.Println("moved outside")
	_, err = srv.Add(12, 12, planeSizeLng, planeSizeLat, p.id, p)
	if err != nil {
		panic(err)
	}
	time.Sleep(2 * time.Second)

	listener.Unsubscribe()

	wg.Wait()
}
