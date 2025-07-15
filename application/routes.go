package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/fohlarbee/orders-api/application/handler"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to Orders API"))
	});

	router.Route("/orders", loadOrderRoutes)

	// router.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte("Hello, World!"))
	// })

	// router.Post("/hello", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte("Hello from POST!"))
	// })

	return router
}

func loadOrderRoutes(router chi.Router) {
	orderHandler := &handler.Order{}

	router.Post("/", orderHandler.CreateOrder);
	router.Get("/", orderHandler.GetOrders);
	router.Get("/{id}", orderHandler.GetById);
	router.Put("/{id}", orderHandler.UpdateByOrderId);
	router.Delete("/{id}", orderHandler.DeleteByOrderId);
	  
}
