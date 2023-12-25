package tests

import (
	"fmt"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestMssqlDatabase() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("SQL Database", "")
	fmt.Println("SQL Database ingested")

	usg, err := ts.getUsage("../../testdata/azure/mssql_database/usage.yaml")
	require.NoError(ts.T(), err)

	state := ts.getDirCosts("../../testdata/azure/mssql_database", *usg)
	cost, err := state.Cost()
	require.NoError(ts.T(), err)
	ts.Equal(31585.371, cost.Decimal.InexactFloat64())
}
