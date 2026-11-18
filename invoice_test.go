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
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildInvoice(t *testing.T) {
	config.Now = func() time.Time {
		return time.Date(2024, 8, 2, 23, 30, 30, 0, time.UTC) // already 3 of August in Portugal Mainland
	}

	t.Run("Empty cases", func(t *testing.T) {
		i, err := saft.BuildInvoice(nil, nil)
		require.NoError(t, err)
		assert.Nil(t, i)

		i, err = saft.BuildInvoice(&bill.Invoice{}, nil)
		require.NoError(t, err)
		assert.Equal(t, "", i.InvoiceNo)
	})

	t.Run("General case", func(t *testing.T) {
		inv := &bill.Invoice{
			Regime:   tax.WithRegime(l10n.PT.Tax()),
			Currency: currency.EUR,
			Customer: &org.Party{
				Name: "Innovatech Lda",
				TaxID: &tax.Identity{
					Code:    "123456789",
					Country: "PT",
				},
			},
			Series:    "FT SERIES-A",
			Code:      "1234",
			IssueDate: cal.MakeDate(2024, 8, 3),
			ValueDate: cal.NewDate(2024, 8, 2),
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeStandard,
					addon.ExtKeySource:      "P",
				}),
			},
			Lines: []*bill.Line{
				{
					Index: 1,
					Item: &org.Item{
						Price: num.NewAmount(800041, 4),
						Unit:  "kg",
					},
					Quantity: num.MakeAmount(12345, 3),
					Total:    num.NewAmount(98765, 2),
					Taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Ext: tax.ExtensionsOf(cbc.CodeMap{
								pt.ExtKeyRegion:     pt.RegionAzores,
								addon.ExtKeyTaxRate: addon.TaxRateReduced,
							}),
							Percent: num.NewPercentage(40, 3),
						},
					},
				},
				{
					Index: 2,
					Item:  &org.Item{},
				},
			},
			Payment: &bill.PaymentDetails{
				Advances: []*pay.Record{
					{
						Date:   cal.NewDate(2024, 2, 29),
						Amount: num.MakeAmount(12345, 2),
						Ext: tax.ExtensionsOf(cbc.CodeMap{
							addon.ExtKeyPaymentMeans: "NU",
						}),
					},
				},
			},
			Totals: &bill.Totals{
				Taxes: &tax.Total{
					Categories: []*tax.CategoryTotal{
						{
							Code:   tax.CategoryVAT,
							Amount: num.MakeAmount(3333, 2),
						},
					},
				},
				Total:        num.MakeAmount(98765, 2),
				TotalWithTax: num.MakeAmount(95432, 2),
			},
		}

		i, err := saft.BuildInvoice(inv, &saft.BuildParams{
			ATCUD:    "AAJFJHNC82-1234",
			Customer: &saft.Customer{CustomerID: "CUSTID"},
			Products: []*saft.Product{
				{
					ProductCode: "PRODUCTCODE1",
					ProductDesc: "Product Description 1",
				},
				{},
			},
			SourceID:          "SOURCEID",
			Hash:              "aX1zuk8FOTLY2UCFH8",
			SigningKeyVersion: "4",
			SystemEntryDate:   "2024-03-01T10:20:40",
		})
		require.NoError(t, err)

		assert.Equal(t, "FT SERIES-A/1234", i.InvoiceNo)
		assert.Equal(t, "AAJFJHNC82-1234", i.ATCUD)
		if ds := i.DocumentStatus; assert.NotNil(t, ds) {
			assert.Equal(t, "N", ds.InvoiceStatus)
			assert.Equal(t, "2024-03-01T10:20:40", ds.InvoiceStatusDate)
			assert.Equal(t, "SOURCEID", ds.SourceID)
			assert.Equal(t, "P", ds.SourceBilling)
		}
		assert.Equal(t, i.Hash, "aX1zuk8FOTLY2UCFH8")
		assert.Equal(t, "4", i.HashControl)
		assert.Equal(t, 8, i.Period)
		assert.Equal(t, "2024-08-03", i.InvoiceDate)
		assert.Equal(t, "FT", i.InvoiceType)
		if sr := i.SpecialRegimes; assert.NotNil(t, sr) {
			assert.Equal(t, 0, sr.SelfBillingIndicator)
			assert.Equal(t, 0, sr.CashVATSchemeIndicator)
			assert.Equal(t, 0, sr.ThirdPartiesBillingIndicator)
		}
		assert.Equal(t, "SOURCEID", i.SourceID)
		assert.Equal(t, "2024-03-01T10:20:40", i.SystemEntryDate)
		assert.Equal(t, "CUSTID", i.CustomerID)
		if ls := i.Lines; assert.Len(t, ls, 2) {
			assert.Equal(t, 1, ls[0].LineNumber)
			assert.Equal(t, "PRODUCTCODE1", ls[0].ProductCode)
			assert.Equal(t, "Product Description 1", ls[0].ProductDesc)
			assert.Equal(t, "12.345", ls[0].Quantity)
			assert.Equal(t, "kg", ls[0].UnitOfMeasure)
			assert.Equal(t, "80.0041", ls[0].UnitPrice)
			assert.Equal(t, "2024-08-02", ls[0].TaxPointDate)
			assert.Nil(t, ls[0].References)
			assert.Equal(t, "Product Description 1", ls[0].Description)
			assert.Empty(t, ls[0].DebitAmount)
			assert.Equal(t, "987.65", ls[0].CreditAmount)
			if tax := ls[0].Tax; assert.NotNil(t, ls) {
				assert.Equal(t, "IVA", tax.TaxType)
				assert.Equal(t, "PT-AC", tax.TaxCountryRegion)
				assert.Equal(t, "RED", tax.TaxCode)
				assert.Equal(t, "4.0", tax.TaxPercentage)
			}
			assert.Empty(t, ls[0].TaxExemptionReason)
			assert.Empty(t, ls[0].TaxExemptionCode)
			assert.Empty(t, ls[0].SettlementAmount)
		}
		if dt := i.DocumentTotals; assert.NotNil(t, dt) {
			assert.Equal(t, "33.33", dt.TaxPayable)
			assert.Equal(t, "987.65", dt.NetTotal)
			assert.Equal(t, "954.32", dt.GrossTotal)
			if ps := dt.Payments; assert.Len(t, ps, 1) {
				assert.Equal(t, "2024-02-29", ps[0].PaymentDate)
				assert.Equal(t, "123.45", ps[0].PaymentAmount)
				assert.Equal(t, "NU", ps[0].PaymentMechanism)
			}
		}
	})

	t.Run("Cash-VAT scheme", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeStandard,
					addon.ExtKeyCashVAT:     "1",
				}),
			},
		}

		i, err := saft.BuildInvoice(inv, nil)
		require.NoError(t, err)
		if sr := i.SpecialRegimes; assert.NotNil(t, sr) {
			assert.Equal(t, 1, sr.CashVATSchemeIndicator)
		}
	})

	t.Run("Cash-VAT scheme disabled", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeStandard,
					addon.ExtKeyCashVAT:     "0",
				}),
			},
		}

		i, err := saft.BuildInvoice(inv, nil)
		require.NoError(t, err)
		if sr := i.SpecialRegimes; assert.NotNil(t, sr) {
			assert.Equal(t, 0, sr.CashVATSchemeIndicator)
		}
	})

	t.Run("Preceding present", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeCreditNote,
				}),
			},
			Preceding: []*org.DocumentRef{
				{
					Series: "FT SERIES-A",
					Code:   "1233",
					Reason: "Preceding reason",
				},
			},
			Lines: []*bill.Line{{Item: &org.Item{}}},
		}

		i, err := saft.BuildInvoice(inv, &saft.BuildParams{
			Products: []*saft.Product{{}},
		})
		require.NoError(t, err)

		require.Len(t, i.Lines, 1)
		require.Len(t, i.Lines[0].References, 1)
		assert.Equal(t, "FT SERIES-A/1233", i.Lines[0].References[0].Reference)
		assert.Equal(t, "Preceding reason", i.Lines[0].References[0].Reason)
	})

	t.Run("Credit note", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeCreditNote,
				}),
			},
			Preceding: []*org.DocumentRef{
				{
					Series: "FT SERIES-A",
					Code:   "1233",
					Reason: "Credit note reason",
				},
			},
			Lines: []*bill.Line{
				{
					Item:  &org.Item{},
					Total: num.NewAmount(10000, 2),
				},
			},
		}

		i, err := saft.BuildInvoice(inv, &saft.BuildParams{
			Products: []*saft.Product{{}},
		})
		require.NoError(t, err)

		assert.Equal(t, "NC", i.InvoiceType)
		require.Len(t, i.Lines, 1)
		assert.Equal(t, "100.00", i.Lines[0].DebitAmount)
		assert.Empty(t, i.Lines[0].CreditAmount)
	})

	t.Run("Tax exempt", func(t *testing.T) {
		inv := &bill.Invoice{
			Lines: []*bill.Line{
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
			},
		}

		i, err := saft.BuildInvoice(inv, &saft.BuildParams{
			Products: []*saft.Product{{}},
		})
		require.NoError(t, err)

		require.Len(t, i.Lines, 1)
		if tax := i.Lines[0].Tax; assert.NotNil(t, tax) {
			assert.Equal(t, "IVA", tax.TaxType)
			assert.Equal(t, "PT", tax.TaxCountryRegion)
			assert.Equal(t, "ISE", tax.TaxCode)
			assert.Equal(t, "0", tax.TaxPercentage)
		}
		assert.LessOrEqual(t, len(i.Lines[0].TaxExemptionReason), 60)
		assert.Equal(t, "Regime da margem de lucro - Objetos de coleção e antiguida", i.Lines[0].TaxExemptionReason)
		assert.Equal(t, "M15", i.Lines[0].TaxExemptionCode)
	})

	t.Run("Discounts present", func(t *testing.T) {
		inv := &bill.Invoice{
			Lines: []*bill.Line{
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
			},
		}

		i, err := saft.BuildInvoice(inv, &saft.BuildParams{
			Products: []*saft.Product{{}},
		})
		require.NoError(t, err)

		require.Len(t, i.Lines, 1)
		assert.Equal(t, "18.018", i.Lines[0].SettlementAmount)
	})

	t.Run("Unit price without discounts", func(t *testing.T) {
		inv := &bill.Invoice{
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(35, 1),
					Item: &org.Item{
						Price: num.NewAmount(345, 2),
					},
					Total: num.NewAmount(1208, 2),
				},
			},
		}

		i, err := saft.BuildInvoice(inv, &saft.BuildParams{
			Products: []*saft.Product{{}},
		})
		require.NoError(t, err)

		require.Len(t, i.Lines, 1)
		assert.Equal(t, "3.45", i.Lines[0].UnitPrice)
	})

	t.Run("Unit price with discounts", func(t *testing.T) {
		inv := &bill.Invoice{
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(35, 1),
					Item: &org.Item{
						Price: num.NewAmount(345, 2),
					},
					Discounts: []*bill.LineDiscount{
						{
							Amount: num.MakeAmount(121, 2),
						},
					},
					Total: num.NewAmount(1087, 2),
				},
			},
		}

		i, err := saft.BuildInvoice(inv, &saft.BuildParams{
			Products: []*saft.Product{{}},
		})
		require.NoError(t, err)

		require.Len(t, i.Lines, 1)
		assert.Equal(t, "3.1057", i.Lines[0].UnitPrice)
	})

	t.Run("Unit price with full discount and currency rounding", func(t *testing.T) {
		inv := &bill.Invoice{
			Lines: []*bill.Line{
				{
					Quantity: num.MakeAmount(12000, 3),
					Item: &org.Item{
						Price: num.NewAmount(20325000, 8),
					},
					Discounts: []*bill.LineDiscount{
						{
							Amount: num.MakeAmount(244, 2),
						},
					},
					Total: num.NewAmount(0, 2),
				},
			},
		}

		i, err := saft.BuildInvoice(inv, &saft.BuildParams{
			Products: []*saft.Product{{}},
		})
		require.NoError(t, err)

		require.Len(t, i.Lines, 1)
		assert.Equal(t, "0.0000", i.Lines[0].UnitPrice)
	})

	t.Run("Zero quantity", func(t *testing.T) {
		inv := &bill.Invoice{
			Lines: []*bill.Line{
				{
					Item: &org.Item{
						Price: num.NewAmount(1234, 2),
					},
					Quantity: num.MakeAmount(0, 3),
					Total:    num.NewAmount(0, 2),
				},
			},
		}

		i, err := saft.BuildInvoice(inv, &saft.BuildParams{
			Products: []*saft.Product{{}},
		})
		require.NoError(t, err)

		require.Len(t, i.Lines, 1)
		assert.Equal(t, "0.000", i.Lines[0].Quantity)
		assert.Equal(t, "12.34", i.Lines[0].UnitPrice)
	})

	t.Run("Order reference present", func(t *testing.T) {
		inv := &bill.Invoice{
			Lines: []*bill.Line{
				{
					Item:  &org.Item{},
					Order: "ORDER-123",
				},
			},
		}

		i, err := saft.BuildInvoice(inv, &saft.BuildParams{
			Products: []*saft.Product{{}},
		})
		require.NoError(t, err)

		require.Len(t, i.Lines, 1)
		require.Len(t, i.Lines[0].OrderReferences, 1)
		assert.Equal(t, "ORDER-123", i.Lines[0].OrderReferences[0].OriginatingON)
		assert.Empty(t, i.Lines[0].OrderReferences[0].OrderDate)
	})

	t.Run("Foreign currency", func(t *testing.T) {
		env, err := test.LoadEnvelope("invoice-b2b.json")
		require.NoError(t, err)

		inv := env.Extract().(*bill.Invoice)
		inv.Currency = currency.USD
		inv.ExchangeRates = []*currency.ExchangeRate{
			{
				From:   currency.USD,
				To:     currency.EUR,
				Amount: num.MakeAmount(123, 2),
			},
		}

		i, err := saft.BuildInvoice(inv, &saft.BuildParams{
			Products: []*saft.Product{{}},
		})
		require.NoError(t, err)

		if dt := i.DocumentTotals; assert.NotNil(t, dt) {
			assert.Equal(t, "125.80", dt.GrossTotal)
			if curr := i.DocumentTotals.Currency; assert.NotNil(t, curr) {
				assert.Equal(t, "USD", curr.CurrencyCode)
				assert.Equal(t, "102.28", curr.CurrencyAmount)
				assert.Equal(t, "1.23", curr.ExchangeRate)
			}
		}
	})

	t.Run("Credit sale", func(t *testing.T) {
		env, err := test.LoadEnvelope("invoice-b2b.json")
		require.NoError(t, err)

		inv := env.Extract().(*bill.Invoice)
		inv.Payment.Terms = &pay.Terms{
			DueDates: []*pay.DueDate{
				{
					Date:   cal.NewDate(2024, 12, 24),
					Amount: num.NewAmount(1234, 2),
				},
				{
					Date:   cal.NewDate(2024, 12, 25),
					Amount: num.NewAmount(1000, 2),
					Notes:  "This is a note",
				},
			},
		}

		i, err := saft.BuildInvoice(inv, &saft.BuildParams{
			Products: []*saft.Product{{}},
		})
		require.NoError(t, err)

		if assert.Len(t, i.DocumentTotals.Settlements, 2) {
			assert.Equal(t, "12.34 EUR por pagar com vencimento em 2024-12-24", i.DocumentTotals.Settlements[0].PaymentTerms)
			assert.Equal(t, "This is a note", i.DocumentTotals.Settlements[1].PaymentTerms)
		}
	})

	t.Run("Manual document", func(t *testing.T) {
		env, err := test.LoadEnvelope("invoice-b2b.json")
		require.NoError(t, err)

		inv := env.Extract().(*bill.Invoice)
		inv.Tax.Ext = tax.ExtensionsOf(cbc.CodeMap{
			addon.ExtKeySource:    addon.SourceBillingManual,
			addon.ExtKeySourceRef: "FTM abc/123",
		})

		i, err := saft.BuildInvoice(inv, &saft.BuildParams{
			SigningKeyVersion: "23",
			Products:          []*saft.Product{{}},
		})
		require.NoError(t, err)

		assert.Equal(t, "M", i.DocumentStatus.SourceBilling)
		assert.Equal(t, "23-FTM abc/123", i.HashControl)
	})

	t.Run("1 decimal gross total", func(t *testing.T) {
		inv := &bill.Invoice{
			Totals: &bill.Totals{
				TotalWithTax: num.MakeAmount(12345, 1),
			},
		}

		i, err := saft.BuildInvoice(inv, nil)
		require.NoError(t, err)

		assert.Equal(t, "1234.50", i.DocumentTotals.GrossTotal)
	})

	t.Run("No decimals gross total", func(t *testing.T) {
		inv := &bill.Invoice{
			Totals: &bill.Totals{
				TotalWithTax: num.MakeAmount(12345, 0),
			},
		}

		i, err := saft.BuildInvoice(inv, nil)
		require.NoError(t, err)

		assert.Equal(t, "12345.00", i.DocumentTotals.GrossTotal)
	})

	t.Run("Withholding taxes", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Notes: []*tax.Note{
					{
						Category: pt.TaxCategoryIRS,
						Text:     "Art. 101.º CIRS",
					},
				},
			},
			Totals: &bill.Totals{
				Taxes: &tax.Total{
					Categories: []*tax.CategoryTotal{
						{
							Code:     pt.TaxCategoryIRS,
							Retained: true,
							Amount:   num.MakeAmount(2300, 2),
						},
					},
				},
				TotalWithTax: num.MakeAmount(95432, 2),
			},
		}

		i, err := saft.BuildInvoice(inv, nil)
		require.NoError(t, err)

		if assert.Len(t, i.WithholdingTaxes, 1) {
			wt := i.WithholdingTaxes[0]
			assert.Equal(t, "IRS", wt.WithholdingTaxType)
			assert.Equal(t, "Art. 101.º CIRS", wt.WithholdingTaxDescription)
			assert.Equal(t, "23.00", wt.WithholdingTaxAmount)
		}
		assert.Equal(t, "954.32", i.DocumentTotals.GrossTotal)
	})
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("Unsupported cases", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeySource: addon.SourceBillingIntegrated,
				}),
			},
			Charges:   []*bill.Charge{{}},
			Discounts: []*bill.Discount{{}},
			Lines: []*bill.Line{
				{
					Charges: []*bill.LineCharge{{}},
				},
			},
		}

		_, err := saft.BuildInvoice(inv, nil)
		assert.ErrorContains(t, err, "tax: (ext: pt-saft-source: value 'I' not allowed")
		assert.ErrorContains(t, err, "charges: not supported")
		assert.ErrorContains(t, err, "discounts: not supported")
		assert.ErrorContains(t, err, "lines: (0: (charges: not supported")
	})

	t.Run("Valid simplified invoice under 1000 EUR", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeSimplified,
				}),
			},
			Totals: &bill.Totals{
				Payable: num.MakeAmount(100000, 2), // Max amount affected by the rule
			},
		}

		_, err := saft.BuildInvoice(inv, nil)
		require.NoError(t, err)

		inv.Totals.Payable = num.MakeAmount(1, 2) // Min amount affected by the rule
		_, err = saft.BuildInvoice(inv, nil)
		require.NoError(t, err)
	})

	t.Run("Valid simplified invoice under 1000 EUR for identified customer", func(t *testing.T) {
		inv := &bill.Invoice{
			Customer: &org.Party{
				TaxID: &tax.Identity{
					Code:    "123456789",
					Country: "PT",
				},
			},
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeSimplified,
				}),
			},
			Lines: []*bill.Line{
				{
					Item: &org.Item{
						Ext: tax.ExtensionsOf(cbc.CodeMap{
							addon.ExtKeyProductType: addon.ProductTypeService,
						}),
					},
				},
			},
			Totals: &bill.Totals{
				Payable: num.MakeAmount(100000, 2), // Max amount affected by the rule
			},
		}

		_, err := saft.BuildInvoice(inv, nil)
		require.NoError(t, err)

		inv.Totals.Payable = num.MakeAmount(10001, 2) // Min amount affected by the rule
		_, err = saft.BuildInvoice(inv, nil)
		require.NoError(t, err)
	})

	t.Run("Valid simplified invoice under 1000 EUR for service", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeSimplified,
				}),
			},
			Lines: []*bill.Line{
				{
					Item: &org.Item{
						Ext: tax.ExtensionsOf(cbc.CodeMap{
							addon.ExtKeyProductType: addon.ProductTypeService,
						}),
					},
				},
			},
			Totals: &bill.Totals{
				Payable: num.MakeAmount(100000, 2), // Max amount affected by the rule
			},
		}

		_, err := saft.BuildInvoice(inv, nil)
		require.NoError(t, err)

		inv.Totals.Payable = num.MakeAmount(10001, 2) // Min amount affected by the rule
		_, err = saft.BuildInvoice(inv, nil)
		require.NoError(t, err)
	})

	t.Run("Valid simplified invoice under 1000 EUR for goods", func(t *testing.T) {
		inv := &bill.Invoice{
			Customer: nil,
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeSimplified,
				}),
			},
			Lines: []*bill.Line{
				{
					Item: &org.Item{
						Ext: tax.ExtensionsOf(cbc.CodeMap{
							addon.ExtKeyProductType: addon.ProductTypeGoods,
						}),
					},
				},
			},
			Totals: &bill.Totals{
				Payable: num.MakeAmount(100000, 2), // Max amount affected by the rule
			},
		}

		_, err := saft.BuildInvoice(inv, nil)
		require.NoError(t, err)

		inv.Totals.Payable = num.MakeAmount(10001, 2) // Min amount affected by the rule
		_, err = saft.BuildInvoice(inv, nil)
		require.NoError(t, err)
	})

	t.Run("Invalid simplified invoice over 1000 EUR", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeSimplified,
				}),
			},
			Lines: []*bill.Line{
				{
					Item: &org.Item{
						Ext: tax.ExtensionsOf(cbc.CodeMap{
							addon.ExtKeyProductType: addon.ProductTypeGoods,
						}),
					},
				},
			},
			Totals: &bill.Totals{
				Payable: num.MakeAmount(100001, 2), // Min amount affected by the rule
			},
		}

		_, err := saft.BuildInvoice(inv, nil)
		assert.ErrorContains(t, err, "payable: too high according to Artigo 40.º do CIVA")
	})

	t.Run("Invalid simplified invoice over 1000 EUR with exemption", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeSimplified,
				}),
			},
			Lines: []*bill.Line{
				{
					Taxes: tax.Set{
						{
							Category: tax.CategoryVAT,
							Ext: tax.ExtensionsOf(cbc.CodeMap{
								addon.ExtKeyTaxRate:   addon.TaxRateExempt,
								addon.ExtKeyExemption: "M10",
							}),
						},
					},
				},
			},
			Totals: &bill.Totals{
				Payable: num.MakeAmount(100001, 2), // Min amount over the 1000 EUR limit
			},
			Notes: []*org.Note{
				{
					Key:  org.NoteKeyLegal,
					Src:  addon.ExtKeyExemption,
					Code: "M10",
					Text: "IVA - regime de isenção / Artigo 57.º do CIVA",
				},
			},
		}

		_, err := saft.BuildInvoice(inv, nil)
		assert.ErrorContains(t, err, "payable: too high according to Artigo 40.º do CIVA")
	})

	t.Run("Missing supplier registration", func(t *testing.T) {
		inv := &bill.Invoice{
			Supplier: &org.Party{
				Name: "Test Supplier",
				TaxID: &tax.Identity{
					Country: l10n.PT.Tax(),
					Code:    "513912622",
				},
			},
		}

		_, err := saft.BuildInvoice(inv, nil)
		assert.ErrorContains(t, err, "registration: cannot be blank")
	})

	t.Run("Missing supplier registration office", func(t *testing.T) {
		inv := &bill.Invoice{
			Supplier: &org.Party{
				Name: "Test Supplier",
				TaxID: &tax.Identity{
					Country: l10n.PT.Tax(),
					Code:    "513912622",
				},
				Registration: &org.Registration{
					Capital: num.NewAmount(100000, 2),
				},
			},
		}

		_, err := saft.BuildInvoice(inv, nil)
		assert.ErrorContains(t, err, "office: cannot be blank")
	})

	t.Run("Missing supplier registration capital", func(t *testing.T) {
		inv := &bill.Invoice{
			Supplier: &org.Party{
				Name: "Test Supplier",
				TaxID: &tax.Identity{
					Country: l10n.PT.Tax(),
					Code:    "513912622",
				},
				Registration: &org.Registration{
					Office: "Test Office",
				},
			},
		}

		_, err := saft.BuildInvoice(inv, nil)
		assert.ErrorContains(t, err, "capital: cannot be blank")
	})

	t.Run("Supplier with complete registration", func(t *testing.T) {
		inv := &bill.Invoice{
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
		}

		_, err := saft.BuildInvoice(inv, nil)
		require.NoError(t, err)
	})

	t.Run("Supplier registration required by NIF prefix", func(t *testing.T) {
		tests := []struct {
			name     string
			nif      string
			required bool
		}{
			{"pessoa coletiva", "500000000", true},
			{"individual", "270044060", false},
			{"non-resident individual", "450000001", false},
			{"public body", "600000001", false},
			{"herança indivisa", "740000004", false},
			{"condomínio", "900000007", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				inv := &bill.Invoice{
					Supplier: &org.Party{
						Name: "Test Supplier",
						TaxID: &tax.Identity{
							Country: l10n.PT.Tax(),
							Code:    cbc.Code(tt.nif),
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
				}

				_, err := saft.BuildInvoice(inv, nil)
				if tt.required {
					assert.ErrorContains(t, err, "registration: cannot be blank")
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("Credit note with missing preceding reason", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeCreditNote,
				}),
			},
			Preceding: []*org.DocumentRef{
				{
					Series: "FT SERIES-A",
					Code:   "1233",
					// Missing Reason field
				},
			},
		}

		_, err := saft.BuildInvoice(inv, nil)
		assert.ErrorContains(t, err, "preceding: (0: (reason: cannot be blank")
	})

	t.Run("Debit note with missing preceding reason", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeDebitNote,
				}),
			},
			Preceding: []*org.DocumentRef{
				{
					Series: "FT SERIES-A",
					Code:   "1233",
					// Missing Reason field
				},
			},
		}

		_, err := saft.BuildInvoice(inv, nil)
		assert.ErrorContains(t, err, "preceding: (0: (reason: cannot be blank")
	})

	t.Run("Credit note with valid preceding reason", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeCreditNote,
				}),
			},
			Preceding: []*org.DocumentRef{
				{
					Series: "FT SERIES-A",
					Code:   "1233",
					Reason: "Credit note reason",
				},
			},
		}

		_, err := saft.BuildInvoice(inv, nil)
		require.NoError(t, err)
	})

	t.Run("Debit note with valid preceding reason", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeDebitNote,
				}),
			},
			Preceding: []*org.DocumentRef{
				{
					Series: "FT SERIES-A",
					Code:   "1233",
					Reason: "Debit note reason",
				},
			},
		}

		_, err := saft.BuildInvoice(inv, nil)
		require.NoError(t, err)
	})

	t.Run("Standard invoice with non-transport preceding should fail", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeStandard,
				}),
			},
			Preceding: []*org.DocumentRef{
				{
					Series: "FT SERIES-A", // FT is an invoice type, not a transport doc
					Code:   "1",
				},
			},
		}

		_, err := saft.BuildInvoice(inv, nil)
		assert.ErrorContains(t, err, "document type (FT) is not supported")
	})

	t.Run("Standard invoice with transport preceding should succeed", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeStandard,
				}),
			},
			Preceding: []*org.DocumentRef{
				{
					Series:    "GT SERIES-A",
					Code:      "1",
					IssueDate: cal.NewDate(2024, 8, 3),
				},
			},
		}

		_, err := saft.BuildInvoice(inv, nil)
		require.NoError(t, err)
	})

	t.Run("Standard invoice with transport preceding generates order references on lines", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeStandard,
				}),
			},
			Preceding: []*org.DocumentRef{
				{
					Series:    "GT SERIES-A",
					Code:      "1",
					IssueDate: cal.NewDate(2024, 8, 3),
				},
			},
			Lines: []*bill.Line{
				{Item: &org.Item{}, Order: "PO-2024-001"},
				{Item: &org.Item{}},
			},
		}

		i, err := saft.BuildInvoice(inv, &saft.BuildParams{
			Products: []*saft.Product{{}, {}},
		})
		require.NoError(t, err)
		require.Len(t, i.Lines, 2)

		// First line: per-line order + GT preceding
		require.Len(t, i.Lines[0].OrderReferences, 2)
		assert.Equal(t, "PO-2024-001", i.Lines[0].OrderReferences[0].OriginatingON)
		assert.Empty(t, i.Lines[0].OrderReferences[0].OrderDate)
		assert.Equal(t, "GT SERIES-A/1", i.Lines[0].OrderReferences[1].OriginatingON)
		assert.Equal(t, "2024-08-03", i.Lines[0].OrderReferences[1].OrderDate)

		// Second line: only GT preceding (no per-line order)
		require.Len(t, i.Lines[1].OrderReferences, 1)
		assert.Equal(t, "GT SERIES-A/1", i.Lines[1].OrderReferences[0].OriginatingON)
		assert.Equal(t, "2024-08-03", i.Lines[1].OrderReferences[0].OrderDate)
	})

	t.Run("FT with transport preceding and line order generates combined order references", func(t *testing.T) {
		// invoice-with-transport-preceding.json: FT with preceding GT/1 + line1.Order="NE SERIES-A/2"
		env, err := test.LoadEnvelope("invoice-with-transport-preceding.json")
		require.NoError(t, err)

		ginv := env.Extract().(*bill.Invoice)
		require.Len(t, ginv.Preceding, 1)

		prods := make([]*saft.Product, len(ginv.Lines))
		for i := range prods {
			prods[i] = &saft.Product{}
		}

		inv, err := saft.BuildInvoice(ginv, &saft.BuildParams{Products: prods})
		require.NoError(t, err)
		require.NotNil(t, inv)
		require.Len(t, inv.Lines, 2)

		// Line 1: line.Order ("NE SERIES-A/2") first, then document-level preceding (GT/1)
		line1 := inv.Lines[0]
		require.Len(t, line1.OrderReferences, 2)
		assert.Equal(t, "NE SERIES-A/2", line1.OrderReferences[0].OriginatingON)
		assert.Empty(t, line1.OrderReferences[0].OrderDate)
		assert.Equal(t, "GT SERIES-A/1", line1.OrderReferences[1].OriginatingON)
		assert.Equal(t, "2024-08-03", line1.OrderReferences[1].OrderDate)

		// Line 2: only document-level preceding (GT/1), no per-line order
		line2 := inv.Lines[1]
		require.Len(t, line2.OrderReferences, 1)
		assert.Equal(t, "GT SERIES-A/1", line2.OrderReferences[0].OriginatingON)
		assert.Equal(t, "2024-08-03", line2.OrderReferences[0].OrderDate)

		// FT transport preceding should NOT appear in References XML
		assert.Empty(t, line1.References)
		assert.Empty(t, line2.References)
	})

	t.Run("Credit note with transport preceding should fail", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeCreditNote,
				}),
			},
			Preceding: []*org.DocumentRef{
				{
					Series: "GT SERIES-A",
					Code:   "1",
					Reason: "Credit note reason",
				},
			},
		}

		_, err := saft.BuildInvoice(inv, nil)
		assert.ErrorContains(t, err, "document type (GT) is not supported")
	})

	t.Run("Credit note with no preceding documents should fail", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeCreditNote,
				}),
			},
			Lines: []*bill.Line{
				{
					Item: &org.Item{},
				},
			},
		}

		_, err := saft.BuildInvoice(inv, nil)
		assert.ErrorContains(t, err, "preceding: cannot be blank")
	})

	t.Run("Debit note with no preceding documents should fail", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeDebitNote,
				}),
			},
			Lines: []*bill.Line{
				{
					Item: &org.Item{},
				},
			},
		}

		_, err := saft.BuildInvoice(inv, nil)
		assert.ErrorContains(t, err, "preceding: cannot be blank")
	})

	t.Run("Corrective invoice with unsupported document type in preceding should fail", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeCreditNote,
				}),
			},
			Preceding: []*org.DocumentRef{
				{
					Series: "ND SERIES-A", // ND is not a correctable document type
					Code:   "1233",
					Reason: "Credit note reason",
				},
			},
		}

		_, err := saft.BuildInvoice(inv, nil)
		assert.ErrorContains(t, err, "document type (ND) is not supported")
	})

	t.Run("Debit note with unsupported document type in preceding should fail", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeDebitNote,
				}),
			},
			Preceding: []*org.DocumentRef{
				{
					Series: "RG SERIES-A", // RG is not a correctable document type
					Code:   "1233",
					Reason: "Debit note reason",
				},
			},
		}

		_, err := saft.BuildInvoice(inv, nil)
		assert.ErrorContains(t, err, "document type (RG) is not supported")
	})

	t.Run("Valid credit note", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeCreditNote,
				}),
			},
			Preceding: []*org.DocumentRef{
				{
					Series: "FT SERIES-A",
					Code:   "1233",
					Reason: "Credit note reason",
				},
			},
		}

		_, err := saft.BuildInvoice(inv, nil)
		require.NoError(t, err)
	})

	t.Run("Valid debit note", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeDebitNote,
				}),
			},
			Preceding: []*org.DocumentRef{
				{
					Series: "FR SERIES-A",
					Code:   "1233",
					Reason: "Debit note reason",
				},
			},
		}

		_, err := saft.BuildInvoice(inv, nil)
		require.NoError(t, err)
	})
}

