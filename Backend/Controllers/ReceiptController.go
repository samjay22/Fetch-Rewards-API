package Controllers

import (
	"Fetch-Rewards-API/Backend/Interfaces"
	Structs2 "Fetch-Rewards-API/Shared/Structs"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"net/http"
	"strconv"
	"sync"
)

type receiptController struct {
	ReceiptService Interfaces.ReceiptService
	NetworkClient  *echo.Echo
	Logger         *zerolog.Logger

	eventHandlers map[string]EventDelegate
	eventQueue    chan Structs2.Event
	mu            sync.RWMutex
}

type EventDelegate func(data interface{})

func RegisterReceiptController(logger *zerolog.Logger, echoClient *echo.Echo, dataService Interfaces.ReceiptService) {
	controller := &receiptController{
		ReceiptService: dataService,
		NetworkClient:  echoClient,
		Logger:         logger,
		eventHandlers:  make(map[string]EventDelegate),
		eventQueue:     make(chan Structs2.Event, 100), // Buffer size can be adjusted based on expected load
	}

	controller.RegisterEventHandler("processReceipt", controller.ProcessReceipt)

	go controller.processQueue()

	controller.NetworkClient.GET("/receipts", controller.GetReceipts)
	controller.NetworkClient.POST("/receipts/process", controller.QueueReceipt)
	controller.NetworkClient.GET("/receipts/:id/points", controller.QueuePoints)
}

func (r *receiptController) RegisterEventHandler(eventType string, handler EventDelegate) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.eventHandlers[eventType] = handler
}

func (r *receiptController) DispatchEvent(eventType string, data interface{}) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if handler, ok := r.eventHandlers[eventType]; ok {
		handler(data)
	} else {
		r.Logger.Warn().Str("event", eventType).Msg("No handler registered for this event")
	}
}

func (r *receiptController) QueueReceipt(c echo.Context) error {
	receiptID := uuid.New().String()
	var req Structs2.Receipt
	req.Id = receiptID
	err := json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		r.Logger.Error().Err(err).Msg("Error decoding request body")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON request")
	}

	r.DispatchEvent("processReceipt", &req)

	return c.JSON(http.StatusOK, &ProcessReceiptRequestReturn{Id: req.Id})
}

func (r *receiptController) QueuePoints(c echo.Context) error {
	receiptID := c.Param("id")

	points, err := r.ReceiptService.GetPointsForReceiptById(receiptID)
	if err != nil {
		r.Logger.Error().Err(err).Msg("Error retrieving points for receipt")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve points for receipt")
	}

	return c.JSON(http.StatusOK, map[string]int64{"points": points})
}

func (r *receiptController) ProcessReceipt(data interface{}) {
	if receipt, ok := data.(*Structs2.Receipt); ok {
		err := r.ReceiptService.ProcessReceipt(receipt)
		if err != nil {
			r.Logger.Error().Err(err).Msg("Error processing ProcessReceipt request")
		}
	} else {
		r.Logger.Error().Msg("Invalid data type received for processReceipt event")
	}
}

func (r *receiptController) ProcessPointsRequest(data interface{}) {
	if receiptID, ok := data.(string); ok {
		_, err := r.ReceiptService.GetPointsForReceiptById(receiptID)
		if err != nil {
			r.Logger.Error().Err(err).Msg("Error processing ProcessReceipt request")
		}
	} else {
		r.Logger.Error().Msg("Invalid data type received for processReceipt event")
	}
}

// GetReceipts handles fetching receipts based on search and pagination
func (r *receiptController) GetReceipts(c echo.Context) error {
	// Retrieve query parameters for search and pagination
	id := c.QueryParam("id")
	retailer := c.QueryParam("retailer")
	purchaseDate := c.QueryParam("purchaseDate")
	purchaseTime := c.QueryParam("purchaseTime")
	total := c.QueryParam("total")
	points := c.QueryParam("points")
	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1 // Default to page 1 if invalid or missing page parameter
	}

	// Create a query filter based on search terms
	queryFilter := createQueryFilter(id, retailer, purchaseDate, purchaseTime, total, points)

	// Implement logic to fetch receipts with search and pagination from the service
	queryFilter.ApplyDefaults()
	receipts, err := r.ReceiptService.GetReceipts(context.Background(), queryFilter, page)
	if err != nil {
		r.Logger.Error().Err(err).Msg("Error retrieving receipts")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve receipts")
	}

	// Return the receipts as JSON response
	return c.JSON(http.StatusOK, receipts)
}

// createQueryFilter parses the search terms and constructs a filter rule
func createQueryFilter(id, retailer, purchaseDate, purchaseTime, total, points string) *Interfaces.ReceiptFilterRule {
	// Initialize an empty filter rule
	queryFilter := &Interfaces.ReceiptFilterRule{}

	// Example logic: set filter fields based on provided search terms
	if id != "" {
		queryFilter.Id = id
	}

	if retailer != "" {
		queryFilter.Retailer = retailer
	}

	if purchaseDate != "" {
		queryFilter.PurchaseDate = purchaseDate
	}

	if purchaseTime != "" {
		queryFilter.PurchaseTime = purchaseTime
	}

	if total != "" {
		queryFilter.Total = total
	}

	if points != "" {
		queryFilter.Points = points
	}

	// Add more conditions for other query parameters as needed

	return queryFilter
}

func (r *receiptController) processQueue() {
	for {
		select {
		case event := <-r.eventQueue:
			r.DispatchEvent(event.Type, event.Data)
		}
	}
}

type ProcessReceiptRequestReturn struct {
	Id string `json:"id"`
}
