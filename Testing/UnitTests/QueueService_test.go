package Testing

import (
	Testing "Fetch-Rewards-API/Testing/Mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestQueueService_DispatchEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueueService := Testing.NewMockQueueService(ctrl)

	// Prepare mock expectations
	expectedData := "test data"
	mockQueueService.EXPECT().DispatchEvent("event_type", expectedData).Return("result", nil)

	// Call the method being tested
	result, err := mockQueueService.DispatchEvent("event_type", expectedData)

	// Verify expectations and assertions
	assert.NoError(t, err)
	assert.Equal(t, "result", result)
}

func TestQueueService_ProcessQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueueService := Testing.NewMockQueueService(ctrl)

	// Prepare mock expectations
	mockQueueService.EXPECT().ProcessQueue()

	// Call the method being tested
	mockQueueService.ProcessQueue()
}

func TestQueueService_QueueEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueueService := Testing.NewMockQueueService(ctrl)

	// Prepare mock expectations
	expectedData := "test data"
	mockQueueService.EXPECT().QueueEvent("event_type", expectedData).Return("result", nil)

	// Call the method being tested
	result, err := mockQueueService.QueueEvent("event_type", expectedData)

	// Verify expectations and assertions
	assert.NoError(t, err)
	assert.Equal(t, "result", result)
}

func TestQueueService_RegisterEventHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueueService := Testing.NewMockQueueService(ctrl)

	// Prepare mock expectations
	mockHandler := func(event interface{}) (interface{}, error) {
		return event, nil
	}
	mockQueueService.EXPECT().RegisterEventHandler("event_type", gomock.Any())

	// Call the method being tested
	mockQueueService.RegisterEventHandler("event_type", mockHandler)
}
