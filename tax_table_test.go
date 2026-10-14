package saft_test

import (
	"testing"

	saft "github.com/invopop/gobl.pt.saft"
	"github.com/invopop/gobl.pt.saft/addon"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestBuildTaxTableEntry(t *testing.T) {
	t.Run("General case", func(t *testing.T) {
		rt := &tax.RateTotal{
			Percent: num.NewPercentage(230, 3),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				pt.ExtKeyRegion:     pt.RegionMainland,
				addon.ExtKeyTaxRate: addon.TaxRateNormal,
			}),
		}

		tte := saft.BuildTaxTableEntry(rt)

		assert.Equal(t, "IVA", tte.TaxType)
		assert.Equal(t, "PT", tte.TaxCountryRegion)
		assert.Equal(t, "NOR", tte.TaxCode)
		assert.Equal(t, "NOR-PT-23.0", tte.Description)
		assert.Equal(t, "23.0", tte.TaxPercentage)
		assert.Equal(t, "8820a49814467b454a84e72ff7665a", tte.TaxTableEntryID)
	})
}
