package config

import (
	"context"
	"flag"
	"os"
	"strconv"
	"time"

	"log/slog"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type configuration struct {
	version        string
	restApiPort    string
	grpcServerAddr string
	env            string
	redisConfig    struct {
		address  string
		password string
		db       int
	}
}

func initializeGrpcConn(cfg configuration) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(cfg.grpcServerAddr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func initializeRedisClient(cfg configuration) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.redisConfig.address,
		Password: cfg.redisConfig.password, // no password set
		DB:       cfg.redisConfig.db,       // use default DB
	})

	// Check if Redis is accessible
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return redisClient, nil
}

func InitConfig() (*redis.Client, *grpc.ClientConn, string) {
	config := configuration{}

	err := godotenv.Load()
	if err != nil {
		slog.Error("failed loading environment variables", err)
	}
	slog.Info("environment variables loaded")

	//?APP
	flag.StringVar(&config.restApiPort, "restApiPort", os.Getenv("REST_API_PORT"), "REST API server port")
	flag.StringVar(&config.version, "version", os.Getenv("APP_VERSION"), "application version")
	flag.StringVar(&config.env, "env", os.Getenv("APP_ENVIRONMENT"), "application environment")

	//?REDIS
	flag.StringVar(&config.redisConfig.address, "redisAddress", os.Getenv("REDIS_ADDRESS"), "Redis address")
	flag.StringVar(&config.redisConfig.password, "redisPassword", os.Getenv("REDIS_PASSWORD"), "Redis password")
	redisDbStr := os.Getenv("REDIS_DB")
	redisDb, err := strconv.Atoi(redisDbStr)
	if err != nil {
		slog.Error("failed convering REDIS_DB envvar to INT", err)
	}
	flag.IntVar(&config.redisConfig.db, "redisDb", redisDb, "Redis DB")

	//?GRPC
	flag.StringVar(&config.grpcServerAddr, "grpcServerAddr", os.Getenv("GRPC_SERVER_ADDR"), "GRPC server port")

	flag.Parse()
	slog.Info("command line variables loaded")

	redisClient, err := initializeRedisClient(config)
	if err != nil {
		slog.Error("failed initializing Redis client", err)
	}
	slog.Info("Redis client initialized")

	grpcConn, err := initializeGrpcConn(config)
	if err != nil {
		slog.Error("failed initializing gRPC connection", err)
	}
	slog.Info("gRPC connection established")

	return redisClient, grpcConn, config.restApiPort
}
