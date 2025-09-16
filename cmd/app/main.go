package main

import (
	"fmt"

	"tesla-app/internal/app/config"
	"tesla-app/internal/app/dsn"
	"tesla-app/internal/app/handler"
	"tesla-app/internal/app/repository"
	"tesla-app/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	router := gin.Default()
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println("DSN:", postgresString)

	repo, err := repository.New(postgresString)
	if err != nil {
		logrus.Fatalf("error initializing repository: %v", err)
	}

	hand := handler.NewHandler(repo)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}
