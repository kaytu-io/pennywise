package resources

import (
	"github.com/kaytu-io/pennywise/server/resource"
	"github.com/mitchellh/mapstructure"
)

type PrivateDNSMXRecord struct {
	provider *Provider

	location string
	// Usage
	// receive monthly number of DNS queries
	monthlyQueries *int64
}

type privateDNSMXRecordValues struct {
	ResourceGroupName ResourceGroupNameStruct `mapstructure:"resource_group_name"`

	Usage struct {
		MonthlyQueries int64 `mapstructure:"monthly_queries"`
	} `mapstructure:"pennywise_usage"`
}

func (p *Provider) newPrivateDNSMXRecord(vals privateDNSMXRecordValues) *PrivateDNSMXRecord {
	inst := &PrivateDNSMXRecord{
		provider:       p,
		location:       vals.ResourceGroupName.Values.Location,
		monthlyQueries: &vals.Usage.MonthlyQueries,
	}
	return inst
}

func decoderPrivateDnsMXRecord(tfVals map[string]interface{}) (privateDNSMXRecordValues, error) {
	var v privateDNSMXRecordValues
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

func (inst *PrivateDNSMXRecord) component() []resource.Component {
	region := getLocationName(inst.location)
	costComponents := privateDNSMXRecord(inst.provider.key, region, inst.monthlyQueries)
	GetCostComponentNamesAndSetLogger(costComponents, inst.provider.logger)
	return costComponents
}

func privateDNSMXRecord(key, region string, monthlyQueries *int64) []resource.Component {
	return DNSQueriesCostComponent(key, region, monthlyQueries)
}
