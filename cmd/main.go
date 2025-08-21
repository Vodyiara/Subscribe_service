package main

import (
	"Test_project_Effective_Mobile"
	"Test_project_Effective_Mobile/pkg/handler"
	"Test_project_Effective_Mobile/pkg/repository"
	"Test_project_Effective_Mobile/pkg/service"
	"context"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

// @title Subscribe_service
// @version 1.0
// @description API server for Subscribe_service

// @host localhost:8000
// @BasePath /

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := InitConfig(); err != nil {
		logrus.Fatalf("init config error: %v", err)
	}

	if err := godotenv.Load(); err != nil {
		logrus.Warn("no .env file found, continuing without it")
	}

	db, err := repository.NewPostgreDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("init postgre db error: %v", err)
	}
	logrus.Infof("connected to db %s:%s/%s", viper.GetString("db.host"), viper.GetString("db.port"), viper.GetString("db.dbname"))

	repos := repository.NewRepository(db)
	service := service.NewService(repos)
	handlers := handler.NewHandler(service)

	srv := new(Test_project_Effective_Mobile.Server)
	go func() {
		port := viper.GetString("port")

		logrus.Infof("starting server on port %s", port)

		if err := srv.Run(port, handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occurred running http server: %s", err.Error())
		}
	}()

	logrus.Info("App started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sig := <-quit
	logrus.Infof("received signal: %s. Shutting down...", sig)

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("server shutdown error: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("db close error: %s", err.Error())
	}

	logrus.Info("App stopped gracefully")
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
