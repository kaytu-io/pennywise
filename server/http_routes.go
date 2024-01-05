package main

import (
	awsrg "github.com/kaytu-io/pennywise/server/aws/region"
	awsres "github.com/kaytu-io/pennywise/server/aws/resources"
	azurermres "github.com/kaytu-io/pennywise/server/azurerm/resources"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/resource"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (h *HttpHandler) Register(e *echo.Echo) {
	v1 := e.Group("/api/v1")

	v1.PUT("/ingest", h.IngestTables)
	v1.GET("/ingest/jobs/:id", h.GetIngestionJob)
	v1.GET("/ingest/jobs", h.ListIngestionJobs)

	cost := v1.Group("/cost")
	cost.GET("/resource", h.GetResourceCost)
	cost.GET("/state", h.GetStateCost)
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

func (h *HttpHandler) GetStateCost(ctx echo.Context) error {
	var req resource.State
	var qResources []query.Resource
	if err := bindValidate(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	for _, res := range req.Resources {
		if res.ProviderName == resource.AzureProvider {
			provider, err := azurermres.NewProvider(azurermres.ProviderName, h.logger)
			if err != nil {
				return err
			}
			resources := make(map[string]resource.Resource)
			resources[res.Address] = res
			components := provider.ResourceComponents(resources, res)
			qResources = append(qResources, query.Resource{
				Address:    res.Address,
				Provider:   res.ProviderName,
				Type:       res.Type,
				Components: components,
			})
		} else if res.ProviderName == resource.AWSProvider {
			provider, err := awsres.NewProvider(awsres.ProviderName, awsrg.Code(res.RegionCode), h.logger)
			if err != nil {
				return err
			}
			resources := make(map[string]resource.Resource)
			resources[res.Address] = res
			components := provider.ResourceComponents(resources, res)
			qResources = append(qResources, query.Resource{
				Address:    res.Address,
				Provider:   res.ProviderName,
				Type:       res.Type,
				Components: components,
			})
		}
	}

	state, err := cost.NewState(ctx.Request().Context(), h.backend, qResources, h.logger)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, state)
}

func (h *HttpHandler) GetResourceCost(ctx echo.Context) error {
	var req resource.Resource
	var qResource query.Resource
	if err := bindValidate(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if req.ProviderName == resource.AzureProvider {
		provider, err := azurermres.NewProvider(azurermres.ProviderName, h.logger)
		if err != nil {
			return err
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
	} else if req.ProviderName == resource.AWSProvider {
		provider, err := awsres.NewProvider(awsres.ProviderName, awsrg.Code(req.RegionCode), h.logger)
		if err != nil {
			return err
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
	state, err := cost.NewState(ctx.Request().Context(), h.backend, []query.Resource{qResource}, h.logger)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, state)
}

// IngestTables adds an ingestion job to receive pricing and store in the database
// Params: provider (query param), service (query param), region (query param)
func (h *HttpHandler) IngestTables(ctx echo.Context) error {
	provider := ctx.QueryParam("provider")
	service := ctx.QueryParam("service")
	region := ctx.QueryParam("region")
	lastId, err := h.scheduler.MakeJob(provider, service, region)
	if err != nil {
		return err
	}
	job, err := h.scheduler.GetJobById(int32(lastId))
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, job)
}

// ListIngestionJobs returns list of ingestion jobs with the provided filters
// Params: provider (query param), service (query param), region (query param), status (query param)
func (h *HttpHandler) ListIngestionJobs(ctx echo.Context) error {
	provider := ctx.QueryParam("provider")
	status := ctx.QueryParam("status")
	service := ctx.QueryParam("service")
	location := ctx.QueryParam("region")
	jobs, err := h.scheduler.GetJobs(status, provider, service, location)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, jobs)
}

// GetIngestionJob returns an ingestion job with the provided id
// Params: id (route param)
func (h *HttpHandler) GetIngestionJob(ctx echo.Context) error {
	idStr := ctx.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return err
	}
	id := int32(id64)
	job, err := h.scheduler.GetJobById(id)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, *job)
}
