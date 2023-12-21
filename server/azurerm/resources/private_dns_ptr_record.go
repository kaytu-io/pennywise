package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
)

type PrivateDNSPTRRecord struct {
	provider *Provider

	location string
	// Usage
	// receive monthly number of DNS queries
	monthlyQueries *int64
}

type privateDNSPTRRecordValues struct {
	ResourceGroupName ResourceGroupNameStruct `mapstructure:"resource_group_name"`

	Usage struct {
		MonthlyQueries int64 `mapstructure:"monthly_queries"`
	} `mapstructure:"pennywise_usage"`
}

func (p *Provider) newPrivateDNSPTRRecord(vals privateDNSPTRRecordValues) *PrivateDNSPTRRecord {
	inst := &PrivateDNSPTRRecord{
		provider:       p,
		location:       vals.ResourceGroupName.Values.Location,
		monthlyQueries: &vals.Usage.MonthlyQueries,
	}
	return inst
}

func decoderPrivateDnsPTRRecord(tfVals map[string]interface{}) (privateDNSPTRRecordValues, error) {
	var v privateDNSPTRRecordValues
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

func (inst *PrivateDNSPTRRecord) component() []query.Component {
	region := getLocationName(inst.location)
	return privateDNSPTRRecord(inst.provider.key, region, inst.monthlyQueries)
}

func privateDNSPTRRecord(key, region string, monthlyQueries *int64) []query.Component {
	return DNSQueriesCostComponent(key, region, monthlyQueries)
}
