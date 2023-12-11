package resources

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"strings"
)

type VirtualNetworkGatewayConnection struct {
	provider *Provider
	typeName *string
	sku      *string
	location string
}

type VirtualNetworkGatewayId struct {
	Values struct {
		Sku string `mapstructure:"sku"`
	} `mapstructure:"values"`
}

type VirtualNetworkGatewayConnectionValue struct {
	VirtualNetworkGatewayId []VirtualNetworkGatewayId `mapstructure:"virtual_network_gateway_id"`
	Location                string                    `mapstructure:"location"`
	Type                    string                    `mapstructure:"type"`
}

func decodeVirtualNetworkGatewayConnection(tfVals map[string]interface{}) (VirtualNetworkGatewayConnectionValue, error) {
	var v VirtualNetworkGatewayConnectionValue
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
	fmt.Printf("tfvalue : %v \n ", tfVals)
	fmt.Printf("value : %v \n", v)
	return v, nil
}

func (p *Provider) newVirtualNetworkGatewayConnection(vals VirtualNetworkGatewayConnectionValue) *VirtualNetworkGatewayConnection {
	var sku string
	if len(vals.VirtualNetworkGatewayId) > 0 {
		sku = vals.VirtualNetworkGatewayId[0].Values.Sku
	}
	inst := &VirtualNetworkGatewayConnection{
		sku:      &sku,
		location: vals.Location,
		typeName: &vals.Type,
	}
	return inst
}

func (inst *VirtualNetworkGatewayConnection) Component() []query.Component {
	sku := "Basic"
	if inst.sku != nil {
		sku = *inst.sku
	}

	region := inst.location
	if strings.ToLower(sku) == "basic" {
		return nil
	}

	costComponents := make([]query.Component, 0)
	if inst.typeName != nil {
		if strings.ToLower(*inst.typeName) == "ipsec" {
			costComponents = append(costComponents, vpnGatewayS2S(region, sku))
		}
	}

	return costComponents
}

func vpnGatewayS2S(region, sku string) query.Component {
	return query.Component{
		Name:           fmt.Sprintf("VPN gateway (%s)", sku),
		Unit:           "hours",
		HourlyQuantity: decimal.NewFromInt(1),
		ProductFilter: &product.Filter{
			Location: util.StringPtr(region),
			Service:  util.StringPtr("VPN Gateway"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "sku_name", Value: util.StringPtr(sku)},
				{Key: "meter_name", ValueRegex: util.StringPtr(fmt.Sprintf("%s", "S2S Connection"))},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	}
}
