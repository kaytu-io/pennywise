package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
)

type DNSSRVRecord struct {
	provider *Provider

	location       string
	monthlyQueries *int64
}

type dnsSRVRecordValues struct {
	Location string `mapstructure:"location"`
	// TODO:we should get MonthlyQueries field from user
	MonthlyQueries int64 `mapstructure:"monthly_queries"`
}

func decoderDNSSRVRecord(tfVals map[string]interface{}) (dnsSRVRecordValues, error) {
	var v dnsSRVRecordValues
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

func (p *Provider) newDNSSRVRecord(vals dnsSRVRecordValues) *DNSSRVRecord {
	inst := &DNSSRVRecord{
		location:       vals.Location,
		monthlyQueries: &vals.MonthlyQueries,
	}
	return inst
}

func (inst *DNSSRVRecord) component() []query.Component {
	return DNSQueriesCostComponent(inst.location, inst.monthlyQueries)
}
