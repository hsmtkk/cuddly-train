package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/hsmtkk/cuddly-train/env"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	port, err := env.Port()
	if err != nil {
		log.Fatal(err)
	}

	hdl := newHandler()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hdl.index)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

type handler struct{}

func newHandler() *handler {
	return &handler{}
}

type responseFormat struct {
	Sum     int `json:"sum"`
	Product int `json:"product"`
}

func (h *handler) index(ectx echo.Context) error {
	num1, err := strconv.Atoi(ectx.QueryParam("num1"))
	if err != nil {
		return fmt.Errorf("strconv.Atoi failed; %w", err)
	}
	num2, err := strconv.Atoi(ectx.QueryParam("num2"))
	if err != nil {
		return fmt.Errorf("strconv.Atoi failed; %w", err)
	}
	resp := responseFormat{
		Sum:     num1 + num2,
		Product: num1 * num2,
	}
	return ectx.JSON(http.StatusOK, resp)
}
