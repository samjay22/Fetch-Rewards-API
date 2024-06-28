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

func loadConfig(logger *zerolog.Logger) (*Structs.Config, error) {
	f, err := os.Open("Backend/Configs/ENV.yml")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to open config.yml")
		return nil, err
	}
	defer f.Close()

	var cfg Structs.Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to parse config.yml")
		return nil, err
	}

	return &cfg, nil
}

type ServicesContainer struct {
	dataService    Interfaces2.DatabaseService
	receiptService Interfaces2.ReceiptService
}

func initServices(logger *zerolog.Logger, cfg *Structs.Config) (*ServicesContainer, error) {

	args := Services2.NewDatabaseServiceArgs{
		Logger: logger,
		Cfg:    cfg,
		Delegate: func(db *sql.DB) {
			db.SetConnMaxLifetime(time.Minute * 5)
			db.SetConnMaxIdleTime(time.Minute)
			db.SetMaxOpenConns(10)
		},
	}

	dataService := Services2.NewDatabaseService(&args)
	recService := Services2.NewReceiptService(logger, dataService)
	return &ServicesContainer{
		dataService:    dataService,
		receiptService: recService,
	}, nil
}

func registerMiddleware(echoClient *echo.Echo, logger *zerolog.Logger) {
	echoClient.Use(middleware.RequestID())
	echoClient.Use(middleware.AddTrailingSlash())
	echoClient.Use(middleware.Recover())
	echoClient.Use(Middleware.Logger(logger))
}

// Uncle Bob AKA Robert C. Martin suggested to avoid having more than 3 arguments in a function, always compress if possible.
type controllerArgs struct {
	echoClient *echo.Echo
	logger     *zerolog.Logger
	cfg        *Structs.Config
	container  *ServicesContainer
}

func registerControllers(args *controllerArgs) error {
	Controllers.RegisterReceiptController(args.logger, args.echoClient, args.container.receiptService)
	return nil
}

func main() {
	logger := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()

	logger.Log().Msg("Logging Active!")

	//We already logged the error
	config, err := loadConfig(&logger)
	if err != nil {
		return
	}

	logger.Log().Msg("Loaded Configs!")

	serviceContainer, err := initServices(&logger, config)
	if err != nil {
		return
	}

	logger.Log().Msg("Loaded Services!")

	//Register API controllers, middleware, and static pages
	networkClient := echo.New()
	networkClient.Static("/", "Frontend/Pages")
	registerMiddleware(networkClient, &logger)
	logger.Log().Msg("Loaded Middleware!")

	args := controllerArgs{
		echoClient: networkClient,
		logger:     &logger,
		cfg:        config,
		container:  serviceContainer,
	}

	err = registerControllers(&args)
	if err != nil {
		logger.Fatal().Err(err)
		return
	}

	logger.Log().Msg("Loaded Controllers!")

	hostString := fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
	logger.Fatal().Err(networkClient.StartTLS(hostString, config.Server.SSLCert, config.Server.SSLKey))
}
