package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
)

type PrivateDNSAAAARecord struct {
	provider *Provider

	location string

	// Usage
	monthlyQueries *int64
}

type privateDNSAAAARecordValues struct {
	ResourceGroupName ResourceGroupNameStruct `mapstructure:"resource_group_name"`

	Usage struct {
		MonthlyQueries int64 `mapstructure:"monthly_queries"`
	} `mapstructure:"tc_usage"`
}

func (p *Provider) newprivateDNSAAAARecord(vals privateDNSAAAARecordValues) *PrivateDNSAAAARecord {
	inst := &PrivateDNSAAAARecord{

		provider:       p,
		location:       vals.ResourceGroupName.Values.Location,
		monthlyQueries: &vals.Usage.MonthlyQueries,
	}
	return inst
}

func decoderPrivateDnsAAAARecord(tfVals map[string]interface{}) (privateDNSAAAARecordValues, error) {
	var v privateDNSAAAARecordValues
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

func (inst *PrivateDNSAAAARecord) component() []query.Component {
	region := getLocationName(inst.location)
	return privateDNSAAAARecord(inst.provider.key, region, inst.monthlyQueries)
}

func privateDNSAAAARecord(key, region string, monthlyQueries *int64) []query.Component {
	return DNSQueriesCostComponent(key, region, monthlyQueries)
}
