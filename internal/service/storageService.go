package service

import (
	"context"
	"fmt"
	"gotus/internal/repository"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func RunService() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	dataCh := make(chan fmt.Stringer)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(3)

	repository.LoadAllFromCSV()

	go func() {
		defer wg.Done()
		GenerateAndStoreData(ctx, dataCh)
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
		LogUpdatesWorker(ctx)
	}()

	// Ожидаем сигнал завершения
	<-sigs
	log.Println("Получен сигнал завершения. Завершаем выполнение...")

	cancel() // Отправка сигнала контексту

	wg.Wait() // Ожидание завершения всех горутин
	log.Println("Все горутины завершены. Приложение остановлено.")
}
