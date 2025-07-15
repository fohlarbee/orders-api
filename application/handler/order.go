package handler

import (
	"fmt"
	"net/http"

	
)

type Order struct {
		
}

func (o *Order) CreateOrder(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Order created successfully")
}

func (o *Order) GetOrders(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Orders retrieved successfully")
}
func (o *Order) GetById(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Order retrieved by ID successfully")
}
func (o *Order) UpdateByOrderId(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Order updated successfully")
}
func (o *Order) DeleteByOrderId(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Order deleted successfully")
}	
