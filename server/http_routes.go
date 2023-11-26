package main

import (
	"github.com/kaytu.io/pennywise/server/aws"
	"github.com/kaytu.io/pennywise/server/azurerm"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *HttpHandler) Register(e *echo.Echo) {
	v1 := e.Group("/api/v1")

	v1.PUT("/ingest/azure", h.IngestAzureTables)
	v1.PUT("/ingest/aws", h.IngestAwsTables)

}

func (h *HttpHandler) IngestAwsTables(ctx echo.Context) error {

	ingester, err := azurerm.NewIngester(ctx.Request().Context(), "Virtual Machines", "eastus")
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	err = IngestPricing(ctx.Request().Context(), h.backend, ingester)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, "Tables successfully ingested")
}

func (h *HttpHandler) IngestAzureTables(ctx echo.Context) error {

	ingester, err := aws.NewIngester("Virtual Machines", "eastus")
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	err = IngestPricing(ctx.Request().Context(), h.backend, ingester)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, "Tables successfully ingested")
}
