package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"vehicles/config"
	"vehicles/infra/conn"
	logger "vehicles/infra/logger"
)

var (
	RootCmd = &cobra.Command{
		Use: "vehicles",
	}
)

func init() {
	RootCmd.AddCommand(serveCmd)
}

// Execute executes the root command
func Execute() {
	// load application configuration
	if err := config.Load(); err != nil {
		//log.Error(err)
		os.Exit(1)
	}

	logger.InitLogger()
	conn.ConnectCache()
	conn.ConnectDb()
	conn.InitRabbitMQ()

	defer conn.CloseRabbitMQ()

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
