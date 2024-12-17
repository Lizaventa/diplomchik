package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/liza/labwork_45/internal/handlers/mocks"
	"github.com/liza/labwork_45/internal/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	mockCourierServiceInterface *mocks.CourierServiceInterface

	mockDeliveryInstance = &model.Delivery{
		Id:              uuid.New(),
		DeliveryDate:    "2024-13-12",
		DeliveryStatus:  "test_delivered",
		DeliveryComment: "test_comment",
	}
)

func TestCreateDelivery(t *testing.T) {
	mockCourierServiceInterface.On("CreateDelivery", mock.Anything, mock.AnythingOfType("*model.Delivery")).Return(nil)
	err := mockCourierServiceInterface.CreateDelivery(context.Background(), mockDeliveryInstance)
	require.NoError(t, err)
}

func TestCreateDeliveryWithError(t *testing.T) {
	mockCourierServiceInterface.On("CreateDelivery", mock.Anything, mock.AnythingOfType("*model.Delivery")).Return(errors.New("create delivery error"))
	err := mockCourierServiceInterface.CreateDelivery(context.Background(), mockDeliveryInstance)
	require.Error(t, err)

	mockCourierServiceInterface.AssertCalled(t, "CreateDelivery", mock.Anything, mockDeliveryInstance)
}
