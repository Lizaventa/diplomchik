package handlers

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/liza/labwork_45/internal/handlers/mocks"
	"github.com/liza/labwork_45/internal/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	mockAuthApiService *mocks.AuthApiSerivce

	mockUserEntity = &model.User{
		ID:           uuid.New(),
		Login:        "test_login",
		Password:     []byte("test_password"),
		Username:     "test_username",
		Role:         "test_role",
		RefreshToken: []byte("test_refresh_token"),
	}

	testTokens = &model.Tokens{
		AccessToken:  "test_access_token",
		RefreshToken: "test_refresh_token",
	}
	testLogin = &model.Login{
		Login:    "test_login",
		Password: "test_password",
	}
	TestSignUpEntity = &model.SignUp{
		Login:    "test_login",
		Password: "test_password",
		Username: "test_username",
	}
	// testUpdateToken = &model.UpdateTokens{
	//     RefreshToken: []byte("test_token"),
	// }
	// testProfile = &model.Profile{
	//     ID:        uuid.New(),
	//     Login:    "test_login",
)

// TestMain function starts the tests
func TestMain(m *testing.M) {
	mockAuthApiService = new(mocks.AuthApiSerivce)
	mockCourierServiceInterface = new(mocks.CourierServiceInterface)
	exitVal := m.Run()
	os.Exit(exitVal)
}

// TestGetAll tests GetAll function mocking GetAll function from Service Interface
func TestGetAll(t *testing.T) {
	mockAuthApiService.On("GetAll", mock.Anything).Return([]*model.User{}, nil).Twice()
	result, err := mockAuthApiService.GetAll(context.Background())
	require.NoError(t, err)
	require.NotNil(t, result)
}

// TestGetPersonalInfo tests GetPersonalInfo function mocking GetPersonalInfo function from Service Interface
func TestGetPersonalInfo(t *testing.T) {
	mockAuthApiService.On("GetPersonalInfo", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(&model.User{}, nil).Twice()
	result, err := mockAuthApiService.GetPersonalInfo(context.Background(), mockUserEntity.ID)
	require.NoError(t, err)
	require.NotNil(t, result)
}

// TestGetPersonalInfoError tests GetPersonalInfo function mocking GetPersonalInfo function from Service Interface
func TestGetPersonalInfoError(t *testing.T) {
	mockAuthApiService.On("GetPersonalInfo", mock.Anything, mockUserEntity.ID).Return(nil, errors.New("get user error"))
	result, err := mockAuthApiService.GetPersonalInfo(context.Background(), mockUserEntity.ID)
	require.Error(t, err)
	require.Nil(t, result)

	mockAuthApiService.AssertCalled(t, "GetPersonalInfo", mock.Anything, mockUserEntity.ID)
}

// TestLogin tests Login function mocking LoginUser function from Service Interface
func TestLogin(t *testing.T) {
	mockAuthApiService.On("LoginUser", mock.Anything, mock.Anything).Return(testTokens.AccessToken, testTokens.RefreshToken, nil)
	accessToken, refreshToken, err := mockAuthApiService.LoginUser(context.Background(), testLogin)
	require.NoError(t, err)
	require.NotNil(t, accessToken)
	require.NotNil(t, refreshToken)
}

// TestLoginError tests Login function mocking LoginUser function from Service Interface
func TestLoginError(t *testing.T) {
	mockAuthApiService.On("LoginUser", mock.Anything, mock.Anything).Return("", "", errors.New("login error"))
	accessToken, refreshToken, err := mockAuthApiService.LoginUser(context.Background(), testLogin)
	require.Error(t, err)
	require.Empty(t, accessToken)
	require.Empty(t, refreshToken)

	mockAuthApiService.AssertCalled(t, "LoginUser", mock.Anything, testLogin)
}

// TestSignUp tests SignUp function mocking SugnUpUser function from Service Interface
func TestSignUp(t *testing.T) {
	mockAuthApiService.On("SignUpUser", mock.Anything, mock.Anything).Return(mockUserEntity.ID, nil)
	id, err := mockAuthApiService.SignUpUser(context.Background(), TestSignUpEntity)
	require.NoError(t, err)
	require.NotNil(t, id)

	mockAuthApiService.AssertCalled(t, "SignUpUser", mock.Anything, TestSignUpEntity)
}

// TestSignUp tests SignUp function mocking SugnUpUser function from Service Interface
func TestSignUpError(t *testing.T) {
	mockAuthApiService.On("SignUpUser", mock.Anything, mock.Anything).Return(uuid.Nil, errors.New("signup error"))
	id, err := mockAuthApiService.SignUpUser(context.Background(), TestSignUpEntity)
	require.Error(t, err)
	require.Empty(t, id)

	mockAuthApiService.AssertCalled(t, "SignUpUser", mock.Anything, TestSignUpEntity)
}

// TestUpdate tests Update function mocking UpdateUser function from Service Interface
func TestDelete(t *testing.T) {
	mockAuthApiService.On("DeleteUserByID", mock.Anything, mockUserEntity.ID).Return(nil)
	err := mockAuthApiService.DeleteUserByID(context.Background(), mockUserEntity.ID)
	require.NoError(t, err)

	mockAuthApiService.AssertCalled(t, "DeleteUserByID", mock.Anything, mockUserEntity.ID)
}

// TestUpdate tests Update function mocking UpdateUser function from Service Interface
func TestDeleteError(t *testing.T) {
	mockAuthApiService.On("DeleteUserByID", mock.Anything, mockUserEntity.ID).Return(errors.New("delete error"))
	err := mockAuthApiService.DeleteUserByID(context.Background(), mockUserEntity.ID)
	require.Error(t, err)

	mockAuthApiService.AssertCalled(t, "DeleteUserByID", mock.Anything, mockUserEntity.ID)
}
