package tests

import (
	"fmt"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestFunctionWindowsApp() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Functions", "eastus")
	fmt.Println("Functions ingested")

	usg, err := ts.getUsage("../../testdata/azure/function_windows_app/usage.yml")
	require.NoError(ts.T(), err)

	state := ts.getDirCosts("../../testdata/azure/function_windows_app", *usg)
	stateCost, err := state.Cost()
	require.NoError(ts.T(), err)
	ts.Equal(13366.00129344, stateCost.Decimal.InexactFloat64())
}
