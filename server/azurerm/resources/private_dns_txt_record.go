package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
)

type PrivateDNSTXTRecord struct {
	provider *Provider

	location       string
	monthlyQueries *int64
}

type privateDNSTXTRecordValues struct {
	Value struct {
		Location string `mapstructure:"location"`
	} `mapstructure:"value"`
	MonthlyQueries int64 `mapstructure:"monthly_queries"`
}

func (p *Provider) newPrivateDNSTXTRecord(vals privateDNSTXTRecordValues) *PrivateDNSTXTRecord {
	inst := &PrivateDNSTXTRecord{
		location:       vals.Value.Location,
		monthlyQueries: &vals.MonthlyQueries,
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
	return privateDNSTXTRecord(inst.location, inst.monthlyQueries)
}

func privateDNSTXTRecord(region string, monthlyQueries *int64) []query.Component {
	return DNSQueriesCostComponent(region, monthlyQueries)
}
