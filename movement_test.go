package saft_test

import (
	"testing"
	"time"

	saft "github.com/invopop/gobl.pt.saft"
	"github.com/invopop/gobl.pt.saft/addon"
	"github.com/invopop/gobl.pt.saft/config"
	"github.com/invopop/gobl.pt.saft/test"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildStockMovement(t *testing.T) {
	// Set a fixed time for testing
	config.Now = func() time.Time {
		return time.Date(2024, 12, 23, 23, 30, 30, 0, time.UTC)
	}

	t.Run("Empty cases", func(t *testing.T) {
		sm, err := saft.BuildStockMovement(nil, nil)
		require.NoError(t, err)
		assert.Nil(t, sm)

		sm, err = saft.BuildStockMovement(emptyDelivery(), nil)
		require.NoError(t, err)
		assert.Equal(t, "", sm.DocumentNumber)
	})

	t.Run("Unsupported cases", func(t *testing.T) {
		dlv := &bill.Delivery{
			Supplier:  &org.Party{},
			Currency:  "USD",
			Charges:   []*bill.Charge{{}},
			Discounts: []*bill.Discount{{}},
			Lines: []*bill.Line{
				{
					Charges: []*bill.LineCharge{{}},
				},
			},
		}

		_, err := saft.BuildStockMovement(dlv, nil)
		assert.ErrorContains(t, err, "currency: must be EUR")
		assert.ErrorContains(t, err, "charges: not supported")
		assert.ErrorContains(t, err, "discounts: not supported")
		assert.ErrorContains(t, err, "lines: (0: (charges: not supported")
	})

	t.Run("General case", func(t *testing.T) {
		dlv := &bill.Delivery{
			Series:    "GT SERIES-A",
			Code:      "123",
			IssueDate: cal.MakeDate(2024, 12, 23),
			Type:      bill.DeliveryTypeNote,
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyMovementType: addon.MovementTypeWaybill,
				}),
			},
			Supplier: &org.Party{
				Name: "Test Supplier",
				TaxID: &tax.Identity{
					Country: "PT",
					Code:    "123456789",
				},
				Addresses: []*org.Address{
					{
						Street:   "Avenida da Liberdade",
						Locality: "Lisboa",
						Code:     "1250-146",
						Country:  "PT",
					},
				},
			},
			Customer: &org.Party{
				Name: "Test Customer",
				TaxID: &tax.Identity{
					Country: "PT",
					Code:    "987654321",
				},
				Addresses: []*org.Address{
					{
						Street:   "123 Delivery St",
						Locality: "Lisbon",
						Code:     "1000-100",
						Country:  "PT",
					},
				},
			},
			Receiver: &org.Party{
				Name: "Test Receiver",
				Addresses: []*org.Address{
					{
						Street:   "123 Delivery St",
						Locality: "Lisbon",
						Code:     "1000-100",
						Country:  "PT",
					},
				},
			},
			Despatcher: &org.Party{
				Name: "Test Despatcher",
				Addresses: []*org.Address{
					{
						Street:   "456 Warehouse St",
						Locality: "Porto",
						Code:     "4000-100",
						Country:  "PT",
					},
				},
			},
			DespatchDate: cal.NewDate(2024, 12, 23),
			ReceiveDate:  cal.NewDate(2024, 12, 24),
			Notes: []*org.Note{
				{
					Text: "Test delivery note comments",
				},
			},
			Lines: []*bill.Line{
				{
					Index:    1,
					Quantity: num.MakeAmount(2, 0),
					Item: &org.Item{
						Price: num.NewAmount(100, 2),
						Name:  "Test Item",
						Unit:  "EA",
					},
					Taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Rate:     "standard",
							Percent:  num.NewPercentage(230, 3),
							Ext: tax.ExtensionsOf(cbc.CodeMap{
								pt.ExtKeyRegion:       "PT",
								addon.ExtKeyTaxRate:   "NOR",
								addon.ExtKeyExemption: "",
							}),
						},
					},
					Total: num.NewAmount(200, 2),
				},
			},
			Totals: &bill.Totals{
				Total:        num.MakeAmount(200, 2),
				Tax:          num.MakeAmount(46, 2),
				TotalWithTax: num.MakeAmount(246, 2),
				Taxes: &tax.Total{
					Categories: []*tax.CategoryTotal{
						{
							Code:   tax.CategoryVAT,
							Amount: num.MakeAmount(46, 2),
						},
					},
					Sum: num.MakeAmount(46, 2),
				},
			},
		}

		customer := &saft.Customer{
			CustomerID: "CUST-987654321",
		}
		products := []*saft.Product{
			{
				ProductCode: "PROD1",
				ProductDesc: "Test Item",
			},
		}

		// Build the stock movement
		sm, err := saft.BuildStockMovement(dlv, &saft.BuildParams{
			ATCUD:             "AAJFJHNC82-123",
			Customer:          customer,
			Products:          products,
			SourceID:          "SOURCE123",
			Hash:              "HASH123",
			SigningKeyVersion: "1",
			SystemEntryDate:   "2024-03-01T10:20:40",
		})

		require.NoError(t, err)
		assert.NotNil(t, sm)
		assert.Equal(t, "GT SERIES-A/123", sm.DocumentNumber)
		assert.Equal(t, "AAJFJHNC82-123", sm.ATCUD)
		assert.Equal(t, "HASH123", sm.Hash)
		assert.Equal(t, "1", sm.HashControl)
		assert.Equal(t, 12, sm.Period)
		assert.Equal(t, "2024-12-23", sm.MovementDate)
		assert.Equal(t, "GT", sm.MovementType)
		assert.Equal(t, "SOURCE123", sm.SourceID)
		assert.Equal(t, "2024-03-01T10:20:40", sm.SystemEntryDate)
		assert.Equal(t, "CUST-987654321", sm.CustomerID)
		assert.Equal(t, "Test delivery note comments", sm.MovementComments)
		assert.Equal(t, "2024-12-23T23:59:59", sm.MovementStartTime)
		assert.Equal(t, "2024-12-24T23:59:59", sm.MovementEndTime)

		// Check ShipTo shipping point
		if shipTo := sm.ShipTo; assert.NotNil(t, shipTo) {
			assert.Equal(t, "2024-12-24", shipTo.DeliveryDate)

			// Check address within shipping point
			if addr := shipTo.Address; assert.NotNil(t, addr) {
				assert.Equal(t, "123 Delivery St", addr.AddressDetail)
				assert.Equal(t, "Lisbon", addr.City)
				assert.Equal(t, "1000-100", addr.PostalCode)
				assert.Equal(t, "PT", addr.Country)
			}
		}

		// Check ShipFrom shipping point
		if shipFrom := sm.ShipFrom; assert.NotNil(t, shipFrom) {
			assert.Equal(t, "2024-12-23", shipFrom.DeliveryDate)

			// Check address within shipping point
			if addr := shipFrom.Address; assert.NotNil(t, addr) {
				assert.Equal(t, "456 Warehouse St", addr.AddressDetail)
				assert.Equal(t, "Porto", addr.City)
				assert.Equal(t, "4000-100", addr.PostalCode)
				assert.Equal(t, "PT", addr.Country)
			}
		}

		if ds := sm.DocumentStatus; assert.NotNil(t, ds) {
			assert.Equal(t, saft.StatusNormal, ds.MovementStatus)
			assert.Equal(t, "2024-03-01T10:20:40", ds.MovementStatusDate)
			assert.Equal(t, "SOURCE123", ds.SourceID)
			assert.Equal(t, "P", ds.SourceBilling)
		}

		if lines := sm.Lines; assert.Len(t, lines, 1) {
			line := lines[0]
			assert.Equal(t, 1, line.LineNumber)
			assert.Equal(t, "PROD1", line.ProductCode)
			assert.Equal(t, "Test Item", line.ProductDesc)
			assert.Equal(t, "2", line.Quantity)
			assert.Equal(t, "EA", line.UnitOfMeasure)
			assert.Equal(t, "1.00", line.UnitPrice) // 2.00 / 2 = 1.00
			assert.Equal(t, "", line.TaxPointDate)  // Not supported by SAF-T PT in Movements
			assert.Equal(t, "Test Item", line.Description)
			assert.Empty(t, line.DebitAmount)
			assert.Equal(t, "2.00", line.CreditAmount)

			// Check tax
			if tax := line.Tax; assert.NotNil(t, tax) {
				assert.Equal(t, "IVA", tax.TaxType)
				assert.Equal(t, "PT", tax.TaxCountryRegion)
				assert.Equal(t, "NOR", tax.TaxCode)
				assert.Equal(t, "23.0", tax.TaxPercentage)
			}
		}

		if dt := sm.DocumentTotals; assert.NotNil(t, dt) {
			assert.Equal(t, "0.46", dt.TaxPayable)
			assert.Equal(t, "2.00", dt.NetTotal)
			assert.Equal(t, "2.46", dt.GrossTotal)
		}
	})

	t.Run("Delivery note not valued", func(t *testing.T) {
		dlv := emptyDelivery()
		sm, err := saft.BuildStockMovement(dlv, nil)
		require.NoError(t, err)
		assert.NotNil(t, sm)
		assert.Equal(t, "0.00", sm.DocumentTotals.TaxPayable)
		assert.Equal(t, "0.00", sm.DocumentTotals.NetTotal)
		assert.Equal(t, "0.00", sm.DocumentTotals.GrossTotal)
	})

	t.Run("Tax exempt", func(t *testing.T) {
		dlv := emptyDelivery()
		dlv.Lines = []*bill.Line{
			{
				Item: &org.Item{},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Ext: tax.ExtensionsOf(cbc.CodeMap{
							pt.ExtKeyRegion:       pt.RegionMainland,
							addon.ExtKeyExemption: "M15",
							addon.ExtKeyTaxRate:   addon.TaxRateExempt,
						}),
					},
				},
				Notes: []*org.Note{
					{
						Key:  org.NoteKeyLegal,
						Src:  addon.ExtKeyExemption,
						Code: "M15",
						Text: "Regime da margem de lucro - Objetos de coleção e antiguida",
					},
				},
			},
		}

		sm, err := saft.BuildStockMovement(dlv, &saft.BuildParams{
			Products: []*saft.Product{{}},
		})
		require.NoError(t, err)

		require.Len(t, sm.Lines, 1)
		if tax := sm.Lines[0].Tax; assert.NotNil(t, tax) {
			assert.Equal(t, "IVA", tax.TaxType)
			assert.Equal(t, "PT", tax.TaxCountryRegion)
			assert.Equal(t, "ISE", tax.TaxCode)
			assert.Equal(t, "0", tax.TaxPercentage)
		}
		assert.LessOrEqual(t, len(sm.Lines[0].TaxExemptionReason), 60)
		assert.Equal(t, "Regime da margem de lucro - Objetos de coleção e antiguida", sm.Lines[0].TaxExemptionReason)
		assert.Equal(t, "M15", sm.Lines[0].TaxExemptionCode)
	})

	t.Run("Discounts present", func(t *testing.T) {
		dlv := emptyDelivery()
		dlv.Lines = []*bill.Line{
			{
				Item: &org.Item{},
				Discounts: []*bill.LineDiscount{
					{
						Amount: num.MakeAmount(1234, 2),
					},
					{
						Amount: num.MakeAmount(5678, 3),
					},
				},
			},
		}

		sm, err := saft.BuildStockMovement(dlv, &saft.BuildParams{
			Products: []*saft.Product{{}},
		})
		require.NoError(t, err)

		require.Len(t, sm.Lines, 1)
		assert.Equal(t, "18.018", sm.Lines[0].SettlementAmount)
	})

	t.Run("Movement type return", func(t *testing.T) {
		dlv := emptyDelivery()
		dlv.Tax = &bill.Tax{
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				addon.ExtKeyMovementType: addon.MovementTypeReturn,
			}),
		}
		dlv.Lines = []*bill.Line{
			{
				Item:  &org.Item{},
				Total: num.NewAmount(10000, 2),
			},
		}

		sm, err := saft.BuildStockMovement(dlv, &saft.BuildParams{
			Products: []*saft.Product{{}},
		})
		require.NoError(t, err)

		require.Len(t, sm.Lines, 1)
		assert.Equal(t, "100.00", sm.Lines[0].DebitAmount)
		assert.Empty(t, sm.Lines[0].CreditAmount)
	})

	t.Run("1 decimal gross total", func(t *testing.T) {
		dlv := emptyDelivery()
		dlv.Totals = &bill.Totals{
			TotalWithTax: num.MakeAmount(12345, 1),
		}

		sm, err := saft.BuildStockMovement(dlv, nil)
		require.NoError(t, err)

		assert.Equal(t, "1234.50", sm.DocumentTotals.GrossTotal)
	})

	t.Run("No decimals gross total", func(t *testing.T) {
		dlv := emptyDelivery()
		dlv.Totals = &bill.Totals{
			TotalWithTax: num.MakeAmount(12345, 0),
		}

		sm, err := saft.BuildStockMovement(dlv, nil)
		require.NoError(t, err)

		assert.Equal(t, "12345.00", sm.DocumentTotals.GrossTotal)
	})

	t.Run("Global transport guide without customer", func(t *testing.T) {
		env, err := test.LoadEnvelope("delivery-waybill-global.json")
		require.NoError(t, err)

		gdlv := env.Extract().(*bill.Delivery)
		require.Nil(t, gdlv.Customer, "GT global fixture must not declare a customer")

		// BuildCustomer(nil) auto-fills the generic "Consumidor final" entry as
		// required by the AT (NIF 999999990) and "Desconhecido" placeholders for
		// the billing address.
		cus := saft.BuildCustomer(gdlv.Customer)
		require.NotNil(t, cus)
		assert.Equal(t, saft.GenericCustomerTaxID, cus.CustomerTaxID)
		assert.Equal(t, saft.GenericCustomerID, cus.CompanyName)
		if assert.NotNil(t, cus.BillingAddress) {
			assert.Equal(t, saft.ValueUnknown, cus.BillingAddress.AddressDetail)
			assert.Equal(t, saft.ValueUnknown, cus.BillingAddress.City)
			assert.Equal(t, saft.ValueUnknown, cus.BillingAddress.PostalCode)
			assert.Equal(t, saft.ValueUnknown, cus.BillingAddress.Country)
		}

		// The stock movement itself converts cleanly without a GOBL customer.
		prods := make([]*saft.Product, len(gdlv.Lines))
		for i := range prods {
			prods[i] = &saft.Product{}
		}
		sm, err := saft.BuildStockMovement(gdlv, &saft.BuildParams{Products: prods})
		require.NoError(t, err)
		require.NotNil(t, sm)
		assert.Equal(t, "GT", sm.MovementType)
	})

	t.Run("Return delivery normalizes to GD", func(t *testing.T) {
		env, err := test.LoadEnvelope("delivery-return.json")
		require.NoError(t, err)

		gdlv := env.Extract().(*bill.Delivery)
		prods := make([]*saft.Product, len(gdlv.Lines))
		for i := range prods {
			prods[i] = &saft.Product{}
		}

		sm, err := saft.BuildStockMovement(gdlv, &saft.BuildParams{Products: prods})
		require.NoError(t, err)
		require.NotNil(t, sm)

		assert.Equal(t, "GD", sm.MovementType)
		assert.Equal(t, saft.DocumentClassTransport, sm.GetDocumentClass())

		// Returns map line totals to DebitAmount (inbound), not CreditAmount.
		require.Len(t, sm.Lines, 1)
		assert.Equal(t, "42.50", sm.Lines[0].DebitAmount)
		assert.Empty(t, sm.Lines[0].CreditAmount)
	})

	t.Run("Preceding GT cross-reference fans out to lines", func(t *testing.T) {
		env, err := test.LoadEnvelope("delivery-with-preceding.json")
		require.NoError(t, err)

		gdlv := env.Extract().(*bill.Delivery)
		require.Len(t, gdlv.Preceding, 1)

		prods := make([]*saft.Product, len(gdlv.Lines))
		for i := range prods {
			prods[i] = &saft.Product{}
		}

		sm, err := saft.BuildStockMovement(gdlv, &saft.BuildParams{Products: prods})
		require.NoError(t, err)
		require.NotNil(t, sm)
		require.Len(t, sm.Lines, 2)

		// First line keeps its per-line gline.Order ("PO-2024-001") and gains the
		// document-level GT preceding reference. Second line only gains the GT.
		line1 := sm.Lines[0]
		require.Len(t, line1.OrderReferences, 2)
		assert.Equal(t, "PO-2024-001", line1.OrderReferences[0].OriginatingON)
		assert.Empty(t, line1.OrderReferences[0].OrderDate)
		assert.Equal(t, "GT SERIES-A/1", line1.OrderReferences[1].OriginatingON)
		assert.Equal(t, "2024-08-03", line1.OrderReferences[1].OrderDate)

		line2 := sm.Lines[1]
		require.Len(t, line2.OrderReferences, 1)
		assert.Equal(t, "GT SERIES-A/1", line2.OrderReferences[0].OriginatingON)
		assert.Equal(t, "2024-08-03", line2.OrderReferences[0].OrderDate)
	})

	t.Run("Preceding with non-transport doc type rejected", func(t *testing.T) {
		dlv := emptyDelivery()
		dlv.Series = "GR SERIES-A"
		dlv.Tax = &bill.Tax{
			Ext: tax.MakeExtensions().Set(addon.ExtKeyMovementType, addon.MovementTypeDeliveryNote),
		}
		dlv.Preceding = []*org.DocumentRef{
			{
				Series:    "FT SERIES-A", // FT is an invoice series, not a transport doc
				Code:      "1",
				IssueDate: cal.NewDate(2024, 8, 3),
			},
		}

		_, err := saft.BuildStockMovement(dlv, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "document type (FT) is not supported")
	})
}

