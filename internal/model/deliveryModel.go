package model

import (
	"github.com/google/uuid"
)

type Delivery struct {
	Id              uuid.UUID `json:"id"`
	CourierId       uuid.UUID `json:"courier_id"`
	DeliveryDate    string    `json:"delivery_date"`
	DeliveryStatus  string    `json:"delivery_status"`
	DeliveryComment string    `json:"delivery_comment"`
}
type DeliveryGet struct {
	Id              uuid.UUID `json:"id"`
	DeliveryDate    string    `json:"delivery_date"`
	DeliveryStatus  string    `json:"delivery_status"`
	DeliveryComment string    `json:"delivery_comment"`
}

type DeliveryStatus struct {
	Id             uuid.UUID `json:"id"`
	DeliveryStatus string    `json:"delivery_status"`
}

type DeliveryId struct {
	Id uuid.UUID `json:"id"`
}
