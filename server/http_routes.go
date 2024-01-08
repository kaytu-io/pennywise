package main

import (
	awsrg "github.com/kaytu-io/pennywise/server/aws/region"
	awsres "github.com/kaytu-io/pennywise/server/aws/resources"
	azurermres "github.com/kaytu-io/pennywise/server/azurerm/resources"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/kaytu-io/pennywise/server/resource"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (h *HttpHandler) Register(e *echo.Echo) {
	v1 := e.Group("/api/v1")

	v1.PUT("/ingestion/jobs", h.IngestTables)
	v1.GET("/ingestion/jobs/:id", h.GetIngestionJob)
	v1.GET("/ingestion/jobs", h.ListIngestionJobs)

	cost := v1.Group("/cost")
	cost.GET("/state", h.GetTerraformCost)
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

// GetTerraformCost godoc
//
//	@Summary	Returns breakdown costs for a terraform file
//	@Tags		cost
//	@Param		request	body	resource.State	true	"Terraform state"
//	@Produce	json
//	@Success	200	{object}	cost.State
//	@Router		/api/v1/cost/state [get]
func (h *HttpHandler) GetTerraformCost(ctx echo.Context) error {
	var req resource.State
	var qResources []resource.Resource
	if err := bindValidate(ctx, &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	for _, res := range req.Resources {
		if res.ProviderName == resource.AzureProvider {
			provider, err := azurermres.NewProvider(azurermres.ProviderName, h.logger)
			if err != nil {
				return err
			}
			resources := make(map[string]resource.ResourceDef)
			resources[res.Address] = res
			components := provider.ResourceComponents(resources, res)
			qResources = append(qResources, resource.Resource{
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
			resources := make(map[string]resource.ResourceDef)
			resources[res.Address] = res
			components := provider.ResourceComponents(resources, res)
			qResources = append(qResources, resource.Resource{
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

// IngestTables godoc
//
//	@Summary	Adds an ingestion job to receive pricing and store in the database
//	@Tags		ingestion
//	@Param		provider	query	string	true	"provider"
//	@Param		service		query	string	true	"service"
//	@Param		region		query	string	false	"region"
//	@Produce	json
//	@Success	200	{object}	ingester.IngestionJob
//	@Router		/api/v1/ingestion/jobs [put]
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

// ListIngestionJobs godoc
//
//	@Summary	Returns list of ingestion jobs with the provided filters
//	@Tags		ingestion
//	@Param		provider	query	string	false	"provider"
//	@Param		service		query	string	false	"service"
//	@Param		region		query	string	false	"region"
//	@Param		status		query	string	false	"status"
//	@Produce	json
//	@Success	200	{object}	[]ingester.IngestionJob
//	@Router		/api/v1/ingestion/jobs [get]
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

// GetIngestionJob godoc
//
//	@Summary	Returns an ingestion job with the provided id
//	@Tags		ingestion
//	@Param		provider	path	string	true	"provider"
//	@Produce	json
//	@Success	200	{object}	ingester.IngestionJob
//	@Router		/api/v1/ingestion/jobs/{id} [get]
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
