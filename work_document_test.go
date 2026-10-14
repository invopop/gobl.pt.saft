package saft_test

import (
	"testing"
	"time"

	saft "github.com/invopop/gobl.pt.saft"
	"github.com/invopop/gobl.pt.saft/addon"
	"github.com/invopop/gobl.pt.saft/config"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildWorkDocumentFromInvoice(t *testing.T) {
	// Set a fixed time for testing
	config.Now = func() time.Time {
		return time.Date(2024, 12, 23, 23, 30, 30, 0, time.UTC)
	}

	t.Run("Empty cases", func(t *testing.T) {
		wd, err := saft.BuildWorkDocumentFromInvoice(nil, nil)
		require.NoError(t, err)
		assert.Nil(t, wd)

		wd, err = saft.BuildWorkDocumentFromInvoice(&bill.Invoice{}, nil)
		require.NoError(t, err)
		assert.Equal(t, "", wd.DocumentNumber)
	})

	t.Run("Unsupported cases", func(t *testing.T) {
		inv := &bill.Invoice{
			Charges:   []*bill.Charge{{}},
			Discounts: []*bill.Discount{{}},
			Lines: []*bill.Line{
				{
					Charges: []*bill.LineCharge{{}},
				},
			},
		}

		_, err := saft.BuildWorkDocumentFromInvoice(inv, nil)
		assert.ErrorContains(t, err, "charges: not supported")
		assert.ErrorContains(t, err, "discounts: not supported")
		assert.ErrorContains(t, err, "lines: (0: (charges: not supported")
	})

	t.Run("General case", func(t *testing.T) {
		inv := &bill.Invoice{
			Series:    "PF SERIES-A",
			Code:      "123",
			IssueDate: cal.MakeDate(2024, 12, 23),
			ValueDate: cal.NewDate(2024, 12, 2),
			Type:      bill.InvoiceTypeProforma,
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeySource:   addon.SourceBillingProduced,
					addon.ExtKeyWorkType: "PF",
				}),
			},
			Supplier: &org.Party{
				Name: "Test Supplier",
				TaxID: &tax.Identity{
					Country: l10n.PT.Tax(),
					Code:    "513912622",
				},
				Addresses: []*org.Address{
					{
						Street:   "Avenida da Liberdade",
						Locality: "Lisboa",
						Code:     "1250-146",
						Country:  "PT",
					},
				},
				Registration: &org.Registration{
					Office:  "Test Office",
					Capital: num.NewAmount(100000, 2),
				},
			},
			Customer: &org.Party{
				Name: "Test Customer",
				TaxID: &tax.Identity{
					Country: "PT",
					Code:    "987654321",
				},
			},
			Lines: []*bill.Line{
				{
					Index:    1,
					Quantity: num.MakeAmount(2, 0),
					Item: &org.Item{
						Name:  "Test Item",
						Price: num.NewAmount(100, 2),
					},
					Taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Rate:     "standard",
							Percent:  num.NewPercentage(230, 3),
							Ext: tax.ExtensionsOf(cbc.CodeMap{
								pt.ExtKeyRegion:     "PT",
								addon.ExtKeyTaxRate: "NOR",
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
				},
			},
		}

		customer := &saft.Customer{
			CustomerID: "987654321",
		}
		products := []*saft.Product{
			{
				ProductCode: "PROD1",
				ProductDesc: "Test Item",
			},
		}

		// Build the work document
		wd, err := saft.BuildWorkDocumentFromInvoice(inv, &saft.BuildParams{
			ATCUD:             "AAJFJHNC82-123",
			Customer:          customer,
			Products:          products,
			SourceID:          "SOURCE123",
			Hash:              "HASH123",
			SigningKeyVersion: "1",
			SystemEntryDate:   "2024-03-01T10:20:40",
		})

		require.NoError(t, err)
		assert.NotNil(t, wd)
		assert.Equal(t, "PF SERIES-A/123", wd.DocumentNumber)
		assert.Equal(t, "AAJFJHNC82-123", wd.ATCUD)
		assert.Equal(t, "HASH123", wd.Hash)
		assert.Equal(t, "1", wd.HashControl)
		assert.Equal(t, 12, wd.Period)
		assert.Equal(t, "2024-12-23", wd.WorkDate)
		assert.Equal(t, "PF", wd.WorkType)
		assert.Equal(t, "SOURCE123", wd.SourceID)
		assert.Equal(t, "2024-03-01T10:20:40", wd.SystemEntryDate)
		assert.Equal(t, "987654321", wd.CustomerID)
		if ds := wd.DocumentStatus; assert.NotNil(t, ds) {
			assert.Equal(t, saft.StatusNormal, ds.WorkStatus)
			assert.Equal(t, "2024-03-01T10:20:40", ds.WorkStatusDate)
			assert.Equal(t, "SOURCE123", ds.SourceID)
			assert.Equal(t, "P", ds.SourceBilling)
		}
		if lines := wd.Lines; assert.Len(t, lines, 1) {
			line := lines[0]
			assert.Equal(t, 1, line.LineNumber)
			assert.Equal(t, "PROD1", line.ProductCode)
			assert.Equal(t, "Test Item", line.ProductDesc)
			assert.Equal(t, "2", line.Quantity)
			assert.Equal(t, "2024-12-02", line.TaxPointDate)
			assert.Equal(t, "Test Item", line.Description)
			assert.Equal(t, "", line.DebitAmount)
			assert.Equal(t, "2.00", line.CreditAmount)

			// Check tax
			if tax := line.Tax; assert.NotNil(t, tax) {
				assert.Equal(t, "IVA", tax.TaxType)
				assert.Equal(t, "PT", tax.TaxCountryRegion)
				assert.Equal(t, "NOR", tax.TaxCode)
				assert.Equal(t, "23.0", tax.TaxPercentage)
			}
		}
		if dt := wd.DocumentTotals; assert.NotNil(t, dt) {
			assert.Equal(t, "0.46", dt.TaxPayable)
			assert.Equal(t, "2.00", dt.NetTotal)
			assert.Equal(t, "2.46", dt.GrossTotal)
		}
	})

	t.Run("Consignment credit", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyWorkType: addon.WorkTypeConsignmentCredit,
				}),
			},
			Lines: []*bill.Line{
				{
					Item:  &org.Item{},
					Total: num.NewAmount(10000, 2),
					Taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Ext: tax.ExtensionsOf(cbc.CodeMap{
								pt.ExtKeyRegion:     pt.RegionMainland,
								addon.ExtKeyTaxRate: addon.TaxRateNormal,
							}),
							Percent: num.NewPercentage(23, 2),
						},
					},
				},
			},
		}

		wd, err := saft.BuildWorkDocumentFromInvoice(inv, &saft.BuildParams{
			Products: []*saft.Product{{}},
		})
		require.NoError(t, err)

		assert.Equal(t, "CC", wd.WorkType)
		require.Len(t, wd.Lines, 1)
		assert.Equal(t, "100.00", wd.Lines[0].DebitAmount)
		assert.Empty(t, wd.Lines[0].CreditAmount)
	})

	t.Run("Foreign currency", func(t *testing.T) {
		// Create a basic order with foreign currency
		inv := &bill.Invoice{
			Currency: currency.USD,
			ExchangeRates: []*currency.ExchangeRate{
				{
					From:   currency.USD,
					To:     currency.EUR,
					Amount: num.MakeAmount(123, 2),
				},
			},
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(2, 0),
					Item: &org.Item{
						Price: num.NewAmount(5228, 2),
					},
					Taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Percent:  num.NewPercentage(230, 3),
						},
					},
				},
			},
		}

		// Build the work document
		wd, err := saft.BuildWorkDocumentFromInvoice(inv, nil)

		require.NoError(t, err)
		assert.NotNil(t, wd)

		// Check currency information
		if dt := wd.DocumentTotals; assert.NotNil(t, dt) {
			assert.Equal(t, "158.19", dt.GrossTotal)
			if curr := dt.Currency; assert.NotNil(t, curr) {
				assert.Equal(t, "USD", curr.CurrencyCode)
				assert.Equal(t, "128.61", curr.CurrencyAmount)
				assert.Equal(t, "1.23", curr.ExchangeRate)
			}
		}
	})

	t.Run("1 decimal gross total from invoice", func(t *testing.T) {
		inv := &bill.Invoice{
			Totals: &bill.Totals{
				TotalWithTax: num.MakeAmount(12345, 1),
			},
		}

		wd, err := saft.BuildWorkDocumentFromInvoice(inv, nil)
		require.NoError(t, err)

		assert.Equal(t, "1234.50", wd.DocumentTotals.GrossTotal)
	})

	t.Run("No decimals gross total from invoice", func(t *testing.T) {
		inv := &bill.Invoice{
			Totals: &bill.Totals{
				TotalWithTax: num.MakeAmount(12345, 0),
			},
		}

		wd, err := saft.BuildWorkDocumentFromInvoice(inv, nil)
		require.NoError(t, err)

		assert.Equal(t, "12345.00", wd.DocumentTotals.GrossTotal)
	})
}

