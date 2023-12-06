package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
)

type DNSCAARecord struct {
	provider *Provider

	location       string
	monthlyQueries *int64
}

type dnsCAARecordValues struct {
	Location string `mapstructure:"location"`
	// TODO:we should get MonthlyQueries field from user
	MonthlyQueries int64 `mapstructure:"monthly_queries"`
}

func decoderDNSCAARecord(tfVals map[string]interface{}) (dnsCAARecordValues, error) {
	var v dnsCAARecordValues
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

func (p *Provider) newDNSCAARecord(vals dnsCAARecordValues) *DNSCAARecord {
	inst := &DNSCAARecord{
		location:       vals.Location,
		monthlyQueries: &vals.MonthlyQueries,
	}
	return inst
}

func (inst *DNSCAARecord) component() []query.Component {
	return DNSQueriesCostComponent(inst.location, inst.monthlyQueries)
}
