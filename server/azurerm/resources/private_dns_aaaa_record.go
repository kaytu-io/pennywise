package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
)

type PrivateDNSAAAARecord struct {
	provider *Provider

	location       string
	monthlyQueries *int64
}

type privateDNSAAAARecordValues struct {
	Value struct {
		Location string `mapstructure:"location"`
	} `mapstructure:"value"`
	MonthlyQueries int64 `mapstructure:"monthly_queries"`
}

func (p *Provider) newprivateDNSAAAARecord(vals privateDNSAAAARecordValues) *PrivateDNSAAAARecord {
	inst := &PrivateDNSAAAARecord{
		location:       vals.Value.Location,
		monthlyQueries: &vals.MonthlyQueries,
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
	return privateDNSAAAARecord(inst.location, inst.monthlyQueries)
}

func privateDNSAAAARecord(region string, monthlyQueries *int64) []query.Component {
	return DNSQueriesCostComponent(region, monthlyQueries)
}
