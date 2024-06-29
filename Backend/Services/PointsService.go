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

type pointsService struct {
	cfg *Structs2.Config
}

func NewPointsService(cfg *Structs2.Config) Interfaces2.PointsService {
	return &pointsService{
		cfg: cfg,
	}
}

func (ps *pointsService) CalculatePoints(receipt *Structs2.Receipt) int {
	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name
	points += len(regexp.MustCompile("[^a-zA-Z0-9]+").ReplaceAllString(receipt.Retailer, ""))

	// Rule 2: 50 points if the total is a round dollar amount with no cents
	totalFloat, err := strconv.ParseFloat(receipt.Total, 64)
	if err == nil && math.Mod(totalFloat, 1) == 0 {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25
	if totalFloat > 0 && math.Mod(totalFloat, 0.25) == 0 {
		points += 25
	}

	// Rule 4: 5 points for every two items on the receipt
	points += len(receipt.Items) / 2 * 5

	// Rule 5: If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer
	for _, item := range receipt.Items {
		trimmedLength := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLength > 0 && trimmedLength%3 == 0 {
			priceFloat, _ := strconv.ParseFloat(item.Price, 64)
			additionalPoints := int(math.Ceil(priceFloat * 0.2))
			points += additionalPoints
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd
	purchaseDateTime, err := time.Parse("2006-01-02 15:04", receipt.PurchaseDate+" "+receipt.PurchaseTime)
	if err == nil && purchaseDateTime.Day()%2 != 0 {
		points += 6
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm
	purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
	startTime, _ := time.Parse("15:04", "14:00")
	endTime, _ := time.Parse("15:04", "16:00")
	if purchaseTime.After(startTime) && purchaseTime.Before(endTime) {
		points += 10
	}

	return points
}
