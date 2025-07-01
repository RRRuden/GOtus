// @title GOtus API
// @version 1.0
// @description REST API для библиотеки GOtus.

package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "gotus/cmd/main/docs" // Импорт с побочным эффектом, чтобы инициализировать Swagger
	router "gotus/internal/api"
	"gotus/internal/config"
	grpc "gotus/internal/grpc/server"
	"gotus/internal/logger"
	"gotus/internal/repository"
	"gotus/internal/service/booking"
	emailsender "gotus/internal/service/email_sender"
	emailtemplater "gotus/internal/service/email_templater"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/lib/pq"
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

	wg.Add(3)

	config := config.LoadConfig("././config/config.yaml")

	pgCfg := config.PostgreSQL
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		pgCfg.Host, pgCfg.Port, pgCfg.User, pgCfg.Password, pgCfg.DBName, pgCfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("PostgreSQL не отвечает: %v", err)
	}
	defer db.Close()

	bookRepo := repository.NewBookPostgresRepo(db)
	bookInstanceRepo := repository.NewBookInstancePostgresRepo(db)
	reservationRepo := repository.NewReservationPostgresRepo(db)
	userRepo := repository.NewUserPostgresRepo(db)
	messageRepo := repository.NewEmailMessagePostgresRepo(db)
	logger := logger.NewRedisLogger(config.Redis.Addr, config.Redis.Password, config.Redis.DB)
	emailSender := emailsender.NewSMTPEmailSender(config.SMTP, messageRepo)
	templater := emailtemplater.NewEmailTemplater(reservationRepo, bookInstanceRepo, bookRepo, config.TemplateDir)
	bookingSerivce := booking.NewBookingService(userRepo, bookRepo, bookInstanceRepo, reservationRepo, logger, templater, emailSender)

	bookingServerAddr := config.BookingServer.Host + ":" + config.BookingServer.Port

	go func() {
		defer wg.Done()
		if err := grpc.RunGRPCServer(bookingSerivce, ":"+config.BookingServer.Port); err != nil {
			log.Fatalf("gRPC сервер ошибка: %v", err)
		}
	}()

	// HTTP-сервер
	srv := &http.Server{
		Addr:    config.HTTPServer.Host + ":" + config.HTTPServer.Port,
		Handler: router.NewRouter(bookRepo, bookInstanceRepo, reservationRepo, userRepo, bookingServerAddr),
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

	wg.Wait() // Ожидание завершения всех горутин
	log.Println("Все горутины завершены. Приложение остановлено.")
}

func GetMongoDb(config config.MongoDBConfig, ctx context.Context) *mongo.Database {
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(config.URI))
	if err != nil {
		log.Fatalf("Ошибка создания Mongo клиента: %v", err)
	}
	if err := mongoClient.Connect(ctx); err != nil {
		log.Fatalf("Ошибка подключения к MongoDB: %v", err)
	}
	mongoDB := mongoClient.Database(config.Database)
	return mongoDB
}
