package main

import (
	"fmt"
	"context"
    "github.com/fohlarbee/orders-api/application"
)

func main() {
	app := application.New();

	err := app.Start(context.TODO())

	if err != nil {
		fmt.Println("Error starting application:", err);
	}
}

