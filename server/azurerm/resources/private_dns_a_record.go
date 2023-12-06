package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
)

type PrivateDNSARecord struct {
	provider *Provider

	location       string
	monthlyQueries *int64
}

type privateDNSARecordValues struct {
	Value struct {
		Location string `mapstructure:"location"`
	} `mapstructure:"value"`
	MonthlyQueries int64 `mapstructure:"monthly_queries"`
}

func (p *Provider) newPrivateDnsARecord(vals privateDNSARecordValues) *PrivateDNSARecord {
	inst := &PrivateDNSARecord{
		location:       vals.Value.Location,
		monthlyQueries: &vals.MonthlyQueries,
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
	return privateDNSARecord(inst.location, inst.monthlyQueries)
}

func privateDNSARecord(region string, monthlyQueries *int64) []query.Component {
	return DNSQueriesCostComponent(region, monthlyQueries)
}
