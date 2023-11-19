package main

import (
	"log/slog"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"

	"github.com/niewolinsky/tw_employee_data_processor/config"

	pb "github.com/niewolinsky/tw_employee_data_service/grpc/employee"
)

type application struct {
	redisClient        *redis.Client
	grpcEmployeeClient pb.EmployeeServiceClient
	validator          *validator.Validate
	wait_group         sync.WaitGroup
}

func main() {
	redisClient, grpcConn, restApiPort := config.InitConfig()
	defer redisClient.Close()

	app := &application{
		redisClient:        redisClient,
		grpcEmployeeClient: pb.NewEmployeeServiceClient(grpcConn),
		validator:          validator.New(),
	}

	err := app.serveREST(restApiPort)
	if err != nil {
		slog.Error("failed starting HTTP server", err)
	}

	slog.Info("stopped server")
}
