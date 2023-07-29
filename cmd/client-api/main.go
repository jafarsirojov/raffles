package main

import (
	"crm/cmd/client-api/handlers"
	"crm/cmd/client-api/router"
	"crm/internal"
	"crm/pkg/db"
	"crm/pkg/logger"
	"crm/pkg/repo/client"
	"fmt"
	"go.uber.org/fx"
)

func main() {
	fmt.Println("Initializing client-api modules: ")
	fx.New(
		handlers.Module,
		router.Module,
		internal.Module,
		client.Module,
		db.Module,
		logger.Module,
	).Run()
}
