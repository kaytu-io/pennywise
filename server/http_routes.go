package main

import (
	"github.com/kaytu-io/pennywise/server/aws"
	awsrg "github.com/kaytu-io/pennywise/server/aws/region"
	awstf "github.com/kaytu-io/pennywise/server/aws/terraform"
	"github.com/kaytu-io/pennywise/server/azurerm"
	azuretf "github.com/kaytu-io/pennywise/server/azurerm/terraform"
	"github.com/kaytu-io/pennywise/server/cost"
	ingester2 "github.com/kaytu-io/pennywise/server/internal/ingester"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/resource"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *HttpHandler) Register(e *echo.Echo) {
	v1 := e.Group("/api/v1")

	v1.PUT("/ingest/azure", h.IngestAzureTables)
	v1.PUT("/ingest/aws", h.IngestAwsTables)

	cost := e.Group("/cost")
	cost.GET("/resource", h.GetResourceCost)
}

func bindValidate(ctx echo.Context, i any) error {
	if err := ctx.Bind(i); err != nil {
		return err
	}

	if err := ctx.Validate(i); err != nil {
		return err
	}

	return nil
}

func (h *HttpHandler) GetResourceCost(ctx echo.Context) error {
	var req resource.Resource
	var qResource query.Resource
	if err := bindValidate(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if req.ProviderName == "azurerm" {
		provider, err := azuretf.NewProvider(azuretf.ProviderName)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		resources := make(map[string]resource.Resource)
		resources[req.Address] = req
		components := provider.ResourceComponents(resources, req)
		qResource = query.Resource{
			Address:    req.Address,
			Provider:   req.ProviderName,
			Type:       req.Type,
			Components: components,
		}
	} else if req.ProviderName == "aws" {
		provider, err := awstf.NewProvider(awstf.ProviderName, awsrg.Code(req.RegionCode))
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		resources := make(map[string]resource.Resource)
		resources[req.Address] = req
		components := provider.ResourceComponents(resources, req)
		qResource = query.Resource{
			Address:    req.Address,
			Provider:   req.ProviderName,
			Type:       req.Type,
			Components: components,
		}
	}
	state, err := cost.NewState(ctx.Request().Context(), h.backend, []query.Resource{qResource})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	cost, err := state.Cost()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, cost)
}

// IngestAwsTables run the ingester to receive pricing and store in the database for aws services
// Params: service (query param), region (query param)
func (h *HttpHandler) IngestAwsTables(ctx echo.Context) error {
	service := ctx.QueryParam("service")
	region := ctx.QueryParam("region")
	ingester, err := azurerm.NewIngester(ctx.Request().Context(), service, region)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	err = ingester2.IngestPricing(ctx.Request().Context(), h.backend, ingester)
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
	err = ingester2.IngestPricing(ctx.Request().Context(), h.backend, ingester)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, "Tables successfully ingested")
}
