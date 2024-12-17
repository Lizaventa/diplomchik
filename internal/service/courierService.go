package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/liza/labwork_45/internal/model"
)

type CourierService struct {
	rps CourierRepository
}

func NewCourierService(rps CourierRepository) *CourierService {
	return &CourierService{rps: rps}
}

type CourierRepository interface {
	UpdateCourierInfo(context.Context, uuid.UUID, *model.Courier) error
	InsertDelivery(context.Context, *model.Delivery) error
	GetAllDeliveries(context.Context) ([]*model.DeliveryGet, error)
	UpdateDeliveryCourier(ctx context.Context, deliveryId uuid.UUID, courieerId uuid.UUID) error
	GetCourierByUserID(context.Context, uuid.UUID) (*model.Courier, error)
	UpdateStatus(context.Context, *model.DeliveryStatus) error
}

func (srv *CourierService) UpdateCourier(ctx context.Context, userId uuid.UUID, courier *model.Courier) error {
	err := srv.rps.UpdateCourierInfo(ctx, userId, courier)
	if err != nil {
		return fmt.Errorf("UpdateStatus: %w", err)
	}
	return nil
}

func (srv *CourierService) CreateDelivery(ctx context.Context, delivery *model.Delivery) error {
	err := srv.rps.InsertDelivery(ctx, delivery)
	if err != nil {
		return fmt.Errorf("InsertDelivery: %w", err)
	}
	return nil
}

func (srv *CourierService) GetAllDeliveries(ctx context.Context) ([]*model.DeliveryGet, error) {
	deliveries, err := srv.rps.GetAllDeliveries(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetAllDeliveries: %w", err)
	}
	return deliveries, nil
}

func (srv *CourierService) AssignCourierToDelivery(ctx context.Context, deliveryId uuid.UUID, userId uuid.UUID) error {
	courier, err := srv.rps.GetCourierByUserID(ctx, userId)
	if err != nil {
		return fmt.Errorf("GetCourierByUserID: %w", err)
	}

	err = srv.rps.UpdateDeliveryCourier(ctx, deliveryId, courier.Id)
	if err != nil {
		return fmt.Errorf("UpdateDeliveryCourier: %w", err)
	}
	return nil
}

func (srv *CourierService) UpdateDeliveryStatus(ctx context.Context, delivery *model.DeliveryStatus) error {
	err := srv.rps.UpdateStatus(ctx, delivery)
	if err != nil {
		return fmt.Errorf("UpdateStatus: %w", err)
	}
	return nil
}
