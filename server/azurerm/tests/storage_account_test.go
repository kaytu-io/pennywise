package tests

import (
	"fmt"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestStorageAccount() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Storage", "eastus")
	fmt.Println("Storage ingested")

	usg, err := ts.getUsage("../../testdata/azure/storage_account/usage.yaml")
	require.NoError(ts.T(), err)

	state := ts.getDirCosts("../../testdata/azure/storage_account", *usg)
	cost, err := state.Cost()
	require.NoError(ts.T(), err)
	ts.Equal(12483120.2124, cost.Decimal.InexactFloat64())
}
