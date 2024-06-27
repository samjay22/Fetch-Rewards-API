package Services

import (
	"Fetch-Rewards-API/Interfaces"
	"Fetch-Rewards-API/Structs"
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
	EventBus Interfaces.EventBus
	Delegate func(database *sql.DB)
}

// Use factory pattern
// NewDatabaseService method decorator for extended usages. We use a generic, each type may require different db settings
func NewDatabaseService(args *NewDatabaseServiceArgs) Interfaces.DatabaseService {
	dbConfig := args.Cfg.Database
	initDb, err := sql.Open("sqlite", fmt.Sprintf("file:%s/%s.db", dbConfig.HomeDir, dbConfig.FileName))
	if err != nil {
		args.Logger.Fatal().Err(err).Msg("Failed to connect to database")
		return nil
	}

	args.Delegate(initDb)

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

		// Execute the query within a transaction
		tx, err := initDb.Begin()
		if err != nil {
			args.Logger.Fatal().Err(err)
		}

		_, err = tx.Exec(query)
		if err != nil {
			tx.Rollback()
			args.Logger.Fatal().Err(err)
		}

		err = tx.Commit()
		if err != nil {
			args.Logger.Fatal().Err(err)
		}

		fmt.Printf("Table %s created successfully.\n", tableDef.TableName)
	}

	return &DatabaseService{
		db:     initDb,
		logger: args.Logger,
	}
}

func (d *DatabaseService) openConnection(ctx context.Context) (*interface{}, error) {
	//TODO implement me

	con, err := d.db.Conn(ctx)
	if err != nil {
		d.logger.Error().Err(err).Msg("Failed to connect to database")
		return nil, err
	}

	returnValue := interface{}(con)
	return &returnValue, nil
}

func (d *DatabaseService) closeConnection(connection *sql.Conn, ctx context.Context) error {
	return connection.Close()
}

func (d *DatabaseService) GetEntityByFilterRule(ctx context.Context, filterRule func(connection interface{}) (*interface{}, error)) (*interface{}, error) {
	conn, err := d.openConnection(ctx)
	if err != nil {
		d.logger.Error().Err(err).Msg("Failed to connect to database")
		return nil, err
	}

	go func() {
		_, err := filterRule(conn)
		if err != nil {
			d.logger.Error().Err(err).Ctx(ctx).Msg("Failed to execute filter rule for get entity")
		}
	}()

	return nil, err
}

func (d *DatabaseService) UpdateEntityByFilterRule(ctx context.Context, filterRule func(connection interface{}) (bool, error)) (bool, error) {
	return true, nil
}

func (d *DatabaseService) DeleteEntityByFilterRule(ctx context.Context, filterRule func(connection interface{}) (bool, error)) (bool, error) {
	return true, nil
}

func (d *DatabaseService) AddEntity(ctx context.Context, del func(dbConnection interface{}) error) error {
	go func() {
		err := del(d.db)
		if err != nil {
			d.logger.Error().Err(err)
		}
	}()

	return nil
}
