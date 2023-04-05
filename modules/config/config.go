package config

import (
	"github.com/go-playground/validator/v10"
	"log"
)

type AppTools struct {
	ErrorLogger log.Logger
	InfoLogger  log.Logger
	Validator   *validator.Validate
}
