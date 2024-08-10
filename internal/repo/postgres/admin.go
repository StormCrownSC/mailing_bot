package postgresrepo

import (
	"MailingBot/internal/entity"
	"MailingBot/pkg/storage/postgres"
	"MailingBot/pkg/utils"
	"context"
	"fmt"
)

type AdminRepo struct {
	db postgres.PgxPool
}

func NewAdminRepo(db postgres.PgxPool) *AdminRepo {
	return &AdminRepo{db: db}
}

func (r *AdminRepo) AddAdmin(ctx context.Context, admin entity.Admin) error {
	query := `
		INSERT INTO admin.admins
		(telegram_id, name, is_deleted)
		VALUES ($1, $2, false)
	`

	_, err := r.db.Exec(ctx, query, admin.TelegramID, admin.Name)
	if err != nil {
		return fmt.Errorf("%s; error: %w", utils.CurrentFunction(), err)
	}

	return nil
}

func (r *AdminRepo) RemoveAdmin(ctx context.Context, admin entity.Admin) error {
	query := `
		UPDATE admin.admins
		SET is_deleted = true
		WHERE telegram_id = $1
	`

	_, err := r.db.Exec(ctx, query, admin.TelegramID)
	if err != nil {
		return fmt.Errorf("%s; error: %w", utils.CurrentFunction(), err)
	}

	return nil
}

func (r *AdminRepo) RecoverAdmin(ctx context.Context, admin entity.Admin) error {
	query := `
		UPDATE admin.admins
		SET is_deleted = false
		WHERE telegram_id = $1
	`

	_, err := r.db.Exec(ctx, query, admin.TelegramID)
	if err != nil {
		return fmt.Errorf("%s; error: %w", utils.CurrentFunction(), err)
	}

	return nil
}

func (r *AdminRepo) GetAdmins(ctx context.Context) ([]entity.Admin, error) {
	query := `
		SELECT 
			telegram_id, 
			name,
			is_deleted
		FROM admin.admins
	`
	result := []entity.Admin{}

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return result, fmt.Errorf("%s; error: %w", utils.CurrentFunction(), err)
	}
	defer rows.Close()

	for rows.Next() {
		var admin entity.Admin
		err = rows.Scan(&admin.TelegramID, &admin.Name, &admin.IsDeleted)
		if err != nil {
			return result, fmt.Errorf("%s; error: %w", utils.CurrentFunction(), err)
		}
		result = append(result, admin)
	}

	return result, nil
}
