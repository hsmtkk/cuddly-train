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
	e.GET("/", hdl.health)
	e.GET("/order/:id", hdl.order)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

type handler struct {
}

func newHandler() *handler {
	return &handler{}
}

type responseFormat struct {
	Items []string `json:"items"`
}

var items = [][]string{
	{"apple", "book", "chocolate"},
	{"dogfood", "egg", "fruit"},
	{"ham", "ink", "juice"},
}

func (h *handler) order(ectx echo.Context) error {
	idStr := ectx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ectx.String(http.StatusBadRequest, fmt.Sprintf("strconv.Atoi failed; %v", err.Error()))
	}
	resp := responseFormat{
		Items: items[id%len(items)],
	}
	return ectx.JSON(http.StatusOK, resp)
}

func (h *handler) health(ectx echo.Context) error {
	return ectx.String(http.StatusOK, "OK")
}
