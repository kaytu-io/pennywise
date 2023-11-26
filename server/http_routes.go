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

	//cost := e.Group("/cost")
	//cost.GET("/azure/resource", h.GetResourceCost)
}

//func (h *HttpHandler) GetResourceCost(ctx echo.Context) error {
//	provider, err := azuretf.NewProvider(azuretf.ProviderName)
//	if err != nil {
//		return ctx.JSON(http.StatusInternalServerError, err.Error())
//	}
//	provider.ResourceComponents()
//	return nil
//}

// IngestAwsTables run the ingester to receive pricing and store in the database for aws services
// Params: service (query param), region (query param)
func (h *HttpHandler) IngestAwsTables(ctx echo.Context) error {
	service := ctx.QueryParam("service")
	region := ctx.QueryParam("region")
	ingester, err := azurerm.NewIngester(ctx.Request().Context(), service, region)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	err = IngestPricing(ctx.Request().Context(), h.backend, ingester)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, "Tables successfully ingested")
}

// IngestAzureTables run the ingester to receive pricing and store in the database for azure services
// Params: service (query param), region (query param)
func (h *HttpHandler) IngestAzureTables(ctx echo.Context) error {
	service := ctx.QueryParam("service")
	region := ctx.QueryParam("region")
	ingester, err := aws.NewIngester(service, region)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	err = IngestPricing(ctx.Request().Context(), h.backend, ingester)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, "Tables successfully ingested")
}
