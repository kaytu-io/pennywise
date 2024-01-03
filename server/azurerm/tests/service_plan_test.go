package tests

import (
	"fmt"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestServicePlan() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")

	ts.IngestService("Azure App Service", "")
	fmt.Println("Container Registry data ingested")

	state := ts.getDirCosts("../../testdata/azure/service_plan", nil)
	stateCost, err := state.Cost()
	require.NoError(ts.T(), err)
	ts.Equal(611320.98, stateCost.Decimal.InexactFloat64())
}
