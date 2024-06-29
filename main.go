package main

import (
	"Fetch-Rewards-API/Backend/Controllers"
	Interfaces2 "Fetch-Rewards-API/Backend/Interfaces"
	"Fetch-Rewards-API/Backend/Middleware"
	Services2 "Fetch-Rewards-API/Backend/Services"
	"Fetch-Rewards-API/Shared/Structs"
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type ServicesContainer struct {
	dataService    Interfaces2.DatabaseService
	cacheService   Interfaces2.CacheService
	receiptService Interfaces2.ReceiptService
	itemService    Interfaces2.ItemService
	pointsService  Interfaces2.PointsService
	queueService   Interfaces2.QueueService
}

type App struct {
	logger           zerolog.Logger
	config           *Structs.Config
	serviceContainer *ServicesContainer
	echoClient       *echo.Echo
}

func main() {
	app := App{}

	err := app.initialize()
	if err != nil {
		app.logger.Fatal().Err(err).Msg("Failed to initialize application")
	}

	hostString := fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.Port)
	app.logger.Fatal().Err(app.echoClient.StartTLS(hostString, app.config.Server.SSLCert, app.config.Server.SSLKey))
}

func (app *App) initialize() error {
	// Initialize logger
	app.logger = zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()
	app.logger.Log().Msg("Logging Active!")

	// Load configuration
	config, err := app.loadConfig()
	if err != nil {
		return err
	}
	app.config = config
	app.logger.Log().Msg("Loaded Configs!")

	// Initialize services
	serviceContainer, err := app.initServices(config)
	if err != nil {
		return err
	}
	app.serviceContainer = serviceContainer
	app.logger.Log().Msg("Loaded Services!")

	// Initialize Echo server
	app.echoClient = echo.New()
	app.setupMiddleware()
	app.logger.Log().Msg("Loaded Middleware!")

	// Register controllers
	err = app.registerControllers()
	if err != nil {
		return err
	}
	app.logger.Log().Msg("Loaded Controllers!")

	return nil
}

func (app *App) loadConfig() (*Structs.Config, error) {
	f, err := os.Open("Backend/Configs/ENV.yml")
	if err != nil {
		app.logger.Fatal().Err(err).Msg("Failed to open config.yml")
		return nil, err
	}
	defer f.Close()

	var cfg Structs.Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		app.logger.Fatal().Err(err).Msg("Failed to parse config.yml")
		return nil, err
	}

	return &cfg, nil
}

func (app *App) initServices(cfg *Structs.Config) (*ServicesContainer, error) {
	logger := &app.logger

	cacheService := Services2.NewMemoryCacheService()

	dataServiceArgs := Services2.NewDatabaseServiceArgs{
		Logger: logger,
		Cfg:    cfg,
		Delegate: func(db *sql.DB) {
			dbConfig := cfg.Database
			db.SetConnMaxLifetime(time.Minute * 5)
			db.SetConnMaxIdleTime(time.Minute)
			db.SetMaxOpenConns(-1)

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

				query += ") WITHOUT ROWID"

				// Execute the query
				_, err := db.Exec(query)
				if err != nil {
					logger.Fatal().Err(err).Msgf("Failed to create table %s", tableDef.TableName)
				}

				fmt.Printf("Table %s created successfully.\n", tableDef.TableName)
			}
		},
		ConnectionString: fmt.Sprintf("file:%s/%s.db?cache=shared", cfg.Database.HomeDir, cfg.Database.FileName),
	}
	dataService := Services2.NewDatabaseService(&dataServiceArgs)

	queueService := Services2.NewQueueService(logger)
	itemService := Services2.NewItemService(dataService, cacheService, logger)
	pointsService := Services2.NewPointsService(cfg)

	recServiceArgs := &Services2.NewReceiptServiceArgs{
		Logger:        logger,
		Cfg:           cfg,
		DataService:   dataService,
		ItemService:   itemService,
		PointsService: pointsService,
		CacheService:  cacheService,
	}

	recService := Services2.NewReceiptService(recServiceArgs)

	return &ServicesContainer{
		dataService:    dataService,
		receiptService: recService,
		itemService:    itemService,
		pointsService:  pointsService,
		queueService:   queueService,
		cacheService:   cacheService,
	}, nil
}

func (app *App) setupMiddleware() {
	app.echoClient.Use(middleware.RequestID())
	app.echoClient.Use(middleware.AddTrailingSlash())
	app.echoClient.Use(middleware.Recover())
	app.echoClient.Use(Middleware.Logger(&app.logger))
}

func (app *App) registerControllers() error {
	receiptControllerArgs := &Controllers.ReceiptControllerArgs{
		Logger:       &app.logger,
		EchoClient:   app.echoClient,
		DataService:  app.serviceContainer.receiptService,
		QueueService: app.serviceContainer.queueService,
	}

	app.echoClient.Static("/", "Frontend/Pages")

	Controllers.RegisterReceiptController(receiptControllerArgs)
	return nil
}