func TestIsCorrection(t *testing.T) {
	t.Run("Standard invoice should not be correction", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeStandard,
				}),
			},
		}

		// We can't test the private function directly, but we can test its behavior
		// through the validation logic
		_, err := saft.BuildInvoice(inv, &saft.BuildParams{
			Products: []*saft.Product{{}},
		})
		require.NoError(t, err)
	})

	t.Run("Simplified invoice should not be correction", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeSimplified,
				}),
			},
			Lines: []*bill.Line{
				{
					Item: &org.Item{},
				},
			},
		}

		_, err := saft.BuildInvoice(inv, &saft.BuildParams{
			Products: []*saft.Product{{}},
		})
		require.NoError(t, err)
	})

	t.Run("Credit note should be correction", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeCreditNote,
				}),
			},
			Preceding: []*org.DocumentRef{
				{
					Series: "FT SERIES-A",
					Code:   "1233",
					Reason: "Credit note reason",
				},
			},
			Lines: []*bill.Line{
				{
					Item: &org.Item{},
				},
			},
		}

		_, err := saft.BuildInvoice(inv, &saft.BuildParams{
			Products: []*saft.Product{{}},
		})
		require.NoError(t, err)
	})

	t.Run("Debit note should be correction", func(t *testing.T) {
		inv := &bill.Invoice{
			Tax: &bill.Tax{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyInvoiceType: addon.InvoiceTypeDebitNote,
				}),
			},
			Preceding: []*org.DocumentRef{
				{
					Series: "FT SERIES-A",
					Code:   "1233",
					Reason: "Debit note reason",
				},
			},
			Lines: []*bill.Line{
				{
					Item: &org.Item{},
				},
			},
		}

		_, err := saft.BuildInvoice(inv, &saft.BuildParams{
			Products: []*saft.Product{{}},
		})
		require.NoError(t, err)
	})

	t.Run("Invoice without tax should not be correction", func(t *testing.T) {
		inv := &bill.Invoice{
			Lines: []*bill.Line{
				{
					Item: &org.Item{},
				},
			},
		}

		_, err := saft.BuildInvoice(inv, &saft.BuildParams{
			Products: []*saft.Product{{}},
		})
		require.NoError(t, err)
	})
}

