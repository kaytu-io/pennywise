package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
)

type PrivateDNSPTRRecord struct {
	provider *Provider

	location string
	// Usage
	monthlyQueries *int64
}

type privateDNSPTRRecordValues struct {
	Location string `mapstructure:"location"`

	Usage struct {
		MonthlyQueries int64 `mapstructure:"monthly_queries"`
	} `mapstructure:"tc_usage"`
}

func (p *Provider) newPrivateDNSPTRRecord(vals privateDNSPTRRecordValues) *PrivateDNSPTRRecord {
	inst := &PrivateDNSPTRRecord{
		location:       vals.Location,
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
	return privateDNSPTRRecord(inst.location, inst.monthlyQueries)
}

func privateDNSPTRRecord(region string, monthlyQueries *int64) []query.Component {
	return DNSQueriesCostComponent(region, monthlyQueries)
}
