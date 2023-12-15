package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
)

type DNSPTRRecord struct {
	provider *Provider

	location string
	// Usage
	monthlyQueries *int64
}

type dnsPTRRecordValues struct {
	ResourceGroupName ResourceGroupNameStruct `mapstructure:"resource_group_name"`

	Usage struct {
		MonthlyQueries int64 `mapstructure:"monthly_queries"`
	} `mapstructure:"tc_usage"`
}

func decoderDNSPTRRecord(tfVals map[string]interface{}) (dnsPTRRecordValues, error) {
	var v dnsPTRRecordValues
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

func (p *Provider) newDNSPTRRecord(vals dnsPTRRecordValues) *DNSPTRRecord {
	inst := &DNSPTRRecord{
		provider:       p,
		location:       vals.ResourceGroupName.Values.Location,
		monthlyQueries: &vals.Usage.MonthlyQueries,
	}
	return inst
}

func (inst *DNSPTRRecord) component() []query.Component {
	region := getLocationName(inst.location)
	return DNSQueriesCostComponent(inst.provider.key, region, inst.monthlyQueries)
}
