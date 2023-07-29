package interfaces

import (
	"context"
	"crm/internal/structs"
)

type LeadAdminRepo interface {
	GetLeadList(ctx context.Context, offset, limit int) (list []structs.Lead, err error)
	GetLeadListByStatus(ctx context.Context, offset, limit int, status string) (list []structs.Lead, err error)
	GetLeadByID(ctx context.Context, id int) (list structs.Lead, err error)
	UpdateLeadStatus(ctx context.Context, id int, status string) error
}

type LeadClientRepo interface {
	InsertLead(ctx context.Context, request structs.Lead) error
}
