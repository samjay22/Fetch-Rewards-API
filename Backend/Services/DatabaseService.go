package Services

import (
	"Fetch-Rewards-API/Backend/Interfaces"
	"Fetch-Rewards-API/Shared/Structs"
	"context"
	"database/sql"
	"fmt"
	"github.com/rs/zerolog"
	_ "modernc.org/sqlite"
)

type DatabaseService struct {
	db     *sql.DB
	logger *zerolog.Logger
}

type NewDatabaseServiceArgs struct {
	Logger   *zerolog.Logger
	Cfg      *Structs.Config
	Delegate func(db *sql.DB)
}

// NewDatabaseService creates a new instance of DatabaseService.
func NewDatabaseService(args *NewDatabaseServiceArgs) Interfaces.DatabaseService {
	dbConfig := args.Cfg.Database
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s/%s.db?cache=shared", dbConfig.HomeDir, dbConfig.FileName))
	if err != nil {
		args.Logger.Fatal().Err(err).Msg("Failed to connect to database")
		return nil
	}

	for _, tableDef := range dbConfig.TableDef {
		query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", tableDef.TableName)

		// Construct column definitions
		for i, rowData := range tableDef.TableRows {
			keyType := ""
			if rowData.IsPrimaryKey {
				keyType = "PRIMARY KEY"
			}

			nullState := "NOT NULL"
			if rowData.IsNull {
				nullState = ""
			}

			// Add column definition to query
			query += fmt.Sprintf("%s %s %s %s", rowData.RowId, rowData.DataType, keyType, nullState)
			if i < len(tableDef.TableRows)-1 {
				query += ","
			}
		}

		query += ")"

		// Execute the query
		_, err := db.Exec(query)
		if err != nil {
			args.Logger.Fatal().Err(err).Msgf("Failed to create table %s", tableDef.TableName)
		}

		fmt.Printf("Table %s created successfully.\n", tableDef.TableName)
	}

	return &DatabaseService{
		db:     db,
		logger: args.Logger,
	}
}

func (d *DatabaseService) GetEntityByFilterRule(ctx context.Context, filterRule func(connection interface{}) (interface{}, error)) (interface{}, error) {
	// Execute the filter rule to fetch data from database
	obj, err := filterRule(d.db)
	if err != nil {
		d.logger.Error().Err(err).Msg("Failed to execute filter rule")
		return nil, err
	}

	return obj, nil
}

func (d *DatabaseService) UpdateEntityByFilterRule(ctx context.Context, filterRule func(connection interface{}) (bool, error)) (bool, error) {
	result, err := filterRule(d.db)
	if err != nil {
		d.logger.Error().Err(err).Msg("Failed to execute update rule")
		return false, err
	}

	return result, nil
}

func (d *DatabaseService) DeleteEntityByFilterRule(ctx context.Context, filterRule func(connection interface{}) (bool, error)) (bool, error) {
	result, err := filterRule(d.db)
	if err != nil {
		d.logger.Error().Err(err).Msg("Failed to execute delete rule")
		return false, err
	}

	return result, nil
}

func (d *DatabaseService) AddEntity(ctx context.Context, del func(dbConnection interface{}) error) error {
	err := del(d.db)
	if err != nil {
		d.logger.Error().Err(err).Msg("Failed to add entity")
		return err
	}

	return nil
}

// Close method to close the database connection
func (d *DatabaseService) Close() error {
	err := d.db.Close()
	if err != nil {
		d.logger.Error().Err(err).Msg("Failed to close database connection")
		return err
	}
	return nil
}
