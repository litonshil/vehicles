package routes

import (
	"github.com/labstack/echo/v4"
	"vehicles/internal/http/controllers"
	"vehicles/internal/http/middlewares"
)

type Routes struct {
	echo              *echo.Echo
	userController    *controllers.UserController
	vehicleController *controllers.VehicleController
}

func New(
	e *echo.Echo,
	userController *controllers.UserController,
	vehicleController *controllers.VehicleController,
) *Routes {
	return &Routes{
		echo:              e,
		userController:    userController,
		vehicleController: vehicleController,
	}
}

func (r *Routes) Init() {
	e := r.echo
	middlewares.Init(e)

	g := e.Group("/v1")
	g.POST("/users", r.userController.CreateUser)
	g.GET("/users", r.userController.GetUsers)

	// vehicle makes
	g.POST("/vehicle-brands", r.vehicleController.CreateBrands)
	g.GET("/vehicle-brands", r.vehicleController.ReadBrands)

	// vehicle models
	g.POST("/vehicle-models", r.vehicleController.CreateVehicleModel)
	g.GET("/vehicle-models", r.vehicleController.ReadVehicleModels)

	// vehicle
	g.POST("/vehicles", r.vehicleController.CreateVehicle)
	g.PATCH("/vehicles/status/:id", r.vehicleController.UpdateVehicleStatus)

}
