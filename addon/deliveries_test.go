package addon_test

import (
	"testing"

	"github.com/invopop/gobl.pt.saft/addon"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeliveryValidation(t *testing.T) {
	t.Run("valid delivery", func(t *testing.T) {
		dlv := validDelivery()
		require.NoError(t, rules.Validate(dlv, withAddonContext()))
	})

	t.Run("missing movement type", func(t *testing.T) {
		dlv := validDelivery()

		dlv.Tax = nil
		assert.ErrorContains(t, rules.Validate(dlv, withAddonContext()), "tax requires 'pt-saft-movement-type' extension")

		dlv.Tax = new(bill.Tax)
		assert.ErrorContains(t, rules.Validate(dlv, withAddonContext()), "tax requires 'pt-saft-movement-type' extension")
	})

	t.Run("missing despatch date", func(t *testing.T) {
		dlv := validDelivery()
		dlv.DespatchDate = nil
		assert.ErrorContains(t, rules.Validate(dlv, withAddonContext()), "despatch date is required")
	})

	t.Run("invalid series format", func(t *testing.T) {
		dlv := validDelivery()

		dlv.Series = "SERIES-A"
		assert.ErrorContains(t, rules.Validate(dlv, withAddonContext()), "series format must be valid")
	})

	t.Run("invalid code format", func(t *testing.T) {
		dlv := validDelivery()

		dlv.Code = "ABCD"
		assert.ErrorContains(t, rules.Validate(dlv, withAddonContext()), "code format must be valid")
	})

	t.Run("valid full code", func(t *testing.T) {
		dlv := validDelivery()

		dlv.Series = ""
		dlv.Code = "GR SERIES-A/123"
		assert.NoError(t, rules.Validate(dlv, withAddonContext()))
	})

	t.Run("invalid full code", func(t *testing.T) {
		dlv := validDelivery()

		dlv.Series = ""
		dlv.Code = "ABCDEF"
		assert.ErrorContains(t, rules.Validate(dlv, withAddonContext()), "code format must be valid")
	})

	t.Run("missing supplier tax ID", func(t *testing.T) {
		dlv := validDelivery()

		dlv.Supplier.TaxID = nil
		assert.ErrorContains(t, rules.Validate(dlv, withAddonContext()), "supplier tax ID is required")

		dlv.Supplier.TaxID = &tax.Identity{
			Country: "PT",
			Code:    "",
		}
		assert.ErrorContains(t, rules.Validate(dlv, withAddonContext()), "supplier tax ID code is required")

		// dlv.Supplier = nil is caught by core GOBL rules (supplier is required)
	})

	t.Run("waybill without customer", func(t *testing.T) {
		dlv := validDelivery()
		dlv.Type = bill.DeliveryTypeWaybill
		dlv.Series = "GT SERIES-A"
		dlv.Tax.Ext = tax.ExtensionsOf(cbc.CodeMap{
			addon.ExtKeyMovementType: addon.MovementTypeWaybill,
		})
		dlv.Customer = nil
		require.NoError(t, rules.Validate(dlv, withAddonContext()))
	})

	t.Run("nil preceding", func(t *testing.T) {
		dlv := validDelivery()
		dlv.Preceding = nil
		require.NoError(t, rules.Validate(dlv, withAddonContext()))
	})

	t.Run("valid preceding", func(t *testing.T) {
		dlv := validDelivery()
		dlv.Preceding = []*org.DocumentRef{
			{
				Series:    "GR SERIES-A",
				Code:      "1",
				IssueDate: cal.NewDate(2023, 1, 1),
			},
		}
		require.NoError(t, rules.Validate(dlv, withAddonContext()))
	})

	t.Run("missing series", func(t *testing.T) {
		dlv := validDelivery()
		dlv.Preceding = []*org.DocumentRef{
			{
				Code:      "1",
				IssueDate: cal.NewDate(2023, 1, 1),
			},
		}
		assert.ErrorContains(t, rules.Validate(dlv, withAddonContext()), "preceding series is required")
	})

	t.Run("missing code", func(t *testing.T) {
		dlv := validDelivery()
		dlv.Preceding = []*org.DocumentRef{
			{
				Series:    "GR SERIES-A",
				IssueDate: cal.NewDate(2023, 1, 1),
			},
		}
		assert.ErrorContains(t, rules.Validate(dlv, withAddonContext()), "preceding code is required")
	})

	t.Run("valid post codes", func(t *testing.T) {
		dlv := validDelivery()
		dlv.Supplier.Addresses = []*org.Address{{Code: "1000-100", Country: "PT"}}
		dlv.Customer.Addresses = []*org.Address{{Code: "2000-200"}}
		dlv.Despatcher = &org.Party{
			Name:      "Test Despatcher",
			Addresses: []*org.Address{{Code: "3000-300", Country: "PT"}},
		}
		dlv.Receiver = &org.Party{
			Name:      "Test Receiver",
			Addresses: []*org.Address{{Code: "4000-400", Country: "PT"}},
		}
		require.NoError(t, rules.Validate(dlv, withAddonContext()))
	})

	t.Run("invalid supplier post code", func(t *testing.T) {
		dlv := validDelivery()
		dlv.Supplier.Addresses = []*org.Address{{Code: "2050", Country: "PT"}}
		assert.ErrorContains(t, rules.Validate(dlv, withAddonContext()), "delivery supplier post code must be in the 'NNNN-NNN' format")
	})

	t.Run("invalid customer post code", func(t *testing.T) {
		dlv := validDelivery()
		dlv.Customer.Addresses = []*org.Address{{Code: "2050", Country: "PT"}}
		assert.ErrorContains(t, rules.Validate(dlv, withAddonContext()), "delivery customer post code must be in the 'NNNN-NNN' format")
	})

	t.Run("invalid despatcher post code", func(t *testing.T) {
		dlv := validDelivery()
		dlv.Despatcher = &org.Party{
			Name:      "Test Despatcher",
			Addresses: []*org.Address{{Code: "2050", Country: "PT"}},
		}
		faults := rules.Validate(dlv, withAddonContext())
		require.Error(t, faults)
		assert.ErrorContains(t, faults, "delivery despatcher post code must be in the 'NNNN-NNN' format")
		assert.True(t, faults.HasPath("$.despatcher.addresses[0].code"))
		assert.True(t, faults.HasCode("GOBL-PT-SAFT-BILL-DELIVERY-12"))
	})

	t.Run("invalid receiver post code", func(t *testing.T) {
		dlv := validDelivery()
		dlv.Receiver = &org.Party{
			Name:      "Test Receiver",
			Addresses: []*org.Address{{Code: "1000 100", Country: "PT"}},
		}
		assert.ErrorContains(t, rules.Validate(dlv, withAddonContext()), "delivery receiver post code must be in the 'NNNN-NNN' format")
	})

	t.Run("invalid post code without country", func(t *testing.T) {
		dlv := validDelivery()
		dlv.Supplier.Addresses = []*org.Address{{Code: "2050"}}
		assert.ErrorContains(t, rules.Validate(dlv, withAddonContext()), "delivery supplier post code must be in the 'NNNN-NNN' format")
	})

	t.Run("foreign post code", func(t *testing.T) {
		dlv := validDelivery()
		dlv.Customer.Addresses = []*org.Address{{Code: "28001", Country: "ES"}}
		require.NoError(t, rules.Validate(dlv, withAddonContext()))
	})

	t.Run("missing post code", func(t *testing.T) {
		dlv := validDelivery()
		dlv.Supplier.Addresses = []*org.Address{{Country: "PT"}}
		require.NoError(t, rules.Validate(dlv, withAddonContext()))
	})

	t.Run("several addresses", func(t *testing.T) {
		dlv := validDelivery()
		dlv.Supplier.Addresses = []*org.Address{
			{Code: "1000-100", Country: "PT"},
			{Code: "2050", Country: "PT"},
		}
		faults := rules.Validate(dlv, withAddonContext())
		require.Error(t, faults)
		assert.True(t, faults.HasPath("$.supplier.addresses[1].code"))
	})

	t.Run("several preceding documents", func(t *testing.T) {
		dlv := validDelivery()
		dlv.Preceding = []*org.DocumentRef{
			{
				Series:    "GR SERIES-A",
				Code:      "1",
				IssueDate: cal.NewDate(2023, 1, 1),
			},
			{
				Series:    "GR SERIES-A",
				Code:      "2",
				IssueDate: cal.NewDate(2023, 1, 1),
			},
		}
		assert.ErrorContains(t, rules.Validate(dlv, withAddonContext()), "preceding must have at most one entry")
	})
}

func TestDeliveryNormalization(t *testing.T) {
	t.Run("note type", func(t *testing.T) {
		dlv := &bill.Delivery{
			Type: bill.DeliveryTypeNote,
		}
		norm.Normalize(dlv, tax.AddonContext(addon.V1))
		require.NotNil(t, dlv.Tax)
		require.NotNil(t, dlv.Tax.Ext)
		assert.Equal(t, addon.MovementTypeDeliveryNote, dlv.Tax.Ext.Get(addon.ExtKeyMovementType))
	})

	t.Run("waybill type", func(t *testing.T) {
		dlv := &bill.Delivery{
			Type: bill.DeliveryTypeWaybill,
		}
		norm.Normalize(dlv, tax.AddonContext(addon.V1))
		require.NotNil(t, dlv.Tax)
		require.NotNil(t, dlv.Tax.Ext)
		assert.Equal(t, addon.MovementTypeWaybill, dlv.Tax.Ext.Get(addon.ExtKeyMovementType))
	})

	t.Run("return tag", func(t *testing.T) {
		dlv := &bill.Delivery{
			Type: bill.DeliveryTypeNote,
		}
		dlv.SetTags(addon.TagReturn)
		norm.Normalize(dlv, tax.AddonContext(addon.V1))
		require.NotNil(t, dlv.Tax)
		require.NotNil(t, dlv.Tax.Ext)
		assert.Equal(t, addon.MovementTypeReturn, dlv.Tax.Ext.Get(addon.ExtKeyMovementType))
	})

	t.Run("respect existing value", func(t *testing.T) {
		dlv := &bill.Delivery{
			Type: bill.DeliveryTypeNote,
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyMovementType: addon.MovementTypeFixedAssets,
				}),
			},
		}
		norm.Normalize(dlv, tax.AddonContext(addon.V1))
		assert.Equal(t, addon.MovementTypeFixedAssets, dlv.Tax.Ext.Get(addon.ExtKeyMovementType))
	})
}

func validDelivery() *bill.Delivery {
	date := cal.NewDate(2023, 1, 1)

	return &bill.Delivery{
		Type:      bill.DeliveryTypeNote,
		IssueDate: *date,
		Series:    "GR SERIES-A",
		Code:      "123",
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: "PT",
				Code:    "123456789",
			},
			Name: "Test Supplier",
		},
		Customer: &org.Party{
			Name: "Test Customer",
		},
		DespatchDate: date,
		Lines: []*bill.Line{
			{
				Index:    1,
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name: "Test Item",
					Unit: "one",
					Ext: tax.ExtensionsOf(cbc.CodeMap{
						addon.ExtKeyProductType: addon.ProductTypeService,
					}),
				},
			},
		},
		Tax: &bill.Tax{
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				addon.ExtKeyMovementType: addon.MovementTypeDeliveryNote,
			}),
		},
		Totals: &bill.Totals{},
	}
}
