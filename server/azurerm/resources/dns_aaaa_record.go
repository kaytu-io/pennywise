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
	ResourceGroupName ResourceGroupNameStruct `mapstructure:"resource_group_name"`

	Usage struct {
		// receives monthly number of DNS queries
		MonthlyQueries int64 `mapstructure:"monthly_queries"`
	} `mapstructure:"pennywise_usage"`
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
		provider:       p,
		location:       vals.ResourceGroupName.Values.Location,
		monthlyQueries: &vals.Usage.MonthlyQueries,
	}
	return inst
}

func (inst *DNSAAAARecord) component() []query.Component {
	region := getLocationName(inst.location)
	return DNSQueriesCostComponent(inst.provider.key, region, inst.monthlyQueries)
}
