// @title GOtus API
// @version 1.0
// @description REST API для библиотеки GOtus.

package main

import (
	"context"
	_ "gotus/cmd/main/docs" // Импорт с побочным эффектом, чтобы инициализировать Swagger
	router "gotus/internal/api"
	"gotus/internal/config"
	grpc "gotus/internal/grpc/server"
	"gotus/internal/repository"
	"gotus/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	RunService()
}

func RunService() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(3)

	config := config.LoadConfig("././config/config.yaml")
	bookRepo := repository.NewBookRepository(config.StoragePath)
	bookInstanceRepo := repository.NewBookInstanceRepository(config.StoragePath)
	reservationRepo := repository.NewReservationRepository(config.StoragePath)
	userRepo := repository.NewUserRepository(config.StoragePath)
	bookingSerivce := service.NewBookingService(userRepo, bookRepo, bookInstanceRepo, reservationRepo)

	storage := repository.NewStorage(bookRepo, bookInstanceRepo, userRepo, reservationRepo)
	storage.LoadAllFromCSV()

	grpcAddr := config.BookingServer.Host + ":" + config.BookingServer.Port

	go func() {
		defer wg.Done()
		if err := grpc.RunGRPCServer(bookingSerivce, ":"+config.BookingServer.Port); err != nil {
			log.Fatalf("gRPC сервер ошибка: %v", err)
		}
	}()

	// HTTP-сервер
	srv := &http.Server{
		Addr:    config.HTTPServer.Host + ":" + config.HTTPServer.Port,
		Handler: router.NewRouter(bookRepo, bookInstanceRepo, reservationRepo, userRepo, grpcAddr),
	}

	go func() {
		defer wg.Done()
		log.Printf("Запуск HTTP сервера на %s...\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		service.LogUpdatesWorker(ctx, storage)
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
