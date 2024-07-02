package Services

import (
	Interfaces2 "Fetch-Rewards-API/Backend/Interfaces"
	"Fetch-Rewards-API/Backend/ServerUtility"
	Structs2 "Fetch-Rewards-API/Shared/Structs"
	"context"
	"database/sql"
	"fmt"
	"math"

	"github.com/rs/zerolog"
)

type receiptService struct {
	logger        *zerolog.Logger
	cfg           *Structs2.Config
	dataService   Interfaces2.DatabaseService
	itemService   Interfaces2.ItemService
	pointsService Interfaces2.PointsService
	cacheService  Interfaces2.CacheService
}

type NewReceiptServiceArgs struct {
	Logger        *zerolog.Logger
	Cfg           *Structs2.Config
	DataService   Interfaces2.DatabaseService
	ItemService   Interfaces2.ItemService
	PointsService Interfaces2.PointsService
	CacheService  Interfaces2.CacheService
}

func NewReceiptService(args *NewReceiptServiceArgs) Interfaces2.ReceiptService {
	return &receiptService{
		logger:        args.Logger,
		dataService:   args.DataService,
		itemService:   args.ItemService,
		pointsService: args.PointsService,
		cfg:           args.Cfg,
		cacheService:  args.CacheService,
	}
}

func (rt *receiptService) GetPointsForReceiptById(id string) (int64, error) {
	entity, _ := rt.cacheService.Get(id)
	if entity != nil {
		if receipt, ok := entity.(*Structs2.Receipt); ok {
			rt.logger.Info().Str("id", id).Msg("Retrieved receipt from cache")
			return receipt.Points, nil
		}
	}

	receipt, err := rt.getReceiptById(id)
	if err != nil {
		return 0, err
	}

	if err = rt.cacheService.Set(id, receipt); err != nil {
		rt.logger.Error().Err(err).Msg("Failed to cache receipt")
		return 0, err
	}

	return receipt.Points, nil
}

func (rt *receiptService) ProcessReceipt(receiptEntity *Structs2.Receipt) error {
	//We are working with a ref, we can set id's without returning the struct
	rt.itemService.GenerateItemIds(receiptEntity.Items)
	points := rt.pointsService.CalculatePoints(receiptEntity)

	if err := rt.insertReceiptAndItems(receiptEntity, points); err != nil {
		return fmt.Errorf("failed to process receipt: %w", err)
	}

	return nil
}

func (rt *receiptService) GetReceipts(ctx context.Context, filterBy *Interfaces2.ReceiptFilterRule, page int) (*Interfaces2.SearchPagePayload, error) {
	const pageSize = 15
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	cacheKey := fmt.Sprintf("receipts_%s_%d", filterBy.ToString(), page)
	if pageData, _ := rt.cacheService.Get(cacheKey); pageData != nil {
		if obj, ok := pageData.(*Interfaces2.SearchPagePayload); ok {
			rt.logger.Info().Str("cacheKey", cacheKey).Msg("Retrieved receipts from cache")
			return obj, nil
		}
	}

	filterFunc := rt.buildReceiptsFilterFunc(ctx, filterBy, pageSize, offset)
	data, err := rt.dataService.GetEntityByFilterRule(ctx, filterFunc)
	if err != nil {
		rt.logger.Error().Err(err).Msg("Failed to retrieve receipts from database")
		return nil, err
	}

	receipts, ok := data.(*Interfaces2.SearchPagePayload)
	if !ok {
		return nil, fmt.Errorf("unexpected type for response data")
	}

	if err = rt.cacheService.Set(cacheKey, receipts); err != nil {
		rt.logger.Error().Err(err).Msg("Failed to cache receipts")
		return nil, err
	}

	return receipts, nil
}

