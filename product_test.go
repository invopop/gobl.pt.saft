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

func TestBuildDeductionProduct(t *testing.T) {
	t.Run("General case", func(t *testing.T) {
		p := saft.BuildDeductionProduct(&bill.Discount{
			Key:    addon.DiscountKeyDeduction,
			Reason: "Deposit for 1 hotel room",
		})

		assert.Equal(t, "O", p.ProductType)
		assert.Equal(t, "Deposit for 1 hotel room", p.ProductDesc)
		assert.NotEmpty(t, p.ProductCode)
		assert.Equal(t, p.ProductCode, p.ProductNumberCode)
	})

	t.Run("Without reason", func(t *testing.T) {
		p := saft.BuildDeductionProduct(&bill.Discount{Key: addon.DiscountKeyDeduction})
		assert.Equal(t, "Adiantamento", p.ProductDesc)
	})

	t.Run("Only deductions", func(t *testing.T) {
		inv := &bill.Invoice{
			Discounts: []*bill.Discount{
				{Key: addon.DiscountKeyDeduction, Reason: "A"},
				{Reason: "Promo"},
				{Key: addon.DiscountKeyDeduction, Reason: "B"},
			},
		}
		ps := saft.BuildDeductionProducts(inv)
		require.Len(t, ps, 2)
		assert.Equal(t, "A", ps[0].ProductDesc)
		assert.Equal(t, "B", ps[1].ProductDesc)
		assert.Nil(t, saft.BuildDeductionProducts(nil))
	})
}
