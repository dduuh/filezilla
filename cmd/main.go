package main

import (
	"context"
	"filezilla/internal/configs"
	"filezilla/internal/handlers"
	"filezilla/internal/repository"
	"filezilla/internal/server"
	"filezilla/internal/service"
	"filezilla/pkg/mong"
	"filezilla/pkg/psql"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const configDir = "configs"

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	cfg, err := configs.Init(configDir)
	if err != nil {
		logrus.Fatal(err)
	}

	srv := &server.Server{}

	db, err := psql.NewPostgreSQL()
	if err != nil {
		logrus.Fatal(err)
	}

	mongoDb, err := mong.NewMongoDB(cfg, context.Background())
	if err != nil {
		logrus.Fatal(err)
	}

	repo := repository.NewRepository(db, mongoDb)
	service := service.NewService(*repo)
	handler := handlers.NewHandler(*service)

	logrus.Printf("Server is listening on port %s", cfg.HttpCfg.Port)

	if err := srv.Run(cfg.HttpCfg.Port, handler.InitHandler()); err != nil {
		logrus.Fatal(err.Error())
	}

}
