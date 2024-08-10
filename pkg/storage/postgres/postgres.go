package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Host     string `env:"POSTGRES_HOST" env-required:"true"`
	Port     string `env:"POSTGRES_PORT" env-required:"true"`
	User     string `env:"POSTGRES_USER" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	DBName   string `env:"POSTGRES_DB" env-required:"true"`
	SSLMode  string `yaml:"ssl_mode"`
	MaxConns int32  `yaml:"max_conns"`
	MinConns int32  `yaml:"min_conns"`
}

type PgxPool interface {
	Close()
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Ping(ctx context.Context) error
}

type Storage struct {
	DB  PgxPool
	cfg *Postgres
}

func NewStorage(cfg *Postgres) *Storage {
	return &Storage{cfg: cfg}
}

func (s *Storage) Connect(ctx context.Context) error {
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		s.cfg.Host,
		s.cfg.Port,
		s.cfg.User,
		s.cfg.Password,
		s.cfg.DBName,
		s.cfg.SSLMode,
	)

	pgxConf, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return err
	}

	pgxConf.MaxConns = s.cfg.MaxConns
	pgxConf.MinConns = s.cfg.MinConns

	db, err := pgxpool.NewWithConfig(ctx, pgxConf)
	if err != nil {
		return fmt.Errorf("error creating new pgx pool: %w", err)
	}

	err = db.Ping(ctx)
	if err != nil {
		return fmt.Errorf("error connecting pgx pool: %w", err)
	}

	s.DB = db

	return nil
}
