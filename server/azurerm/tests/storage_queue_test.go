package tests

import (
	"fmt"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestStorageQueue() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Storage", "westus")
	fmt.Println("Storage ingested")

	usg, err := ts.getUsage("../../testdata/azure/storage_queue/usage.yaml")
	require.NoError(ts.T(), err)

	state := ts.getDirCosts("../../testdata/azure/storage_queue", *usg)
	cost, err := state.Cost()
	require.NoError(ts.T(), err)
	ts.Equal(1802.432, cost.Decimal.InexactFloat64())
}
