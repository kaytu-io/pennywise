package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
)

type PrivateDNSCNAMERecord struct {
	provider *Provider

	location string

	// Usage
	monthlyQueries *int64
}

type privateDNSCNAMERecordValues struct {
	Location string `mapstructure:"location"`

	Usage struct {
		MonthlyQueries int64 `mapstructure:"monthly_queries"`
	} `mapstructure:"tc_usage"`
}

func (p *Provider) newprivateDNSCNAMERecord(vals privateDNSCNAMERecordValues) *PrivateDNSCNAMERecord {
	inst := &PrivateDNSCNAMERecord{
		location:       vals.Location,
		monthlyQueries: &vals.Usage.MonthlyQueries,
	}
	return inst
}

func decoderPrivateDnsCNAMERecord(tfVals map[string]interface{}) (privateDNSCNAMERecordValues, error) {
	var v privateDNSCNAMERecordValues
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

func (inst *PrivateDNSCNAMERecord) component() []query.Component {
	return privateDNSCNAMERecord(inst.location, inst.monthlyQueries)
}

func privateDNSCNAMERecord(region string, monthlyQueries *int64) []query.Component {
	return DNSQueriesCostComponent(region, monthlyQueries)
}
