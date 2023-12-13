package resources

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/tier_request"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
)

type VirtualNetworkGateway struct {
	provider *Provider

	sku      *string
	location string
	// usage
	p2sConnection         *int64
	monthlyDataTransferGB *int64
}
type VirtualNetworkGatewayValue struct {
	Sku      string `mapstructure:"sku"`
	Location string `mapstructure:"location"`

	Usage struct {
		P2sConnection         int64 `mapstructure:"p2s_connection"`
		MonthlyDataTransferGB int64 `mapstructure:"monthly_data_transfer_gb"`
	} `mapstructure:"pennywise_usage"`
}

func decodeVirtualNetworkGateway(tfVals map[string]interface{}) (VirtualNetworkGatewayValue, error) {
	var v VirtualNetworkGatewayValue
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

func (p *Provider) newVirtualNetworkGateway(vals VirtualNetworkGatewayValue) *VirtualNetworkGateway {
	inst := &VirtualNetworkGateway{
		sku:                   &vals.Sku,
		location:              vals.Location,
		p2sConnection:         &vals.Usage.P2sConnection,
		monthlyDataTransferGB: &vals.Usage.MonthlyDataTransferGB,
	}
	return inst
}

func (inst *VirtualNetworkGateway) Components() []query.Component {
	var connection, dataTransfers *decimal.Decimal
	sku := "Basic"
	region := inst.location
	zone := regionToZone(region)

	if inst.sku != nil {
		sku = *inst.sku
	}
	meterName := sku

	costComponents := make([]query.Component, 0)

	if sku == "Basic" {
		meterName = "Basic Gateway"
	}

	costComponents = append(costComponents, vpnGateway(region, sku, meterName))

	if inst.p2sConnection != nil {
		connection = decimalPtr(decimal.NewFromInt(*inst.p2sConnection))
		if connection != nil {
			connectionLimits := []int{128}
			connectionValues := tier_request.CalculateTierBuckets(*connection, connectionLimits)
			if connectionValues[1].GreaterThan(decimal.Zero) {
				costComponents = append(costComponents, vpnGatewayP2S(region, sku, &connectionValues[1]))
			}
		}
	} else {
		costComponents = append(costComponents, vpnGatewayP2S(region, sku, connection))
	}

	if inst.monthlyDataTransferGB != nil {
		dataTransfers = decimalPtr(decimal.NewFromInt(*inst.monthlyDataTransferGB))
		if dataTransfers != nil {
			costComponents = append(costComponents, vpnGatewayDataTransfers(zone, sku, dataTransfers))
		}
	} else {
		costComponents = append(costComponents, vpnGatewayDataTransfers(zone, sku, dataTransfers))
	}

	return costComponents
}

func vpnGateway(region, sku, meterName string) query.Component {
	return query.Component{
		Name:           fmt.Sprintf("VPN gateway (%s)", sku),
		Unit:           "hours",
		HourlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Location: util.StringPtr(region),
			Service:  util.StringPtr("VPN Gateway"),
			Family:   util.StringPtr("Networking"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "meter_name", Value: util.StringPtr(meterName)},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

func vpnGatewayP2S(region, sku string, connection *decimal.Decimal) query.Component {
	return query.Component{
		Name:           "VPN gateway P2S tunnels (over 128)",
		Unit:           "tunnel",
		HourlyQuantity: *connection,
		ProductFilter: &product.Filter{
			Location: util.StringPtr(region),
			Service:  util.StringPtr("VPN Gateway"),
			Family:   util.StringPtr("Networking"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "sku_name", Value: util.StringPtr(sku)},
				{Key: "meter_name", ValueRegex: util.StringPtr(fmt.Sprintf("%s", "P2S Connection"))},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}

func vpnGatewayDataTransfers(zone, sku string, dataTransfers *decimal.Decimal) query.Component {
	return query.Component{
		Name:            "VPN gateway data tranfer",
		Unit:            "GB",
		MonthlyQuantity: *dataTransfers,
		ProductFilter: &product.Filter{
			Location: util.StringPtr(zone),
			Service:  util.StringPtr("VPN Gateway"),
			Family:   util.StringPtr("Networking"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "product_name", ValueRegex: util.StringPtr(fmt.Sprintf("%s", "VPN Gateway Bandwidth"))},
				{Key: "meter_name", ValueRegex: util.StringPtr(fmt.Sprintf("%s", "Inter-Virtual Network Data Transfer Out"))},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}
