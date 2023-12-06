package resources

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"strings"
)

type VirtualNetworkGatewayConnection struct {
	provider                *Provider
	typeName                *string
	virtualNetworkGatewayId map[string]string
	location                string
}
type VirtualNetworkGatewayConnectionValue struct {
	VirtualNetworkGatewayId map[string]string `mapstructure:"virtual_network_gateway_id"`
	location                string            `mapstructure:"location"`
	Type                    string            `mapstructure:"type"`
}

func decoderVirtualNetworkGatewayConnection(tfVals map[string]interface{}) (VirtualNetworkGatewayConnectionValue, error) {
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
	return v, nil
}

func (p *Provider) newVirtualNetworkGatewayConnection(vals VirtualNetworkGatewayConnectionValue) *VirtualNetworkGatewayConnection {
	inst := &VirtualNetworkGatewayConnection{
		virtualNetworkGatewayId: vals.VirtualNetworkGatewayId,
		location:                vals.location,
		typeName:                &vals.Type,
	}
	return inst
}

func (inst *VirtualNetworkGatewayConnection) Component() []query.Component {
	sku := "Basic"
	if len(inst.virtualNetworkGatewayId) > 0 {
		for k, v := range inst.virtualNetworkGatewayId {
			if k == "sku" {
				sku = v
			}
		}
	}

	region := inst.location
	if strings.ToLower(sku) == "basic" {
		//return &schema.Resource{
		//	Name:      d.Address,
		//	NoPrice:   true,
		//	IsSkipped: true,
		//}
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
				{Key: "skuName", Value: util.StringPtr(sku)},
				{Key: "meterName", ValueRegex: util.StringPtr(fmt.Sprintf("/%s/i", "S2S Connection"))},
			},
		},
	}
}
