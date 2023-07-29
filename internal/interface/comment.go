package interfaces

import (
	"context"
	"crm/internal/structs"
)

type CommentRepo interface {
	GetCommentsByLeadID(ctx context.Context, leadID int) (comments []structs.Comment, err error)
	AddComment(ctx context.Context, comment structs.Comment) error
}
