package main

import (
	"context"
	"fmt"

	"tesla-app/internal/app/config"
	"tesla-app/internal/app/dsn"
	"tesla-app/internal/app/handler"
	"tesla-app/internal/app/repository"
	"tesla-app/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

func main() {
	router := gin.Default()

	// Инициализация Minio клиента
	minioClient, err := minio.New("minio:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		logrus.Fatalf("Error creating Minio client: %v", err)
	}

	// Проверка подключения к Minio
	ctx := context.Background()
	_, err = minioClient.ListBuckets(ctx)
	if err != nil {
		logrus.Fatalf("Error connecting to Minio: %v", err)
	}
	logrus.Info("Successfully connected to Minio")

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

	// Передаем Minio клиент в handler
	hand := handler.NewHandler(repo, minioClient)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}
