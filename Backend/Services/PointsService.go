package Services

import (
	Interfaces2 "Fetch-Rewards-API/Backend/Interfaces"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	Structs2 "Fetch-Rewards-API/Shared/Structs"
)

// pointsService implements the Interfaces2.PointsService interface
type pointsService struct {
	cfg      *Structs2.Config
	ruleList []func(*Structs2.Receipt) int
}

// NewPointsService creates a new instance of PointsService
func NewPointsService(cfg *Structs2.Config) Interfaces2.PointsService {
	ps := &pointsService{
		cfg: cfg,
	}

	//The logic for calculating points is more extendable if we use the strategy pattern
	//This change also ensure open closed, and SRP
	ps.ruleList = []func(*Structs2.Receipt) int{
		ps.pointsForRetailerName,
		ps.pointsForRoundTotal,
		ps.pointsForMultipleOfQuarter,
		ps.pointsForItems,
		ps.pointsForItemDescriptions,
		ps.pointsForOddDay,
		ps.pointsForAfternoonPurchase,
	}

	return ps
}

// CalculatePoints calculates the points for a given receipt based on various rules
func (ps *pointsService) CalculatePoints(receipt *Structs2.Receipt) int {
	points := 0
	for _, rule := range ps.ruleList {
		points += rule(receipt)
	}
	return points
}

// pointsForRetailerName calculates points based on the retailer name
func (ps *pointsService) pointsForRetailerName(receipt *Structs2.Receipt) int {
	return len(regexp.MustCompile("[^a-zA-Z0-9]+").ReplaceAllString(receipt.Retailer, ""))
}

// pointsForRoundTotal calculates points if the total is a round dollar amount
func (ps *pointsService) pointsForRoundTotal(receipt *Structs2.Receipt) int {
	totalFloat, err := strconv.ParseFloat(receipt.Total, 64)
	if err == nil && math.Mod(totalFloat, 1) == 0 {
		return 50
	}
	return 0
}

// pointsForMultipleOfQuarter calculates points if the total is a multiple of 0.25
func (ps *pointsService) pointsForMultipleOfQuarter(receipt *Structs2.Receipt) int {
	totalFloat, err := strconv.ParseFloat(receipt.Total, 64)
	if err == nil && math.Mod(totalFloat, 0.25) == 0 {
		return 25
	}
	return 0
}

// pointsForItems calculates points for the number of items on the receipt
func (ps *pointsService) pointsForItems(receipt *Structs2.Receipt) int {
	return len(receipt.Items) / 2 * 5
}

// pointsForItemDescriptions calculates points based on the item descriptions
func (ps *pointsService) pointsForItemDescriptions(receipt *Structs2.Receipt) int {
	points := 0
	for _, item := range receipt.Items {
		trimmedLength := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLength > 0 && trimmedLength%3 == 0 {
			priceFloat, _ := strconv.ParseFloat(item.Price, 64)
			additionalPoints := int(math.Ceil(priceFloat * 0.2))
			points += additionalPoints
		}
	}
	return points
}

// pointsForOddDay calculates points if the purchase date is an odd day
func (ps *pointsService) pointsForOddDay(receipt *Structs2.Receipt) int {
	purchaseDateTime, err := time.Parse("2006-01-02 15:04", receipt.PurchaseDate+" "+receipt.PurchaseTime)
	if err == nil && purchaseDateTime.Day()%2 != 0 {
		return 6
	}
	return 0
}

// pointsForAfternoonPurchase calculates points if the purchase time is between 2:00pm and 4:00pm
func (ps *pointsService) pointsForAfternoonPurchase(receipt *Structs2.Receipt) int {
	purchaseTimeParsed, _ := time.Parse("15:04", receipt.PurchaseTime)
	startTime, _ := time.Parse("15:04", "14:00")
	endTime, _ := time.Parse("15:04", "16:00")
	if purchaseTimeParsed.After(startTime) && purchaseTimeParsed.Before(endTime) {
		return 10
	}
	return 0
}