func TestBuildWorkDocumentFromOrder(t *testing.T) {
	// Set a fixed time for testing
	config.Now = func() time.Time {
		return time.Date(2024, 8, 2, 23, 30, 30, 0, time.UTC) // already 3 of August in Portugal Mainland
	}

	t.Run("Recovered order reports its own source", func(t *testing.T) {
		ord := &bill.Order{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyWorkType:  addon.WorkTypeOther,
					addon.ExtKeySource:    addon.SourceBillingManual,
					addon.ExtKeySourceRef: "OUM SERIES-A/1",
				}),
			},
		}

		wd, err := saft.BuildWorkDocumentFromOrder(ord, &saft.BuildParams{
			SigningKeyVersion: "4",
		})
		require.NoError(t, err)

		assert.Equal(t, "M", wd.DocumentStatus.SourceBilling)
		assert.Equal(t, "4-OUM SERIES-A/1", wd.HashControl)
	})

	t.Run("Integrated order is rejected", func(t *testing.T) {
		ord := &bill.Order{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyWorkType: addon.WorkTypeOther,
					addon.ExtKeySource:   addon.SourceBillingIntegrated,
				}),
			},
		}

		_, err := saft.BuildWorkDocumentFromOrder(ord, nil)
		assert.ErrorContains(t, err, "pt-saft-source")
	})

	t.Run("Empty cases", func(t *testing.T) {
		wd, err := saft.BuildWorkDocumentFromOrder(nil, nil)
		require.NoError(t, err)
		assert.Nil(t, wd)

		wd, err = saft.BuildWorkDocumentFromOrder(&bill.Order{}, nil)
		require.NoError(t, err)
		assert.Equal(t, "", wd.DocumentNumber)
	})

	t.Run("Unsupported cases", func(t *testing.T) {
		ord := &bill.Order{
			Charges:   []*bill.Charge{{}},
			Discounts: []*bill.Discount{{}},
			Lines: []*bill.Line{
				{
					Charges: []*bill.LineCharge{{}},
				},
			},
		}

		_, err := saft.BuildWorkDocumentFromOrder(ord, nil)
		assert.ErrorContains(t, err, "charges: not supported")
		assert.ErrorContains(t, err, "discounts: not supported")
		assert.ErrorContains(t, err, "lines: (0: (charges: not supported")
	})

	t.Run("General case", func(t *testing.T) {
		ord := &bill.Order{
			Regime:    tax.WithRegime(l10n.PT.Tax()),
			Type:      bill.OrderTypePurchase,
			Series:    "NE SERIES-A",
			Code:      "123",
			IssueDate: cal.MakeDate(2024, 8, 3),
			ValueDate: cal.NewDate(2024, 8, 2),
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeySource:   addon.SourceBillingProduced,
					addon.ExtKeyWorkType: "NE",
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
			},
			Lines: []*bill.Line{
				{
					Index:    1,
					Quantity: num.MakeAmount(2, 0),
					Item: &org.Item{
						Name:  "Test Item",
						Price: num.NewAmount(100, 2),
					},
					Taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Percent:  num.NewPercentage(230, 3),
							Ext: tax.ExtensionsOf(cbc.CodeMap{
								pt.ExtKeyRegion:     "PT",
								addon.ExtKeyTaxRate: "NOR",
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
				},
			},
		}

		customer := &saft.Customer{
			CustomerID: "987654321",
		}
		products := []*saft.Product{
			{
				ProductCode: "PROD1",
				ProductDesc: "Test Item",
			},
		}

		// Build the work document
		wd, err := saft.BuildWorkDocumentFromOrder(ord, &saft.BuildParams{
			ATCUD:             "AAJFJHNC82-123",
			Customer:          customer,
			Products:          products,
			SourceID:          "SOURCE123",
			Hash:              "HASH123",
			SigningKeyVersion: "1",
			SystemEntryDate:   "2024-03-01T10:20:30",
		})

		require.NoError(t, err)
		assert.NotNil(t, wd)
		assert.Equal(t, "NE SERIES-A/123", wd.DocumentNumber)
		assert.Equal(t, "AAJFJHNC82-123", wd.ATCUD)
		assert.Equal(t, "HASH123", wd.Hash)
		assert.Equal(t, "1", wd.HashControl)
		assert.Equal(t, 8, wd.Period)
		assert.Equal(t, "2024-08-03", wd.WorkDate)
		assert.Equal(t, "NE", wd.WorkType)
		assert.Equal(t, "SOURCE123", wd.SourceID)
		assert.Equal(t, "2024-03-01T10:20:30", wd.SystemEntryDate)
		assert.Equal(t, "987654321", wd.CustomerID)
		if ds := wd.DocumentStatus; assert.NotNil(t, ds) {
			assert.Equal(t, saft.StatusNormal, ds.WorkStatus)
			assert.Equal(t, "2024-03-01T10:20:30", ds.WorkStatusDate)
			assert.Equal(t, "SOURCE123", ds.SourceID)
			assert.Equal(t, "P", ds.SourceBilling)
		}
		if lines := wd.Lines; assert.Len(t, lines, 1) {
			line := lines[0]
			assert.Equal(t, 1, line.LineNumber)
			assert.Equal(t, "PROD1", line.ProductCode)
			assert.Equal(t, "Test Item", line.ProductDesc)
			assert.Equal(t, "2", line.Quantity)
			assert.Equal(t, "2024-08-02", line.TaxPointDate)
			assert.Equal(t, "Test Item", line.Description)
			assert.Equal(t, "", line.DebitAmount)
			assert.Equal(t, "2.00", line.CreditAmount)

			// Check tax
			if tax := line.Tax; assert.NotNil(t, tax) {
				assert.Equal(t, "IVA", tax.TaxType)
				assert.Equal(t, "PT", tax.TaxCountryRegion)
				assert.Equal(t, "NOR", tax.TaxCode)
				assert.Equal(t, "23.0", tax.TaxPercentage)
			}
		}
		if dt := wd.DocumentTotals; assert.NotNil(t, dt) {
			assert.Equal(t, "0.46", dt.TaxPayable)
			assert.Equal(t, "2.00", dt.NetTotal)
			assert.Equal(t, "2.46", dt.GrossTotal)
		}
	})

	t.Run("Foreign currency", func(t *testing.T) {
		// Create a basic order with foreign currency
		ord := &bill.Order{
			Currency: currency.USD,
			ExchangeRates: []*currency.ExchangeRate{
				{
					From:   currency.USD,
					To:     currency.EUR,
					Amount: num.MakeAmount(123, 2),
				},
			},
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(2, 0),
					Item: &org.Item{
						Price: num.NewAmount(5228, 2),
					},
					Taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Percent:  num.NewPercentage(230, 3),
						},
					},
				},
			},
		}

		// Build the work document
		wd, err := saft.BuildWorkDocumentFromOrder(ord, nil)

		require.NoError(t, err)
		assert.NotNil(t, wd)

		// Check currency information
		if dt := wd.DocumentTotals; assert.NotNil(t, dt) {
			assert.Equal(t, "158.19", dt.GrossTotal)
			if curr := dt.Currency; assert.NotNil(t, curr) {
				assert.Equal(t, "USD", curr.CurrencyCode)
				assert.Equal(t, "128.61", curr.CurrencyAmount)
				assert.Equal(t, "1.23", curr.ExchangeRate)
			}
		}
	})

	t.Run("1 decimal gross total from order", func(t *testing.T) {
		ord := &bill.Order{
			Totals: &bill.Totals{
				TotalWithTax: num.MakeAmount(12345, 1),
			},
		}

		wd, err := saft.BuildWorkDocumentFromOrder(ord, nil)
		require.NoError(t, err)

		assert.Equal(t, "1234.50", wd.DocumentTotals.GrossTotal)
	})

	t.Run("No decimals gross total from order", func(t *testing.T) {
		ord := &bill.Order{
			Totals: &bill.Totals{
				TotalWithTax: num.MakeAmount(12345, 0),
			},
		}

		wd, err := saft.BuildWorkDocumentFromOrder(ord, nil)
		require.NoError(t, err)

		assert.Equal(t, "12345.00", wd.DocumentTotals.GrossTotal)
	})
}

