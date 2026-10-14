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
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildPayment(t *testing.T) {
	config.Now = func() time.Time {
		return time.Date(2024, 12, 23, 23, 30, 30, 0, time.UTC)
	}

	t.Run("Empty cases", func(t *testing.T) {
		p, err := saft.BuildPayment(nil, nil)
		require.NoError(t, err)
		assert.Nil(t, p)

		p, err = saft.BuildPayment(&bill.Payment{}, nil)
		require.NoError(t, err)
		assert.Equal(t, "", p.PaymentRefNo)
	})

	t.Run("Unsupported cases", func(t *testing.T) {
		r := &bill.Payment{
			Currency: "USD",
			Lines: []*bill.PaymentLine{
				{
					Amount: num.MakeAmount(12345, 2),
				},
			},
		}

		_, err := saft.BuildPayment(r, nil)
		assert.ErrorContains(t, err, "currency: must be EUR")
	})

	t.Run("General case", func(t *testing.T) {
		r := &bill.Payment{
			Identify:  uuid.IdentifyParse("b48d3ef4-df19-11ef-8d60-325096b39f47"),
			Series:    "RG SERIES-A",
			Code:      "1234",
			IssueDate: cal.MakeDate(2024, 12, 23),
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				addon.ExtKeySource:      addon.SourceBillingProduced,
				addon.ExtKeyPaymentType: addon.PaymentTypeOther,
			}),
			Notes: []*org.Note{
				{
					Key:  org.NoteKeyGeneral,
					Text: "Description of the payment",
				},
			},
			Methods: []*pay.Record{
				{
					Key:    "cash",
					Amount: num.MakeAmount(13667, 2),
					Ext: tax.ExtensionsOf(cbc.CodeMap{
						addon.ExtKeyPaymentMeans: "NU",
					}),
				},
			},
			Lines: []*bill.PaymentLine{
				{
					Index: 1,
					Document: &org.DocumentRef{
						Series:      "FT SERIES-A",
						Code:        "123",
						IssueDate:   cal.NewDate(2024, 2, 29),
						Description: "Description of the document",
					},
					Tax: &tax.Total{
						Categories: []*tax.CategoryTotal{
							{
								Code: tax.CategoryVAT,
								Rates: []*tax.RateTotal{
									{
										Ext: tax.ExtensionsOf(cbc.CodeMap{
											addon.ExtKeyTaxRate: "NOR",
											pt.ExtKeyRegion:     "PT",
										}),
										Base:    num.MakeAmount(11111, 2),
										Percent: num.NewPercentage(23, 2),
										Amount:  num.MakeAmount(2556, 2),
									},
								},
								Amount: num.MakeAmount(2556, 2),
							},
						},
						Sum: num.MakeAmount(2556, 2),
					},
					Amount: num.MakeAmount(13667, 2),
				},
			},
			Total: num.MakeAmount(13667, 2),
		}

		p, err := saft.BuildPayment(r, &saft.BuildParams{
			ATCUD:           "AAJFJHNC99-1234",
			SourceID:        "SOURCEID",
			Customer:        &saft.Customer{CustomerID: "CUSTID"},
			SystemEntryDate: "2024-03-01T10:20:40",
		})
		require.NoError(t, err)

		assert.Equal(t, "RG SERIES-A/1234", p.PaymentRefNo)
		assert.Equal(t, "AAJFJHNC99-1234", p.ATCUD)
		assert.Equal(t, 12, p.Period)
		assert.Equal(t, "2024-12-23", p.TransactionDate)
		assert.Equal(t, "RG", p.PaymentType)
		assert.Equal(t, "Description of the payment", p.Description)
		assert.Equal(t, "b48d3ef4-df19-11ef-8d60-325096b39f47", p.SystemID)
		if ds := p.DocumentStatus; assert.NotNil(t, ds) {
			assert.Equal(t, "N", ds.PaymentStatus)
			assert.Equal(t, "2024-03-01T10:20:40", ds.PaymentStatusDate)
			assert.Equal(t, "SOURCEID", ds.SourceID)
			assert.Equal(t, "P", ds.SourcePayment)
		}
		if assert.Len(t, p.PaymentMethods, 1) {
			pm := p.PaymentMethods[0]
			assert.Equal(t, "NU", pm.PaymentMechanism)
			assert.Equal(t, "136.67", pm.PaymentAmount)
			assert.Equal(t, "2024-12-23", pm.PaymentDate)
		}
		assert.Equal(t, "SOURCEID", p.SourceID)
		assert.Equal(t, "2024-03-01T10:20:40", p.SystemEntryDate)
		if assert.Len(t, p.Lines, 1) {
			l := p.Lines[0]
			assert.Equal(t, 1, l.LineNumber)
			if assert.Len(t, l.SourceDocumentIDs, 1) {
				sd := p.Lines[0].SourceDocumentIDs[0]
				assert.Equal(t, "FT SERIES-A/123", sd.OriginatingON)
				assert.Equal(t, "2024-02-29", sd.InvoiceDate)
				assert.Equal(t, "Description of the document", sd.Description)
			}
			assert.Equal(t, "111.11", l.CreditAmount)
			if tax := l.Tax; assert.NotNil(t, tax) {
				assert.Equal(t, "IVA", tax.TaxType)
				assert.Equal(t, "PT", tax.TaxCountryRegion)
				assert.Equal(t, "NOR", tax.TaxCode)
				assert.Equal(t, "23", tax.TaxPercentage)
			}
		}
		if dt := p.DocumentTotals; assert.NotNil(t, dt) {
			assert.Equal(t, "136.67", dt.GrossTotal)
			assert.Equal(t, "111.11", dt.NetTotal)
			assert.Equal(t, "25.56", dt.TaxPayable)
		}
	})

	t.Run("Advanced VAT rates case", func(t *testing.T) {
		r := &bill.Payment{
			Lines: []*bill.PaymentLine{
				{
					Amount: num.MakeAmount(13667, 2),
					Tax: &tax.Total{
						Categories: []*tax.CategoryTotal{
							{
								Code: tax.CategoryVAT,
								Rates: []*tax.RateTotal{
									{
										Base:    num.MakeAmount(11111, 2),
										Percent: num.NewPercentage(23, 2),
										Amount:  num.MakeAmount(2556, 2),
									},
								},
								Amount: num.MakeAmount(2556, 2),
							},
						},
						Sum: num.MakeAmount(2556, 2),
					},
				},
				{
					Amount: num.MakeAmount(22222, 2),
					Tax: &tax.Total{
						Categories: []*tax.CategoryTotal{
							{
								Code: tax.CategoryVAT,
								Rates: []*tax.RateTotal{
									{
										Ext: tax.ExtensionsOf(cbc.CodeMap{
											addon.ExtKeyExemption: "M15",
										}),
										Base:   num.MakeAmount(22222, 2),
										Amount: num.AmountZero,
									},
								},
								Amount: num.AmountZero,
							},
						},
						Sum: num.AmountZero,
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

		p, err := saft.BuildPayment(r, nil)
		require.NoError(t, err)

		if assert.Len(t, p.Lines, 2) {
			assert.Equal(t, "111.11", p.Lines[0].CreditAmount)
			assert.Equal(t, "", p.Lines[0].TaxExemptionCode)
			assert.Equal(t, "", p.Lines[0].TaxExemptionReason)
			assert.Equal(t, "222.22", p.Lines[1].CreditAmount)
			assert.Equal(t, "M15", p.Lines[1].TaxExemptionCode)
			assert.Equal(t, "Regime da margem de lucro - Objetos de coleção e antiguida", p.Lines[1].TaxExemptionReason)
		}
	})

	t.Run("Refund case", func(t *testing.T) {
		r := &bill.Payment{
			Lines: []*bill.PaymentLine{
				{
					Amount: num.MakeAmount(13667, 2),
					Refund: true,
					Tax: &tax.Total{
						Categories: []*tax.CategoryTotal{
							{
								Code: tax.CategoryVAT,
								Rates: []*tax.RateTotal{
									{
										Ext: tax.ExtensionsOf(cbc.CodeMap{
											addon.ExtKeyTaxRate: "NOR",
											pt.ExtKeyRegion:     "PT",
										}),
										Base:    num.MakeAmount(11111, 2),
										Percent: num.NewPercentage(23, 2),
										Amount:  num.MakeAmount(2556, 2),
									},
								},
								Amount: num.MakeAmount(2556, 2),
							},
						},
						Sum: num.MakeAmount(2556, 2),
					},
				},
			},
		}

		p, err := saft.BuildPayment(r, nil)
		require.NoError(t, err)

		assert.Equal(t, "111.11", p.Lines[0].DebitAmount)
	})

	t.Run("1 decimal gross total", func(t *testing.T) {
		r := &bill.Payment{
			Total: num.MakeAmount(12345, 1),
		}

		p, err := saft.BuildPayment(r, nil)
		require.NoError(t, err)

		assert.Equal(t, "1234.50", p.DocumentTotals.GrossTotal)
	})

	t.Run("No decimals gross total", func(t *testing.T) {
		r := &bill.Payment{
			Total: num.MakeAmount(12345, 0),
		}

		p, err := saft.BuildPayment(r, nil)
		require.NoError(t, err)

		assert.Equal(t, "12345.00", p.DocumentTotals.GrossTotal)
	})

	t.Run("Withholding taxes", func(t *testing.T) {
		r := &bill.Payment{
			Lines: []*bill.PaymentLine{
				{
					Amount: num.MakeAmount(13667, 2),
					Tax: &tax.Total{
						Categories: []*tax.CategoryTotal{
							{
								Code: tax.CategoryVAT,
								Rates: []*tax.RateTotal{
									{Base: num.MakeAmount(11111, 2), Percent: num.NewPercentage(23, 2), Amount: num.MakeAmount(2556, 2)},
								},
								Amount: num.MakeAmount(2556, 2),
							},
							{
								Code:     pt.TaxCategoryIRS,
								Retained: true,
								Amount:   num.MakeAmount(2300, 2),
							},
						},
						Sum: num.MakeAmount(2556, 2),
					},
				},
			},
		}

		p, err := saft.BuildPayment(r, nil)
		require.NoError(t, err)

		if assert.Len(t, p.WithholdingTaxes, 1) {
			wt := p.WithholdingTaxes[0]
			assert.Equal(t, "IRS", wt.WithholdingTaxType)
			assert.Equal(t, "", wt.WithholdingTaxDescription) // payments carry no description
			assert.Equal(t, "23.00", wt.WithholdingTaxAmount)
		}
	})
}

func TestPaymentValidation(t *testing.T) {
	t.Run("Missing supplier registration", func(t *testing.T) {
		p := &bill.Payment{
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
			},
		}

		_, err := saft.BuildPayment(p, nil)
		assert.ErrorContains(t, err, "registration: cannot be blank")
	})

	t.Run("Missing supplier registration office", func(t *testing.T) {
		p := &bill.Payment{
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
					Capital: num.NewAmount(100000, 2),
				},
			},
		}

		_, err := saft.BuildPayment(p, nil)
		assert.ErrorContains(t, err, "office: cannot be blank")
	})

	t.Run("Missing supplier registration capital", func(t *testing.T) {
		p := &bill.Payment{
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
					Office: "Test Office",
				},
			},
		}

		_, err := saft.BuildPayment(p, nil)
		assert.ErrorContains(t, err, "capital: cannot be blank")
	})

	t.Run("Rate amount not consistent with base", func(t *testing.T) {
		p := &bill.Payment{
			Lines: []*bill.PaymentLine{
				{
					Amount: num.MakeAmount(13000, 2),
					Tax: &tax.Total{
						Categories: []*tax.CategoryTotal{
							{
								Code: tax.CategoryVAT,
								Rates: []*tax.RateTotal{
									{
										Base:    num.MakeAmount(10000, 2),
										Percent: num.NewPercentage(23, 2),
										Amount:  num.MakeAmount(3000, 2), // should be 23.00, not 30.00
									},
								},
								Amount: num.MakeAmount(3000, 2),
							},
						},
						Sum: num.MakeAmount(3000, 2),
					},
				},
			},
		}

		_, err := saft.BuildPayment(p, nil)
		assert.ErrorContains(t, err, "rate amount is not the percentage of the base")
	})

	t.Run("Line amount not consistent with base and percent", func(t *testing.T) {
		p := &bill.Payment{
			Lines: []*bill.PaymentLine{
				{
					Amount: num.MakeAmount(15000, 2), // should be 123.00, not 150.00
					Tax: &tax.Total{
						Categories: []*tax.CategoryTotal{
							{
								Code: tax.CategoryVAT,
								Rates: []*tax.RateTotal{
									{
										Base:    num.MakeAmount(10000, 2),
										Percent: num.NewPercentage(23, 2),
										Amount:  num.MakeAmount(2300, 2),
									},
								},
								Amount: num.MakeAmount(2300, 2),
							},
						},
						Sum: num.MakeAmount(2300, 2),
					},
				},
			},
		}

		_, err := saft.BuildPayment(p, nil)
		assert.ErrorContains(t, err, "line amount is not consistent with the tax base and percentage")
	})

	t.Run("Valid line amounts", func(t *testing.T) {
		p := &bill.Payment{
			Lines: []*bill.PaymentLine{
				{
					Amount: num.MakeAmount(500, 2),
					Tax: &tax.Total{
						Categories: []*tax.CategoryTotal{
							{
								Code: tax.CategoryVAT,
								Rates: []*tax.RateTotal{
									{
										Base:    num.MakeAmount(407, 2),
										Percent: num.NewPercentage(23, 2),
										Amount:  num.MakeAmount(94, 2),
									},
								},
								Amount: num.MakeAmount(94, 2),
							},
						},
						Sum: num.MakeAmount(94, 2),
					},
				},
			},
		}

		_, err := saft.BuildPayment(p, nil)
		assert.NoError(t, err)
	})

	t.Run("Supplier with complete registration", func(t *testing.T) {
		p := &bill.Payment{
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

		_, err := saft.BuildPayment(p, nil)
		require.NoError(t, err)
	})
}

