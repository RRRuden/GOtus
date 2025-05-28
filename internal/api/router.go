package router

import (
	b_handler "gotus/internal/api/book"
	r_handler "gotus/internal/api/reservation"
	u_handler "gotus/internal/api/user"
	"gotus/internal/repository"
	"net/http"
)

func NewRouter(bookRepo *repository.BookRepository, bookInstanceRepo *repository.BookInstanceRepository, resRepo *repository.ReservationRepository, userRepo *repository.UserRepository) http.Handler {
	mux := http.NewServeMux()

	bookHandler := &b_handler.BookHandler{Repo: bookRepo}
	bookInstanceHandler := &b_handler.BookInstanceHandler{Repo: bookInstanceRepo}
	resHandler := &r_handler.ReservationHandler{Repo: resRepo}
	userHandler := &u_handler.UserHandler{Repo: userRepo}

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
