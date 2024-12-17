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

type CourierHandler struct {
	srv CourierServiceInterface
}

func NewCourierHandler(srv CourierServiceInterface) *CourierHandler {
	return &CourierHandler{srv: srv}
}

type CourierServiceInterface interface {
	UpdateCourier(context.Context, uuid.UUID, *model.Courier) error
	CreateDelivery(context.Context, *model.Delivery) error
	GetAllDeliveries(context.Context) ([]*model.DeliveryGet, error)
	AssignCourierToDelivery(context.Context, uuid.UUID, uuid.UUID) error
	UpdateDeliveryStatus(context.Context, *model.DeliveryStatus) error
}

// UpdateCourier updates courier info for the given Courier instance
// @Summary UpdateCourier
// @Description Updates Courier information
// @Tags Courier Bussiness logic
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body model.Courier true "Courier to update"
// @Success 200 {object} map[string]interface{} "Courier info has been sucessfully updated"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /courier/updatecourier [patch]
func (h *CourierHandler) UpdateCourier(c echo.Context) error {
	userId, err := middleware.GetPayloadFromToken(strings.Split(c.Request().Header.Get("Authorization"), " ")[1])
	if err != nil {
		logrus.WithFields(logrus.Fields{"id": userId}).Errorf("GetPayloadFromToken: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("GetPayloadFromToken: %v", err))
	}
	courier := &model.Courier{}
	err = c.Bind(courier)
	if err != nil {
		logrus.WithFields(logrus.Fields{"courier": courier}).Errorf("Bind: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind: %v", err))
	}
	err = h.srv.UpdateCourier(c.Request().Context(), userId, courier)
	if err != nil {
		logrus.WithFields(logrus.Fields{"userId": userId}).Errorf("CreateCourier: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("CreateCourier: %v", err))
	}
	response := map[string]interface{}{
		"message": "courier info updated!",
	}
	return c.JSON(http.StatusCreated, response)
}

// CreateDelivery creates a new delivery
// @Summary CreateDelivery
// @Description Creates a new delivery instance
// @Tags Courier Bussiness logic
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body model.Delivery true "Delivery to create"
// @Success 200 {object} map[string]interface{} "Delivery has been sucessfully created"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /delivery/create_delivary [post]
func (h *CourierHandler) CreateDelivery(c echo.Context) error {
	delivery := &model.Delivery{}
	err := c.Bind(delivery)
	if err != nil {
		logrus.WithFields(logrus.Fields{"delivery": delivery}).Errorf("Bind: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind: %v", err))
	}
	err = h.srv.CreateDelivery(c.Request().Context(), delivery)
	if err != nil {
		logrus.WithFields(logrus.Fields{"delivery": delivery}).Errorf("CreateDelivery: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("CreateDelivery: %v", err))
	}
	response := map[string]interface{}{
		"message": "delivery info updated!",
	}
	return c.JSON(http.StatusCreated, response)
}

// CreateDelivery creates a new delivery
// @Summary GetAlldeliveries
// @Description GetAlldeliveries
// @Tags Courier Bussiness logic
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} model.DeliveryGet "Delivery has been sucessfully created"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /courier/getalldeliveries [get]
func (h *CourierHandler) GetAlldeliveries(c echo.Context) error {
	deliveries, err := h.srv.GetAllDeliveries(c.Request().Context())
	if err != nil {
		logrus.WithFields(logrus.Fields{"id": c.Param("userId")}).Errorf("GetAllDeliveries: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("GetAllDeliveries: %v", err))
	}
	return c.JSON(http.StatusOK, deliveries)
}

// CreateDelivery creates a new delivery
// @Summary ChooseAvailibleDelivery
// @Description allows the courier to choose delivery
// @Tags Courier Bussiness logic
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body model.DeliveryId true "Delivery to create"
// @Success 200 {string} string "Delivery has been sucessfully choosed by courier"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /courier/choose_availible_delivery [patch]
func (h *CourierHandler) ChooseAvailibleDelivery(c echo.Context) error {
	userId, err := middleware.GetPayloadFromToken(strings.Split(c.Request().Header.Get("Authorization"), " ")[1])
	if err != nil {
		logrus.WithFields(logrus.Fields{"id": userId}).Errorf("GetPayloadFromToken: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("GetPayloadFromToken: %v", err))
	}
	Id := &model.DeliveryId{}
	err = c.Bind(Id)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Id": Id}).Errorf("Bind: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind: %v", err))
	}

	err = h.srv.AssignCourierToDelivery(c.Request().Context(), Id.Id, userId)
	if err != nil {
		logrus.WithFields(logrus.Fields{"userId": userId}).Errorf("AssignCourierToDelivery: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("AssignCourierToDelivery: %v", err))
	}
	return c.JSON(http.StatusOK, "OK, let's go!")
}

// CreateDelivery creates a new delivery
// @Summary UpdateDeliveryStatus
// @Description allows the courier to update delivery status
// @Tags Courier Bussiness logic
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body model.DeliveryStatus true "Delivery status to update"
// @Success 200 {string} string "Delivery status has been sucessfully updated by courier"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /courier/update_delivery_status [patch]
func (h *CourierHandler) UpdateDeliveryStatus(c echo.Context) error {
	delivery := &model.DeliveryStatus{}
	err := c.Bind(delivery)
	if err != nil {
		logrus.WithFields(logrus.Fields{"delivery": delivery}).Errorf("Bind: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Bind: %v", err))
	}
	err = h.srv.UpdateDeliveryStatus(c.Request().Context(), delivery)
	if err != nil {
		logrus.WithFields(logrus.Fields{"deliveryId": delivery}).Errorf("UpdateDeliveryStatus: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("UpdateDeliveryStatus: %v", err))
	}
	return c.JSON(http.StatusOK, "Status has been changed successfully")

}
