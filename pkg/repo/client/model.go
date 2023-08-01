package client

import (
	"crm/pkg/repo/client/estate"
	"crm/pkg/repo/client/lead"
	"crm/pkg/repo/client/text"
	"go.uber.org/fx"
)

var Module = fx.Options(
	lead.Module,
	estate.Module,
	text.Module,
)
