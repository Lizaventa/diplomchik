package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/liza/labwork_45/internal/model"
)

// type PsqlConnection struct {
// 	pool *pgxpool.Pool
// }

// func NewPsqlConnection(pool *pgxpool.Pool) *PsqlConnection {
// 	return &PsqlConnection{pool: pool}
// }

func (db *PsqlConnection) InsertUser(ctx context.Context, user *model.SaveUser) (uuid.UUID, error) {
	id := uuid.New()
	insertQuery := "INSERT INTO labwork.user (id, login, password, username, role) VALUES ($1, $2, $3, $4, $5)"
	query, err := db.pool.Exec(ctx, insertQuery, id, user.Login, user.Password, user.Username, user.Role)
	if err != nil && !query.Insert() {
		return uuid.Nil, fmt.Errorf("Exec(): %w", err)
	}

	if user.Role == "Courier" {
		createCourier := "INSERT INTO labwork.courier (id, userid) VALUES ($1, $2)"
		courierQuery, err := db.pool.Exec(ctx, createCourier, uuid.New(), id)
		if err != nil && !courierQuery.Insert() {
			return uuid.Nil, fmt.Errorf("Exec(): %w", err)
		}
	}
	return id, nil
}

func (db *PsqlConnection) GetUserByLogin(ctx context.Context, login string) (*model.HashedLogin, error) {
	selectedUser := &model.HashedLogin{}
	queryString := "SELECT id, login, password, role FROM labwork.user WHERE login=$1"
	err := db.pool.QueryRow(ctx, queryString, login).Scan(&selectedUser.ID, &selectedUser.Login, &selectedUser.Password, &selectedUser.Role)
	if err != nil {
		return nil, fmt.Errorf("Exec(): %w", err)
	}

	return selectedUser, nil
}

func (db *PsqlConnection) GetAll(ctx context.Context) ([]*model.User, error) {
	rows, err := db.pool.Query(ctx, "SELECT id, login, password, username, refresh_token, role FROM labwork.user")
	if err != nil {
		return nil, fmt.Errorf("Query(): %w", err)
	}
	defer rows.Close()

	// Create slice to store data from our SQL request
	var results []*model.User

	// go;) through each line
	for rows.Next() {
		person := &model.User{}
		err := rows.Scan(&person.ID, &person.Login, &person.Password, &person.Username, &person.RefreshToken, &person.Role)
		if err != nil {
			return nil, fmt.Errorf("Scan(): %w", err) // Returning error message
		}
		results = append(results, person)
	}
	return results, rows.Err()
}

func (db *PsqlConnection) SaveRefreshToken(ctx context.Context, ID uuid.UUID, refreshToken []byte) error {
	//err := db.pool.QueryRow(ctx, "SELECT id FROM labwork.user WHERE id=$1", ID).Scan(id)
	_, err := db.pool.Exec(ctx, "UPDATE labwork.user SET refresh_token=$1 WHERE id=$2", refreshToken, ID)
	if err != nil {
		return fmt.Errorf("Exec(): %w", err)
	}
	return nil
}

func (db *PsqlConnection) GetRefreshTokenByID(ctx context.Context, ID uuid.UUID) ([]byte, error) {
	refreshToken := make([]byte, 0)
	err := db.pool.QueryRow(ctx, "SELECT refresh_token FROM labwork.user WHERE id=$1", ID).Scan(&refreshToken)
	if err != nil {
		return nil, fmt.Errorf("Exec(): %w", err)
	}
	return refreshToken, nil
}

func (db *PsqlConnection) GetUserByID(ctx context.Context, ID uuid.UUID) (*model.User, error) {
	userInfo := &model.User{}
	queryString := "SELECT id, login, password, username, refresh_token, role FROM labwork.user WHERE id=$1"
	err := db.pool.QueryRow(ctx, queryString, ID).Scan(&userInfo.ID, &userInfo.Login, &userInfo.Password, &userInfo.Username, &userInfo.RefreshToken, &userInfo.Role)
	if err != nil {
		return nil, fmt.Errorf("Exec(): %w", err)
	}
	return userInfo, nil
}

func (db *PsqlConnection) DeleteUserByID(ctx context.Context, ID uuid.UUID) error {
	_, err := db.pool.Exec(ctx, "DELETE FROM labwork.user WHERE id=$1", ID)
	if err != nil {
		return fmt.Errorf("Exec(): %w", err)
	}
	return nil
}
