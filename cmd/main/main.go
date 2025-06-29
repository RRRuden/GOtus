// @title GOtus API
// @version 1.0
// @description REST API для библиотеки GOtus.

package main

import (
	"context"
	_ "gotus/cmd/main/docs" // Импорт с побочным эффектом, чтобы инициализировать Swagger
	router "gotus/internal/api"
	"gotus/internal/config"
	"gotus/internal/generator"
	grpc "gotus/internal/grpc/server"
	"gotus/internal/logger"
	"gotus/internal/repository"
	"gotus/internal/service/booking"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	RunService()
}

func RunService() {
	var wg sync.WaitGroup

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(2)

	config := config.LoadConfig("././config/config.yaml")
	ctx := context.Background()
	// MongoDB
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(config.MongoDB.URI))
	if err != nil {
		log.Fatalf("Ошибка создания Mongo клиента: %v", err)
	}
	if err := mongoClient.Connect(ctx); err != nil {
		log.Fatalf("Ошибка подключения к MongoDB: %v", err)
	}
	mongoDB := mongoClient.Database(config.MongoDB.Database)

	bookRepo := repository.NewMongoBookRepository(mongoDB)
	bookInstanceRepo := repository.NewMongoBookInstanceRepository(mongoDB)
	reservationRepo := repository.NewMongoReservationRepository(mongoDB)
	userRepo := repository.NewMongoUserRepository(mongoDB)
	logger := logger.NewRedisLogger(config.Redis.Addr, config.Redis.Password, config.Redis.DB)
	bookingSerivce := booking.NewBookingService(userRepo, bookRepo, bookInstanceRepo, reservationRepo, logger)

	generator.GenerateInitialDataForMongoDb(userRepo, bookRepo, bookInstanceRepo, reservationRepo)

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

	// Ожидаем сигнал завершения
	<-sigs
	log.Println("Получен сигнал завершения. Завершаем выполнение...")

	// Завершаем http-сервер
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("Ошибка при остановке сервера: %v", err)
	}

	// Завершаем gRPC-сервер
	grpc.StopGRPCServer()

	mongoClient.Disconnect(ctx)

	wg.Wait() // Ожидание завершения всех горутин
	log.Println("Все горутины завершены. Приложение остановлено.")
}
