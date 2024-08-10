package postgresrepo

import (
	"MailingBot/internal/entity"
	"MailingBot/pkg/storage/postgres"
	"MailingBot/pkg/utils"
	"context"
	"fmt"
)

type TextRepo struct {
	db postgres.PgxPool
}

func NewTextRepo(db postgres.PgxPool) *TextRepo {
	return &TextRepo{db: db}
}

func (r *TextRepo) AddText(ctx context.Context, text entity.Text) error {
	query := `
		INSERT INTO mailing.texts
		(text)
		VALUES ($1)
	`

	_, err := r.db.Exec(ctx, query, text.Text)
	if err != nil {
		return fmt.Errorf("%s; error: %w", utils.CurrentFunction(), err)
	}

	return nil
}

func (r *TextRepo) RemoveText(ctx context.Context, text entity.Text) error {
	query := `
		UPDATE mailing.texts
		SET is_deleted = true
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, text.ID)
	if err != nil {
		return fmt.Errorf("%s; error: %w", utils.CurrentFunction(), err)
	}

	return nil
}

func (r *TextRepo) RecoverText(ctx context.Context, text entity.Text) error {
	query := `
		UPDATE mailing.texts
		SET is_deleted = false
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, text.ID)
	if err != nil {
		return fmt.Errorf("%s; error: %w", utils.CurrentFunction(), err)
	}

	return nil
}

func (r *TextRepo) GetTexts(ctx context.Context) ([]entity.Text, error) {
	query := `
		SELECT 
			id, 
			text,
			is_deleted
		FROM mailing.texts
		ORDER BY is_deleted
	`
	result := []entity.Text{}

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return result, fmt.Errorf("%s; error: %w", utils.CurrentFunction(), err)
	}
	defer rows.Close()

	for rows.Next() {
		var text entity.Text
		err = rows.Scan(&text.ID, &text.Text, &text.IsDeleted)
		if err != nil {
			return result, fmt.Errorf("%s; error: %w", utils.CurrentFunction(), err)
		}
		result = append(result, text)
	}

	return result, nil
}

func (r *TextRepo) GetText(ctx context.Context, id uint32) (entity.Text, error) {
	query := `
        SELECT 
            id, 
            text,
            is_deleted
        FROM mailing.texts
        WHERE id = $1
    `
	result := entity.Text{}

	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(&result.ID, &result.Text, &result.IsDeleted)
	if err != nil {
		return result, fmt.Errorf("%s; error: %w", utils.CurrentFunction(), err)
	}

	return result, nil
}
