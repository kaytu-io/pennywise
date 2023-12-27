package resources

import (
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/mitchellh/mapstructure"
)

type AutomationDSCNodeConfiguration struct {
	provider                *Provider
	location                string
	nonAzureConfigNodeCount *int64
}

type AutomationDSCNodeConfigurationValue struct {
	ResourceGroupName ResourceGroupName `mapstructure:"resource_group_name"`
	Usage             struct {
		NonAzureConfigNodeCount *int64 `mapstructure:"non_azure_config_node_count"`
	} `mapstructure:"pennywise_usage"`
}

func (p *Provider) newAutomationDSCNodeConfiguration(vals AutomationDSCNodeConfigurationValue) *AutomationDSCNodeConfiguration {
	inst := &AutomationDSCNodeConfiguration{
		provider:                p,
		location:                vals.ResourceGroupName.Values.Location,
		nonAzureConfigNodeCount: vals.Usage.NonAzureConfigNodeCount,
	}
	return inst
}

func decodeAutomationDSCNodeConfiguration(tfVals map[string]interface{}) (AutomationDSCNodeConfigurationValue, error) {
	var v AutomationDSCNodeConfigurationValue
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

func (inst *AutomationDSCNodeConfiguration) Component() []query.Component {
	costComponent := automationDSCNodesCostComponent(inst.location, inst.nonAzureConfigNodeCount)

	GetCostComponentNamesAndSetLogger(costComponent, inst.provider.logger)
	return costComponent
}
