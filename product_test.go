package saft_test

import (
	"testing"

	saft "github.com/invopop/gobl.pt.saft"
	"github.com/invopop/gobl.pt.saft/addon"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildProduct(t *testing.T) {
	t.Run("General case", func(t *testing.T) {
		item := &org.Item{
			Name: "Product Name",
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				addon.ExtKeyProductType: "S",
			}),
		}

		p := saft.BuildProduct(&bill.Line{Item: item})

		assert.Equal(t, "S", p.ProductType)
		assert.Equal(t, "082fa0838490ddae50c893b5fe1694", p.ProductCode)
		assert.Equal(t, "Product Name", p.ProductDesc)
		assert.Equal(t, "082fa0838490ddae50c893b5fe1694", p.ProductNumberCode)
	})

	t.Run("Explicit EAN code", func(t *testing.T) {
		item := &org.Item{
			Identities: []*org.Identity{
				{
					Key:  org.IdentityKeyEAN,
					Code: "4006381333931",
				},
			},
		}

		p := saft.BuildProduct(&bill.Line{Item: item})

		require.Equal(t, "4006381333931", p.ProductNumberCode)
	})
}
