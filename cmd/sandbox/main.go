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
	srv := spatial.New(25, 50, 100)
	bounds, err := spatial.MakeBounds(-10.0, -10.0, 10.0, 10.0)
	if err != nil {
		panic(err)
	}

	updates := srv.Subscribe(bounds)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for msg := range updates.Updates() {
			fmt.Println(msg)
		}
		wg.Done()
		fmt.Println("done reading updates")
	}()

	p := &plane{"U-AQ7E"}

	_, err = srv.Add(0, 0, planeSizeLng, planeSizeLat, p.id, p)
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

	updates.Unsubscribe()
	wg.Wait()

}
