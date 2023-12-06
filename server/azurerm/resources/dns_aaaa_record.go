package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
)

type DNSAAAARecord struct {
	provider *Provider

	location string

	// Usage
	monthlyQueries *int64
}

type dnsAAAARecordValues struct {
	Location string `mapstructure:"location"`

	Usage struct {
		MonthlyQueries int64 `mapstructure:"monthly_queries"`
	} `mapstructure:"tc_usage"`
}

func decoderDNSAAAARecord(tfVals map[string]interface{}) (dnsAAAARecordValues, error) {
	var v dnsAAAARecordValues
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

func (p *Provider) newDNSAAAARecord(vals dnsAAAARecordValues) *DNSAAAARecord {
	inst := &DNSAAAARecord{
		location:       vals.Location,
		monthlyQueries: &vals.Usage.MonthlyQueries,
	}
	return inst
}

func (inst *DNSAAAARecord) component() []query.Component {
	return DNSQueriesCostComponent(inst.location, inst.monthlyQueries)
}