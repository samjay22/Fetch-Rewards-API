package Controllers

import (
	"Fetch-Rewards-API/Backend/Interfaces"
	Structs2 "Fetch-Rewards-API/Shared/Structs"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"net/http"
	"strconv"
)

//Controller Struct
type receiptController struct {
	ReceiptService Interfaces.ReceiptService
	NetworkClient  *echo.Echo
	Logger         *zerolog.Logger
	QueueService   Interfaces.QueueService
}

// ReceiptControllerArgs this struct cleans up the arguments to the controller init function.
type ReceiptControllerArgs struct {
	Logger       *zerolog.Logger
	EchoClient   *echo.Echo
	DataService  Interfaces.ReceiptService
	QueueService Interfaces.QueueService
}

// RegisterReceiptController constructor
func RegisterReceiptController(args *ReceiptControllerArgs) {
	controller := &receiptController{
		ReceiptService: args.DataService,
		NetworkClient:  args.EchoClient,
		Logger:         args.Logger,
		QueueService:   args.QueueService,
	}

	//Register the handler for the event topic
	controller.QueueService.RegisterEventHandler("processReceipt", controller.ProcessReceipt)
	go controller.QueueService.ProcessQueue()

	//Register endpoints for controller
	controller.NetworkClient.GET("/receipts", controller.GetReceipts)
	controller.NetworkClient.POST("/receipts/process", controller.QueueReceipt)
	controller.NetworkClient.GET("/receipts/:id/points", controller.QueuePoints)
}

// QueueReceipt queue an event to process a receipt
func (r *receiptController) QueueReceipt(c echo.Context) error {
	receiptID := uuid.New().String()
	var req Structs2.Receipt
	req.Id = receiptID
	err := json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		r.Logger.Error().Err(err).Msg("Error decoding request body")
		return echo.NewHTTPError(http.StatusBadRequest, "The receipt is invalid")
	}

	// Submit the event for the topic
	response, err := r.QueueService.QueueEvent("processReceipt", &req)
	if err != nil {
		r.Logger.Error().Err(err).Msg("Error processing receipt")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to process receipt")
	}

	return c.JSON(http.StatusOK, response)
}

// QueuePoints we keep the naming convention, this gets points for a given ID
func (r *receiptController) QueuePoints(c echo.Context) error {
	receiptID := c.Param("id")

	points, err := r.ReceiptService.GetPointsForReceiptById(receiptID)
	if err != nil {
		r.Logger.Error().Err(err).Msg("Error retrieving points for receipt")
		return echo.NewHTTPError(http.StatusNotFound, "No receipt found for that Id")
	}

	return c.JSON(http.StatusOK, map[string]int64{"points": points})
}

//ProcessReceipt this the a handler for the event topic to process client receipts 
func (r *receiptController) ProcessReceipt(data interface{}) (interface{}, error) {
	if receipt, ok := data.(*Structs2.Receipt); ok {
		err := r.ReceiptService.ProcessReceipt(receipt)
		if err != nil {
			r.Logger.Error().Err(err).Msg("Error processing ProcessReceipt request")
			return nil, err
		}
		return &ProcessReceiptRequestReturn{Id: receipt.Id}, nil
	} else {
		r.Logger.Error().Msg("Invalid data type received for processReceipt event")
		return nil, fmt.Errorf("invalid data type")
	}
}

// GetReceipts get all recipts based on specified filter critiera
func (r *receiptController) GetReceipts(c echo.Context) error {
	id := c.QueryParam("id")
	retailer := c.QueryParam("retailer")
	purchaseDate := c.QueryParam("purchaseDate")
	purchaseTime := c.QueryParam("purchaseTime")
	total := c.QueryParam("total")
	points := c.QueryParam("points")
	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	queryFilter := createQueryFilter(id, retailer, purchaseDate, purchaseTime, total, points)
	receipts, err := r.ReceiptService.GetReceipts(context.Background(), queryFilter, page)
	if err != nil {
		r.Logger.Error().Err(err).Msg("Error retrieving receipts")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve receipts")
	}

	return c.JSON(http.StatusOK, receipts)
}

func createQueryFilter(id, retailer, purchaseDate, purchaseTime, total, points string) *Interfaces.ReceiptFilterRule {
	return &Interfaces.ReceiptFilterRule{
		Id:           id,
		Retailer:     retailer,
		PurchaseDate: purchaseDate,
		PurchaseTime: purchaseTime,
		Total:        total,
		Points:       points,
	}
}

type ProcessReceiptRequestReturn struct {
	Id string `json:"id"`
}