func TestStockMovementSetStatus(t *testing.T) {
	// Set a fixed time for testing
	config.Now = func() time.Time {
		return time.Date(2024, 3, 1, 10, 20, 30, 0, time.UTC)
	}

	sm := &saft.StockMovement{}
	sm.SetStatus(saft.StatusNormal, "Reason", "SOURCE123")

	if ds := sm.DocumentStatus; assert.NotNil(t, ds) {
		assert.Equal(t, saft.StatusNormal, ds.MovementStatus)
		assert.Equal(t, "Reason", ds.Reason)
		assert.Equal(t, "2024-03-01T10:20:30", ds.MovementStatusDate)
		assert.Equal(t, "SOURCE123", ds.SourceID)
		assert.Equal(t, "P", ds.SourceBilling)
	}
}

func TestStockMovementGetStatus(t *testing.T) {
	t.Run("Empty case", func(t *testing.T) {
		sm := &saft.StockMovement{}
		assert.Equal(t, "", sm.GetStatus())
	})

	t.Run("With status", func(t *testing.T) {
		sm := &saft.StockMovement{
			DocumentStatus: &saft.MovementStatus{
				MovementStatus: saft.StatusNormal,
			},
		}
		assert.Equal(t, saft.StatusNormal, sm.GetStatus())
	})
}

func TestStockMovementGetStatusDate(t *testing.T) {
	t.Run("Empty case", func(t *testing.T) {
		sm := &saft.StockMovement{}
		assert.Equal(t, "", sm.GetStatusDate())
	})

	t.Run("With status date", func(t *testing.T) {
		sm := &saft.StockMovement{
			DocumentStatus: &saft.MovementStatus{
				MovementStatusDate: "2024-03-01T10:20:30",
			},
		}
		assert.Equal(t, "2024-03-01T10:20:30", sm.GetStatusDate())
	})
}

func emptyDelivery() *bill.Delivery {
	// TODO: return a real empty delivery when the related validations are moved to the addon
	return &bill.Delivery{
		Customer: &org.Party{
			Addresses: []*org.Address{
				{
					Street:   "123 Delivery St",
					Locality: "Lisbon",
					Code:     "1000-100",
					Country:  "PT",
				},
			},
		},
		Despatcher: &org.Party{
			Addresses: []*org.Address{
				{
					Street:   "456 Warehouse St",
					Locality: "Porto",
					Code:     "4000-100",
					Country:  "PT",
				},
			},
		},
	}
}
