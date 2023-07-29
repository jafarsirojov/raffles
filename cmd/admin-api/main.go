package main

import (
	"crm/cmd/admin-api/handlers"
	"crm/cmd/admin-api/router"
	"crm/internal"
	"crm/pkg/db"
	"crm/pkg/logger"
	"crm/pkg/repo/admin"
	"fmt"
	"go.uber.org/fx"
)

func main() {

	fmt.Println("Initializing admin-api modules: ")
	fx.New(
		handlers.Module,
		router.Module,
		internal.Module,
		admin.Module,
		db.Module,
		logger.Module,
	).Run()
}
