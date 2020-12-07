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
	rects := spatial.GetBoundingBoxes(spatial.MapBounds{-10.0, -10.0, 10.0, 10.0})

	listeners := make([]*spatial.Listener, len(rects))
	for i, rect := range rects {
		listeners[i] = srv.Subscribe(rect)
	}

	for i := range listeners {
		lst := listeners[i]
		wg.Add(1)
		go func() {
			for msg := range lst.Updates() {
				fmt.Println(msg)
			}
			fmt.Println("done reading updates")
			wg.Done()
		}()
	}

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

	for _, lst := range listeners {
		lst.Unsubscribe()
	}

	wg.Wait()
}
