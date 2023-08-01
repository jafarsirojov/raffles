package interfaces

import (
	"context"
	"crm/internal/structs"
)

type TextAdminRepo interface {
	GetTexts(ctx context.Context) (texts []structs.Text, err error)
	UpdateText(ctx context.Context, text structs.Text) error
}

type TextClientRepo interface {
	GetTexts(ctx context.Context) (texts []structs.Text, err error)
}
