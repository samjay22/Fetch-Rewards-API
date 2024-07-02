package Services

import (
	"Fetch-Rewards-API/Backend/Interfaces"
	"Fetch-Rewards-API/Shared/Structs"
	"context"
	"database/sql"
	"github.com/rs/zerolog"
	_ "modernc.org/sqlite"
)

type databaseService struct {
	db     *sql.DB
	logger *zerolog.Logger
}

type NewDatabaseServiceArgs struct {
	Logger           *zerolog.Logger
	Cfg              *Structs.Config
	Delegate         func(db *sql.DB)
	ConnectionString string
}

// NewDatabaseService creates a new instance of DatabaseService.
func NewDatabaseService(args *NewDatabaseServiceArgs) Interfaces.DatabaseService {
	db, err := sql.Open("sqlite", args.ConnectionString)
	if err != nil {
		args.Logger.Fatal().Err(err).Msg("Failed to connect to database")
		return nil
	}

	args.Delegate(db)

	return &databaseService{
		db:     db,
		logger: args.Logger,
	}
}

func (d *databaseService) GetEntityByFilterRule(ctx context.Context, filterRule func(connection interface{}) (interface{}, error)) (interface{}, error) {
	obj, err := filterRule(d.db)
	if err != nil {
		d.logger.Error().Err(err).Msg("Failed to execute filter rule")
		return nil, err
	}
	return obj, nil
}

func (d *databaseService) UpdateEntityByFilterRule(ctx context.Context, filterRule func(connection interface{}) (bool, error)) (bool, error) {
	result, err := filterRule(d.db)
	if err != nil {
		d.logger.Error().Err(err).Msg("Failed to execute update rule")
		return false, err
	}
	return result, nil
}

func (d *databaseService) DeleteEntityByFilterRule(ctx context.Context, filterRule func(connection interface{}) (bool, error)) (bool, error) {
	result, err := filterRule(d.db)
	if err != nil {
		d.logger.Error().Err(err).Msg("Failed to execute delete rule")
		return false, err
	}
	return result, nil
}

func (d *databaseService) AddEntity(ctx context.Context, delegate func(connection interface{}) error) error {
	if err := delegate(d.db); err != nil {
		d.logger.Error().Err(err).Msg("Failed to add entity")
		return err
	}
	return nil
}

func (d *databaseService) Close() error {
	if err := d.db.Close(); err != nil {
		d.logger.Error().Err(err).Msg("Failed to close database connection")
		return err
	}
	return nil
}
