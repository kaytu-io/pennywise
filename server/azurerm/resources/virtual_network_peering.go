package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/kaytu-io/pennywise/server/resource"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"strings"
)

// VirtualNetworkPeering is the entity that holds the logic to calculate price
// of the azure_network_virtualnetwork_peering
type VirtualNetworkPeering struct {
	provider *Provider

	sourceLocation      string
	destinationLocation string

	// Usage
	monthlyDataTransferGB decimal.Decimal
}

type DestinationLocationStruct struct {
	Values struct {
		Location string `mapstructure:"location"`
	} `mapstructure:"Values"`
}

type SourceLocationStruct struct {
	Values struct {
		Location string `mapstructure:"location"`
	} `mapstructure:"Values"`
}

// virtualNetworkPeeringValues is holds the values that we need to be able
// to calculate the price of the Virtual Network Peering Values
type virtualNetworkPeeringValues struct {
	SourceLocation      SourceLocationStruct      `mapstructure:"virtual_network_name"`
	DestinationLocation DestinationLocationStruct `mapstructure:"remote_virtual_network_id"`

	Usage struct {
		// receive monthly inbound/outbound data transferred by the VNET peering in GB.
		MonthlyDataTransferGB float64 `mapstructure:"monthly_data_transfer_gb"`
	} `mapstructure:"pennywise_usage"`
}

// decodeVirtualNetworkPeeringValues decodes and returns computeInstanceValues from a Terraform values map.
func decodeVirtualNetworkPeeringValues(tfVals map[string]interface{}) (virtualNetworkPeeringValues, error) {
	var v virtualNetworkPeeringValues
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
	if v.Usage.MonthlyDataTransferGB == 0 {
		v.Usage.MonthlyDataTransferGB = 100
	}
	return v, nil
}

func (p *Provider) newVirtualNetworkPeering(vals virtualNetworkPeeringValues) *VirtualNetworkPeering {
	sourceLocation := getLocationName(vals.SourceLocation.Values.Location)
	destinationLocation := getLocationName(vals.DestinationLocation.Values.Location)
	inst := &VirtualNetworkPeering{
		provider: p,

		sourceLocation:        sourceLocation,
		destinationLocation:   destinationLocation,
		monthlyDataTransferGB: decimal.NewFromFloat(vals.Usage.MonthlyDataTransferGB),
	}
	return inst
}

func (inst *VirtualNetworkPeering) Components() []resource.Component {
	firstQuery := inst.egressDataProcessedCostComponent(inst.provider.key)
	secondQuery := inst.ingressDataProcessedCostComponent(inst.provider.key)
	components := []resource.Component{
		firstQuery,
		secondQuery,
	}
	GetCostComponentNamesAndSetLogger(components, inst.provider.logger)

	return components
}

func (inst *VirtualNetworkPeering) egressDataProcessedCostComponent(key string) resource.Component {
	if inst.sourceLocation == inst.destinationLocation {
		return resource.Component{
			Name:            "Outbound data transfer",
			Unit:            "GB",
			MonthlyQuantity: inst.monthlyDataTransferGB,
			ProductFilter: &product.Filter{
				Provider: util.StringPtr(key),
				Service:  util.StringPtr("Virtual Network"),
				Family:   util.StringPtr("Networking"),
				Location: util.StringPtr("Global"),
				AttributeFilters: []*product.AttributeFilter{
					{Key: "meter_name", Value: util.StringPtr("Intra-Region Egress")},
				},
			},
			PriceFilter: &price.Filter{
				AttributeFilters: []*price.AttributeFilter{
					{Key: "type", Value: util.StringPtr("Consumption")},
				},
			},
		}
	}

	return resource.Component{
		Name:            "Outbound data transfer",
		Unit:            "GB",
		MonthlyQuantity: inst.monthlyDataTransferGB,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(key),
			Location: util.StringPtr(virtualNetworkPeeringConvertRegion(inst.sourceLocation)),
			Family:   util.StringPtr("Networking"),
			Service:  util.StringPtr("VPN Gateway"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", ValueRegex: util.StringPtr("VPN Gateway Bandwidth")},
				{Key: "meter_name", ValueRegex: util.StringPtr("Inter-Virtual Network Data Transfer Out")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

func (inst *VirtualNetworkPeering) ingressDataProcessedCostComponent(key string) resource.Component {
	if inst.sourceLocation == inst.destinationLocation {
		return resource.Component{
			Name:            "Inbound data transfer",
			Unit:            "GB",
			MonthlyQuantity: inst.monthlyDataTransferGB,
			ProductFilter: &product.Filter{
				Provider: util.StringPtr(key),
				Location: util.StringPtr("Global"),
				Service:  util.StringPtr("Virtual Network"),
				Family:   util.StringPtr("Networking"),
				AttributeFilters: []*product.AttributeFilter{
					{Key: "meter_name", Value: util.StringPtr("Intra-Region Ingress")},
				},
			},
			PriceFilter: &price.Filter{
				AttributeFilters: []*price.AttributeFilter{
					{Key: "type", Value: util.StringPtr("Consumption")},
				},
			},
		}
	}

	return resource.Component{
		Name:            "Inbound data transfer",
		Unit:            "GB",
		MonthlyQuantity: inst.monthlyDataTransferGB,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(key),
			Location: util.StringPtr(virtualNetworkPeeringConvertRegion(inst.destinationLocation)),
			Service:  util.StringPtr("VPN Gateway"),
			Family:   util.StringPtr("Networking"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", ValueRegex: util.StringPtr("VPN Gateway Bandwidth")},
				{Key: "meter_name", ValueRegex: util.StringPtr("Inter-Virtual Network Data Transfer Out")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

func virtualNetworkPeeringConvertRegion(region string) string {
	zone := regionToZone(region)
	if strings.HasPrefix(strings.ToLower(region), "usgov") {
		zone = "US Gov Zone 1"
	}
	if strings.HasPrefix(strings.ToLower(region), "germany") {
		zone = "DE Zone 1"
	}
	if strings.HasPrefix(strings.ToLower(region), "china") {
		zone = "CN Zone 1"
	}
	return zone
}
