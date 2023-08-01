package text

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

func NewRepo(params Params) interfaces.TextAdminRepo {
	return &repo{
		db:     params.DB,
		logger: params.Logger,
	}
}

type repo struct {
	db     db.Querier
	logger *zap.Logger
}

func (r *repo) GetTexts(ctx context.Context) (texts []structs.Text, err error) {
	rows, err := r.db.Query(ctx, `
SELECT 
    key, 
    value
FROM text
WHERE 1=1;`)
	if err != nil {
		r.logger.Error("pkg.repo.admin.text.GetTexts r.db.Query", zap.Error(err))
		return nil, err
	}

	for rows.Next() {
		var t structs.Text
		err = rows.Scan(
			&t.Key,
			&t.Value,
		)
		if err != nil {
			r.logger.Error("pkg.repo.admin.text.GetTexts rows.Scan()", zap.Error(err))
			return nil, err
		}

		texts = append(texts, t)
	}

	if len(texts) == 0 {
		return nil, errors.ErrNotFound
	}

	return texts, nil
}

func (r *repo) UpdateText(ctx context.Context, text structs.Text) error {
	_, err := r.db.Exec(ctx, `UPDATE text SET value = $2 WHERE key = $1;`, text.Key, text.Value)
	if err != nil {
		r.logger.Error("pkg.repo.admin.text.UpdateText r.db.Exec", zap.Error(err))
		return err
	}

	return nil
}
