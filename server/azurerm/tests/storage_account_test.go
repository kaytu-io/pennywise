package tests

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestStorageAccount() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Storage", "westus")
	fmt.Println("Storage ingested")

	usg, err := ts.getUsage("../../testdata/azure/storage_account/usage.yml")
	require.NoError(ts.T(), err)

	state := ts.getDirCosts("../../testdata/azure/storage_account", *usg)
	costComponents := state.GetCostComponents()
	expectedCostComponents := []cost.Component{}

	ts.Equal(len(expectedCostComponents), len(costComponents))
	for _, comp := range expectedCostComponents {
		ts.True(componentExists(comp, costComponents), fmt.Sprintf("Could not match component %s: %v", comp.Name, comp))
	}
	fmt.Println(costComponents)
	fmt.Println(state.CostString())
}
