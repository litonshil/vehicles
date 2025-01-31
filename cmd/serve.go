package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"vehicles/infra/conn"
	rabbitmq "vehicles/infra/rabbitmq"
	"vehicles/internal/http/controllers"
	httpRoutes "vehicles/internal/http/routes"
	httpServer "vehicles/internal/http/server"
	"vehicles/internal/repositories/db"
	userservice "vehicles/internal/services/user"
	vehicleservice "vehicles/internal/services/vehicle"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: serve,
}

func serve(cmd *cobra.Command, args []string) {
	// base context
	baseContext := context.Background()

	mdbClient := conn.Db()
	rmq := rabbitmq.RMQ()
	redisClient := conn.NewCacheClient()
	dbRepo := db.NewRepository(mdbClient)
	userSvc := userservice.NewUserService(dbRepo, redisClient)
	vehicleSvc := vehicleservice.NewVehicleService(dbRepo, rmq)

	// HttpServer
	var HttpServer = httpServer.New()

	userController := controllers.NewUserController(
		baseContext,
		userSvc,
	)

	vehicleController := controllers.NewVehicleController(
		baseContext,
		vehicleSvc,
	)

	var Routes = httpRoutes.New(
		HttpServer.Echo,
		userController,
		vehicleController,
	)

	// Spooling
	Routes.Init()
	HttpServer.Start()
}
