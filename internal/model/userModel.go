package model

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Login        string    `json:"login"`
	Password     []byte    `json:"password"`
	Username     string    `json:"username"`
	Role         string    `json:"role"`
	RefreshToken []byte    `json:"refresh_token"`
}

type Login struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type HashedLogin struct {
	ID       uuid.UUID `json:"id"`
	Login    string    `json:"login"`
	Password []byte    `json:"password"`
	Role     string    `json:"role"`
}

type SignUp struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type SaveUser struct {
	Login    string `json:"login"`
	Password []byte `json:"password"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// Tokens struct for Access/Refresh tokens
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
