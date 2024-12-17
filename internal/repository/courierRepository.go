package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/liza/labwork_45/internal/model"
)

func (db *PsqlConnection) GetCourierByUserID(ctx context.Context, userId uuid.UUID) (*model.Courier, error) {
	courier := &model.Courier{}
	query := "SELECT id, userid, name, surname, status, performance_indicator FROM labwork.courier WHERE userid=$1"
	err := db.pool.QueryRow(ctx, query, userId).Scan(&courier.Id, &courier.UserId, &courier.Name, &courier.Surname, &courier.Status, &courier.Perfomance_indicator)
	if err != nil {
		return nil, fmt.Errorf("QueryRow(): %w", err)
	}
	return courier, nil
}

func (db *PsqlConnection) UpdateCourierInfo(ctx context.Context, userId uuid.UUID, courier *model.Courier) error {
	query := "UPDATE labwork.courier SET name=$1, surname=$2, status=$3, performance_indicator=$4 WHERE userid=$5"
	update, err := db.pool.Exec(ctx, query, courier.Name, courier.Surname, courier.Status, courier.Perfomance_indicator, userId)
	if err != nil && update.RowsAffected() == 0 {
		return fmt.Errorf("Exec(): %w", err)
	}
	return nil
}

func (db *PsqlConnection) UpdateStatus(ctx context.Context, delivery *model.DeliveryStatus) error {
	_, err := db.pool.Exec(ctx, "UPDATE labwork.delivery SET delivery_status=$1 WHERE id=$2", delivery.DeliveryStatus, delivery.Id)
	if err != nil {
		return err
	}
	return nil
}

func (db *PsqlConnection) UpdateDeliveryCourier(ctx context.Context, deliveryId uuid.UUID, courieerId uuid.UUID) error {
	_, err := db.pool.Exec(ctx, "UPDATE labwork.delivery SET courier_id=$1 WHERE id=$2", courieerId, deliveryId)
	if err != nil {
		return err
	}
	return nil
}

func (db *PsqlConnection) InsertDelivery(ctx context.Context, delivery *model.Delivery) error {
	id := uuid.New()
	insert := "INSERT INTO labwork.delivery (id, delivery_date, delivery_status, delivery_comment) VALUES ($1, $2, $3, $4)"
	query, err := db.pool.Exec(ctx, insert, id, delivery.DeliveryDate, delivery.DeliveryStatus, delivery.DeliveryComment)
	if err != nil && !query.Insert() {
		return fmt.Errorf("Exec(): %w", err)
	}
	return nil
}

func (db *PsqlConnection) GetAllDeliveries(ctx context.Context) ([]*model.DeliveryGet, error) {
	rows, err := db.pool.Query(ctx, "SELECT id, delivery_date, delivery_status, delivery_comment FROM labwork.delivery")
	if err != nil {
		return nil, fmt.Errorf("Query(): %w", err)
	}
	defer rows.Close()

	var result []*model.DeliveryGet

	for rows.Next() {
		delivery := &model.DeliveryGet{}
		err := rows.Scan(&delivery.Id, &delivery.DeliveryDate, &delivery.DeliveryStatus, &delivery.DeliveryComment)
		if err != nil {
			return nil, fmt.Errorf("Scan(): %w", err)
		}
		result = append(result, delivery)
	}
	return result, nil
}
