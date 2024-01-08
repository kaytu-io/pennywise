package testutil

import (
	"fmt"
	"github.com/kaytu-io/pennywise/server/resource"
	"testing"

	"github.com/stretchr/testify/assert"
)

// EqualQueryComponents will compare the components but the MonthlyQuantity will be
// compared via String and the rest with assert.Equal
func EqualQueryComponents(t *testing.T, eqcs, aqcs []resource.Component) {
	t.Helper()

	for i, eqc := range eqcs {
		if eqc.MonthlyQuantity.String() == aqcs[i].MonthlyQuantity.String() {
			eqc.MonthlyQuantity = aqcs[i].MonthlyQuantity
		} else {
			assert.Fail(t, fmt.Sprintf("Expected MonthlyQuantity to be %q but was %q", eqc.MonthlyQuantity.String(), aqcs[i].MonthlyQuantity.String()))
			continue
		}
		assert.Equal(t, eqc, aqcs[i])
	}
}
