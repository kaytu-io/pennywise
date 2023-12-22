package resources

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
	"strings"
)

// ServicePlan is the entity that holds the logic to calculate price
// of the azurerm_service_plan
type ServicePlan struct {
	provider *Provider

	location    string
	skuName     string
	workerCount int64
	osType      string
}

// servicePlanValues is holds the values that we need to be able
// to calculate the price of the ServicePlan
type servicePlanValues struct {
	Location    string `mapstructure:"location"`
	SkuName     string `mapstructure:"sku_name"`
	WorkerCount *int64 `mapstructure:"worker_count"`
	OsType      string `mapstructure:"os_type"`
}

// decodeServicePlanValues decodes and returns servicePlanValues from a Terraform values map.
func decodeServicePlanValues(tfVals map[string]interface{}) (servicePlanValues, error) {
	var v servicePlanValues
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &v,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return v, err
	}

	if err := decoder.Decode(tfVals); err != nil {
		return v, err
	}
	return v, nil
}

// newServiceEnvironment initializes a new ServicePlan from the provider
func (p *Provider) newServicePlan(vals servicePlanValues) *ServicePlan {
	workerCount := int64(1)
	if vals.WorkerCount != nil {
		workerCount = *vals.WorkerCount
	}
	inst := &ServicePlan{
		provider: p,

		location:    vals.Location,
		skuName:     vals.SkuName,
		workerCount: workerCount,
		osType:      vals.OsType,
	}
	return inst
}

func (inst *ServicePlan) Components() []query.Component {
	var components []query.Component

	productName := "Standard Plan"
	sku := inst.skuName

	if len(inst.skuName) < 2 || strings.ToLower(inst.skuName[:2]) == "ep" || strings.ToLower(inst.skuName[:2]) == "ws" || strings.ToLower(inst.skuName[:2]) == "y1" {
		return components
	}

	firstLetter := strings.ToLower(inst.skuName[:1])
	os := strings.ToLower(inst.osType)
	var additionalAttributeFilters []*product.AttributeFilter

	switch firstLetter {
	case "s":
		sku = "S" + inst.skuName[1:]
	case "b":
		sku = "B" + inst.skuName[1:]
		productName = "Basic Plan"
	case "f":
		productName = "Free Plan"
	case "d":
		sku = "Shared"
		productName = "Shared Plan"
	case "p", "i":
		sku, productName, additionalAttributeFilters = getVersionedAppServicePlanSKU(sku, os)
	}

	if strings.ToLower(inst.skuName) == "shared" {
		sku = "Shared"
		productName = "Shared Plan"
	}

	if os == "linux" && productName != "Isolated Plan" && productName != "Premium Plan" && productName != "Shared Plan" {
		productName += " - Linux"
	}

	components = append(components, servicePlanCostComponent(
		inst.location,
		fmt.Sprintf("Instance usage (%s)", inst.skuName),
		productName,
		sku,
		inst.workerCount,
		additionalAttributeFilters...))

	return components
}
