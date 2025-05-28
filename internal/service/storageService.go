package service

import (
	"context"
	router "gotus/internal/api"
	"gotus/internal/config"
	"gotus/internal/repository"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func RunService() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	//dataCh := make(chan fmt.Stringer)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(2)

	config := config.LoadConfig("./config/config.yaml")
	bookRepo := repository.NewBookRepository(config.StoragePath)
	bookInstanceRepo := repository.NewBookInstanceRepository(config.StoragePath)
	reservationRepo := repository.NewReservationRepository(config.StoragePath)
	userRepo := repository.NewUserRepository(config.StoragePath)

	storage := repository.NewStorage(bookRepo, bookInstanceRepo, userRepo, reservationRepo)
	storage.LoadAllFromCSV()

	// HTTP-сервер
	srv := &http.Server{
		Addr:    config.HTTPServer.Host + ":" + config.HTTPServer.Port,
		Handler: router.NewRouter(bookRepo, bookInstanceRepo, reservationRepo, userRepo),
	}

	go func() {
		defer wg.Done()
		log.Printf("Запуск HTTP сервера на %s...\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	/*go func() {
		defer wg.Done()
		GenerateAndStoreData(ctx, dataCh)
		close(dataCh)
	}()

	go func() {
		defer wg.Done()
		for item := range dataCh {
			storage.Store(item)
		}
	}()*/

	go func() {
		defer wg.Done()
		LogUpdatesWorker(ctx, storage)
	}()

	// Ожидаем сигнал завершения
	<-sigs
	log.Println("Получен сигнал завершения. Завершаем выполнение...")

	cancel() // Отправка сигнала контексту

	// Завершаем http-сервер
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("Ошибка при остановке сервера: %v", err)
	}

	wg.Wait() // Ожидание завершения всех горутин
	log.Println("Все горутины завершены. Приложение остановлено.")
}
