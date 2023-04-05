package main

import (
	"context"
	"encoding/gob"
	"github.com/firat.atmaca/go-mongo/core"
	"github.com/firat.atmaca/go-mongo/handlers"
	"github.com/firat.atmaca/go-mongo/modules/config"
	. "github.com/firat.atmaca/go-mongo/modules/database/mongo"
	"github.com/firat.atmaca/go-mongo/modules/database/storage"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
)

var app config.AppTools

func main() {
	gob.Register(map[string]interface{}{})
	gob.Register(primitive.NewObjectID())

	InfoLogger := log.New(os.Stderr, " ", log.LstdFlags|log.Lshortfile)
	ErrorLogger := log.New(os.Stderr, " ", log.LstdFlags|log.Lshortfile)

	validator := validator.New()
	app.InfoLogger = *InfoLogger
	app.ErrorLogger = *ErrorLogger
	app.Validator = validator

	err := godotenv.Load()
	if err != nil {
		app.ErrorLogger.Fatal("No .env file available")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		app.ErrorLogger.Fatalln("mongodb uri string not found : ")
	}

	client, err := NewClient(uri)
	if err != nil {
		app.ErrorLogger.Fatalln("mongodb client is not ready : ")
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			app.ErrorLogger.Fatal(err)
			return
		}
	}()
	appRouter := gin.New()

	userStorage := storage.NewUserStorage(&app, client)
	userHandler := handlers.NewUserHandler(&app, userStorage)
	core.Routes(appRouter, userHandler)

	err = appRouter.Run()
	if err != nil {
		app.ErrorLogger.Fatal(err)
	}
}
