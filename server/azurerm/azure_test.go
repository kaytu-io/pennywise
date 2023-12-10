package azurerm

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kaytu-io/pennywise/cli/parser/hcl"
	"github.com/kaytu-io/pennywise/cli/usage"
	resources2 "github.com/kaytu-io/pennywise/server/azurerm/resources"
	"github.com/kaytu-io/pennywise/server/cost"
	ingester2 "github.com/kaytu-io/pennywise/server/internal/ingester"
	"github.com/kaytu-io/pennywise/server/internal/mysql"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/resource"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
	"os"
	"testing"
)

var (
	MySQLHost     = os.Getenv("MYSQL_HOST")
	MySQLPort     = os.Getenv("MYSQL_PORT")
	MySQLDb       = os.Getenv("MYSQL_DB")
	MySQLUser     = os.Getenv("MYSQL_USERNAME")
	MySQLPassword = os.Getenv("MYSQL_PASSWORD")
)

type AzureTestSuite struct {
	suite.Suite

	backend *mysql.Backend
}

func TestAzure(t *testing.T) {
	suite.Run(t, &AzureTestSuite{})
}

func (ts *AzureTestSuite) SetupSuite() {
	//db, _, err := sqlmock.New()
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true", MySQLUser, MySQLPassword, MySQLHost, MySQLPort, MySQLDb)
	db, err := sql.Open("mysql", dataSource)

	require.NoError(ts.T(), err)
	err = mysql.Migrate(context.Background(), db, "pricing_migrations")

	ts.backend = mysql.NewBackend(db)
}

func (ts *AzureTestSuite) IngestService(service, region string) {
	ingester, err := NewIngester(service, region)
	require.NoError(ts.T(), err)

	err = ingester2.IngestPricing(context.Background(), ts.backend, ingester)
	require.NoError(ts.T(), err)

}

func (ts *AzureTestSuite) getDirCosts(projectDir string, usage usage.Usage) *cost.State {
	providerName, hclResources, err := hcl.ParseHclResources(projectDir, usage)
	require.NoError(ts.T(), err)

	var qResources []query.Resource
	resources := make(map[string]resource.Resource)
	provider, err := resources2.NewProvider(resources2.ProviderName)
	require.NoError(ts.T(), err)

	for _, rs := range hclResources {
		res := rs.ToResource(providerName, nil)
		resources[res.Address] = res
	}

	for _, res := range resources {
		components := provider.ResourceComponents(resources, res)
		qResource := query.Resource{
			Address:    res.Address,
			Provider:   res.ProviderName,
			Type:       res.Type,
			Components: components,
		}
		qResources = append(qResources, qResource)
	}

	state, err := cost.NewState(context.Background(), ts.backend, qResources)
	require.NoError(ts.T(), err)

	return state
}

func (ts *AzureTestSuite) TestLoadBalancer() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Load Balancer", "")
	fmt.Println("Load Balancer data ingested")

	lbUsage := usage.Usage{"azurerm_lb": map[string]interface{}{
		"monthly_data_proceed": 1000,
	}}
	cost := ts.getDirCosts("../testdata/azure/load_balancer", lbUsage)
	fmt.Println(cost.CostString())

}

func (ts *AzureTestSuite) TestPublicIp() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Virtual Network", "")
	fmt.Println("Virtual Network data ingested")

	usg := usage.Usage{}
	cost := ts.getDirCosts("../testdata/azure/public_ip", usg)
	fmt.Println(cost.CostString())

}

func (ts *AzureTestSuite) TestPublicIpPrefix() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Virtual Network", "")
	fmt.Println("Virtual Network data ingested")

	usg := usage.Usage{}
	cost := ts.getDirCosts("../testdata/azure/public_ip_prefix", usg)
	fmt.Println(cost.CostString())

}

func (ts *AzureTestSuite) TestPrivateEndpoint() {
	ts.SetupSuite()
	fmt.Println("Suite Setup")
	ts.IngestService("Virtual Network", "")
	fmt.Println("Virtual Network data ingested")

	usg := usage.Usage{"azurerm_private_endpoint": map[string]interface{}{
		"monthly_inbound_data_processed_gb":  100,
		"monthly_outbound_data_processed_gb": 100,
	}}
	cost := ts.getDirCosts("../testdata/azure/private_endpoint", usg)
	fmt.Println(cost.CostString())

}