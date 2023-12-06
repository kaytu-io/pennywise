package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
)

type DNSTEXTRecord struct {
	provider *Provider

	location       string
	monthlyQueries *int64
}

type dnsTEXTRecordValues struct {
	Location string `mapstructure:"location"`
	Usage    struct {
		MonthlyQueries int64 `mapstructure:"monthly_queries"`
	} `mapstructure:"tc_usage"`
}

func decoderDNSTXTRecord(tfVals map[string]interface{}) (dnsTEXTRecordValues, error) {
	var v dnsTEXTRecordValues
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

func (p *Provider) newDNSTXTRecord(vals dnsTEXTRecordValues) *DNSTEXTRecord {
	inst := &DNSTEXTRecord{
		location:       vals.Location,
		monthlyQueries: &vals.Usage.MonthlyQueries,
	}
	return inst
}

func (inst *DNSTEXTRecord) component() []query.Component {
	return DNSQueriesCostComponent(inst.location, inst.monthlyQueries)
}
