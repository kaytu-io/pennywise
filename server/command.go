package main

import (
	"fmt"
	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"gopkg.in/go-playground/validator.v9"
)

type Config struct {
	Mysql struct {
		Host     string `koanf:"host"`
		Port     string `koanf:"port"`
		DB       string `koanf:"db"`
		Username string `koanf:"username"`
		Password string `koanf:"password"`
	} `koanf:"mysql"`
	Http struct {
		Address string `koanf:"address"`
	} `koanf:"http"`
}

type customValidator struct {
	validate *validator.Validate
}

func (v customValidator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

func Command() *cobra.Command {
	return &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			return start(cmd.Context())
		},
	}
}

func start(ctx context.Context) error {
	config := load(true)
	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("new logger: %w", err)
	}
	handler, err := InitializeHttpHandler(logger, config.Mysql.Username, config.Mysql.Password, config.Mysql.Host,
		config.Mysql.Port, config.Mysql.DB)
	if err != nil {
		return fmt.Errorf("init http handler: %w", err)
	}

	if HttpAddress == "" {
		HttpAddress = ":8080"
	}
	return registerAndStart(logger, HttpAddress, handler)
}

func registerAndStart(logger *zap.Logger, address string, handler *HttpHandler) error {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Recover())
	e.Use(echozap.ZapLogger(logger))
	e.Validator = customValidator{
		validate: validator.New(),
	}
	handler.Register(e)
	return e.Start(address)
}
