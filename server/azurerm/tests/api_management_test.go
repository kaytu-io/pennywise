package tests

import (
	"fmt"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestApiManagement() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")

	ts.IngestService("API Management", "eastus")
	fmt.Println("API Management data ingested")

	usg, err := ts.getUsage("../../testdata/azure/api_management/usage.yml")
	require.NoError(ts.T(), err)

	state := ts.getDirCosts("../../testdata/azure/api_management", *usg)
	stateCost, err := state.Cost()
	require.NoError(ts.T(), err)

	ts.Equal(24798.863, stateCost.InexactFloat64())
}