func (rt *receiptService) getReceiptById(id string) (*Structs2.Receipt, error) {
	r, err := rt.dataService.GetEntityByFilterRule(context.Background(), func(dbI interface{}) (interface{}, error) {
		receipt := &Structs2.Receipt{}
		db, ok := dbI.(*sql.DB)
		if !ok {
			return nil, fmt.Errorf("unexpected type for database connection")
		}

		row := db.QueryRow("SELECT Id, Retailer, PurchaseDate, PurchaseTime, Total, Points FROM Receipts WHERE Id = ?", id)
		if err := row.Scan(&receipt.Id, &receipt.Retailer, &receipt.PurchaseDate, &receipt.PurchaseTime, &receipt.Total, &receipt.Points); err != nil {
			return nil, err
		}

		items, err := rt.itemService.GetItemsForReceipt(db, id)
		if err != nil {
			return nil, err
		}
		receipt.Items = items

		return receipt, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error retrieving receipt from database: %w", err)
	}

	receipt, ok := r.(*Structs2.Receipt)
	if !ok {
		return nil, fmt.Errorf("unexpected type for receipt")
	}

	return receipt, nil
}

func (rt *receiptService) insertReceiptAndItems(receiptEntity *Structs2.Receipt, points int) error {
	//Invalidate cache
	err := rt.cacheService.Purge()
	if err != nil {
		return fmt.Errorf("error with purge cache request!")
	}

	err = rt.dataService.AddEntity(context.Background(), func(i interface{}) error {
		db, ok := i.(*sql.DB)
		if !ok {
			return fmt.Errorf("invalid database connection")
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer func() {
			if err != nil {
				tx.Rollback()
				rt.logger.Error().Err(err).Msg("Transaction rolled back")
			} else {
				tx.Commit()
				rt.logger.Info().Msg("Transaction committed successfully")
			}
		}()

		_, err = tx.ExecContext(context.Background(), "INSERT INTO receipts (Id, Retailer, PurchaseDate, PurchaseTime, Total, Points) VALUES (?, ?, ?, ?, ?, ?)",
			receiptEntity.Id, receiptEntity.Retailer, receiptEntity.PurchaseDate, receiptEntity.PurchaseTime, receiptEntity.Total, points)
		if err != nil {
			return err
		}

		if err = rt.itemService.InsertItems(tx, receiptEntity.Id, receiptEntity.Items); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		rt.logger.Error().Err(err).Msg("Failed to insert receipt into database")
		return fmt.Errorf("failed to insert receipt into database: %w", err)
	}

	return nil
}

func (rt *receiptService) buildReceiptsFilterFunc(ctx context.Context, filterBy *Interfaces2.ReceiptFilterRule, pageSize, offset int) func(interface{}) (interface{}, error) {
	return func(db interface{}) (interface{}, error) {
		dbInstance, ok := db.(*sql.DB)
		if !ok {
			return nil, fmt.Errorf("unexpected database instance type")
		}

		// Initialize MySQLQueryBuilder
		queryBuilder := ServerUtility.NewMySQLQueryBuilder()

		// Specify fields to select dynamically, entity framework has built in libraries, I needed to do this manually
		fields := []string{"Id", "Retailer", "PurchaseDate", "PurchaseTime", "Total", "Points"}
		queryBuilder.
			Where("Id", filterBy.Id+"%").
			Where("Retailer", filterBy.Retailer+"%").
			Where("PurchaseDate", filterBy.PurchaseDate+"%").
			Where("PurchaseTime", filterBy.PurchaseTime+"%").
			Where("Total", filterBy.Total+"%").
			Where("Points", filterBy.Points+"%").
			Order("Id", "DESC").
			Limit(pageSize).
			Offset(offset).
			SelectFields(fields)

		query := queryBuilder.BuildFullQueryOn("Receipts")

		rows, err := dbInstance.Query(query)
		if err != nil {
			return nil, fmt.Errorf("error executing query: %w", err)
		}
		defer rows.Close()

		// Process retrieved rows into receipt objects
		var receipts []Structs2.Receipt
		for rows.Next() {
			var receipt Structs2.Receipt
			if err := rows.Scan(&receipt.Id, &receipt.Retailer, &receipt.PurchaseDate, &receipt.PurchaseTime, &receipt.Total, &receipt.Points); err != nil {
				return nil, fmt.Errorf("error scanning row: %w", err)
			}

			// Retrieve items associated with each receipt
			items, err := rt.itemService.GetItemsForReceipt(dbInstance, receipt.Id)
			if err != nil {
				return nil, fmt.Errorf("error fetching items for receipt %s: %w", receipt.Id, err)
			}
			receipt.Items = items

			receipts = append(receipts, receipt)
		}

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating over rows: %w", err)
		}

		// Query total count of matching receipts for pagination
		countQuery := queryBuilder.Build()
		var totalRows int
		if err := dbInstance.QueryRow(fmt.Sprintf("SELECT COUNT(Id) FROM Receipts WHERE %s", countQuery)).Scan(&totalRows); err != nil {
			return nil, fmt.Errorf("error getting total count: %w", err)
		}

		totalPages := int(math.Ceil(float64(totalRows) / float64(pageSize)))

		return &Interfaces2.SearchPagePayload{
			Receipts: receipts,
			MaxPages: totalPages,
		}, nil
	}
}
