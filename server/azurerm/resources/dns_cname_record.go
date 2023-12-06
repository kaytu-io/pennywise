package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
)

type DNSCNAMERecord struct {
	provider *Provider

	location       string
	monthlyQueries *int64
}

type dnsCNAMERecordValues struct {
	Location string `mapstructure:"location"`
	// TODO:we should get MonthlyQueries field from user
	MonthlyQueries int64 `mapstructure:"monthly_queries"`
}

func decoderDNSCNAMERecord(tfVals map[string]interface{}) (dnsCNAMERecordValues, error) {
	var v dnsCNAMERecordValues
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

func (p *Provider) newDNSCNAMERecord(vals dnsCNAMERecordValues) *DNSCNAMERecord {
	inst := &DNSCNAMERecord{
		location:       vals.Location,
		monthlyQueries: &vals.MonthlyQueries,
	}
	return inst
}

func (inst *DNSCNAMERecord) component() []query.Component {
	return DNSQueriesCostComponent(inst.location, inst.monthlyQueries)
}
