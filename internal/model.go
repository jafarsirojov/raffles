package internal

import (
	"crm/internal/admin"
	"crm/internal/auth"
	"crm/internal/client"
	"go.uber.org/fx"
)

var Module = fx.Options(
	client.Module,
	admin.Module,
	auth.Module,
)
