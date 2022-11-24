package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	profileAddress, err := env.RequiredVar("PROFILE_ADDRESS")
	if err != nil {
		log.Fatal(err)
	}
	orderAddress, err := env.RequiredVar("ORDER_ADDRESS")
	if err != nil {
		log.Fatal(err)
	}

	hdl := newHandler(profileAddress, orderAddress)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hdl.health)
	e.GET("/customer/:id", hdl.customer)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

type handler struct {
	profileAddress string
	orderAddress   string
}

func newHandler(profileAddress, orderAddress string) *handler {
	return &handler{profileAddress, orderAddress}
}

type responseFormat struct {
	Profile interface{} `json:"profile"`
	Order   interface{} `json:"order"`
}

func (h *handler) customer(ectx echo.Context) error {
	id, err := strconv.Atoi(ectx.Param("id"))
	if err != nil {
		return ectx.String(http.StatusBadRequest, fmt.Sprintf("strconv.Atoi failed; %v", err.Error()))
	}
	profileURL := fmt.Sprintf("http://%s/profile/%d", h.profileAddress, id)
	orderURL := fmt.Sprintf("http://%s/order/%d", h.orderAddress, id)
	profileResp, err := h.httpGet(profileURL)
	if err != nil {
		return ectx.String(http.StatusInternalServerError, err.Error())
	}
	orderResp, err := h.httpGet(orderURL)
	if err != nil {
		return ectx.String(http.StatusInternalServerError, err.Error())
	}
	var response responseFormat
	if err := json.Unmarshal(profileResp, &response.Profile); err != nil {
		return ectx.String(http.StatusInternalServerError, err.Error())
	}
	if err := json.Unmarshal(orderResp, &response.Order); err != nil {
		return ectx.String(http.StatusInternalServerError, err.Error())
	}
	return ectx.JSON(http.StatusOK, response)
}

func (h *handler) httpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http.Get failed; %s; %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non 200 HTTP status; %d; %s", resp.StatusCode, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll failed; %w", err)
	}
	return body, nil
}

func (h *handler) health(ectx echo.Context) error {
	return ectx.String(http.StatusOK, "OK")
}
