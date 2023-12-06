package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
)

type DNSMXRecord struct {
	provider *Provider

	location string
	// Usage
	monthlyQueries *int64
}

type dnsMAXRecordValues struct {
	Location string `mapstructure:"location"`
	Usage    struct {
		MonthlyQueries int64 `mapstructure:"monthly_queries"`
	} `mapstructure:"tc_usage"`
}

func decoderDNSMXRecord(tfVals map[string]interface{}) (dnsMAXRecordValues, error) {
	var v dnsMAXRecordValues
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

func (p *Provider) newDNSMXRecord(vals dnsMAXRecordValues) *DNSMXRecord {
	inst := &DNSMXRecord{
		location:       vals.Location,
		monthlyQueries: &vals.Usage.MonthlyQueries,
	}
	return inst
}

func (inst *DNSMXRecord) component() []query.Component {
	return DNSQueriesCostComponent(inst.location, inst.monthlyQueries)
}