func TestInvoiceSetStatus(t *testing.T) {
	config.Now = func() time.Time {
		return time.Date(2024, 3, 1, 10, 20, 30, 0, time.UTC)
	}

	i := &saft.Invoice{}
	i.SetStatus("N", "Reason", "SOURCE123")

	if ds := i.DocumentStatus; assert.NotNil(t, ds) {
		assert.Equal(t, "N", ds.InvoiceStatus)
		assert.Equal(t, "Reason", ds.Reason)
		assert.Equal(t, "2024-03-01T10:20:30", ds.InvoiceStatusDate)
		assert.Equal(t, "SOURCE123", ds.SourceID)
		assert.Equal(t, "P", ds.SourceBilling)
	}
}

func TestInvoiceGetStatus(t *testing.T) {
	t.Run("Empty case", func(t *testing.T) {
		i := &saft.Invoice{}
		assert.Equal(t, "", i.GetStatus())
	})

	t.Run("With status", func(t *testing.T) {
		i := &saft.Invoice{
			DocumentStatus: &saft.InvoiceStatus{
				InvoiceStatus: "N",
			},
		}
		assert.Equal(t, "N", i.GetStatus())
	})
}

func TestInvoiceGetStatusDate(t *testing.T) {
	t.Run("Empty case", func(t *testing.T) {
		i := &saft.Invoice{}
		assert.Equal(t, "", i.GetStatusDate())
	})

	t.Run("With status date", func(t *testing.T) {
		i := &saft.Invoice{
			DocumentStatus: &saft.InvoiceStatus{
				InvoiceStatusDate: "2024-03-01T10:20:30",
			},
		}
		assert.Equal(t, "2024-03-01T10:20:30", i.GetStatusDate())
	})
}

