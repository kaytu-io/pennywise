package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
)

type DNSSRVRecord struct {
	provider *Provider

	location string
	// Usage
	monthlyQueries *int64
}

type dnsSRVRecordValues struct {
	ResourceGroupName ResourceGroupNameStruct `mapstructure:"resource_group_name"`

	Usage struct {
		MonthlyQueries int64 `mapstructure:"monthly_queries"`
	} `mapstructure:"tc_usage"`
}

func decoderDNSSRVRecord(tfVals map[string]interface{}) (dnsSRVRecordValues, error) {
	var v dnsSRVRecordValues
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

func (p *Provider) newDNSSRVRecord(vals dnsSRVRecordValues) *DNSSRVRecord {
	inst := &DNSSRVRecord{
		provider:       p,
		location:       vals.ResourceGroupName.Values.Location,
		monthlyQueries: &vals.Usage.MonthlyQueries,
	}
	return inst
}

func (inst *DNSSRVRecord) component() []query.Component {
	region := getLocationName(inst.location)
	return DNSQueriesCostComponent(inst.provider.key, region, inst.monthlyQueries)
}
