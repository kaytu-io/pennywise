package tier_request

import (
	"github.com/shopspring/decimal"
)

func CalculateTierBuckets(requestCount decimal.Decimal, tierLimits []int) []decimal.Decimal {
	overTier := false
	tiers := make([]decimal.Decimal, 0)

	for limit := range tierLimits {
		tier := decimal.NewFromInt(int64(tierLimits[limit]))

		if requestCount.GreaterThanOrEqual(tier) {
			tiers = append(tiers, tier)
			requestCount = requestCount.Sub(tier)
			overTier = true
		} else if requestCount.LessThan(tier) {
			tiers = append(tiers, requestCount)
			requestCount = decimal.Zero
			overTier = false
		}
	}

	if overTier {
		tiers = append(tiers, requestCount)
	} else {
		tiers = append(tiers, decimal.NewFromInt(0))
	}
	return tiers
}
