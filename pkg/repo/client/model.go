package client

import (
	"crm/pkg/repo/client/estate"
	"crm/pkg/repo/client/lead"
	"go.uber.org/fx"
)

var Module = fx.Options(
	lead.Module,
	estate.Module,
)
