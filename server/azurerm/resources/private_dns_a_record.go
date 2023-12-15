package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
)

type PrivateDNSARecord struct {
	provider *Provider

	location string
	// Usage
	monthlyQueries *int64
}

type privateDNSARecordValues struct {
	ResourceGroupName ResourceGroupNameStruct `mapstructure:"resource_group_name"`

	Usage struct {
		MonthlyQueries int64 `mapstructure:"monthly_queries"`
	} `mapstructure:"tc_usage"`
}

func (p *Provider) newPrivateDnsARecord(vals privateDNSARecordValues) *PrivateDNSARecord {
	inst := &PrivateDNSARecord{
		location:       vals.ResourceGroupName.Values.Location,
		monthlyQueries: &vals.Usage.MonthlyQueries,
	}
	return inst
}

func decoderPrivateDnsARecord(tfVals map[string]interface{}) (privateDNSARecordValues, error) {
	var v privateDNSARecordValues
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

func (inst *PrivateDNSARecord) component() []query.Component {
	region := getLocationName(inst.location)
	return privateDNSARecord(inst.provider.key, region, inst.monthlyQueries)
}

func privateDNSARecord(key string, region string, monthlyQueries *int64) []query.Component {
	return DNSQueriesCostComponent(key, region, monthlyQueries)
}
