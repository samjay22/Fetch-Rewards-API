package Interfaces

import Structs2 "Fetch-Rewards-API/Shared/Structs"

type PointsService interface {
	CalculatePoints(receipt *Structs2.Receipt) int
}
