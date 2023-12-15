package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
)

type PrivateDNSTXTRecord struct {
	provider *Provider

	location string

	// Usage
	monthlyQueries *int64
}

type privateDNSTXTRecordValues struct {
	ResourceGroupName ResourceGroupNameStruct `mapstructure:"resource_group_name"`

	Usage struct {
		MonthlyQueries int64 `mapstructure:"monthly_queries"`
	} `mapstructure:"tc_usage"`
}

func (p *Provider) newPrivateDNSTXTRecord(vals privateDNSTXTRecordValues) *PrivateDNSTXTRecord {
	inst := &PrivateDNSTXTRecord{
		provider:       p,
		location:       vals.ResourceGroupName.Values.Location,
		monthlyQueries: &vals.Usage.MonthlyQueries,
	}
	return inst
}

func decoderPrivateDnsTXTRecord(tfVals map[string]interface{}) (privateDNSTXTRecordValues, error) {
	var v privateDNSTXTRecordValues
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

func (inst *PrivateDNSTXTRecord) component() []query.Component {
	region := getLocationName(inst.location)
	return privateDNSTXTRecord(inst.provider.key, region, inst.monthlyQueries)
}

func privateDNSTXTRecord(key, region string, monthlyQueries *int64) []query.Component {
	return DNSQueriesCostComponent(key, region, monthlyQueries)
}