// deductionTestInvoice returns a standard invoice with one line that deducts two advance
// payment invoices and also has a transport document as preceding.
func deductionTestInvoice() *bill.Invoice {
	vat := tax.Set{{
		Category: tax.CategoryVAT,
		Percent:  num.NewPercentage(230, 3),
		Ext: tax.ExtensionsOf(cbc.CodeMap{
			pt.ExtKeyRegion:     "PT",
			addon.ExtKeyTaxRate: addon.TaxRateNormal,
		}),
	}}
	return &bill.Invoice{
		Tax: &bill.Tax{
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				addon.ExtKeyInvoiceType: addon.InvoiceTypeStandard,
			}),
		},
		ValueDate: cal.NewDate(2024, 8, 3),
		Preceding: []*org.DocumentRef{
			{Series: "FT SERIES-A", Code: "1", IssueDate: cal.NewDate(2024, 8, 1)},
			{Series: "GT SERIES-A", Code: "1", IssueDate: cal.NewDate(2024, 8, 2)},
			{Series: "FT SERIES-A", Code: "2", IssueDate: cal.NewDate(2024, 8, 2)},
		},
		Lines: []*bill.Line{
			{
				Index:    1,
				Quantity: num.MakeAmount(3, 0),
				Item:     &org.Item{Price: num.NewAmount(10000, 2)},
				Taxes:    vat,
				Total:    num.NewAmount(30000, 2),
			},
		},
		Discounts: []*bill.Discount{
			{
				Key:    addon.DiscountKeyDeduction,
				Reason: "Deposit for 1 hotel room",
				Amount: num.MakeAmount(5000, 2),
				Taxes:  vat,
				Ext:    tax.ExtensionsOf(cbc.CodeMap{addon.ExtKeyReference: "FT SERIES-A/1"}),
			},
			{
				Key:    addon.DiscountKeyDeduction,
				Reason: "Deposit for 2 hotel rooms",
				Amount: num.MakeAmount(10000, 2),
				Taxes:  vat,
				Ext:    tax.ExtensionsOf(cbc.CodeMap{addon.ExtKeyReference: "FT SERIES-A/2"}),
			},
		},
	}
}

