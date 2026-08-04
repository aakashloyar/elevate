package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/aakashloyar/elevate/user/config"
	httpuser "github.com/aakashloyar/elevate/user/internal/adapter/in/http"
	postgres "github.com/aakashloyar/elevate/user/internal/adapter/out/postgres"
	"github.com/aakashloyar/elevate/user/internal/application/ports/out/system"
	usersvc "github.com/aakashloyar/elevate/user/internal/application/service"
)

func main() {
	ctx := context.Background()
	_ = ctx

	port, err := strconv.Atoi(config.App.Postgres.Port)
	if err != nil {
		log.Fatalf("invalid POSTGRES_PORT: %v", err)
	}

	dbConfig := postgres.Config{
		Host:     config.App.Postgres.Host,
		Port:     port,
		User:     config.App.Postgres.User,
		Password: config.App.Postgres.Password,
		DBName:   config.App.Postgres.DBName,
		SSLMode:  config.App.Postgres.SSLMode,
	}

	db, err := dbConfig.NewDB()
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	defer db.Close()

	userRepo := postgres.NewUserRepository(db)
	clock := system.SystemClock{}
	idGen := system.UUIDGenerator{}

	createUserService := usersvc.NewCreateUserService(userRepo, idGen, clock)
	getUserService := usersvc.NewGetUserService(userRepo)
	deleteUserService := usersvc.NewDeleteUserService(userRepo)
	handler := httpuser.NewHandler(createUserService, getUserService, deleteUserService)

	mux := http.NewServeMux()
	httpuser.RegisterRoutes(mux, handler)

	serverPort := config.App.Server.Port
	log.Printf("user service starting on :%s", serverPort)
	if err := http.ListenAndServe(":"+serverPort, mux); err != nil {
		log.Fatal(err)
	}
}
