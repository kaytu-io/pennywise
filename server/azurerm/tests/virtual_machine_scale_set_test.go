package tests

import (
	"fmt"
	"github.com/stretchr/testify/require"
)

func (ts *AzureTestSuite) TestVirtualMachineScaleSet() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	//ts.IngestService("Virtual Machines", "West Europe")
	//fmt.Println("Virtual Machine data ingested")

	ts.IngestService("Storage", "West Europe")
	fmt.Println("Storage data ingested")

	usg, err := ts.getUsage("../../testdata/azure/virtual_machine_scale_set/usage.yml")
	require.NoError(ts.T(), err)

	stat := ts.getDirCosts("../../testdata/azure/virtual_machine_scale_set", *usg)
	costComponent := stat.GetCostComponents()
	for k, v := range costComponent {
		fmt.Printf("cost component : %v \n", k)
		fmt.Printf("name : %v \n ", v.Name)
		fmt.Printf("unit : %v \n ", v.Unit)
		fmt.Printf("rate : %v \n ", v.Rate)
		fmt.Printf("Details : %v \n ", v.Details)
		fmt.Printf("Usage : %v \n ", v.Usage)
		fmt.Printf("MonthlyQuantity : %v \n ", v.MonthlyQuantity)
		fmt.Printf("HourlyQuantity : %v \n ", v.HourlyQuantity)
		fmt.Printf("Error : %v \n ", v.Error)
		fmt.Printf("\n")
	}
}