func TestWorkDocumentSetStatus(t *testing.T) {
	// Set a fixed time for testing
	config.Now = func() time.Time {
		return time.Date(2024, 3, 1, 10, 20, 30, 0, time.UTC)
	}

	t.Run("preserves the original source billing", func(t *testing.T) {
		wd := &saft.WorkDocument{
			DocumentStatus: &saft.WorkStatus{
				WorkStatus:    saft.StatusNormal,
				SourceBilling: "M",
			},
		}
		wd.SetStatus(saft.StatusCanceled, "Reason", "SOURCE123")

		if ds := wd.DocumentStatus; assert.NotNil(t, ds) {
			assert.Equal(t, saft.StatusCanceled, ds.WorkStatus)
			assert.Equal(t, "2024-03-01T10:20:30", ds.WorkStatusDate)
			assert.Equal(t, "SOURCE123", ds.SourceID)
			assert.Equal(t, "M", ds.SourceBilling)
			assert.Equal(t, "Reason", ds.Reason)
		}
	})

	t.Run("defaults the source billing when there is none", func(t *testing.T) {
		wd := &saft.WorkDocument{}
		wd.SetStatus(saft.StatusNormal, "Reason", "SOURCE123")

		if ds := wd.DocumentStatus; assert.NotNil(t, ds) {
			assert.Equal(t, saft.StatusNormal, ds.WorkStatus)
			assert.Equal(t, "P", ds.SourceBilling)
		}
	})
}

func TestWorkDocumentGetStatus(t *testing.T) {
	t.Run("Empty case", func(t *testing.T) {
		wd := &saft.WorkDocument{}
		assert.Equal(t, "", wd.GetStatus())
	})

	t.Run("With status", func(t *testing.T) {
		wd := &saft.WorkDocument{
			DocumentStatus: &saft.WorkStatus{
				WorkStatus: saft.StatusNormal,
			},
		}
		assert.Equal(t, saft.StatusNormal, wd.GetStatus())
	})
}

func TestWorkDocumentGetStatusDate(t *testing.T) {
	t.Run("Empty case", func(t *testing.T) {
		wd := &saft.WorkDocument{}
		assert.Equal(t, "", wd.GetStatusDate())
	})

	t.Run("With status date", func(t *testing.T) {
		wd := &saft.WorkDocument{
			DocumentStatus: &saft.WorkStatus{
				WorkStatusDate: "2024-03-01T10:20:30",
			},
		}
		assert.Equal(t, "2024-03-01T10:20:30", wd.GetStatusDate())
	})
}
