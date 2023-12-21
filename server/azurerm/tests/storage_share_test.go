package tests

import (
	"fmt"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestStorageShare() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Storage", "westus")
	fmt.Println("Storage ingested")

	usg, err := ts.getUsage("../../testdata/azure/storage_share/usage.yml")
	require.NoError(ts.T(), err)

	state := ts.getDirCosts("../../testdata/azure/storage_share", *usg)
	cost, err := state.Cost()
	require.NoError(ts.T(), err)
	ts.Equal(946944.866, cost.Decimal.InexactFloat64())
}
