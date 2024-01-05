package resources

import (
	"github.com/kaytu-io/pennywise/server/azurerm/resources"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/kaytu-io/pennywise/server/resource"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/kaytu-io/pennywise/server/aws/region"
)

// Volume represents an EBS volume that can be cost-estimated.
type Volume struct {
	provider *Provider
	logger   *zap.Logger
	region   region.Code

	volumeType string
	size       decimal.Decimal
	iops       decimal.Decimal
}

// volumeValues represents the structure of Terraform values for aws_ebs_volume resource.
type volumeValues struct {
	AvailabilityZone string  `mapstructure:"availability_zone"`
	Type             string  `mapstructure:"type"`
	Size             float64 `mapstructure:"size"`
	IOPS             float64 `mapstructure:"iops"`
}

// decodeVolumeValues decodes and returns volumeValues from a Terraform values map.
func decodeVolumeValues(tfVals map[string]interface{}) (volumeValues, error) {
	var v volumeValues
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

// newVolume creates a new Volume from volumeValues.
func (p *Provider) newVolume(vals volumeValues) *Volume {
	v := &Volume{
		provider:   p,
		logger:     p.logger,
		region:     p.region,
		volumeType: "gp3",
		size:       decimal.NewFromInt(8),
		iops:       decimal.NewFromInt(16000),
	}

	if reg := region.NewFromZone(vals.AvailabilityZone); reg.Valid() {
		v.region = reg
	}

	if vals.Type != "" {
		v.volumeType = vals.Type
	}

	if vals.Size > 0 {
		v.size = decimal.NewFromFloat(vals.Size)
	}

	if vals.IOPS > 0 {
		v.iops = decimal.NewFromFloat(vals.IOPS)
	}

	return v
}

// Components returns the price component queries that make up the Volume.
func (v *Volume) Components() []resource.Component {
	comps := []resource.Component{v.storageComponent()}

	if v.volumeType == "io1" || v.volumeType == "io2" {
		comps = append(comps, v.iopsComponent())
	}

	resources.GetCostComponentNamesAndSetLogger(comps, v.logger)
	return comps
}

func (v *Volume) storageComponent() resource.Component {
	return resource.Component{
		Name:            "Storage",
		MonthlyQuantity: v.size,
		Unit:            "GB",
		Details:         []string{v.volumeType},
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(v.provider.key),
			Service:  util.StringPtr("AmazonEC2"),
			Family:   util.StringPtr("Storage"),
			Location: util.StringPtr(v.region.String()),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "VolumeAPIName", Value: util.StringPtr(v.volumeType)},
			},
		},
	}
}

func (v *Volume) iopsComponent() resource.Component {
	return resource.Component{
		Name:            "Provisioned IOPS",
		MonthlyQuantity: v.iops,
		Unit:            "IOPS",
		ProductFilter: &product.Filter{
			Provider: util.StringPtr(v.provider.key),
			Service:  util.StringPtr("AmazonEC2"),
			Family:   util.StringPtr("System Operation"),
			Location: util.StringPtr(v.region.String()),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "VolumeAPIName", Value: util.StringPtr(v.volumeType)},
				{Key: "UsageType", ValueRegex: util.StringPtr("^EBS:VolumeP-IOPS")},
			},
		},
	}
}
