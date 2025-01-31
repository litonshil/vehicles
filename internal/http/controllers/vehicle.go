package controllers

import (
	"context"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"vehicles/internal/domain"
	"vehicles/types"
)

type VehicleController struct {
	baseCtx   context.Context
	vehicleuc domain.VehicleUseCase
}

func NewVehicleController(
	baseCtx context.Context,
	vehicleuc domain.VehicleUseCase,
) *VehicleController {
	return &VehicleController{
		baseCtx:   baseCtx,
		vehicleuc: vehicleuc,
	}
}

func (cntlr *VehicleController) CreateVehicle(c echo.Context) error {
	var vehicle types.VehicleReq

	if err := c.Bind(&vehicle); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	err := cntlr.vehicleuc.CreateVehicle(vehicle)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create vehicle"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Vehicle created successfully"})
}

func (cntlr *VehicleController) CreateBrands(c echo.Context) error {
	var brand types.VehicleBrand

	if err := c.Bind(&brand); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	err := cntlr.vehicleuc.CreateBrand(brand)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create brand"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Vehicle Brand created successfully"})
}

func (cntlr *VehicleController) ReadBrands(c echo.Context) error {
	brands, err := cntlr.vehicleuc.ReadBrands()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read brands"})
	}

	return c.JSON(http.StatusOK, brands)
}

func (cntlr *VehicleController) CreateVehicleModel(c echo.Context) error {
	var vehicleModel types.VehicleModel

	if err := c.Bind(&vehicleModel); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	err := cntlr.vehicleuc.CreateModel(vehicleModel)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create vehicle model"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Vehicle model created successfully"})
}

func (cntlr *VehicleController) ReadVehicleModels(c echo.Context) error {
	models, err := cntlr.vehicleuc.ReadModels()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read models"})
	}

	return c.JSON(http.StatusOK, models)
}

func (cntlr *VehicleController) UpdateVehicleStatus(c echo.Context) error {

	id := c.Param("id")
	status := c.QueryParam("status")

	vehicleID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid vehicle ID"})
	}

	req := types.UpdateStatusReq{
		ID:     vehicleID,
		Status: status,
	}

	err = cntlr.vehicleuc.UpdateVehicleStatus(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update vehicle status"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Vehicle status updated successfully"})
}
