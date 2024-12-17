package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/liza/labwork_45/internal/middleware"
	"github.com/liza/labwork_45/internal/model"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=authApiHandler.go -destination=mock/authApiHandler.go

type authApiHandler struct {
	srv AuthApiSerivce
}

func NewAuthApiHandler(srv AuthApiSerivce) *authApiHandler {
	return &authApiHandler{srv: srv}
}

type AuthApiSerivce interface {
	GetAll(ctx context.Context) ([]*model.User, error)
	LoginUser(ctx context.Context, login *model.Login) (string, string, error)
	SignUpUser(ctx context.Context, user *model.SignUp) (uuid.UUID, error)
	GetPersonalInfo(ctx context.Context, ID uuid.UUID) (*model.User, error)
	DeleteUserByID(ctx context.Context, ID uuid.UUID) error
}

// GetAll function returns list of all users in database(test function)
// @Summary GetAll
// @tags temporary methods
// @Description GetAll function returns list of all users in database(test function)
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} model.User "All users in system"
// @Failure 404 {string} string "Error message"
// @Router /auth/getall [get]
func (handler *authApiHandler) GetAll(c echo.Context) error {
	req, err := handler.srv.GetAll(c.Request().Context())
	if err != nil {
		logrus.WithFields(logrus.Fields{"request": req}).Errorf("GetAll: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("GetAll: %v", err))
	}
	return c.JSON(http.StatusOK, req)
}

// Login function handles the login request and returns user's access and refresh tokens
// @Summary Login
// @tags Authentication methods
// @Description Logs in a user and returns access and refresh tokens
// @Accept json
// @Produce json
// @Param input body model.Login true "Login details"
// @Success 200 {object} map[string]interface{} " Generating access and refresh tokens"
// @Failure 404 {string} string "Error message"
// @Router /auth/login [post]
func (handler *authApiHandler) Login(c echo.Context) error {
	login := &model.Login{}
	err := c.Bind(login)
	if err != nil {
		logrus.WithFields(logrus.Fields{"login": login}).Errorf("Bind: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind: %v", err))
	}
	accessToken, refreshToken, err := handler.srv.LoginUser(c.Request().Context(), login)
	if err != nil {
		logrus.WithFields(logrus.Fields{"login": login}).Errorf("LoginUser: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("LoginUser: %v", err))
	}
	response := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	return c.JSON(http.StatusOK, response)
}

// GetPersonalInfo function receives GET request from client
// @Summary Get Personal Info
// @Description Fetch personal information of the active user based on the access token provided in the Authorization header.
// @Tags User methods
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} model.User "User's personal information"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /auth/getpersonalinfo [get]
func (handler *authApiHandler) GetPersonalInfo(c echo.Context) error {
	id, err := middleware.GetPayloadFromToken(strings.Split(c.Request().Header.Get("Authorization"), " ")[1])
	if err != nil {
		logrus.WithFields(logrus.Fields{"id": id}).Errorf("GetPayloadFromToken: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("GetPayloadFromToken: %v", err))
	}
	userInfo, err := handler.srv.GetPersonalInfo(c.Request().Context(), id)
	if err != nil {
		logrus.WithFields(logrus.Fields{"id": id}).Errorf("GetPersonalInfo: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("GetPersonalInfo: %v", err))
	}
	return c.JSON(http.StatusOK, userInfo)
}

// SignUp function receives POST reauest from client to register user in system
// @Summary SignUp
// @Description Creates a new user in the system
// @Tags Authentication methods
// @Accept json
// @Produce json
// @Param input body model.SignUp true "Sign up details"
// @Success 200 {object} map[string]interface{} "User has been registered successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /auth/signup [post]
func (handler *authApiHandler) SignUp(c echo.Context) error {
	user := &model.SignUp{}
	err := c.Bind(user)
	if err != nil {
		logrus.WithFields(logrus.Fields{"user": user}).Errorf("Bind: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind: %v", err))
	}
	id, err := handler.srv.SignUpUser(c.Request().Context(), user)
	if err != nil {
		logrus.WithFields(logrus.Fields{"user": user}).Errorf("SignUpUser: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("SignUpUser: %v", err))
	}
	response := map[string]interface{}{
		"message": "user created!",
		"id":      id,
	}

	return c.JSON(http.StatusOK, response)
}

// DeleteUser deletes a user from the database
// @Summary Delete
// @tags User methods
// @Description Delete a user from the database
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {string} string "All users in system"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /auth/delete [delete]
func (handler *authApiHandler) DeleteUser(c echo.Context) error {
	id, err := middleware.GetPayloadFromToken(strings.Split(c.Request().Header.Get("Authorization"), " ")[1])
	if err != nil {
		logrus.WithFields(logrus.Fields{"id": id}).Errorf("GetPayloadFromToken: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("GetPayloadFromToken: %v", err))
	}
	err = handler.srv.DeleteUserByID(c.Request().Context(), id)
	if err != nil {
		logrus.WithFields(logrus.Fields{"id": id}).Errorf("DeleteUserByID: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("DeleteUserByID: %v", err))
	}
	return c.JSON(http.StatusOK, "User has been deleted from the system")
}

func (handler *authApiHandler) RefreshTokenPair(c echo.Context) error {
	return c.JSON(http.StatusOK, "gotovo")
}
