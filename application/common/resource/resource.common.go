package resource

import (
	"antrein/bc-dashboard/model/config"
	"context"
	_ "database/sql"
	"io"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

type CommonResource struct {
	Db   *sqlx.DB
	Vld  *validator.Validate
	GRPC *grpc.ClientConn
}

func NewCommonResource(cfg *config.Config, ctx context.Context) (*CommonResource, error) {
	db, err := sqlx.Open("postgres", cfg.Database.PostgreDB.Host)
	if err != nil {
		return nil, err
	}

	err = migrateDb(ctx, db)
	if err != nil {
		return nil, err
	}

	vld := validator.New()

	grpcClient, err := grpc.Dial(cfg.GRPCConfig.DashboardQueue, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	rsc := CommonResource{
		Db:   db,
		Vld:  vld,
		GRPC: grpcClient,
	}
	return &rsc, nil
}

func migrateDb(ctx context.Context, db *sqlx.DB) error {
	filePath := "./files/migrations/migrate.sql"
	migrationFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer migrationFile.Close()

	migration, err := io.ReadAll(migrationFile)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, string(migration))
	if err != nil {
		return err
	}

	return nil
}
