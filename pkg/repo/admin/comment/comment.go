package comment

import (
	"context"
	interfaces "crm/internal/interface"
	"crm/internal/structs"
	"crm/pkg/db"
	"crm/pkg/errors"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Provide(NewRepo)

type Params struct {
	fx.In
	DB     db.Querier
	Logger *zap.Logger
}

func NewRepo(params Params) interfaces.CommentRepo {
	return &repo{
		db:     params.DB,
		logger: params.Logger,
	}
}

type repo struct {
	db     db.Querier
	logger *zap.Logger
}

func (r *repo) GetCommentsByLeadID(ctx context.Context, leadID int) (comments []structs.Comment, err error) {
	rows, err := r.db.Query(ctx, `
SELECT 
    id,
    lead_id,
    admin_id, 
    admin_login, 
    text, 
    to_char(created_at AT TIME ZONE 'Asia/Dubai', 'DD-MM HH24:MI')
FROM comment
WHERE lead_id = $1
ORDER BY id DESC;`, leadID)
	if err != nil {
		r.logger.Error("pkg.repo.comment.GetCommentsByLeadID r.db.Query", zap.Error(err))
		return nil, err
	}

	for rows.Next() {
		var c structs.Comment
		err = rows.Scan(
			&c.ID,
			&c.LeadID,
			&c.AdminID,
			&c.AdminLogin,
			&c.Text,
			&c.CreatedAt,
		)
		if err != nil {
			r.logger.Error("pkg.repo.comment.GetCommentsByLeadID rows.Scan()", zap.Error(err))
			return nil, err
		}

		comments = append(comments, c)
	}

	if len(comments) == 0 {
		return nil, errors.ErrNotFound
	}

	r.logger.Info("pkg.repo.comment.GetCommentsByLeadID ", zap.Any("comments", comments))

	return comments, nil
}

func (r *repo) AddComment(ctx context.Context, comment structs.Comment) error {
	_, err := r.db.Exec(ctx, `insert into comment (lead_id, admin_id, admin_login, text)values($1, $2, $3, $4);`,
		comment.LeadID,
		comment.AdminID,
		comment.AdminLogin,
		comment.Text,
	)
	if err != nil {
		r.logger.Error("pkg.repo.comment.AddComment r.db.Exec", zap.Error(err))
		return err
	}

	return nil
}
