package Testing

import (
	Testing "Fetch-Rewards-API/Testing/Mocks"
	"testing"

	Structs "Fetch-Rewards-API/Shared/Structs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPointsService_CalculatePoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPointsService := Testing.NewMockPointsService(ctrl)

	receipt := &Structs.Receipt{
		Id:    "receipt123",
		Total: "100.0",
	}
	expectedPoints := 50
	mockPointsService.EXPECT().CalculatePoints(receipt).Return(expectedPoints)

	points := mockPointsService.CalculatePoints(receipt)

	assert.Equal(t, expectedPoints, points)
}
