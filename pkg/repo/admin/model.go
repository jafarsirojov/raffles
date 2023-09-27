package admin

import (
	"crm/pkg/repo/admin/auth"
	"crm/pkg/repo/admin/comment"
	"crm/pkg/repo/admin/estate"
	"crm/pkg/repo/admin/landing"
	"crm/pkg/repo/admin/lead"
	"crm/pkg/repo/admin/text"
	"go.uber.org/fx"
)

var Module = fx.Options(
	lead.Module,
	auth.Module,
	comment.Module,
	estate.Module,
	text.Module,
	landing.Module,
)
