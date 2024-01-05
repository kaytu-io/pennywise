package resources

import (
	"fmt"
	"github.com/kaytu-io/infracost/external/usage"
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
	"github.com/kaytu-io/pennywise/server/internal/query"
	"github.com/kaytu-io/pennywise/server/internal/util"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

type SearchService struct {
	provider       *Provider
	location       string
	sku            string
	partitionCount *int64
	replicaCount   *int64
	// usage
	monthlyImagesExtracted *int64 `mapstructure:"monthly_images_extracted"`
}

type SearchServiceValue struct {
	Location       string `mapstructure:"location"`
	Sku            string `mapstructure:"sku"`
	PartitionCount *int64 `mapstructure:"partition_count"`
	ReplicaCount   *int64 `mapstructure:"replica_count"`

	Usage struct {
		MonthlyImagesExtracted *int64 `mapstructure:"monthly_images_extracted"`
	} `mapstructure:"pennywise_usage"`
}

func decodesearchService(tfVals map[string]interface{}) (SearchServiceValue, error) {
	var v SearchServiceValue
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &v,
	}
	fmt.Printf("tfvalue : %v \n ", tfVals)
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return v, err
	}

	if err := decoder.Decode(tfVals); err != nil {
		return v, err
	}
	return v, nil
}

func (p *Provider) newSearchService(value SearchServiceValue) *SearchService {
	inst := &SearchService{
		provider:               p,
		location:               value.Location,
		sku:                    value.Sku,
		partitionCount:         value.PartitionCount,
		replicaCount:           value.ReplicaCount,
		monthlyImagesExtracted: value.Usage.MonthlyImagesExtracted,
	}
	return inst
}

func (inst SearchService) Component() []query.Component {
	region := getLocationName(inst.location)

	var costComponents []query.Component

	sku := inst.sku

	if sku == "free" {
		return nil
	}

	if strings.HasPrefix(sku, "standard") {
		sku = sku[:len(sku)-1] + " s" + sku[len(sku)-1:]
	}
	if strings.HasPrefix(sku, "storage") {
		sku = strings.ReplaceAll(sku, "_", " ")
	}
	sku = strings.Title(sku)

	partitionCount := decimal.NewFromInt(1)
	replicaCount := decimal.NewFromInt(1)

	if inst.partitionCount != nil {
		partitionCount = decimal.NewFromInt(*inst.partitionCount)
	}
	if inst.replicaCount != nil {
		replicaCount = decimal.NewFromInt(*inst.replicaCount)
	}
	units := decimalPtr(partitionCount.Mul(replicaCount))

	var skuName string
	skuElems := strings.Split(sku, " ")
	for _, v := range skuElems {
		skuName += cases.Title(language.English).String(v) + " "
	}
	unitName := "unit"
	if units.GreaterThan(decimal.NewFromInt(1)) {
		unitName += "s"
	}
	HourlyQu := decimal.Zero
	if units != nil {
		HourlyQu = *units
	}
	costComponents = append(costComponents, query.Component{
		Name:           fmt.Sprintf("Search usage (%s, %s %s)", skuName[:len(skuName)-1], units.String(), unitName),
		Unit:           "hours",
		HourlyQuantity: HourlyQu,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr("azurerm"),
			Location: util.StringPtr(region),
			Service:  util.StringPtr("Azure Cognitive Search"),
			Family:   util.StringPtr("Web"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "sku_name", Value: util.StringPtr(fmt.Sprintf("%s", sku))},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
			},
		},
	})

	var images *decimal.Decimal
	if inst.monthlyImagesExtracted != nil {
		images = decimalPtr(decimal.NewFromInt(*inst.monthlyImagesExtracted))
		tierLimits := []int{1_000_000, 4_000_000}
		tiers := usage.CalculateTierBuckets(*images, tierLimits)

		type dataTier struct {
			name       string
			startUsage string
		}

		data := []dataTier{
			{name: "Image extraction (first 1M)", startUsage: "0"},
			{name: "Image extraction (next 4M)", startUsage: "1000"},
			{name: "Image extraction (over 5M)", startUsage: "5000"},
		}
		for i, d := range data {
			if tiers[i].GreaterThan(decimal.Zero) {
				costComponents = append(costComponents, searchServiceCostComponent(
					region,
					d.name,
					d.startUsage,
					tiers[i].Div(decimal.NewFromInt(1000))))
			}
		}
	} else {
		monthlyQu := decimal.Zero
		if images != nil {
			monthlyQu = *images
		}

		costComponents = append(costComponents, searchServiceCostComponent(
			region,
			"Image extraction (first 1M)",
			"0",
			monthlyQu))
	}

	return costComponents
}

func searchServiceCostComponent(region, name, startUsage string, qty decimal.Decimal) query.Component {
	return query.Component{
		Name:            name,
		Unit:            "1000 images",
		MonthlyQuantity: qty,
		ProductFilter: &product.Filter{
			Provider: util.StringPtr("azurerm"),
			Location: util.StringPtr(region),
			Service:  util.StringPtr("Azure Cognitive Search"),
			Family:   util.StringPtr("Web"),
			AttributeFilters: []*product.AttributeFilter{
				{Key: "sku_name", Value: util.StringPtr("Document Cracking")},
			},
		},
		PriceFilter: &price.Filter{
			AttributeFilters: []*price.AttributeFilter{
				{Key: "type", Value: util.StringPtr("Consumption")},
				{Key: "tier_minimum_units", Value: util.StringPtr(startUsage)},
			},
		},
	}
}
