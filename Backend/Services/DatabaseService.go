package Services

import (
	"Fetch-Rewards-API/Backend/Interfaces"
	"Fetch-Rewards-API/Shared/Structs"
	"context"
	"database/sql"
	"github.com/rs/zerolog"
	_ "modernc.org/sqlite"
)

// databaseService implements the Interfaces.DatabaseService interface
type databaseService struct {
	db     *sql.DB
	logger *zerolog.Logger
	args   *ScopedConstructorArgs
}

// NewDatabaseServiceArgs contains the arguments required to create a new DatabaseService
type NewDatabaseServiceArgs struct {
	Logger           *zerolog.Logger
	Cfg              *Structs.Config
	Delegate         func(db *sql.DB)
	ConnectionString string
}

// ScopedConstructorArgs contains the arguments required to create a new scoped DatabaseService
type ScopedConstructorArgs struct {
	Args *NewDatabaseServiceArgs
}

// NewDatabaseService creates a new instance of DatabaseService
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
		args:   &ScopedConstructorArgs{Args: args},
	}
}

func (d *databaseService) withNewConnection(action func(db *sql.DB) (interface{}, error)) (interface{}, error) {
	db, err := sql.Open("sqlite", d.args.Args.ConnectionString)
	if err != nil {
		d.logger.Fatal().Err(err).Msg("Failed to connect to database")
		return nil, err
	}
	defer func() {
		if err := db.Close(); err != nil {
			d.logger.Error().Err(err).Msg("Failed to close database connection")
		}
	}()
	return action(db)
}

// GetEntityByFilterRule retrieves an entity based on the provided filter rule
func (d *databaseService) GetEntityByFilterRule(ctx context.Context, filterRule func(connection interface{}) (interface{}, error)) (interface{}, error) {
	return d.withNewConnection(func(db *sql.DB) (interface{}, error) {
		obj, err := filterRule(db)
		if err != nil {
			d.logger.Error().Err(err).Msg("Failed to execute filter rule")
			return nil, err
		}
		return obj, nil
	})
}

// UpdateEntityByFilterRule updates an entity based on the provided filter rule
func (d *databaseService) UpdateEntityByFilterRule(ctx context.Context, filterRule func(connection interface{}) (bool, error)) (bool, error) {
	result, err := d.withNewConnection(func(db *sql.DB) (interface{}, error) {
		return filterRule(db)
	})
	if err != nil {
		return false, err
	}
	return result.(bool), nil
}

// DeleteEntityByFilterRule deletes an entity based on the provided filter rule
func (d *databaseService) DeleteEntityByFilterRule(ctx context.Context, filterRule func(connection interface{}) (bool, error)) (bool, error) {
	result, err := d.withNewConnection(func(db *sql.DB) (interface{}, error) {
		return filterRule(db)
	})
	if err != nil {
		return false, err
	}
	return result.(bool), nil
}

// AddEntity adds a new entity using the provided delegate function
func (d *databaseService) AddEntity(ctx context.Context, delegate func(connection interface{}) error) error {
	_, err := d.withNewConnection(func(db *sql.DB) (interface{}, error) {
		if err := delegate(db); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}
