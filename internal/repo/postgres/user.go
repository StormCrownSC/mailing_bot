package postgresrepo

import (
	"MailingBot/internal/entity"
	"MailingBot/pkg/storage/postgres"
	"MailingBot/pkg/utils"
	"context"
	"fmt"
)

type UserRepo struct {
	db postgres.PgxPool
}

func NewUserRepo(db postgres.PgxPool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) AddUser(ctx context.Context, user entity.User) error {
	query := `
		INSERT INTO mailing.users
		(telegram_id, name, is_deleted)
		VALUES ($1, $2, false)
	`

	_, err := r.db.Exec(ctx, query, user.TelegramID, user.Name)
	if err != nil {
		return fmt.Errorf("%s; error: %w", utils.CurrentFunction(), err)
	}

	return nil
}

func (r *UserRepo) RemoveUser(ctx context.Context, user entity.User) error {
	query := `
		UPDATE mailing.users
		SET is_deleted = true
		WHERE telegram_id = $1
	`

	_, err := r.db.Exec(ctx, query, user.TelegramID)
	if err != nil {
		return fmt.Errorf("%s; error: %w", utils.CurrentFunction(), err)
	}

	return nil
}

func (r *UserRepo) RecoverUser(ctx context.Context, user entity.User) error {
	query := `
		UPDATE mailing.users
		SET is_deleted = false
		WHERE telegram_id = $1
	`

	_, err := r.db.Exec(ctx, query, user.TelegramID)
	if err != nil {
		return fmt.Errorf("%s; error: %w", utils.CurrentFunction(), err)
	}

	return nil
}

func (r *UserRepo) GetUsers(ctx context.Context) ([]entity.User, error) {
	query := `
		SELECT 
			telegram_id, 
			name,
			is_deleted
		FROM mailing.users
	`
	result := []entity.User{}

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return result, fmt.Errorf("%s; error: %w", utils.CurrentFunction(), err)
	}
	defer rows.Close()

	for rows.Next() {
		var user entity.User
		err = rows.Scan(&user.TelegramID, &user.Name, &user.IsDeleted)
		if err != nil {
			return result, fmt.Errorf("%s; error: %w", utils.CurrentFunction(), err)
		}
		result = append(result, user)
	}

	return result, nil
}
