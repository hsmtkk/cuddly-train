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
	e.GET("/profile/:id", hdl.profile)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

type handler struct {
}

func newHandler() *handler {
	return &handler{}
}

type responseFormat struct {
	Name        string `json:"name"`
	MailAddress string `json:"mailaddress"`
}

var names = []string{"Alice", "Bob", "Carol"}
var mails = []string{"alice@example.com", "bob@example.com", "carol@example.com"}

func (h *handler) profile(ectx echo.Context) error {
	id, err := strconv.Atoi(ectx.Param("id"))
	if err != nil {
		ectx.String(http.StatusBadRequest, fmt.Sprintf("strconv.Atoi failed; %v", err.Error()))
	}
	resp := responseFormat{
		Name:        names[id%len(names)],
		MailAddress: mails[id%len(mails)],
	}
	return ectx.JSON(http.StatusOK, resp)
}

func (h *handler) health(ectx echo.Context) error {
	return ectx.String(http.StatusOK, "OK")
}
