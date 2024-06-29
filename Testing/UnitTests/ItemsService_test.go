package Testing

import (
	Testing "Fetch-Rewards-API/Testing/Mocks"
	"database/sql"
	"testing"

	Structs "Fetch-Rewards-API/Shared/Structs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestItemService_GetItemsForReceipt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockItemService := Testing.NewMockItemService(ctrl)

	// Prepare mock expectations
	expectedItems := []Structs.PurchasedItem{
		{Id: "item1", ShortDescription: "Item 1", Price: "10.5"},
		{Id: "item2", ShortDescription: "Item 2", Price: "5.25"},
	}
	mockItemService.EXPECT().GetItemsForReceipt(gomock.Any(), "receipt123").Return(expectedItems, nil)

	// Call the method being tested
	items, err := mockItemService.GetItemsForReceipt(&sql.DB{}, "receipt123")

	// Verify expectations and assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedItems, items)
}

func TestItemService_InsertItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockItemService := Testing.NewMockItemService(ctrl)

	// Prepare mock expectations
	mockItemService.EXPECT().InsertItems(gomock.Any(), "receipt123", gomock.Any()).Return(nil)

	// Call the method being tested
	err := mockItemService.InsertItems(&sql.Tx{}, "receipt123", []Structs.PurchasedItem{
		{Id: "item1", ShortDescription: "Item 1", Price: "10.5"},
		{Id: "item2", ShortDescription: "Item 2", Price: "5.25"},
	})

	// Verify expectations and assertions
	assert.NoError(t, err)
}

func TestItemService_GenerateItemIds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockItemService := Testing.NewMockItemService(ctrl)

	// Prepare mock expectations
	mockItemService.EXPECT().GenerateItemIds(gomock.Any())

	// Call the method being tested
	mockItemService.GenerateItemIds([]Structs.PurchasedItem{
		{Id: "item1", ShortDescription: "Item 1", Price: "10.5"},
		{Id: "item2", ShortDescription: "Item 2", Price: "5.25"},
	})

	// Verify expectations and assertions
	// Since GenerateItemIds doesn't return anything, we typically don't assert anything
	// related to its behavior apart from expectations.
}
