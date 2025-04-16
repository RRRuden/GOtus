package main

import (
	"fmt"
	"gotus/internal/repository"
	"gotus/internal/service"
	"sync"
)

func main() {
	dataCh := make(chan fmt.Stringer)

	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		repository.Store(dataCh)
	}()

	go func() {
		defer wg.Done()
		repository.LogUpdatesWorker()
	}()

	go func() {
		defer wg.Done()
		service.GenerateAndStoreData(dataCh)
		close(dataCh)
	}()

	wg.Wait()
}
