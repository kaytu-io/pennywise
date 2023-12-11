package ingester

import (
	"context"
	"fmt"
	"github.com/kaytu-io/pennywise/server/internal/backend"
	"github.com/kaytu-io/pennywise/server/internal/price"
	"github.com/kaytu-io/pennywise/server/internal/product"
)

//go:generate mockgen -destination=mock/ingester.go -mock_names=Ingester=Ingester -package mock github.com/kaytu-io/pennywise/server Ingester

// Ingester represents a vendor-specific mechanism to load pricing data.
type Ingester interface {
	// Ingest downloads pricing data from a cloud provider and sends prices with their associated products
	// on the returned channel.
	Ingest(ctx context.Context, chSize int) <-chan *price.WithProduct

	// Err returns any potential error.
	Err() error
}

// IngestPricing uses the Ingester to load the pricing data and stores it into the Backend.
func IngestPricing(ctx context.Context, be backend.Backend, ingester Ingester) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	skuProductID := make(map[string]product.ID)
	for pp := range ingester.Ingest(ctx, 8) {
		if id, ok := skuProductID[pp.Product.SKU+pp.Product.MeterId]; ok {
			pp.Product.ID = id
		} else {
			var err error
			pp.Product.ID, err = be.Products().Upsert(ctx, pp.Product)
			if err != nil {
				return fmt.Errorf("failed to upsert product (SKU=%q): %w", pp.Product.SKU, err)
			}
			skuProductID[pp.Product.SKU+pp.Product.MeterId] = pp.Product.ID
		}

		if _, err := be.Prices().Upsert(ctx, pp); err != nil {
			return fmt.Errorf("failed to upsert price (SKU=%q): %w", pp.Product.SKU, err)
		}
	}

	if err := ingester.Err(); err != nil {
		return fmt.Errorf("unexpected ingester error: %w", err)
	}
	return nil
}