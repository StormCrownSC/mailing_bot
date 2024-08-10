package repo

import (
	"MailingBot/internal/entity"
	postgresrepo "MailingBot/internal/repo/postgres"
	"MailingBot/pkg/storage/postgres"
	"context"
)

type Admin interface {
	AddAdmin(ctx context.Context, admin entity.Admin) error
	RemoveAdmin(ctx context.Context, admin entity.Admin) error
	RecoverAdmin(ctx context.Context, admin entity.Admin) error
	GetAdmins(ctx context.Context) ([]entity.Admin, error)
}

type User interface {
	AddUser(ctx context.Context, user entity.User) error
	RemoveUser(ctx context.Context, user entity.User) error
	RecoverUser(ctx context.Context, user entity.User) error
	GetUsers(ctx context.Context) ([]entity.User, error)
}

type Text interface {
	AddText(ctx context.Context, text entity.Text) error
	RemoveText(ctx context.Context, text entity.Text) error
	RecoverText(ctx context.Context, text entity.Text) error
	GetTexts(ctx context.Context) ([]entity.Text, error)
	GetText(ctx context.Context, id uint32) (entity.Text, error)
}

type Repository struct {
	Admin Admin
	User  User
	Text  Text
}

func NewRepository(db postgres.PgxPool) *Repository {
	return &Repository{
		Admin: postgresrepo.NewAdminRepo(db),
		User:  postgresrepo.NewUserRepo(db),
		Text:  postgresrepo.NewTextRepo(db),
	}
}
