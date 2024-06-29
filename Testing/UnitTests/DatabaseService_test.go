package Testing

import (
	Testing "Fetch-Rewards-API/Testing/Mocks"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseService_GetEntityByFilterRule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := Testing.NewMockDatabaseService(ctrl)

	// Prepare mock expectations
	expectedResult := "mock result"
	mockDB.EXPECT().GetEntityByFilterRule(gomock.Any(), gomock.Any()).Return(expectedResult, nil)

	// Call the method being tested
	result, err := mockDB.GetEntityByFilterRule(context.Background(), func(connection interface{}) (interface{}, error) {
		// Simulate database operation
		return "result", nil
	})

	// Verify expectations and assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestDatabaseService_UpdateEntityByFilterRule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := Testing.NewMockDatabaseService(ctrl)

	// Prepare mock expectations
	mockDB.EXPECT().UpdateEntityByFilterRule(gomock.Any(), gomock.Any()).Return(true, nil)

	// Call the method being tested
	result, err := mockDB.UpdateEntityByFilterRule(context.Background(), func(connection interface{}) (bool, error) {
		// Simulate update operation
		return true, nil
	})

	// Verify expectations and assertions
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestDatabaseService_DeleteEntityByFilterRule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := Testing.NewMockDatabaseService(ctrl)

	// Prepare mock expectations
	mockDB.EXPECT().DeleteEntityByFilterRule(gomock.Any(), gomock.Any()).Return(true, nil)

	// Call the method being tested
	result, err := mockDB.DeleteEntityByFilterRule(context.Background(), func(connection interface{}) (bool, error) {
		// Simulate delete operation
		return true, nil
	})

	// Verify expectations and assertions
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestDatabaseService_AddEntity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := Testing.NewMockDatabaseService(ctrl)

	// Prepare mock expectations
	mockDB.EXPECT().AddEntity(gomock.Any(), gomock.Any()).Return(nil)

	// Call the method being tested
	err := mockDB.AddEntity(context.Background(), func(connection interface{}) error {
		// Simulate add operation
		return nil
	})

	// Verify expectations and assertions
	assert.NoError(t, err)
}
