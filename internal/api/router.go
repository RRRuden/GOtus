package router

import (
	_ "gotus/cmd/main/docs"
	b_handler "gotus/internal/api/book"
	booking_handler "gotus/internal/api/booking"
	r_handler "gotus/internal/api/reservation"
	u_handler "gotus/internal/api/user"
	"gotus/internal/repository"
	"gotus/internal/service"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(bookRepo *repository.BookRepository, bookInstanceRepo *repository.BookInstanceRepository, resRepo *repository.ReservationRepository, userRepo *repository.UserRepository) http.Handler {
	mux := http.NewServeMux()

	bookingService := service.NewBookingService(userRepo, bookRepo, bookInstanceRepo, resRepo)

	bookHandler := &b_handler.BookHandler{Repo: bookRepo}
	bookInstanceHandler := &b_handler.BookInstanceHandler{Repo: bookInstanceRepo}
	resHandler := &r_handler.ReservationHandler{Repo: resRepo}
	userHandler := &u_handler.UserHandler{Repo: userRepo}
	bookingHandler := &booking_handler.BookingHandler{Service: bookingService}

	// Book Routes
	mux.HandleFunc("/api/book", method(bookHandler.CreateBook, "POST"))
	mux.HandleFunc("/api/books", method(bookHandler.GetAllBooks, "GET"))
	mux.HandleFunc("/api/book/", bookHandler.BookByISBNHandler)

	// Book Instance Routes
	mux.HandleFunc("/api/bookinstance", method(bookInstanceHandler.CreateBookInstance, "POST"))
	mux.HandleFunc("/api/bookinstances", method(bookInstanceHandler.GetAllBookInstances, "GET"))
	mux.HandleFunc("/api/bookinstance/", bookInstanceHandler.BookInstanceByIDHandler)

	// Reservation Routes
	mux.HandleFunc("/api/reservation", method(resHandler.CreateReservation, "POST"))
	mux.HandleFunc("/api/reservations", method(resHandler.GetAllReservations, "GET"))
	mux.HandleFunc("/api/reservation/", resHandler.ReservationByIDHandler)

	// User Routes
	mux.HandleFunc("/api/user", method(userHandler.CreateUser, "POST"))
	mux.HandleFunc("/api/users", method(userHandler.GetAllUsers, "GET"))
	mux.HandleFunc("/api/user/", userHandler.UserByIDHandler)

	// User Routes
	mux.HandleFunc("/api/booking/create", method(bookingHandler.CreateBooking, "POST"))
	mux.HandleFunc("/api/booking/extend/", method(bookingHandler.ExtendBooking, "PUT"))
	mux.HandleFunc("/api/booking/cancel/", method(bookingHandler.CancelBooking, "POST"))
	mux.HandleFunc("/api/booking/end/", method(bookingHandler.EndBooking, "POST"))

	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	return mux
}

func method(h http.HandlerFunc, allowed string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowed {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h(w, r)
	}
}