func TestBuildInvoiceDeductions(t *testing.T) {
	t.Run("deduction lines", func(t *testing.T) {
		inv := deductionTestInvoice()
		prods := append(saft.BuildProducts(inv.Lines), saft.BuildDeductionProducts(inv)...)

		i, err := saft.BuildInvoice(inv, &saft.BuildParams{Products: prods})
		require.NoError(t, err)
		require.Len(t, i.Lines, 3)

		item := i.Lines[0]
		assert.Equal(t, "300.00", item.CreditAmount)
		assert.Empty(t, item.References)
		require.Len(t, item.OrderReferences, 1)
		assert.Equal(t, "GT SERIES-A/1", item.OrderReferences[0].OriginatingON)

		for n, exp := range []struct{ ref, date, desc, amount string }{
			{"FT SERIES-A/1", "2024-08-01", "Deposit for 1 hotel room", "50.00"},
			{"FT SERIES-A/2", "2024-08-02", "Deposit for 2 hotel rooms", "100.00"},
		} {
			line := i.Lines[n+1]
			assert.Equal(t, n+2, line.LineNumber)
			require.Len(t, line.OrderReferences, 1)
			assert.Equal(t, exp.ref, line.OrderReferences[0].OriginatingON)
			assert.Equal(t, exp.date, line.OrderReferences[0].OrderDate)
			assert.Equal(t, prods[n+1].ProductCode, line.ProductCode)
			assert.Equal(t, exp.desc, line.ProductDesc)
			assert.Equal(t, exp.desc, line.Description)
			assert.Equal(t, "1", line.Quantity)
			assert.Equal(t, "one", line.UnitOfMeasure)
			assert.Equal(t, exp.amount, line.UnitPrice)
			assert.Equal(t, "2024-08-03", line.TaxPointDate)
			assert.Equal(t, exp.amount, line.DebitAmount)
			assert.Empty(t, line.CreditAmount)
			assert.Empty(t, line.References)
			if assert.NotNil(t, line.Tax) {
				assert.Equal(t, "IVA", line.Tax.TaxType)
				assert.Equal(t, "PT", line.Tax.TaxCountryRegion)
				assert.Equal(t, "NOR", line.Tax.TaxCode)
				assert.Equal(t, "23.0", line.Tax.TaxPercentage)
			}
		}
	})

	t.Run("reference defaults to the only preceding document", func(t *testing.T) {
		inv := deductionTestInvoice()
		inv.Preceding = inv.Preceding[:1]
		inv.Discounts = inv.Discounts[:1]
		inv.Discounts[0].Ext = tax.Extensions{}
		prods := append(saft.BuildProducts(inv.Lines), saft.BuildDeductionProducts(inv)...)

		i, err := saft.BuildInvoice(inv, &saft.BuildParams{Products: prods})
		require.NoError(t, err)
		require.Len(t, i.Lines, 2)
		assert.Empty(t, i.Lines[0].References)
		require.Len(t, i.Lines[1].OrderReferences, 1)
		assert.Equal(t, "FT SERIES-A/1", i.Lines[1].OrderReferences[0].OriginatingON)
	})

	t.Run("exempt deduction", func(t *testing.T) {
		inv := deductionTestInvoice()
		inv.Discounts[0].Taxes = tax.Set{{
			Category: tax.CategoryVAT,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				addon.ExtKeyTaxRate:   addon.TaxRateExempt,
				addon.ExtKeyExemption: "M16",
			}),
		}}
		prods := append(saft.BuildProducts(inv.Lines), saft.BuildDeductionProducts(inv)...)

		i, err := saft.BuildInvoice(inv, &saft.BuildParams{Products: prods})
		require.NoError(t, err)
		assert.Equal(t, "M16", i.Lines[1].TaxExemptionCode)
		assert.NotEmpty(t, i.Lines[1].TaxExemptionReason)
	})

	t.Run("products missing for deductions", func(t *testing.T) {
		inv := deductionTestInvoice()

		i, err := saft.BuildInvoice(inv, &saft.BuildParams{Products: saft.BuildProducts(inv.Lines)})
		require.NoError(t, err)
		assert.Nil(t, i.Lines)
	})

	t.Run("convert adds deduction products", func(t *testing.T) {
		config.Now = func() time.Time {
			return time.Date(2024, 8, 3, 0, 0, 0, 0, pt.New().TimeLocation())
		}

		env, err := test.LoadEnvelope("invoice-deductions.json")
		require.NoError(t, err)

		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		})
		require.NoError(t, err)

		res, err := saft.Convert(env, "SOURCE_ID", &config.Config{SigningKeyVersion: "3"})
		require.NoError(t, err)

		require.Len(t, res.Products, 3)
		assert.Equal(t, "Deposit for 1 hotel room", res.Products[1].ProductDesc)
		assert.Equal(t, "O", res.Products[1].ProductType)
		assert.Equal(t, "Deposit for 2 hotel rooms", res.Products[2].ProductDesc)

		inv := res.Document.(*saft.Invoice)
		require.Len(t, inv.Lines, 3)
		assert.Equal(t, "150.00", inv.DocumentTotals.NetTotal)
		assert.Equal(t, "34.50", inv.DocumentTotals.TaxPayable)
		assert.Equal(t, "184.50", inv.DocumentTotals.GrossTotal)
	})
}

func TestInvoiceDeductionValidation(t *testing.T) {
	t.Run("deducted invoices allowed as preceding", func(t *testing.T) {
		_, err := saft.BuildInvoice(deductionTestInvoice(), nil)
		require.NoError(t, err)
	})

	t.Run("invoice preceding not deducted", func(t *testing.T) {
		inv := deductionTestInvoice()
		inv.Discounts = inv.Discounts[:1]

		_, err := saft.BuildInvoice(inv, nil)
		assert.ErrorContains(t, err, "preceding: (2: document type (FT) is not supported.)")
	})

	t.Run("other document discounts", func(t *testing.T) {
		inv := deductionTestInvoice()
		inv.Discounts = append(inv.Discounts, &bill.Discount{Reason: "Promo", Amount: num.MakeAmount(1000, 2)})

		_, err := saft.BuildInvoice(inv, nil)
		assert.ErrorContains(t, err, "discounts: not supported, use line discounts instead")
	})
}
