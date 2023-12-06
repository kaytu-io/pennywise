package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
)

type DNSNSRecord struct {
	provider *Provider

	location string
	// Usage
	monthlyQueries *int64
}

type dnsNSRecordValues struct {
	Location string `mapstructure:"location"`
	Usage    struct {
		MonthlyQueries int64 `mapstructure:"monthly_queries"`
	} `mapstructure:"tc_usage"`
}

func decoderDNSNSRecord(tfVals map[string]interface{}) (dnsNSRecordValues, error) {
	var v dnsNSRecordValues
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

func (p *Provider) newDNSNSRecord(vals dnsNSRecordValues) *DNSNSRecord {
	inst := &DNSNSRecord{
		location:       vals.Location,
		monthlyQueries: &vals.Usage.MonthlyQueries,
	}
	return inst
}

func (inst *DNSNSRecord) component() []query.Component {
	return DNSQueriesCostComponent(inst.location, inst.monthlyQueries)
}
