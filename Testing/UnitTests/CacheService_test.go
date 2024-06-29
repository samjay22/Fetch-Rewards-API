package Testing

import (
	Testing "Fetch-Rewards-API/Testing/Mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCacheService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := Testing.NewMockCacheService(ctrl)

	// Prepare mock expectations
	expectedResult := "mock result"
	mockCache.EXPECT().Get(gomock.Any()).Return(expectedResult, nil)

	// Call the method being tested
	result, err := mockCache.Get("key")

	// Verify expectations and assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestCacheService_Set(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := Testing.NewMockCacheService(ctrl)

	// Prepare mock expectations
	mockCache.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil)

	// Call the method being tested
	err := mockCache.Set("key", "value")

	// Verify expectations and assertions
	assert.NoError(t, err)
}

func TestCacheService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := Testing.NewMockCacheService(ctrl)

	// Prepare mock expectations
	mockCache.EXPECT().Delete(gomock.Any()).Return(nil)

	// Call the method being tested
	err := mockCache.Delete("key")

	// Verify expectations and assertions
	assert.NoError(t, err)
}
