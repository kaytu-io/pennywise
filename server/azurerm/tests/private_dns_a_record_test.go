package tests

import (
	"fmt"
)

func (ts *AzureTestSuite) TestPrivateDNSARecord() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Azure DNS", "West Europe")
	fmt.Println("Azure DNS data ingested")

	//usg, err := ts.getUsage("../../testdata/azure/private_dns_a_record/usage.json")
	//require.NoError(ts.T(), err)

	stat := ts.getDirCosts("../../testdata/azure/private_dns_a_record", nil)
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
