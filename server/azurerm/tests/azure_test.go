package tests

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kaytu-io/pennywise/cli/parser/hcl"
	"github.com/kaytu-io/pennywise/cli/usage"
	"github.com/kaytu-io/pennywise/server/azurerm"
	resources2 "github.com/kaytu-io/pennywise/server/azurerm/resources"
	"github.com/kaytu-io/pennywise/server/cost"
	ingester2 "github.com/kaytu-io/pennywise/server/internal/ingester"
	"github.com/kaytu-io/pennywise/server/internal/mysql"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/resource"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
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
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true", MySQLUser, MySQLPassword, MySQLHost, MySQLPort, MySQLDb)
	db, err := sql.Open("mysql", dataSource)
	require.NoError(ts.T(), err)
	err = mysql.Migrate(context.Background(), db, "pricing_migrations")

	ts.backend = mysql.NewBackend(db)
}

func (ts *AzureTestSuite) IngestService(service, region string) {
	ingester, err := azurerm.NewIngester(service, region)
	require.NoError(ts.T(), err)

	err = ingester2.IngestPricing(context.Background(), ts.backend, ingester)
	require.NoError(ts.T(), err)

}

func (ts *AzureTestSuite) getUsage(usagePath string) (*usage.Usage, error) {
	var usg usage.Usage
	if usagePath != "" {
		usageFile, err := os.Open(usagePath)
		if err != nil {
			return nil, fmt.Errorf("error while reading usage file %s", err)
		}
		defer usageFile.Close()

		ext := filepath.Ext(usagePath)
		switch ext {
		case ".json":
			err = json.NewDecoder(usageFile).Decode(&usg)
		case ".yaml", ".yml":
			err = yaml.NewDecoder(usageFile).Decode(&usg)
		default:
			return nil, fmt.Errorf("unsupported file format %s for usage file", ext)
		}
		if err != nil {
			return nil, fmt.Errorf("error while parsing usage file %s", err)
		}

	} else {
		usg = usage.Default
	}
	return &usg, nil
}

func (ts *AzureTestSuite) getDirCosts(projectDir string, usg usage.Usage) *cost.State {
	providerName, hclResources, err := hcl.ParseHclResources(projectDir, usg)
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

func checkComponents(result, expected cost.Component) bool {
	if result.Name == expected.Name && result.MonthlyQuantity.Equal(expected.MonthlyQuantity) &&
		result.HourlyQuantity.Equal(expected.HourlyQuantity) && result.Unit == expected.Unit && result.Rate.Decimal.Equal(expected.Rate.Decimal) &&
		result.Rate.Currency == expected.Rate.Currency && result.Usage == expected.Usage && result.Error == expected.Error {
		return true
	} else {
		return false
	}
}

func componentExists(component cost.Component, comps []cost.Component) bool {
	for _, comp := range comps {
		if checkComponents(comp, component) {
			return true
		}
	}
	return false
}
