package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "gotus/cmd/main/docs" // Импорт с побочным эффектом, чтобы инициализировать Swagger
	"gotus/internal/config"
	"gotus/internal/grpc/api/booking_api"
	grpcBooking "gotus/internal/grpc/server"
	"gotus/internal/logger"
	"gotus/internal/repository"
	"gotus/internal/service/booking"
	emailsender "gotus/internal/service/email_sender"
	emailtemplater "gotus/internal/service/email_templater"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {

	var wg sync.WaitGroup

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)

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
		if err := grpcBooking.RunGRPCServer(bookingSerivce, ":"+config.BookingServer.Port); err != nil {
			log.Fatalf("gRPC сервер ошибка: %v", err)
		}
	}()

	conn, err := grpc.NewClient(bookingServerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := booking_api.NewBookingServiceClient(conn)

	res, error := client.CreateBooking(context.Background(), &booking_api.CreateBookingRequest{
		UserId: 6,
		Isbn:   "978-3-16-148410-5",
	})

	if error != nil {
		log.Fatalf("error get request: %s", err.Error())
	}

	log.Println("CreateBooking resul", res)

	// Ожидаем сигнал завершения
	<-sigs
	log.Println("Получен сигнал завершения. Завершаем выполнение...")

	// Завершаем gRPC-сервер
	grpcBooking.StopGRPCServer()

	wg.Wait() // Ожидание завершения всех горутин
	log.Println("Все горутины завершены. Приложение остановлено.")
}
