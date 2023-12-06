package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
)

type PrivateDNSMXRecord struct {
	provider *Provider

	location       string
	monthlyQueries *int64
}

type privateDNSMXRecordValues struct {
	Value struct {
		Location string `mapstructure:"location"`
	} `mapstructure:"value"`
	MonthlyQueries int64 `mapstructure:"monthly_queries"`
}

func (p *Provider) newPrivateDNSMXRecord(vals privateDNSMXRecordValues) *PrivateDNSMXRecord {
	inst := &PrivateDNSMXRecord{
		location:       vals.Value.Location,
		monthlyQueries: &vals.MonthlyQueries,
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

func (inst *PrivateDNSMXRecord) component() []query.Component {
	return privateDNSMXRecord(inst.location, inst.monthlyQueries)
}

func privateDNSMXRecord(region string, monthlyQueries *int64) []query.Component {
	return DNSQueriesCostComponent(region, monthlyQueries)
}
