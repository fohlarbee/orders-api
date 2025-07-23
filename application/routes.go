package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/fohlarbee/orders-api/application/handler"
	"github.com/fohlarbee/orders-api/repository/order"
)

func (a *App) loadRoutes()  {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to Orders API"))
	});

	router.Route("/orders", a.loadOrderRoutes)
	a.router = router
}

func (a *App) loadOrderRoutes(router chi.Router) {
	orderHandler := &handler.Order{
		Repo: &order.RedisRepo{
			Client: a.rdb, 
		},
	}

	router.Post("/", orderHandler.CreateOrder);
	router.Get("/", orderHandler.GetOrders);
	router.Get("/{id}", orderHandler.GetById);
	router.Put("/{id}", orderHandler.UpdateByOrderId);
	router.Delete("/{id}", orderHandler.DeleteByOrderId);
	  
}
