package service

import (
	"fmt"
	"gotus/internal/repository"
	"sync"
)

func RunService() {
	var wg sync.WaitGroup

	dataCh := make(chan fmt.Stringer)

	wg.Add(3)

	go func() {
		defer wg.Done()
		GenerateAndStoreData(dataCh)
		close(dataCh)
	}()

	go func() {
		defer wg.Done()
		for item := range dataCh {
			repository.Store(item)
		}
	}()

	go func() {
		defer wg.Done()
		LogUpdatesWorker()
	}()

	wg.Wait()
}
