package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/fohlarbee/orders-api/application"
)

func main() {
	app := application.New(application.LoadConfig());
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	err := app.Start(ctx)

	if err != nil {
		fmt.Println("Error starting application:", err);
	}
}

