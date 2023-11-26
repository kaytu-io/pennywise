package mysql

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/internal/price"
	product2 "github.com/kaytu-io/pennywise/server/internal/product"
	"strings"
)

// Where represents the parts of a SQL WHERE clause.
type Where struct {
	conditions []string
	params     []interface{}
}

// String returns the string of the WHERE clause.
func (w *Where) String() string {
	return strings.Join(w.conditions, " AND ")
}

// Parameters returns the slice of parameters to be passed to the Exec or Query method.
func (w *Where) Parameters() []interface{} {
	return w.params
}

func (w *Where) add(condition string, params ...interface{}) {
	w.conditions = append(w.conditions, condition)
	w.params = append(w.params, params...)
}

func parseProductFilter(filter *product2.Filter) *Where {
	w := &Where{}

	if filter == nil {
		return w
	}

	type fieldMapping struct {
		key string
		val *string
	}
	equalFields := []fieldMapping{
		{key: "provider", val: filter.Provider},
		{key: "location", val: filter.Location},
		{key: "service", val: filter.Service},
		{key: "family", val: filter.Family},
		{key: "sku", val: filter.SKU},
	}

	for _, fm := range equalFields {
		if fm.val != nil {
			w.add(fmt.Sprintf("%s = ?", fm.key), *fm.val)
		}
	}

	for _, f := range filter.AttributeFilters {
		if f.Value != nil {
			w.add(fmt.Sprintf("JSON_UNQUOTE(JSON_EXTRACT(attributes, '$.%s')) = ?", f.Key), *f.Value)
		} else if f.ValueRegex != nil {
			w.add(fmt.Sprintf("JSON_UNQUOTE(JSON_EXTRACT(attributes, '$.%s')) RLIKE ?", f.Key), *f.ValueRegex)
		}
	}

	return w
}

func parsePriceFilter(filter *price.Filter, products []*product2.Product) *Where {
	w := &Where{}

	var productIds []uint32
	for _, p := range products {
		productIds = append(productIds, uint32(p.ID))
	}

	var idInterfaces []interface{}
	for _, id := range productIds {
		idInterfaces = append(idInterfaces, id)
	}

	if len(idInterfaces) > 0 {
		placeholders := make([]string, len(idInterfaces))
		for i := range idInterfaces {
			placeholders[i] = "?"
		}
		inClause := fmt.Sprintf("product_id IN (%s)", strings.Join(placeholders, ","))
		w.add(inClause, idInterfaces...)
	}

	if filter == nil {
		return w
	}

	type fieldMapping struct {
		key string
		val *string
	}
	equalFields := []fieldMapping{
		{key: "unit", val: filter.Unit},
		{key: "currency", val: filter.Currency},
	}

	for _, fm := range equalFields {
		if fm.val != nil {
			w.add(fmt.Sprintf("%s = ?", fm.key), *fm.val)
		}
	}

	for _, f := range filter.AttributeFilters {
		if f.Value != nil {
			w.add(fmt.Sprintf("JSON_UNQUOTE(JSON_EXTRACT(attributes, '$.%s')) = ?", f.Key), *f.Value)
		} else if f.ValueRegex != nil {
			w.add(fmt.Sprintf("JSON_UNQUOTE(JSON_EXTRACT(attributes, '$.%s')) RLIKE ?", f.Key), *f.ValueRegex)
		}
	}

	return w
}
