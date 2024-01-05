package resources

import (
	"github.com/kaytu-io/pennywise/server/resource"
	"github.com/mitchellh/mapstructure"
)

type DNSMXRecord struct {
	provider *Provider

	location string
	// Usage
	monthlyQueries *int64
}

type dnsMAXRecordValues struct {
	ResourceGroupName ResourceGroupNameStruct `mapstructure:"resource_group_name"`

	Usage struct {
		// receives monthly number of DNS queries
		MonthlyQueries int64 `mapstructure:"monthly_queries"`
	} `mapstructure:"pennywise_usage"`
}

func decoderDNSMXRecord(tfVals map[string]interface{}) (dnsMAXRecordValues, error) {
	var v dnsMAXRecordValues
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

func (p *Provider) newDNSMXRecord(vals dnsMAXRecordValues) *DNSMXRecord {
	inst := &DNSMXRecord{
		provider:       p,
		location:       vals.ResourceGroupName.Values.Location,
		monthlyQueries: &vals.Usage.MonthlyQueries,
	}
	return inst
}

func (inst *DNSMXRecord) component() []resource.Component {
	region := getLocationName(inst.location)
	costComponents := DNSQueriesCostComponent(inst.provider.key, region, inst.monthlyQueries)
	GetCostComponentNamesAndSetLogger(costComponents, inst.provider.logger)
	return costComponents
}
