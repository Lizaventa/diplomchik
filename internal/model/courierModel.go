package model

import "github.com/google/uuid"

type Courier struct {
	Id                   uuid.UUID `json:"id"`
	UserId               uuid.UUID `json:"userid"`
	Name                 string    `json:"name"`
	Surname              string    `json:"surname"`
	Status               string    `json:"status"`
	Perfomance_indicator int       `json:"perfomance_indicator"`
}
