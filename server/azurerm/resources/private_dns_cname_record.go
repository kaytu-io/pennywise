package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
)

type PrivateDNSCNAMERecord struct {
	provider *Provider

	location string

	// Usage
	// receive monthly number of DNS queries
	monthlyQueries *int64
}

type privateDNSCNAMERecordValues struct {
	ResourceGroupName ResourceGroupNameStruct `mapstructure:"resource_group_name"`

	Usage struct {
		MonthlyQueries int64 `mapstructure:"monthly_queries"`
	} `mapstructure:"pennywise_usage"`
}

func (p *Provider) newprivateDNSCNAMERecord(vals privateDNSCNAMERecordValues) *PrivateDNSCNAMERecord {
	inst := &PrivateDNSCNAMERecord{
		provider:       p,
		location:       vals.ResourceGroupName.Values.Location,
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
	region := getLocationName(inst.location)
	costComponents := privateDNSCNAMERecord(inst.provider.key, region, inst.monthlyQueries)
	GetCostComponentNamesAndSetLogger(costComponents, inst.provider.logger)
	return costComponents
}

func privateDNSCNAMERecord(key, region string, monthlyQueries *int64) []query.Component {
	return DNSQueriesCostComponent(key, region, monthlyQueries)
}