func TestPaymentSetStatus(t *testing.T) {
	config.Now = func() time.Time {
		return time.Date(2024, 3, 1, 10, 20, 30, 0, time.UTC)
	}

	p := &saft.Payment{}
	p.SetStatus("N", "Reason", "SOURCE123")

	if ds := p.DocumentStatus; assert.NotNil(t, ds) {
		assert.Equal(t, "N", ds.PaymentStatus)
		assert.Equal(t, "Reason", ds.Reason)
		assert.Equal(t, "2024-03-01T10:20:30", ds.PaymentStatusDate)
		assert.Equal(t, "SOURCE123", ds.SourceID)
		assert.Equal(t, "P", ds.SourcePayment)
	}
}

func TestPaymentGetStatus(t *testing.T) {
	t.Run("Empty case", func(t *testing.T) {
		p := &saft.Payment{}
		assert.Equal(t, "", p.GetStatus())
	})

	t.Run("With status", func(t *testing.T) {
		p := &saft.Payment{
			DocumentStatus: &saft.PaymentStatus{
				PaymentStatus: "N",
			},
		}
		assert.Equal(t, "N", p.GetStatus())
	})
}

func TestPaymentGetStatusDate(t *testing.T) {
	t.Run("Empty case", func(t *testing.T) {
		p := &saft.Payment{}
		assert.Equal(t, "", p.GetStatusDate())
	})

	t.Run("With status date", func(t *testing.T) {
		p := &saft.Payment{
			DocumentStatus: &saft.PaymentStatus{
				PaymentStatusDate: "2024-03-01T10:20:30",
			},
		}
		assert.Equal(t, "2024-03-01T10:20:30", p.GetStatusDate())
	})
}
