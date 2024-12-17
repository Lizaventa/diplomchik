package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/liza/labwork_45/internal/model"
	"github.com/stretchr/testify/require"
)

var rps *PsqlConnection

var (
	testProfile = &model.SaveUser{
		Login:    "test_login",
		Password: []byte("test_password"),
		Username: "test_username",
		Role:     "test_role",
	}

	test_refreshToken = []byte("test_refresh_token")
	// testLogin = &model.Login{
	// 	Login:    "test_login",
	// 	Password: "test_password",
	// }
	// testUpdateToken = &model.UpdateTokens{
	// 	RefreshToken: []byte("test_token"),
	// }
)

func TestCreateDeleteProfile(t *testing.T) {
	id, err := CreateTestProfile()
	require.NoError(t, err)
	defer func() {
		err = DeleteTestProfile(id)
		require.NoError(t, err)
	}()
}

func TestGetUserByLogin(t *testing.T) {
	id, err := CreateTestProfile()
	require.NoError(t, err)
	defer func() {
		err = DeleteTestProfile(id)
		require.NoError(t, err)
	}()

	selectedUser, err := rps.GetUserByLogin(context.Background(), testProfile.Login)
	require.NoError(t, err)
	require.Equal(t, testProfile.Login, selectedUser.Login)
	require.Equal(t, testProfile.Password, selectedUser.Password)
	require.Equal(t, testProfile.Role, selectedUser.Role)
}

func TestGetUserByWrongLogin(t *testing.T) {
	id, err := CreateTestProfile()
	require.NoError(t, err)
	defer func() {
		err = DeleteTestProfile(id)
		require.NoError(t, err)
	}()

	_, err = rps.GetUserByLogin(context.Background(), "Wrong_Login_template")
	require.Error(t, err)
}

func TestDeleteUserByID(t *testing.T) {
	id, err := CreateTestProfile()
	require.NoError(t, err)
	defer func() {
		err = DeleteTestProfile(id)
		require.NoError(t, err)
	}()
}

func TestDeleteUserByWrongID(t *testing.T) {
	id, err := CreateTestProfile()
	require.NoError(t, err)
	err = DeleteTestProfile(uuid.New())
	require.NoError(t, err)
	defer func() {
		err = DeleteTestProfile(id)
		require.NoError(t, err)
	}()
}

func TestGetAllUsers(t *testing.T) {
	id, err := CreateTestProfile()
	require.NoError(t, err)
	defer func() {
		err = DeleteTestProfile(id)
		require.NoError(t, err)
	}()

	users, err := rps.GetAll(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, users)
}

func TestGetRefreshToken(t *testing.T) {
	id, err := CreateTestProfile()
	require.NoError(t, err)
	defer func() {
		err = DeleteTestProfile(id)
		require.NoError(t, err)
	}()

	err = rps.SaveRefreshToken(context.Background(), id, test_refreshToken)
	require.NoError(t, err)

	refreshToken, err := rps.GetRefreshTokenByID(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, refreshToken)
	require.Equal(t, test_refreshToken, refreshToken)
}

func TestSaveRefreshToken(t *testing.T) {
	id, err := CreateTestProfile()
	require.NoError(t, err)
	defer func() {
		err = DeleteTestProfile(id)
		require.NoError(t, err)
	}()

	err = rps.SaveRefreshToken(context.Background(), id, test_refreshToken)
	require.NoError(t, err)

	refreshToken, err := rps.GetRefreshTokenByID(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, refreshToken)
	require.Equal(t, test_refreshToken, refreshToken)
}

func TestGetUserByID(t *testing.T) {
	id, err := CreateTestProfile()
	require.NoError(t, err)
	defer func() {
		err = DeleteTestProfile(id)
		require.NoError(t, err)
	}()

	selectedUser, err := rps.GetUserByID(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, testProfile.Login, selectedUser.Login)
	require.Equal(t, testProfile.Password, selectedUser.Password)
	require.Equal(t, testProfile.Role, selectedUser.Role)
}
