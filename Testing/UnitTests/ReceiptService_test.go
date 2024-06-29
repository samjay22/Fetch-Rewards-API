package Testing

import (
	Testing "Fetch-Rewards-API/Testing/Mocks"
	"context"
	"testing"

	Interfaces "Fetch-Rewards-API/Backend/Interfaces"
	Structs "Fetch-Rewards-API/Shared/Structs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestReceiptService_GetPointsForReceiptById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReceiptService := Testing.NewMockReceiptService(ctrl)

	// Prepare mock expectations
	expectedPoints := int64(100)
	mockReceiptService.EXPECT().GetPointsForReceiptById("receipt123").Return(expectedPoints, nil)

	// Call the method being tested
	points, err := mockReceiptService.GetPointsForReceiptById("receipt123")

	// Verify expectations and assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedPoints, points)
}

func TestReceiptService_GetReceipts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReceiptService := Testing.NewMockReceiptService(ctrl)

	// Prepare mock expectations
	expectedPayload := &Interfaces.SearchPagePayload{}
	mockReceiptService.EXPECT().GetReceipts(gomock.Any(), gomock.Any(), 1).Return(expectedPayload, nil)

	// Call the method being tested
	payload, err := mockReceiptService.GetReceipts(context.Background(), nil, 1)

	// Verify expectations and assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedPayload, payload)
}

func TestReceiptService_ProcessReceipt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReceiptService := Testing.NewMockReceiptService(ctrl)

	// Prepare mock expectations
	mockReceipt := &Structs.Receipt{Id: "receipt123"}
	mockReceiptService.EXPECT().ProcessReceipt(mockReceipt).Return(nil)

	// Call the method being tested
	err := mockReceiptService.ProcessReceipt(mockReceipt)

	// Verify expectations and assertions
	assert.NoError(t, err)
}
